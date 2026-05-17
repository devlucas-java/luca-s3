package model

import (
	"time"

	"github.com/devlucas-java/luca-s3/internal/domain/enums"
)

type ResolutionProgress struct {
	Resolution    enums.Resolution `json:"resolution"`
	Status        enums.JobStatus  `json:"status"`
	SegmentsDone  int              `json:"segments_done"`
	TotalSegments int              `json:"total_segments"`
	Percent       int              `json:"percent"`
	ThumbnailDone bool             `json:"thumbnail_done"`
	Error         string           `json:"error,omitempty"`
}

type TranscodeJob struct {
	ID              string               `json:"id"`
	VideoID         string               `json:"video_id"`
	OriginalPath    string               `json:"original_path"`
	RequestedRes    []enums.Resolution   `json:"requested_res"`
	Status          enums.JobStatus      `json:"status"`
	Progress        []ResolutionProgress `json:"progress"`
	OriginalWidth   int                  `json:"original_width"`
	OriginalHeight  int                  `json:"original_height"`
	DurationSeconds float64              `json:"duration_seconds"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
}

func (j *TranscodeJob) OverallPercent() int {
	if len(j.Progress) == 0 {
		return 0
	}
	total := 0
	for _, p := range j.Progress {
		total += p.Percent
	}
	return total / len(j.Progress)
}

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
