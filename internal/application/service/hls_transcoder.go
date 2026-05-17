package service

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/devlucas-java/luca-s3/internal/domain/enums"
	minioclient "github.com/devlucas-java/luca-s3/internal/infrastructure/minio"
	"github.com/devlucas-java/luca-s3/pkg/logger"
)

const (
	segmentDuration  = 2
	segmentTolerance = 2
)

type ProgressFunc func(segmentsDone int)

func transcodeToHLS(
	ctx context.Context,
	minio *minioclient.Client,
	videoID string,
	resolution enums.Resolution,
	inputPath string,
	workDir string,
	startSegment int,
	totalSegments int,
	onProgress ProgressFunc,
	is360 bool,
) (int, error) {
	log := logger.Instance()

	resDir := filepath.Join(workDir, resolution.String())
	if err := os.MkdirAll(resDir, 0755); err != nil {
		return 0, fmt.Errorf("mkdir: %w", err)
	}

	playlistPath := filepath.Join(resDir, "playlist.m3u8")
	segmentPattern := filepath.Join(resDir, "segment_%05d.ts")

	// For 360° video: preserve spherical metadata and use scale2ref to keep
	// the equirectangular projection intact. The key flags are:
	//   -map_metadata 0        — copy all metadata (including spherical box)
	//   -movflags +faststart   — needed for streaming
	//   scale filter with -2   — keeps exact 2:1 ratio for equirectangular
	args := []string{
		"-y",
		"-i", inputPath,
	}

	if is360 {
		// Copy spherical metadata from input to output
		args = append(args, "-map_metadata", "0")
		// Scale keeping exact width, height auto-calculated to preserve 2:1 ratio
		args = append(args,
			"-vf", fmt.Sprintf("scale=%d:-2", resolution.Width()),
		)
	} else {
		args = append(args, "-vf", "scale="+resolution.AspectScale())
	}

	args = append(args,
		"-c:v", "libx264",
		"-crf", "23",
		"-preset", "fast",
		"-c:a", "aac",
		"-b:a", "128k",
		"-hls_time", fmt.Sprintf("%d", segmentDuration),
		"-hls_playlist_type", "vod",
		"-hls_segment_filename", segmentPattern,
		"-start_number", fmt.Sprintf("%d", startSegment),
		"-progress", "pipe:2",
		playlistPath,
	)

	cmd := exec.CommandContext(ctx, "ffmpeg", args...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return 0, fmt.Errorf("stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return 0, fmt.Errorf("ffmpeg start: %w", err)
	}

	uploadedSegments := startSegment
	go func() {
		uploadedSegments = watchAndUploadSegments(
			ctx, minio, videoID, resolution,
			resDir, startSegment, totalSegments, onProgress,
		)
	}()

	go drainFFmpegProgress(stderr, log)

	if err := cmd.Wait(); err != nil {
		return uploadedSegments - startSegment, fmt.Errorf("ffmpeg: %w", err)
	}

	finalCount, err := uploadRemainingHLS(ctx, minio, videoID, resolution, resDir, uploadedSegments)
	if err != nil {
		return finalCount, err
	}

	total := finalCount - startSegment
	log.Infof("hls: %s/%s — %d segments uploaded", videoID, resolution, total)
	return total, nil
}

func watchAndUploadSegments(
	ctx context.Context,
	minio *minioclient.Client,
	videoID string,
	resolution enums.Resolution,
	resDir string,
	startSegment int,
	totalSegments int,
	onProgress ProgressFunc,
) int {
	uploaded := startSegment
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return uploaded
		case <-ticker.C:
			entries, err := os.ReadDir(resDir)
			if err != nil {
				continue
			}
			for _, entry := range entries {
				if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".ts") {
					continue
				}
				numStr := strings.TrimPrefix(entry.Name(), "segment_")
				numStr = strings.TrimSuffix(numStr, ".ts")
				num, err := strconv.Atoi(numStr)
				if err != nil || num < uploaded {
					continue
				}

				localPath := filepath.Join(resDir, entry.Name())
				minioPath := fmt.Sprintf("%s/hls/%s/%s", videoID, resolution.String(), entry.Name())

				if err := minio.UploadFile(ctx, minioPath, localPath, "video/mp2t"); err != nil {
					continue
				}

				uploaded = num + 1
				if onProgress != nil {
					onProgress(uploaded - startSegment)
				}
			}
		}
	}
}

func uploadRemainingHLS(
	ctx context.Context,
	minio *minioclient.Client,
	videoID string,
	resolution enums.Resolution,
	resDir string,
	alreadyUploaded int,
) (int, error) {
	entries, err := os.ReadDir(resDir)
	if err != nil {
		return alreadyUploaded, fmt.Errorf("readdir: %w", err)
	}

	maxSeg := alreadyUploaded
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		localPath := filepath.Join(resDir, entry.Name())
		minioPath := fmt.Sprintf("%s/hls/%s/%s", videoID, resolution.String(), entry.Name())

		contentType := "application/vnd.apple.mpegurl"
		if strings.HasSuffix(entry.Name(), ".ts") {
			contentType = "video/mp2t"
			numStr := strings.TrimPrefix(entry.Name(), "segment_")
			numStr = strings.TrimSuffix(numStr, ".ts")
			if num, err := strconv.Atoi(numStr); err == nil && num >= alreadyUploaded {
				_ = minio.UploadFile(ctx, minioPath, localPath, contentType)
				if num+1 > maxSeg {
					maxSeg = num + 1
				}
				continue
			}
		}
		_ = minio.UploadFile(ctx, minioPath, localPath, contentType)
	}
	return maxSeg, nil
}

func drainFFmpegProgress(r io.Reader, log interface{ Debugf(string, ...any) }) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
	}
}

func getLastSegmentNumber(ctx context.Context, minio *minioclient.Client, videoID string, resolution enums.Resolution) int {
	prefix := fmt.Sprintf("%s/hls/%s/segment_", videoID, resolution.String())
	objects := minio.ListObjects(ctx, prefix)

	max := -1
	for _, obj := range objects {
		name := filepath.Base(obj)
		if !strings.HasPrefix(name, "segment_") || !strings.HasSuffix(name, ".ts") {
			continue
		}
		numStr := strings.TrimSuffix(strings.TrimPrefix(name, "segment_"), ".ts")
		if num, err := strconv.Atoi(numStr); err == nil && num > max {
			max = num
		}
	}
	return max
}

func calcTotalSegments(durationSeconds float64) int {
	return int(math.Ceil(durationSeconds / segmentDuration))
}

func countSegments(ctx context.Context, minio *minioclient.Client, videoID string, resolution enums.Resolution) int {
	prefix := fmt.Sprintf("%s/hls/%s/segment_", videoID, resolution.String())
	objects := minio.ListObjects(ctx, prefix)
	count := 0
	for _, obj := range objects {
		if strings.HasSuffix(obj, ".ts") {
			count++
		}
	}
	return count
}
