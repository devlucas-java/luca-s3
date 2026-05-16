package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/devlucas-java/luca-s3/internal/domain/enums"
	minioclient "github.com/devlucas-java/luca-s3/internal/infrastructure/minio"
	"github.com/devlucas-java/luca-s3/pkg/logger"
)

const segmentDuration = 2 // segundos por segmento

// transcodeToHLS gera segmentos HLS de 2s para uma resolução específica.
// Retoma de onde parou verificando segmentos existentes no MinIO.
func transcodeToHLS(
	ctx context.Context,
	minio *minioclient.Client,
	videoID string,
	resolution enums.Resolution,
	inputPath string,
	workDir string,
	startSegment int,
) (int, error) {
	log := logger.Instance()

	resDir := filepath.Join(workDir, resolution.String())
	if err := os.MkdirAll(resDir, 0755); err != nil {
		return 0, fmt.Errorf("mkdir: %w", err)
	}

	playlistPath := filepath.Join(resDir, "playlist.m3u8")
	segmentPattern := filepath.Join(resDir, "segment_%05d.ts")

	// FFmpeg para gerar HLS com segmentos de 2s
	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-y",
		"-i", inputPath,
		"-vf", "scale="+resolution.AspectScale(),
		"-c:v", "libx264",
		"-crf", "23",
		"-preset", "medium",
		"-c:a", "aac",
		"-b:a", "128k",
		"-hls_time", fmt.Sprintf("%d", segmentDuration),
		"-hls_playlist_type", "vod",
		"-hls_segment_filename", segmentPattern,
		"-start_number", fmt.Sprintf("%d", startSegment),
		playlistPath,
	)

	if out, err := cmd.CombinedOutput(); err != nil {
		return 0, fmt.Errorf("ffmpeg: %w\n%s", err, string(out))
	}

	// Upload playlist e segmentos para MinIO
	entries, err := os.ReadDir(resDir)
	if err != nil {
		return 0, fmt.Errorf("readdir: %w", err)
	}

	segmentCount := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		localPath := filepath.Join(resDir, entry.Name())
		minioPath := fmt.Sprintf("videos/%s/hls/%s/%s", videoID, resolution.String(), entry.Name())

		contentType := "application/vnd.apple.mpegurl"
		if strings.HasSuffix(entry.Name(), ".ts") {
			contentType = "video/mp2t"
			segmentCount++
		}

		if err := minio.UploadFile(ctx, minioPath, localPath, contentType); err != nil {
			return segmentCount, fmt.Errorf("upload %s: %w", entry.Name(), err)
		}
	}

	log.Infof("hls: uploaded %d segments for %s/%s", segmentCount, videoID, resolution)
	return segmentCount, nil
}

// getLastSegmentNumber verifica qual o último segmento já existente no MinIO.
func getLastSegmentNumber(ctx context.Context, minio *minioclient.Client, videoID string, resolution enums.Resolution) int {
	prefix := fmt.Sprintf("videos/%s/hls/%s/segment_", videoID, resolution.String())

	// Lista objetos com o prefixo
	objects := minio.ListObjects(ctx, prefix)

	maxSegment := -1
	for _, obj := range objects {
		// Extrai número do segmento: segment_00042.ts -> 42
		name := filepath.Base(obj)
		if !strings.HasPrefix(name, "segment_") || !strings.HasSuffix(name, ".ts") {
			continue
		}

		numStr := strings.TrimPrefix(name, "segment_")
		numStr = strings.TrimSuffix(numStr, ".ts")

		var num int
		if _, err := fmt.Sscanf(numStr, "%d", &num); err == nil {
			if num > maxSegment {
				maxSegment = num
			}
		}
	}

	return maxSegment
}
