package model

import (
	"time"

	"github.com/devlucas-java/luca-s3/internal/domain/enums"
)

// ResolutionProgress tracks progress for a specific resolution.
type ResolutionProgress struct {
	Resolution    enums.Resolution `json:"resolution"`
	Status        enums.JobStatus  `json:"status"`
	LastSegment   int              `json:"last_segment"`   // Last processed segment
	TotalSegments int              `json:"total_segments"` // Total expected segments
	ThumbnailDone bool             `json:"thumbnail_done"` // Whether thumbnail was generated
	Error         string           `json:"error,omitempty"`
}

// TranscodeJob represents an HLS transcoding job.
type TranscodeJob struct {
	ID              string               `json:"id"`
	VideoID         string               `json:"video_id"`      // Video ID in MinIO
	OriginalPath    string               `json:"original_path"` // Path in MinIO: videos/{id}.mp4
	RequestedRes    []enums.Resolution   `json:"requested_res"` // Requested resolutions
	Status          enums.JobStatus      `json:"status"`
	Progress        []ResolutionProgress `json:"progress"`
	OriginalWidth   int                  `json:"original_width"`   // Original video width
	OriginalHeight  int                  `json:"original_height"`  // Original video height
	DurationSeconds float64              `json:"duration_seconds"` // Duration in seconds
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
}

// Finalize sets the final job status.
func (j *TranscodeJob) Finalize() {
	anyDone := false
	allFailed := true
	for _, p := range j.Progress {
		if p.Status == enums.JobStatusDone {
			anyDone = true
			allFailed = false
		}
		if p.Status != enums.JobStatusFailed {
			allFailed = false
		}
	}

	switch {
	case anyDone:
		j.Status = enums.JobStatusDone
	case allFailed:
		j.Status = enums.JobStatusFailed
	default:
		j.Status = enums.JobStatusProcessing
	}
}
