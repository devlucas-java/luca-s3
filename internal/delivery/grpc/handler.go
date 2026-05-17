package grpc

import (
	"context"
	"fmt"
	"strings"
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
	inspectService   *service.InspectService
	minio            *minioclient.Client
}

func NewHandler(transcodeService *service.TranscodeService, inspectService *service.InspectService, minio *minioclient.Client) *Handler {
	return &Handler{
		transcodeService: transcodeService,
		inspectService:   inspectService,
		minio:            minio,
	}
}

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
		msg := err.Error()
		switch {
		case containsAny(msg, "not found in minio", "no files"):
			return nil, status.Errorf(codes.NotFound, "%v", err)
		case containsAny(msg, "already has an active job"):
			return nil, status.Errorf(codes.AlreadyExists, "%v", err)
		default:
			return nil, status.Errorf(codes.Internal, "transcode: %v", err)
		}
	}

	return &pb.TranscodeVideoResponse{
		JobId:   job.ID,
		VideoId: job.VideoID,
	}, nil
}

func (h *Handler) WatchJob(req *pb.WatchJobRequest, stream pb.VideoService_WatchJobServer) error {
	if req.VideoId == "" {
		return status.Error(codes.InvalidArgument, "video_id is required")
	}

	interval := time.Duration(req.IntervalMs) * time.Millisecond
	if interval < 200*time.Millisecond || interval > 5*time.Second {
		interval = 500 * time.Millisecond
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case <-ticker.C:
			job, err := h.transcodeService.GetJob(stream.Context(), req.VideoId)
			if err != nil {
				return status.Errorf(codes.NotFound, "no job found for video: %v", err)
			}

			if err := stream.Send(jobToProto(job)); err != nil {
				return err
			}

			if job.Status == enums.JobStatusDone || job.Status == enums.JobStatusFailed {
				return nil
			}
		}
	}
}

func (h *Handler) ListJobs(ctx context.Context, req *pb.ListJobsRequest) (*pb.ListJobsResponse, error) {
	page := int(req.Page)
	pageSize := int(req.PageSize)

	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	filterStatus := protoToJobStatus(req.FilterStatus)

	jobs, total, err := h.transcodeService.ListJobs(ctx, filterStatus, page, pageSize)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list jobs: %v", err)
	}

	items := make([]*pb.JobStatusResponse, len(jobs))
	for i, j := range jobs {
		items[i] = jobToProto(j)
	}

	return &pb.ListJobsResponse{
		Jobs:     items,
		Total:    int32(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}

func (h *Handler) GetHLSManifest(ctx context.Context, req *pb.GetHLSManifestRequest) (*pb.GetHLSManifestResponse, error) {
	if req.VideoId == "" {
		return nil, status.Error(codes.InvalidArgument, "video_id is required")
	}

	manifestPath := fmt.Sprintf("%s/hls/master.m3u8", req.VideoId)

	exists, err := h.minio.ObjectExists(ctx, manifestPath)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "check manifest: %v", err)
	}
	if !exists {
		return nil, status.Errorf(codes.NotFound, "HLS manifest not found for video: %s", req.VideoId)
	}

	return &pb.GetHLSManifestResponse{
		Url: h.minio.PublicURL(manifestPath),
	}, nil
}

// InspectVideo scans MinIO and returns a full report of what exists for a video.
func (h *Handler) InspectVideo(ctx context.Context, req *pb.InspectVideoRequest) (*pb.InspectVideoResponse, error) {
	if req.VideoId == "" {
		return nil, status.Error(codes.InvalidArgument, "video_id is required")
	}

	info, err := h.inspectService.InspectVideo(ctx, req.VideoId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "inspect: %v", err)
	}
	if !info.Found {
		return nil, status.Errorf(codes.NotFound, "no files found in MinIO for video_id %q", req.VideoId)
	}

	resInfos := make([]*pb.ResolutionInfo, len(info.Resolutions))
	for i, r := range info.Resolutions {
		resInfos[i] = &pb.ResolutionInfo{
			Resolution:   resolutionToProto(r.Resolution),
			SegmentCount: int32(r.SegmentCount),
			HasPlaylist:  r.HasPlaylist,
			HasThumbnail: r.HasThumbnail,
			SizeBytes:    r.SizeBytes,
		}
	}

	return &pb.InspectVideoResponse{
		VideoId:        info.VideoID,
		Found:          info.Found,
		OriginalPath:   info.OriginalPath,
		OriginalSize:   info.OriginalSize,
		HasMaster:      info.HasMaster,
		Resolutions:    resInfos,
		TotalSizeBytes: info.TotalSizeBytes,
	}, nil
}

func (h *Handler) DeleteVideo(ctx context.Context, req *pb.DeleteVideoRequest) (*pb.DeleteVideoResponse, error) {
	if req.VideoId == "" {
		return nil, status.Error(codes.InvalidArgument, "video_id is required")
	}

	if err := h.transcodeService.DeleteVideo(ctx, req.VideoId); err != nil {
		// Map known error types to proper gRPC codes
		msg := err.Error()
		switch {
		case containsAny(msg, "not found", "no files"):
			return nil, status.Errorf(codes.NotFound, "%v", err)
		case containsAny(msg, "currently processing", "currently pending"):
			return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
		default:
			return nil, status.Errorf(codes.Internal, "%v", err)
		}
	}

	return &pb.DeleteVideoResponse{
		Success: true,
		Message: fmt.Sprintf("video %s and all related files deleted", req.VideoId),
	}, nil
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
