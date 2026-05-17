package grpc

import (
	"time"

	"github.com/devlucas-java/luca-s3/internal/delivery/grpc/pb"
	"github.com/devlucas-java/luca-s3/internal/domain/enums"
	"github.com/devlucas-java/luca-s3/internal/domain/model"
)

func jobToProto(job *model.TranscodeJob) *pb.JobStatusResponse {
	progress := make([]*pb.ResolutionProgress, len(job.Progress))
	for i, p := range job.Progress {
		progress[i] = &pb.ResolutionProgress{
			Resolution:    resolutionToProto(p.Resolution),
			Status:        jobStatusToProto(p.Status),
			SegmentsDone:  int32(p.SegmentsDone),
			TotalSegments: int32(p.TotalSegments),
			Percent:       int32(p.Percent),
			ThumbnailDone: p.ThumbnailDone,
			Error:         p.Error,
		}
	}
	return &pb.JobStatusResponse{
		JobId:           job.ID,
		VideoId:         job.VideoID,
		Status:          jobStatusToProto(job.Status),
		Progress:        progress,
		OverallPercent:  int32(job.OverallPercent()),
		OriginalWidth:   int32(job.OriginalWidth),
		OriginalHeight:  int32(job.OriginalHeight),
		DurationSeconds: job.DurationSeconds,
		CreatedAt:       job.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       job.UpdatedAt.Format(time.RFC3339),
	}
}

func protoToResolution(res pb.Resolution) enums.Resolution {
	switch res {
	case pb.Resolution_RESOLUTION_360P:
		return enums.RESOLUTION_360P
	case pb.Resolution_RESOLUTION_480P:
		return enums.RESOLUTION_480P
	case pb.Resolution_RESOLUTION_720P:
		return enums.RESOLUTION_720P
	case pb.Resolution_RESOLUTION_1080P:
		return enums.RESOLUTION_1080P
	case pb.Resolution_RESOLUTION_1440P:
		return enums.RESOLUTION_1440P
	case pb.Resolution_RESOLUTION_4K:
		return enums.RESOLUTION_4K
	default:
		return ""
	}
}

func resolutionToProto(res enums.Resolution) pb.Resolution {
	switch res {
	case enums.RESOLUTION_360P:
		return pb.Resolution_RESOLUTION_360P
	case enums.RESOLUTION_480P:
		return pb.Resolution_RESOLUTION_480P
	case enums.RESOLUTION_720P:
		return pb.Resolution_RESOLUTION_720P
	case enums.RESOLUTION_1080P:
		return pb.Resolution_RESOLUTION_1080P
	case enums.RESOLUTION_1440P:
		return pb.Resolution_RESOLUTION_1440P
	case enums.RESOLUTION_4K:
		return pb.Resolution_RESOLUTION_4K
	default:
		return pb.Resolution_RESOLUTION_UNKNOWN
	}
}

func jobStatusToProto(s enums.JobStatus) pb.JobStatus {
	switch s {
	case enums.JobStatusPending:
		return pb.JobStatus_JOB_STATUS_PENDING
	case enums.JobStatusProcessing:
		return pb.JobStatus_JOB_STATUS_PROCESSING
	case enums.JobStatusDone:
		return pb.JobStatus_JOB_STATUS_DONE
	case enums.JobStatusFailed:
		return pb.JobStatus_JOB_STATUS_FAILED
	default:
		return pb.JobStatus_JOB_STATUS_UNKNOWN
	}
}

func protoToJobStatus(s pb.JobStatus) enums.JobStatus {
	switch s {
	case pb.JobStatus_JOB_STATUS_PENDING:
		return enums.JobStatusPending
	case pb.JobStatus_JOB_STATUS_PROCESSING:
		return enums.JobStatusProcessing
	case pb.JobStatus_JOB_STATUS_DONE:
		return enums.JobStatusDone
	case pb.JobStatus_JOB_STATUS_FAILED:
		return enums.JobStatusFailed
	default:
		return ""
	}
}
