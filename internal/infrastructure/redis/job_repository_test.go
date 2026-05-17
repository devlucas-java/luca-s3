package redis

import (
	"context"
	"testing"
	"time"

	"github.com/devlucas-java/luca-s3/internal/domain/enums"
	"github.com/devlucas-java/luca-s3/internal/domain/model"
)

// newTestRepo connects to a local Redis and returns a JobRepository.
// Skips the test if Redis is not available.
func newTestRepo(t *testing.T) *JobRepository {
	t.Helper()
	client, err := New("localhost:6379", "", 0)
	if err != nil {
		t.Skip("Redis not available:", err)
	}
	t.Cleanup(func() { client.Close() })
	return NewJobRepository(client)
}

func makeJob(videoID string, status enums.JobStatus) *model.TranscodeJob {
	return &model.TranscodeJob{
		ID:           videoID,
		VideoID:      videoID,
		OriginalPath: videoID + ".mp4",
		RequestedRes: []enums.Resolution{enums.RESOLUTION_720P, enums.RESOLUTION_480P},
		Status:       status,
		Progress: []model.ResolutionProgress{
			{
				Resolution:    enums.RESOLUTION_720P,
				Status:        status,
				SegmentsDone:  23,
				TotalSegments: 23,
				Percent:       100,
				ThumbnailDone: true,
			},
			{
				Resolution:    enums.RESOLUTION_480P,
				Status:        status,
				SegmentsDone:  23,
				TotalSegments: 23,
				Percent:       100,
				ThumbnailDone: true,
			},
		},
		OriginalWidth:   1280,
		OriginalHeight:  720,
		DurationSeconds: 45.7,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func TestJobRepository_SaveAndFind(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	job := makeJob("teste", enums.JobStatusDone)

	t.Cleanup(func() { _ = repo.Delete(ctx, job.ID) })

	if err := repo.Save(ctx, job); err != nil {
		t.Fatalf("Save: %v", err)
	}

	t.Run("FindByID", func(t *testing.T) {
		got, err := repo.FindByID(ctx, job.ID)
		if err != nil {
			t.Fatalf("FindByID: %v", err)
		}
		if got.VideoID != job.VideoID {
			t.Errorf("VideoID: want %s, got %s", job.VideoID, got.VideoID)
		}
		if got.Status != job.Status {
			t.Errorf("Status: want %s, got %s", job.Status, got.Status)
		}
	})

	t.Run("FindByVideoID", func(t *testing.T) {
		got, err := repo.FindByVideoID(ctx, job.VideoID)
		if err != nil {
			t.Fatalf("FindByVideoID: %v", err)
		}
		if got.ID != job.ID {
			t.Errorf("ID: want %s, got %s", job.ID, got.ID)
		}
	})
}

func TestJobRepository_Update(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	job := makeJob("teste-update", enums.JobStatusPending)
	t.Cleanup(func() { _ = repo.Delete(ctx, job.ID) })

	if err := repo.Save(ctx, job); err != nil {
		t.Fatalf("Save: %v", err)
	}

	job.Status = enums.JobStatusDone
	job.UpdatedAt = time.Now()
	if err := repo.Update(ctx, job); err != nil {
		t.Fatalf("Update: %v", err)
	}

	got, err := repo.FindByID(ctx, job.ID)
	if err != nil {
		t.Fatalf("FindByID after update: %v", err)
	}
	if got.Status != enums.JobStatusDone {
		t.Errorf("Status after update: want done, got %s", got.Status)
	}
}

func TestJobRepository_ActiveJobForVideo(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	t.Run("no job returns nil", func(t *testing.T) {
		active, err := repo.ActiveJobForVideo(ctx, "video-nao-existe")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if active != nil {
			t.Error("expected nil for non-existent video")
		}
	})

	t.Run("done job returns nil", func(t *testing.T) {
		job := makeJob("teste-done", enums.JobStatusDone)
		t.Cleanup(func() { _ = repo.Delete(ctx, job.ID) })
		_ = repo.Save(ctx, job)

		active, err := repo.ActiveJobForVideo(ctx, job.VideoID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if active != nil {
			t.Error("done job should not be returned as active")
		}
	})

	t.Run("processing job returns job", func(t *testing.T) {
		job := makeJob("teste-processing", enums.JobStatusProcessing)
		t.Cleanup(func() { _ = repo.Delete(ctx, job.ID) })
		_ = repo.Save(ctx, job)

		active, err := repo.ActiveJobForVideo(ctx, job.VideoID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if active == nil {
			t.Fatal("expected active job for processing video")
		}
		if active.ID != job.ID {
			t.Errorf("ID: want %s, got %s", job.ID, active.ID)
		}
	})

	t.Run("pending job returns job", func(t *testing.T) {
		job := makeJob("teste-pending", enums.JobStatusPending)
		t.Cleanup(func() { _ = repo.Delete(ctx, job.ID) })
		_ = repo.Save(ctx, job)

		active, err := repo.ActiveJobForVideo(ctx, job.VideoID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if active == nil {
			t.Fatal("expected active job for pending video")
		}
	})
}

func TestJobRepository_ListAll(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	jobs := []*model.TranscodeJob{
		makeJob("teste-list-1", enums.JobStatusDone),
		makeJob("teste-list-2", enums.JobStatusFailed),
		makeJob("teste-list-3", enums.JobStatusProcessing),
	}

	for _, j := range jobs {
		j := j
		t.Cleanup(func() { _ = repo.Delete(ctx, j.ID) })
		if err := repo.Save(ctx, j); err != nil {
			t.Fatalf("Save %s: %v", j.ID, err)
		}
	}

	all, err := repo.ListAll(ctx)
	if err != nil {
		t.Fatalf("ListAll: %v", err)
	}

	// At least our 3 jobs must be present (Redis may have others)
	found := map[string]bool{}
	for _, j := range all {
		found[j.ID] = true
	}
	for _, j := range jobs {
		if !found[j.ID] {
			t.Errorf("job %s not found in ListAll", j.ID)
		}
	}

	// Verify sorted newest first
	for i := 1; i < len(all); i++ {
		if all[i].CreatedAt.After(all[i-1].CreatedAt) {
			t.Errorf("ListAll not sorted newest first at index %d", i)
		}
	}
}

func TestJobRepository_Delete(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	job := makeJob("teste-delete", enums.JobStatusDone)
	if err := repo.Save(ctx, job); err != nil {
		t.Fatalf("Save: %v", err)
	}

	if err := repo.Delete(ctx, job.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	_, err := repo.FindByID(ctx, job.ID)
	if err == nil {
		t.Error("expected error after delete, got nil")
	}

	_, err = repo.FindByVideoID(ctx, job.VideoID)
	if err == nil {
		t.Error("expected error finding by video_id after delete, got nil")
	}
}

func TestJobRepository_DeleteByVideoID(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	job := makeJob("teste-delete-video", enums.JobStatusDone)
	if err := repo.Save(ctx, job); err != nil {
		t.Fatalf("Save: %v", err)
	}

	if err := repo.DeleteByVideoID(ctx, job.VideoID); err != nil {
		t.Fatalf("DeleteByVideoID: %v", err)
	}

	_, err := repo.FindByVideoID(ctx, job.VideoID)
	if err == nil {
		t.Error("expected error after DeleteByVideoID, got nil")
	}
}

func TestJobRepository_FindByID_NotFound(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	_, err := repo.FindByID(ctx, "id-que-nao-existe")
	if err == nil {
		t.Error("expected error for non-existent ID, got nil")
	}
}

func TestJobRepository_FindByVideoID_NotFound(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	_, err := repo.FindByVideoID(ctx, "video-que-nao-existe")
	if err == nil {
		t.Error("expected error for non-existent video_id, got nil")
	}
}
