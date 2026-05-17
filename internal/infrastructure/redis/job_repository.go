package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/devlucas-java/luca-s3/internal/domain/enums"
	"github.com/devlucas-java/luca-s3/internal/domain/model"
)

const (
	jobKeyPrefix  = "transcode:job:"
	videoIndexKey = "transcode:video:" // video_id → job_id (latest active)
	jobTTL        = 7 * 24 * time.Hour
)

type JobRepository struct {
	client *Client
}

func NewJobRepository(client *Client) *JobRepository {
	return &JobRepository{client: client}
}

// Save persists a job and updates the video→job index.
func (r *JobRepository) Save(ctx context.Context, job *model.TranscodeJob) error {
	data, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("marshal job: %w", err)
	}

	pipe := r.client.client.Pipeline()
	pipe.Set(ctx, jobKeyPrefix+job.ID, data, jobTTL)
	pipe.Set(ctx, videoIndexKey+job.VideoID, job.ID, jobTTL)
	_, err = pipe.Exec(ctx)
	return err
}

// Update persists updated job state (same as Save).
func (r *JobRepository) Update(ctx context.Context, job *model.TranscodeJob) error {
	return r.Save(ctx, job)
}

// FindByID returns a job by its own ID.
func (r *JobRepository) FindByID(ctx context.Context, id string) (*model.TranscodeJob, error) {
	data, err := r.client.client.Get(ctx, jobKeyPrefix+id).Bytes()
	if err != nil {
		return nil, fmt.Errorf("job not found: %s", id)
	}
	var job model.TranscodeJob
	if err := json.Unmarshal(data, &job); err != nil {
		return nil, fmt.Errorf("unmarshal job: %w", err)
	}
	return &job, nil
}

// FindByVideoID returns the latest job for a given MinIO video ID.
func (r *JobRepository) FindByVideoID(ctx context.Context, videoID string) (*model.TranscodeJob, error) {
	jobID, err := r.client.client.Get(ctx, videoIndexKey+videoID).Result()
	if err != nil {
		return nil, fmt.Errorf("no job for video: %s", videoID)
	}
	return r.FindByID(ctx, jobID)
}

// ActiveJobForVideo returns the active (pending/processing) job for a video, if any.
func (r *JobRepository) ActiveJobForVideo(ctx context.Context, videoID string) (*model.TranscodeJob, error) {
	job, err := r.FindByVideoID(ctx, videoID)
	if err != nil {
		return nil, nil // no job at all — not an error
	}
	if job.Status == enums.JobStatusPending || job.Status == enums.JobStatusProcessing {
		return job, nil
	}
	return nil, nil
}

// ListAll returns all jobs in cache, sorted by CreatedAt descending.
// If multiple jobs exist for the same video_id (e.g. leftover from old UUIDs),
// only the most recent one is returned.
func (r *JobRepository) ListAll(ctx context.Context) ([]*model.TranscodeJob, error) {
	var cursor uint64
	// videoID → most recent job
	seen := map[string]*model.TranscodeJob{}

	for {
		keys, next, err := r.client.client.Scan(ctx, cursor, jobKeyPrefix+"*", 200).Result()
		if err != nil {
			return nil, err
		}

		for _, key := range keys {
			data, err := r.client.client.Get(ctx, key).Bytes()
			if err != nil {
				continue
			}
			var job model.TranscodeJob
			if err := json.Unmarshal(data, &job); err != nil {
				continue
			}
			// Keep only the most recent job per video_id
			if existing, ok := seen[job.VideoID]; !ok || job.CreatedAt.After(existing.CreatedAt) {
				seen[job.VideoID] = &job
			}
		}

		cursor = next
		if cursor == 0 {
			break
		}
	}

	jobs := make([]*model.TranscodeJob, 0, len(seen))
	for _, j := range seen {
		jobs = append(jobs, j)
	}

	// Sort newest first
	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].CreatedAt.After(jobs[j].CreatedAt)
	})

	return jobs, nil
}

// Delete removes a job and its video index entry.
func (r *JobRepository) Delete(ctx context.Context, id string) error {
	job, err := r.FindByID(ctx, id)
	if err == nil {
		r.client.client.Del(ctx, videoIndexKey+job.VideoID)
	}
	return r.client.client.Del(ctx, jobKeyPrefix+id).Err()
}

// DeleteByVideoID removes all jobs for a video.
func (r *JobRepository) DeleteByVideoID(ctx context.Context, videoID string) error {
	job, err := r.FindByVideoID(ctx, videoID)
	if err == nil {
		r.client.client.Del(ctx, jobKeyPrefix+job.ID)
	}
	r.client.client.Del(ctx, videoIndexKey+videoID)
	return nil
}
