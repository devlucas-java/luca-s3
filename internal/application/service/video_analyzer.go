package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/devlucas-java/luca-s3/internal/domain/enums"
)

type VideoProbeInfo struct {
	Width      int
	Height     int
	Duration   float64
	Is360      bool
	Projection string
}

func analyzeVideo(ctx context.Context, videoPath string) (*VideoProbeInfo, error) {
	cmd := exec.CommandContext(ctx, "ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height,duration,side_data_list",
		"-show_entries", "stream_tags=spherical,projection,stereo_mode",
		"-show_entries", "format=duration",
		"-show_entries", "format_tags=spherical,projection",
		"-of", "json",
		videoPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("ffprobe: %w\n%s", err, string(output))
	}

	var result struct {
		Streams []struct {
			Width    int    `json:"width"`
			Height   int    `json:"height"`
			Duration string `json:"duration"`
			Tags     struct {
				Spherical  string `json:"spherical"`
				Projection string `json:"projection"`
			} `json:"tags"`
			SideDataList []struct {
				SideDataType string `json:"side_data_type"`
				Projection   string `json:"projection"`
			} `json:"side_data_list"`
		} `json:"streams"`
		Format struct {
			Duration string `json:"duration"`
			Tags     struct {
				Spherical  string `json:"spherical"`
				Projection string `json:"projection"`
			} `json:"tags"`
		} `json:"format"`
	}

	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("parse ffprobe: %w", err)
	}

	if len(result.Streams) == 0 {
		return nil, fmt.Errorf("no video stream found")
	}

	info := &VideoProbeInfo{
		Width:  result.Streams[0].Width,
		Height: result.Streams[0].Height,
	}

	if result.Streams[0].Duration != "" {
		info.Duration, _ = strconv.ParseFloat(result.Streams[0].Duration, 64)
	} else if result.Format.Duration != "" {
		info.Duration, _ = strconv.ParseFloat(result.Format.Duration, 64)
	}

	for _, sd := range result.Streams[0].SideDataList {
		if strings.Contains(strings.ToLower(sd.SideDataType), "spherical") {
			info.Is360 = true
			info.Projection = sd.Projection
			break
		}
	}

	if !info.Is360 {
		tags := result.Streams[0].Tags
		if strings.EqualFold(tags.Spherical, "true") || tags.Projection != "" {
			info.Is360 = true
			info.Projection = tags.Projection
		}
	}

	if !info.Is360 {
		ftags := result.Format.Tags
		if strings.EqualFold(ftags.Spherical, "true") || ftags.Projection != "" {
			info.Is360 = true
			info.Projection = ftags.Projection
		}
	}

	if !info.Is360 && info.Height > 0 {
		ratio := float64(info.Width) / float64(info.Height)
		if ratio >= 1.95 && ratio <= 2.05 {
			info.Is360 = true
			info.Projection = "equirectangular"
		}
	}

	return info, nil
}

func filterValidResolutions(requested []enums.Resolution, originalHeight int) []enums.Resolution {
	valid := make([]enums.Resolution, 0, len(requested))
	for _, res := range requested {
		if res.Height() <= originalHeight {
			valid = append(valid, res)
		}
	}
	return valid
}
