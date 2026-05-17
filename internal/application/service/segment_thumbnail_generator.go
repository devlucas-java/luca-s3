package service

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/devlucas-java/luca-s3/internal/domain/enums"
	minioclient "github.com/devlucas-java/luca-s3/internal/infrastructure/minio"
)

func generateResolutionThumbnail(
	ctx context.Context,
	minio *minioclient.Client,
	inputPath string,
	workDir string,
	videoID string,
	resolution enums.Resolution,
) error {
	thumbFile := filepath.Join(workDir, fmt.Sprintf("thumb_%s.jpg", resolution.String()))

	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-y",
		"-i", inputPath,
		"-ss", "00:00:05",
		"-vframes", "1",
		"-vf", "scale="+resolution.AspectScale(),
		"-q:v", "2",
		thumbFile,
	)

	if out, err := cmd.CombinedOutput(); err != nil {
		cmd2 := exec.CommandContext(ctx, "ffmpeg",
			"-y",
			"-i", inputPath,
			"-vframes", "1",
			"-vf", "scale="+resolution.AspectScale(),
			"-q:v", "2",
			thumbFile,
		)
		if out2, err2 := cmd2.CombinedOutput(); err2 != nil {
			return fmt.Errorf("%w\n%s", err2, string(out2))
		}
		_ = out
	}

	minioPath := fmt.Sprintf("%s/hls/%s/thumbnail.jpg", videoID, resolution.String())
	return minio.UploadFile(ctx, minioPath, thumbFile, "image/jpeg")
}
