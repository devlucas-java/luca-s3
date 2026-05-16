package service

import (
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/devlucas-java/luca-s3/internal/domain/enums"
	"github.com/devlucas-java/luca-s3/internal/domain/model"
	minioclient "github.com/devlucas-java/luca-s3/internal/infrastructure/minio"
	"github.com/devlucas-java/luca-s3/internal/infrastructure/redis"
	"github.com/devlucas-java/luca-s3/pkg/logger"
	"github.com/google/uuid"
)

type TranscodeService struct {
	minio   *minioclient.Client
	jobRepo *redis.JobRepository
	workDir string
}

func NewTranscodeService(minio *minioclient.Client, jobRepo *redis.JobRepository) *TranscodeService {
	workDir := os.Getenv("FFMPEG_WORK_DIR")
	if workDir == "" {
		workDir = "/tmp/ffmpeg"
	}
	return &TranscodeService{
		minio:   minio,
		jobRepo: jobRepo,
		workDir: workDir,
	}
}

// TranscodeVideo inicia transcodificação HLS para um vídeo no MinIO.
func (s *TranscodeService) TranscodeVideo(
	ctx context.Context,
	videoID string,
	originalPath string,
	resolutions []enums.Resolution,
) (*model.TranscodeJob, error) {
	log := logger.Instance()

	// Verifica se vídeo existe no MinIO
	exists, err := s.minio.ObjectExists(ctx, originalPath)
	if err != nil {
		return nil, fmt.Errorf("check minio: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("video not found in minio: %s", originalPath)
	}

	// Cria diretório temporário
	jobID := uuid.New().String()
	jobDir := filepath.Join(s.workDir, jobID)
	if err := os.MkdirAll(jobDir, 0755); err != nil {
		return nil, fmt.Errorf("mkdir: %w", err)
	}
	defer os.RemoveAll(jobDir)

	// Download do vídeo original
	localInput := filepath.Join(jobDir, "input"+filepath.Ext(originalPath))
	if err := s.minio.DownloadFile(ctx, originalPath, localInput); err != nil {
		return nil, fmt.Errorf("download: %w", err)
	}

	// Analisa vídeo para obter dimensões e duração
	info, err := analyzeVideo(ctx, localInput)
	if err != nil {
		return nil, fmt.Errorf("analyze: %w", err)
	}

	// Filtra resoluções válidas (não maiores que o original)
	validRes := filterValidResolutions(resolutions, info.Height)
	if len(validRes) == 0 {
		return nil, fmt.Errorf("no valid resolutions (original is %dx%d)", info.Width, info.Height)
	}

	// Cria job
	job := &model.TranscodeJob{
		ID:              jobID,
		VideoID:         videoID,
		OriginalPath:    originalPath,
		RequestedRes:    validRes,
		Status:          enums.JobStatusProcessing,
		OriginalWidth:   info.Width,
		OriginalHeight:  info.Height,
		DurationSeconds: info.Duration,
		Progress:        make([]model.ResolutionProgress, 0),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.jobRepo.Save(ctx, job); err != nil {
		return nil, fmt.Errorf("save job: %w", err)
	}

	// Processa cada resolução
	for _, res := range validRes {
		progress := model.ResolutionProgress{
			Resolution: res,
			Status:     enums.JobStatusProcessing,
		}

		// Verifica último segmento existente no MinIO
		lastSegment := getLastSegmentNumber(ctx, s.minio, videoID, res)
		progress.LastSegment = lastSegment
		progress.TotalSegments = int(math.Ceil(info.Duration / segmentDuration))

		log.Infof("transcode: video=%s res=%s resuming from segment %d/%d",
			videoID, res, lastSegment+1, progress.TotalSegments)

		// Transcoda HLS (retoma de onde parou)
		segmentCount, err := transcodeToHLS(ctx, s.minio, videoID, res, localInput, jobDir, lastSegment+1)
		if err != nil {
			log.Errorf("transcode: video=%s res=%s: %v", videoID, res, err)
			progress.Status = enums.JobStatusFailed
			progress.Error = err.Error()
		} else {
			progress.LastSegment = lastSegment + segmentCount
			progress.Status = enums.JobStatusDone
		}

		// Gera thumbnail para esta resolução
		if progress.Status == enums.JobStatusDone {
			thumbErr := generateResolutionThumbnail(ctx, s.minio, localInput, jobDir, videoID, res)
			progress.ThumbnailDone = (thumbErr == nil)
			if thumbErr != nil {
				log.Warnf("transcode: thumbnail video=%s res=%s: %v", videoID, res, thumbErr)
			}
		}

		job.Progress = append(job.Progress, progress)
	}

	// Gera master playlist
	if err := s.generateMasterPlaylist(ctx, videoID, validRes); err != nil {
		log.Errorf("transcode: master playlist video=%s: %v", videoID, err)
	}

	job.Finalize()
	job.UpdatedAt = time.Now()
	if err := s.jobRepo.Update(ctx, job); err != nil {
		log.Errorf("transcode: update job %s: %v", jobID, err)
	}

	log.Infof("transcode: job=%s video=%s status=%s", jobID, videoID, job.Status)
	return job, nil
}

// DeleteVideo deletes all video files from MinIO.
func (s *TranscodeService) DeleteVideo(ctx context.Context, videoID string) error {
	log := logger.Instance()

	// Delete entire video folder
	prefix := fmt.Sprintf("videos/%s/", videoID)
	if err := s.minio.DeleteFolder(ctx, prefix); err != nil {
		return fmt.Errorf("delete folder: %w", err)
	}

	// Delete jobs related to this video
	jobs, err := s.jobRepo.FindByVideoID(ctx, videoID)
	if err == nil {
		for _, job := range jobs {
			_ = s.jobRepo.Delete(ctx, job.ID)
		}
	}

	log.Infof("delete: video=%s deleted", videoID)
	return nil
}

// GetJob retorna o status de um job.
func (s *TranscodeService) GetJob(ctx context.Context, jobID string) (*model.TranscodeJob, error) {
	return s.jobRepo.FindByID(ctx, jobID)
}

// GetJobsByVideo retorna todos os jobs de um vídeo.
func (s *TranscodeService) GetJobsByVideo(ctx context.Context, videoID string) ([]*model.TranscodeJob, error) {
	return s.jobRepo.FindByVideoID(ctx, videoID)
}

// generateMasterPlaylist cria o master playlist HLS.
func (s *TranscodeService) generateMasterPlaylist(ctx context.Context, videoID string, resolutions []enums.Resolution) error {
	tmpFile := filepath.Join(s.workDir, "master.m3u8")
	f, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("create master: %w", err)
	}
	defer f.Close()
	defer os.Remove(tmpFile)

	fmt.Fprintln(f, "#EXTM3U")
	fmt.Fprintln(f, "#EXT-X-VERSION:3")

	for _, res := range resolutions {
		bandwidth := res.Width() * res.Height() * 2 // estimativa
		fmt.Fprintf(f, "#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%dx%d\n",
			bandwidth, res.Width(), res.Height())
		fmt.Fprintf(f, "%s/playlist.m3u8\n", res.String())
	}

	minioPath := fmt.Sprintf("videos/%s/hls/master.m3u8", videoID)
	return s.minio.UploadFile(ctx, minioPath, tmpFile, "application/vnd.apple.mpegurl")
}
