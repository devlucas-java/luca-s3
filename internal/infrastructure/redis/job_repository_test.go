package redis

import (
	"context"
	"testing"
	"time"

	"github.com/devlucas-java/luca-s3/internal/domain/enums"
	"github.com/devlucas-java/luca-s3/internal/domain/model"
)

func TestJobRepository(t *testing.T) {
	// Skip if no Redis available
	client, err := New("localhost:6379", "", 0)
	if err != nil {
		t.Skip("Redis not available:", err)
	}
	defer client.Close()

	repo := NewJobRepository(client)
	ctx := context.Background()

	t.Run("SaveAndFind", func(t *testing.T) {
		job := &model.TranscodeJob{
			ID:              "test-job-1",
			VideoID:         "video-1",
			OriginalPath:    "videos/video-1.mp4",
			RequestedRes:    []enums.Resolution{enums.RESOLUTION_720P, enums.RESOLUTION_480P},
			Status:          enums.JobStatusProcessing,
			OriginalWidth:   1920,
			OriginalHeight:  1080,
			DurationSeconds: 120.5,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		// Save
		err := repo.Save(ctx, job)
		if err != nil {
			t.Fatalf("save failed: %v", err)
		}

		// Find by ID
		found, err := repo.FindByID(ctx, job.ID)
		if err != nil {
			t.Fatalf("find failed: %v", err)
		}

		if found.ID != job.ID {
			t.Fatalf("expected ID %s, got %s", job.ID, found.ID)
		}
		if found.VideoID != job.VideoID {
			t.Fatalf("expected VideoID %s, got %s", job.VideoID, found.VideoID)
		}

		// Delete
		err = repo.Delete(ctx, job.ID)
		if err != nil {
			t.Fatalf("delete failed: %v", err)
		}
	})

	t.Run("FindByVideoID", func(t *testing.T) {
		videoID := "video-2"

		// Create multiple jobs for same video
		for i := 0; i < 3; i++ {
			job := &model.TranscodeJob{
				ID:           "test-job-" + string(rune('a'+i)),
				VideoID:      videoID,
				OriginalPath: "videos/" + videoID + ".mp4",
				Status:       enums.JobStatusPending,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			err := repo.Save(ctx, job)
			if err != nil {
				t.Fatalf("save failed: %v", err)
			}
			defer repo.Delete(ctx, job.ID)
		}

		// Find by video ID
		jobs, err := repo.FindByVideoID(ctx, videoID)
		if err != nil {
			t.Fatalf("find by video id failed: %v", err)
		}

		if len(jobs) < 3 {
			t.Fatalf("expected at least 3 jobs, got %d", len(jobs))
		}
	})
}
