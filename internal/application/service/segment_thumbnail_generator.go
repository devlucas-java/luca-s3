package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/devlucas-java/luca-s3/internal/domain/enums"
	minioclient "github.com/devlucas-java/luca-s3/internal/infrastructure/minio"
)

// generateResolutionThumbnail gera thumbnail para uma resolução específica.
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
		// Fallback: primeiro frame
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

	minioPath := fmt.Sprintf("videos/%s/hls/%s/thumbnail.jpg", videoID, resolution.String())
	return minio.UploadFile(ctx, minioPath, thumbFile, "image/jpeg")
}

// generateSegmentThumbnails gera thumbnails para cada segmento de 2s.
// TODO: Implementar geração de thumbnail por segmento para preview no seek.
func generateSegmentThumbnails(
	ctx context.Context,
	minio *minioclient.Client,
	inputPath string,
	workDir string,
	videoID string,
	resolution enums.Resolution,
	totalSegments int,
) error {
	thumbDir := filepath.Join(workDir, "thumbs", resolution.String())
	if err := os.MkdirAll(thumbDir, 0755); err != nil {
		return fmt.Errorf("mkdir thumbs: %w", err)
	}

	// Gera thumbnail a cada 2 segundos
	for i := 0; i < totalSegments; i++ {
		timestamp := i * segmentDuration
		thumbFile := filepath.Join(thumbDir, fmt.Sprintf("thumb_%05d.jpg", i))

		cmd := exec.CommandContext(ctx, "ffmpeg",
			"-y",
			"-ss", fmt.Sprintf("%d", timestamp),
			"-i", inputPath,
			"-vframes", "1",
			"-vf", "scale=160:-1", // thumbnail pequena
			"-q:v", "5",
			thumbFile,
		)

		if out, err := cmd.CombinedOutput(); err != nil {
			// Ignora erros em thumbnails individuais
			_ = out
			continue
		}

		// Upload para MinIO
		minioPath := fmt.Sprintf("videos/%s/hls/%s/thumbs/thumb_%05d.jpg", videoID, resolution.String(), i)
		if err := minio.UploadFile(ctx, minioPath, thumbFile, "image/jpeg"); err != nil {
			// Ignora erros de upload de thumbnails
			continue
		}
	}

	return nil
}
