package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"

	"github.com/devlucas-java/luca-s3/internal/domain/enums"
)

// VideoInfo contém informações do vídeo original.
type VideoInfo struct {
	Width    int
	Height   int
	Duration float64
}

// analyzeVideo usa ffprobe para obter informações do vídeo.
func analyzeVideo(ctx context.Context, videoPath string) (*VideoInfo, error) {
	cmd := exec.CommandContext(ctx, "ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height,duration",
		"-show_entries", "format=duration",
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
		} `json:"streams"`
		Format struct {
			Duration string `json:"duration"`
		} `json:"format"`
	}

	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("parse ffprobe: %w", err)
	}

	if len(result.Streams) == 0 {
		return nil, fmt.Errorf("no video stream found")
	}

	info := &VideoInfo{
		Width:  result.Streams[0].Width,
		Height: result.Streams[0].Height,
	}

	// Tenta pegar duração do stream, senão do format
	if result.Streams[0].Duration != "" {
		info.Duration, _ = strconv.ParseFloat(result.Streams[0].Duration, 64)
	} else if result.Format.Duration != "" {
		info.Duration, _ = strconv.ParseFloat(result.Format.Duration, 64)
	}

	return info, nil
}

// filterValidResolutions remove resoluções maiores que o original.
func filterValidResolutions(requested []enums.Resolution, originalHeight int) []enums.Resolution {
	valid := make([]enums.Resolution, 0, len(requested))
	for _, res := range requested {
		if res.Height() <= originalHeight {
			valid = append(valid, res)
		}
	}
	return valid
}
