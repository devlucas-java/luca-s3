package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/devlucas-java/luca-s3/internal/domain/enums"
	"github.com/devlucas-java/luca-s3/internal/domain/model"
	minioclient "github.com/devlucas-java/luca-s3/internal/infrastructure/minio"
	"github.com/devlucas-java/luca-s3/internal/infrastructure/redis"
	"github.com/devlucas-java/luca-s3/pkg/logger"
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

func (s *TranscodeService) TranscodeVideo(
	ctx context.Context,
	videoID string,
	originalPath string,
	resolutions []enums.Resolution,
) (*model.TranscodeJob, error) {
	if active, err := s.jobRepo.ActiveJobForVideo(ctx, videoID); err != nil {
		return nil, fmt.Errorf("check active job: %w", err)
	} else if active != nil {
		return nil, fmt.Errorf("video %s already has an active job (%s) — wait for it to finish", videoID, active.Status)
	}

	exists, err := s.minio.ObjectExists(ctx, originalPath)
	if err != nil {
		return nil, fmt.Errorf("check minio: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("video not found in minio: %s", originalPath)
	}

	job := &model.TranscodeJob{
		ID:           videoID,
		VideoID:      videoID,
		OriginalPath: originalPath,
		RequestedRes: resolutions,
		Status:       enums.JobStatusPending,
		Progress:     make([]model.ResolutionProgress, 0),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.jobRepo.Save(ctx, job); err != nil {
		return nil, fmt.Errorf("save job: %w", err)
	}

	go s.process(job)
	return job, nil
}

func (s *TranscodeService) process(job *model.TranscodeJob) {
	log := logger.Instance()
	ctx := context.Background()

	log.Infof("transcode: video=%s started", job.VideoID)

	job.Status = enums.JobStatusProcessing
	job.UpdatedAt = time.Now()
	_ = s.jobRepo.Update(ctx, job)

	jobDir := filepath.Join(s.workDir, job.VideoID)
	if err := os.MkdirAll(jobDir, 0755); err != nil {
		s.failJob(ctx, job, fmt.Errorf("mkdir: %w", err))
		return
	}
	defer os.RemoveAll(jobDir)

	localInput := filepath.Join(jobDir, "input"+filepath.Ext(job.OriginalPath))
	if err := s.minio.DownloadFile(ctx, job.OriginalPath, localInput); err != nil {
		s.failJob(ctx, job, fmt.Errorf("download: %w", err))
		return
	}

	info, err := analyzeVideo(ctx, localInput)
	if err != nil {
		s.failJob(ctx, job, fmt.Errorf("analyze: %w", err))
		return
	}

	job.OriginalWidth = info.Width
	job.OriginalHeight = info.Height
	job.DurationSeconds = info.Duration
	_ = s.jobRepo.Update(ctx, job)

	if info.Is360 {
		log.Infof("transcode: video=%s detected as 360° (%s)", job.VideoID, info.Projection)
	}

	validRes := filterValidResolutions(job.RequestedRes, info.Height)
	if len(validRes) == 0 {
		s.failJob(ctx, job, fmt.Errorf("no valid resolutions for %dx%d", info.Width, info.Height))
		return
	}

	expectedSegments := calcTotalSegments(info.Duration)

	for _, res := range validRes {
		existing := countSegments(ctx, s.minio, job.VideoID, res)
		alreadyDone := existing >= expectedSegments-segmentTolerance

		entry := model.ResolutionProgress{
			Resolution:    res,
			TotalSegments: expectedSegments,
		}

		if alreadyDone {
			entry.Status = enums.JobStatusDone
			entry.SegmentsDone = expectedSegments
			entry.Percent = 100
			log.Infof("transcode: video=%s res=%s already complete (%d/%d segments) — skipping",
				job.VideoID, res, existing, expectedSegments)
		} else {
			entry.Status = enums.JobStatusPending
			entry.SegmentsDone = existing
		}

		job.Progress = append(job.Progress, entry)
	}
	_ = s.jobRepo.Update(ctx, job)

	allDone := true
	for _, p := range job.Progress {
		if p.Status != enums.JobStatusDone {
			allDone = false
			break
		}
	}
	if allDone {
		log.Infof("transcode: video=%s all resolutions already complete — nothing to do", job.VideoID)
		job.Status = enums.JobStatusDone
		job.UpdatedAt = time.Now()
		_ = s.jobRepo.Update(ctx, job)
		return
	}

	for i := range job.Progress {
		p := &job.Progress[i]

		if p.Status == enums.JobStatusDone {
			continue
		}

		startSeg := getLastSegmentNumber(ctx, s.minio, job.VideoID, p.Resolution) + 1

		p.Status = enums.JobStatusProcessing
		_ = s.jobRepo.Update(ctx, job)

		log.Infof("transcode: video=%s res=%s start_seg=%d expected=%d",
			job.VideoID, p.Resolution, startSeg, p.TotalSegments)

		onProgress := func(segsDone int) {
			p.SegmentsDone = startSeg + segsDone
			if p.TotalSegments > 0 {
				p.Percent = min(100, p.SegmentsDone*100/p.TotalSegments)
			}
			job.UpdatedAt = time.Now()
			_ = s.jobRepo.Update(ctx, job)
		}

		count, err := transcodeToHLS(
			ctx, s.minio, job.VideoID, p.Resolution,
			localInput, jobDir, startSeg, p.TotalSegments, onProgress,
			info.Is360,
		)

		if err != nil {
			log.Errorf("transcode: video=%s res=%s: %v", job.VideoID, p.Resolution, err)
			p.Status = enums.JobStatusFailed
			p.Error = err.Error()
		} else {
			p.SegmentsDone = startSeg + count
			p.Percent = 100
			p.Status = enums.JobStatusDone
		}

		if p.Status == enums.JobStatusDone {
			thumbErr := generateResolutionThumbnail(ctx, s.minio, localInput, jobDir, job.VideoID, p.Resolution)
			p.ThumbnailDone = thumbErr == nil
			if thumbErr != nil {
				log.Warnf("transcode: video=%s res=%s thumbnail: %v", job.VideoID, p.Resolution, thumbErr)
			}
		}

		job.UpdatedAt = time.Now()
		_ = s.jobRepo.Update(ctx, job)
	}

	if err := s.generateMasterPlaylist(ctx, job.VideoID, validRes); err != nil {
		log.Errorf("transcode: video=%s master playlist: %v", job.VideoID, err)
	}

	job.Finalize()
	job.UpdatedAt = time.Now()
	_ = s.jobRepo.Update(ctx, job)

	log.Infof("transcode: video=%s finished status=%s overall=%d%%",
		job.VideoID, job.Status, job.OverallPercent())
}

func (s *TranscodeService) failJob(ctx context.Context, job *model.TranscodeJob, err error) {
	logger.Instance().Errorf("transcode: video=%s failed: %v", job.VideoID, err)
	job.Status = enums.JobStatusFailed
	job.UpdatedAt = time.Now()
	_ = s.jobRepo.Update(ctx, job)
}

func (s *TranscodeService) GetJob(ctx context.Context, videoID string) (*model.TranscodeJob, error) {
	return s.jobRepo.FindByVideoID(ctx, videoID)
}

func (s *TranscodeService) ListJobs(ctx context.Context, filterStatus enums.JobStatus, page, pageSize int) ([]*model.TranscodeJob, int, error) {
	all, err := s.jobRepo.ListAll(ctx)
	if err != nil {
		return nil, 0, err
	}

	filtered := all
	if filterStatus != "" {
		filtered = make([]*model.TranscodeJob, 0, len(all))
		for _, j := range all {
			if j.Status == filterStatus {
				filtered = append(filtered, j)
			}
		}
	}

	total := len(filtered)

	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	if page < 1 {
		page = 1
	}
	start := (page - 1) * pageSize
	if start >= total {
		return []*model.TranscodeJob{}, total, nil
	}
	end := start + pageSize
	if end > total {
		end = total
	}

	return filtered[start:end], total, nil
}

func (s *TranscodeService) DeleteVideo(ctx context.Context, videoID string) error {
	log := logger.Instance()

	if active, err := s.jobRepo.ActiveJobForVideo(ctx, videoID); err != nil {
		return fmt.Errorf("check active job: %w", err)
	} else if active != nil {
		return fmt.Errorf("cannot delete: video %s has an active job (%s) — wait for it to finish", videoID, active.Status)
	}

	objects := s.minio.ListObjects(ctx, fmt.Sprintf("%s", videoID))
	if len(objects) == 0 {
		return fmt.Errorf("video not found: no files in MinIO for video_id %q", videoID)
	}

	if err := s.minio.DeleteFolder(ctx, fmt.Sprintf("%s/", videoID)); err != nil {
		return fmt.Errorf("delete minio folder: %w", err)
	}

	for _, obj := range objects {
		_ = s.minio.DeleteObject(ctx, obj)
	}

	if job, err := s.jobRepo.FindByVideoID(ctx, videoID); err == nil {
		_ = s.jobRepo.Delete(ctx, job.ID)
	}

	log.Infof("delete: video=%s deleted (%d objects removed)", videoID, len(objects))
	return nil
}

func (s *TranscodeService) generateMasterPlaylist(ctx context.Context, videoID string, resolutions []enums.Resolution) error {
	tmpFile := filepath.Join(s.workDir, "master_"+videoID+".m3u8")
	f, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("create master: %w", err)
	}
	defer f.Close()
	defer os.Remove(tmpFile)

	fmt.Fprintln(f, "#EXTM3U")
	fmt.Fprintln(f, "#EXT-X-VERSION:3")

	for _, res := range resolutions {
		bandwidth := res.Width() * res.Height() * 2
		fmt.Fprintf(f, "#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%dx%d\n",
			bandwidth, res.Width(), res.Height())
		fmt.Fprintf(f, "%s/playlist.m3u8\n", res.String())
	}

	return s.minio.UploadFile(ctx, fmt.Sprintf("%s/hls/master.m3u8", videoID), tmpFile, "application/vnd.apple.mpegurl")
}
