package service

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/devlucas-java/luca-s3/internal/domain/enums"
	minioclient "github.com/devlucas-java/luca-s3/internal/infrastructure/minio"
)

type ResolutionInfo struct {
	Resolution   enums.Resolution
	SegmentCount int
	HasPlaylist  bool
	HasThumbnail bool
	SizeBytes    int64
}

// VideoInfo describes everything that exists in MinIO for a video_id.
type VideoInfo struct {
	VideoID        string
	Found          bool   // false if nothing exists at all
	OriginalPath   string // e.g. {id}.mp4 — empty if not found
	OriginalSize   int64
	HasMaster      bool // {id}/hls/master.m3u8
	Resolutions    []ResolutionInfo
	TotalSizeBytes int64
}

// InspectService inspects MinIO to report what exists for a given video.
type InspectService struct {
	minio *minioclient.Client
}

func NewInspectService(minio *minioclient.Client) *InspectService {
	return &InspectService{minio: minio}
}

// InspectVideo scans MinIO and returns a full report for the given video_id.
func (s *InspectService) InspectVideo(ctx context.Context, videoID string) (*VideoInfo, error) {
	if videoID == "" {
		return nil, fmt.Errorf("video_id is required")
	}

	info := &VideoInfo{VideoID: videoID}

	// ── 1. Look for original file (any supported extension) ──────────────────
	for _, ext := range []enums.Extension{
		enums.EXTENSION_MP4,
		enums.EXTENSION_MOV,
		enums.EXTENSION_WEBM,
	} {
		path := fmt.Sprintf("%s.%s", videoID, ext.String())
		size, exists := s.minio.StatObject(ctx, path)
		if exists {
			info.OriginalPath = path
			info.OriginalSize = size
			info.TotalSizeBytes += size
			break
		}
	}

	// ── 2. Check master playlist ──────────────────────────────────────────────
	masterPath := fmt.Sprintf("%s/hls/master.m3u8", videoID)
	_, info.HasMaster = s.minio.StatObject(ctx, masterPath)

	// ── 3. Scan all HLS objects under {id}/hls/ ───────────────────────────────
	hlsPrefix := fmt.Sprintf("%s/hls/", videoID)
	objects := s.minio.ListObjects(ctx, hlsPrefix)

	// Group by resolution
	resMap := map[enums.Resolution]*ResolutionInfo{}

	for _, obj := range objects {
		// Skip master playlist
		if strings.HasSuffix(obj, "master.m3u8") {
			continue
		}

		// Path pattern: {id}/hls/{resolution}/{file}
		rel := strings.TrimPrefix(obj, hlsPrefix)
		parts := strings.SplitN(rel, "/", 2)
		if len(parts) != 2 {
			continue
		}

		res := enums.Resolution(parts[0])
		file := parts[1]

		if _, ok := resMap[res]; !ok {
			resMap[res] = &ResolutionInfo{Resolution: res}
		}
		ri := resMap[res]

		switch {
		case strings.HasSuffix(file, ".ts"):
			ri.SegmentCount++
		case file == "playlist.m3u8":
			ri.HasPlaylist = true
		case file == "thumbnail.jpg":
			ri.HasThumbnail = true
		}

		// Accumulate size
		if size, exists := s.minio.StatObject(ctx, obj); exists {
			ri.SizeBytes += size
			info.TotalSizeBytes += size
		}
	}

	// Convert map to sorted slice (by resolution height)
	resOrder := []enums.Resolution{
		enums.RESOLUTION_360P,
		enums.RESOLUTION_480P,
		enums.RESOLUTION_720P,
		enums.RESOLUTION_1080P,
		enums.RESOLUTION_1440P,
		enums.RESOLUTION_4K,
	}
	for _, res := range resOrder {
		if ri, ok := resMap[res]; ok {
			info.Resolutions = append(info.Resolutions, *ri)
		}
	}

	// Also add any unknown resolutions not in the standard list
	for res, ri := range resMap {
		found := false
		for _, r := range resOrder {
			if r == res {
				found = true
				break
			}
		}
		if !found {
			info.Resolutions = append(info.Resolutions, *ri)
		}
	}

	info.Found = info.OriginalPath != "" || info.HasMaster || len(info.Resolutions) > 0

	_ = filepath.Base // keep import
	return info, nil
}
