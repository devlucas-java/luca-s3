package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/devlucas-java/luca-s3/internal/domain/model"
)

const (
	jobKeyPrefix = "transcode:job:"
	jobTTL       = 7 * 24 * time.Hour // 7 days
)

type JobRepository struct {
	client *Client
}

func NewJobRepository(client *Client) *JobRepository {
	return &JobRepository{client: client}
}

func (r *JobRepository) Save(ctx context.Context, job *model.TranscodeJob) error {
	key := jobKeyPrefix + job.ID
	data, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("marshal job: %w", err)
	}

	return r.client.client.Set(ctx, key, data, jobTTL).Err()
}

func (r *JobRepository) Update(ctx context.Context, job *model.TranscodeJob) error {
	return r.Save(ctx, job) // Same as save in Redis
}

func (r *JobRepository) FindByID(ctx context.Context, id string) (*model.TranscodeJob, error) {
	key := jobKeyPrefix + id
	data, err := r.client.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, fmt.Errorf("job not found: %s", id)
	}

	var job model.TranscodeJob
	if err := json.Unmarshal(data, &job); err != nil {
		return nil, fmt.Errorf("unmarshal job: %w", err)
	}

	return &job, nil
}

func (r *JobRepository) FindByVideoID(ctx context.Context, videoID string) ([]*model.TranscodeJob, error) {
	// Scan all job keys and filter by video_id
	var cursor uint64
	var jobs []*model.TranscodeJob

	for {
		keys, nextCursor, err := r.client.client.Scan(ctx, cursor, jobKeyPrefix+"*", 100).Result()
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

			if job.VideoID == videoID {
				jobs = append(jobs, &job)
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return jobs, nil
}

func (r *JobRepository) Delete(ctx context.Context, id string) error {
	key := jobKeyPrefix + id
	return r.client.client.Del(ctx, key).Err()
}
