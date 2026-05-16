package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/devlucas-java/luca-s3/internal/application/service"
	"github.com/devlucas-java/luca-s3/internal/delivery/grpc/pb"
	"github.com/devlucas-java/luca-s3/internal/domain/enums"
	minioclient "github.com/devlucas-java/luca-s3/internal/infrastructure/minio"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	pb.UnimplementedVideoServiceServer
	transcodeService *service.TranscodeService
	minio            *minioclient.Client
}

func NewHandler(transcodeService *service.TranscodeService, minio *minioclient.Client) *Handler {
	return &Handler{
		transcodeService: transcodeService,
		minio:            minio,
	}
}

// TranscodeVideo inicia transcodificação HLS para um vídeo no MinIO.
func (h *Handler) TranscodeVideo(ctx context.Context, req *pb.TranscodeVideoRequest) (*pb.TranscodeVideoResponse, error) {
	if req.VideoId == "" {
		return nil, status.Error(codes.InvalidArgument, "video_id is required")
	}
	if req.OriginalPath == "" {
		return nil, status.Error(codes.InvalidArgument, "original_path is required")
	}
	if len(req.Resolutions) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one resolution is required")
	}

	resolutions := make([]enums.Resolution, len(req.Resolutions))
	for i, r := range req.Resolutions {
		resolutions[i] = protoToResolution(r)
	}

	job, err := h.transcodeService.TranscodeVideo(ctx, req.VideoId, req.OriginalPath, resolutions)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "transcode: %v", err)
	}

	return &pb.TranscodeVideoResponse{
		JobId:   job.ID,
		VideoId: job.VideoID,
	}, nil
}

// GetJobStatus retorna o status de um job.
func (h *Handler) GetJobStatus(ctx context.Context, req *pb.JobStatusRequest) (*pb.JobStatusResponse, error) {
	job, err := h.transcodeService.GetJob(ctx, req.JobId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "job not found: %v", err)
	}

	progress := make([]*pb.ResolutionProgress, len(job.Progress))
	for i, p := range job.Progress {
		progress[i] = &pb.ResolutionProgress{
			Resolution:    resolutionToProto(p.Resolution),
			Status:        jobStatusToProto(p.Status),
			LastSegment:   int32(p.LastSegment),
			TotalSegments: int32(p.TotalSegments),
			ThumbnailDone: p.ThumbnailDone,
			Error:         p.Error,
		}
	}

	return &pb.JobStatusResponse{
		JobId:           job.ID,
		VideoId:         job.VideoID,
		Status:          jobStatusToProto(job.Status),
		Progress:        progress,
		OriginalWidth:   int32(job.OriginalWidth),
		OriginalHeight:  int32(job.OriginalHeight),
		DurationSeconds: job.DurationSeconds,
	}, nil
}

// GetHLSManifest retorna a URL do master playlist HLS.
func (h *Handler) GetHLSManifest(ctx context.Context, req *pb.GetHLSManifestRequest) (*pb.GetHLSManifestResponse, error) {
	expiry := time.Duration(req.ExpiresIn) * time.Second
	if expiry == 0 {
		expiry = 1 * time.Hour
	}

	manifestPath := fmt.Sprintf("videos/%s/hls/master.m3u8", req.VideoId)

	// Verifica se existe
	exists, err := h.minio.ObjectExists(ctx, manifestPath)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "check manifest: %v", err)
	}
	if !exists {
		return nil, status.Errorf(codes.NotFound, "HLS manifest not found for video: %s", req.VideoId)
	}

	url, err := h.minio.PresignedURL(ctx, manifestPath, expiry)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "presigned url: %v", err)
	}

	return &pb.GetHLSManifestResponse{
		Url:       url,
		ExpiresIn: int32(expiry.Seconds()),
	}, nil
}

// DeleteVideo deletes all video files from MinIO and related jobs.
func (h *Handler) DeleteVideo(ctx context.Context, req *pb.DeleteVideoRequest) (*pb.DeleteVideoResponse, error) {
	if err := h.transcodeService.DeleteVideo(ctx, req.VideoId); err != nil {
		return &pb.DeleteVideoResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.DeleteVideoResponse{
		Success: true,
		Message: fmt.Sprintf("video %s deleted", req.VideoId),
	}, nil
}

// ─── Conversores ─────────────────────────────────────────────────────────────

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

func jobStatusToProto(status enums.JobStatus) pb.JobStatus {
	switch status {
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
