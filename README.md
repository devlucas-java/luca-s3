# Luca-S3 Transcode Worker

[🇧🇷 Versão em Português](README.pt-BR.md) | [🇪🇸 Versión en Español](README.es.md)

A high-performance video transcoding worker that converts videos to HLS (HTTP Live Streaming) format with adaptive bitrate streaming. Built with Go and designed to work seamlessly with MinIO object storage.

## 🎯 Overview

Luca-S3 is a specialized microservice focused solely on video transcoding. It's not a complete video management system - it's a **worker** that:

- ✅ Transcodes videos to HLS format with multiple resolutions
- ✅ Generates 2-second segments for smooth streaming
- ✅ Creates thumbnails for each resolution
- ✅ Supports 360° video with spatial metadata
- ✅ Resumes interrupted transcoding jobs
- ✅ Validates resolution constraints (won't upscale)

**What it doesn't do:**
- ❌ Video upload (use MinIO directly)
- ❌ Video metadata management (use your own database)
- ❌ User authentication (implement in your API gateway)
- ❌ Video playback (use HLS player in your frontend)

## 🏗️ Architecture

```
┌─────────────┐      ┌──────────────┐      ┌─────────────┐
│   Client    │─────▶│  gRPC API    │─────▶│   Worker    │
│  (Your App) │      │  (Port 50051)│      │  (FFmpeg)   │
└─────────────┘      └──────────────┘      └─────────────┘
                            │                      │
                            ▼                      ▼
                     ┌──────────────┐      ┌─────────────┐
                     │    Redis     │      │    MinIO    │
                     │ (Job State)  │      │  (Storage)  │
                     └──────────────┘      └─────────────┘
```

### Components

- **gRPC Server**: Exposes transcoding API
- **FFmpeg**: Handles video processing
- **MinIO**: Object storage (shared volume with worker)
- **Redis**: Job state management
- **Temporary Storage**: `/tmp/ffmpeg` for processing (ephemeral)

## 📁 Project Structure

```
luca-s3/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── application/
│   │   └── service/
│   │       ├── transcode_service.go          # Main transcoding logic
│   │       ├── hls_transcoder.go             # HLS generation
│   │       ├── video_analyzer.go             # Video metadata extraction
│   │       └── segment_thumbnail_generator.go # Thumbnail generation
│   ├── delivery/
│   │   └── grpc/
│   │       ├── handler.go          # gRPC handlers
│   │       ├── server.go           # gRPC server setup
│   │       └── pb/                 # Generated protobuf files
│   ├── domain/
│   │   ├── enums/
│   │   │   ├── extension.go        # Supported formats (MP4, MOV, WEBM)
│   │   │   ├── job_status.go       # Job states
│   │   │   ├── resolution.go       # Supported resolutions
│   │   │   └── video_type.go       # Normal/360° video
│   │   └── model/
│   │       └── job.go              # TranscodeJob model
│   └── infrastructure/
│       ├── minio/
│       │   ├── client.go           # MinIO client
│       │   └── client_test.go      # MinIO tests
│       └── redis/
│           ├── client.go           # Redis client
│           ├── client_test.go      # Redis tests
│           ├── job_repository.go   # Job persistence
│           └── job_repository_test.go
├── pkg/
│   └── logger/
│       └── logger.go               # Logging utility
├── proto/
│   └── video.proto                 # gRPC service definition
├── configs/
│   └── config.go                   # Configuration management
├── docker-compose.yaml             # Docker orchestration
├── Dockerfile                      # Worker container
├── .default.env                    # Environment variables
└── README.md
```

## 🚀 Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.23+ (for local development)
- FFmpeg (for local development)

### Running with Docker

```bash
# Start all services
docker-compose up -d

# Check logs
docker-compose logs -f worker

# Stop services
docker-compose down
```

Services will be available at:
- **Worker gRPC**: `localhost:50051`
- **MinIO Console**: `http://localhost:9001` (username/password)
- **MinIO API**: `localhost:9000`
- **Redis**: `localhost:6379`

### Local Development

```bash
# Install dependencies
go mod download

# Run tests
go test ./internal/infrastructure/...

# Generate proto files
protoc --go_out=internal/delivery/grpc/pb \
       --go-grpc_out=internal/delivery/grpc/pb \
       --go_opt=paths=source_relative \
       --go-grpc_opt=paths=source_relative \
       proto/video.proto

# Run worker
go run cmd/server/main.go
```

## 🔧 Configuration

Environment variables (`.default.env`):

```env
# Redis Configuration
REDIS_ADDR=redis:6379
REDIS_PASSWORD=
REDIS_DB=0

# MinIO Configuration
MINIO_ENDPOINT=minio:9000
MINIO_ACCESS_KEY=username
MINIO_SECRET_KEY=password
MINIO_USE_SSL=false

# Server Configuration
SERVER_PORT=50051

# FFmpeg Configuration
FFMPEG_WORK_DIR=/tmp/ffmpeg
```

## 📡 API Reference

### TranscodeVideo

Start HLS transcoding for a video in MinIO.

```protobuf
rpc TranscodeVideo(TranscodeVideoRequest) returns (TranscodeVideoResponse);

message TranscodeVideoRequest {
  string video_id = 1;                    // Video ID in MinIO
  string original_path = 2;               // Path: videos/{id}.mp4
  repeated Resolution resolutions = 3;    // Desired resolutions
}
```

**Example:**
```go
req := &pb.TranscodeVideoRequest{
    VideoId:      "abc123",
    OriginalPath: "videos/abc123.mp4",
    Resolutions:  []pb.Resolution{
        pb.Resolution_RESOLUTION_1080P,
        pb.Resolution_RESOLUTION_720P,
        pb.Resolution_RESOLUTION_480P,
    },
}
```

### GetJobStatus

Query transcoding job status.

```protobuf
rpc GetJobStatus(JobStatusRequest) returns (JobStatusResponse);

message JobStatusResponse {
  string job_id = 1;
  string video_id = 2;
  JobStatus status = 3;
  repeated ResolutionProgress progress = 4;
  int32 original_width = 5;
  int32 original_height = 6;
  double duration_seconds = 7;
}
```

### GetHLSManifest

Get presigned URL for HLS master playlist.

```protobuf
rpc GetHLSManifest(GetHLSManifestRequest) returns (GetHLSManifestResponse);

message GetHLSManifestRequest {
  string video_id = 1;
  int32 expires_in = 2;  // seconds, default 3600
}
```

### DeleteVideo

Delete video and all related files from MinIO.

```protobuf
rpc DeleteVideo(DeleteVideoRequest) returns (DeleteVideoResponse);
```

## 🎬 How It Works

### 1. Video Upload (Your Responsibility)

Upload videos directly to MinIO:

```bash
# Using MinIO CLI
mc cp video.mp4 myminio/videos/abc123.mp4

# Or use MinIO SDK in your application
```

### 2. Transcode Request

Call the worker via gRPC:

```go
job, err := client.TranscodeVideo(ctx, &pb.TranscodeVideoRequest{
    VideoId:      "abc123",
    OriginalPath: "videos/abc123.mp4",
    Resolutions:  []pb.Resolution{pb.Resolution_RESOLUTION_720P},
})
```

### 3. Processing

The worker:
1. Downloads video from MinIO to `/tmp/ffmpeg/{job_id}/`
2. Analyzes video (resolution, duration)
3. Filters invalid resolutions (no upscaling)
4. Checks MinIO for existing segments (resume capability)
5. Generates HLS segments (2s each) with FFmpeg
6. Uploads segments to MinIO: `videos/{id}/hls/{resolution}/`
7. Generates thumbnails
8. Creates master playlist
9. Cleans up temporary files

### 4. Playback (Your Responsibility)

Use the HLS manifest URL in your video player:

```javascript
const player = new Hls();
player.loadSource('https://minio.example.com/videos/abc123/hls/master.m3u8');
player.attachMedia(videoElement);
```

## 🎯 Supported Formats

### Input Formats
- **MP4** (H.264/H.265)
- **MOV** (QuickTime)
- **WEBM** (VP8/VP9)

### Output Resolutions
- **4K** (3840x2160)
- **1440p** (2560x1440)
- **1080p** (1920x1080)
- **720p** (1280x720)
- **480p** (854x480)
- **360p** (640x360)

**Note:** Worker automatically filters out resolutions higher than the source video.

## 🔄 Resume Capability

The worker can resume interrupted jobs:

1. Checks MinIO for existing segments: `videos/{id}/hls/{res}/segment_*.ts`
2. Finds the last segment number
3. Continues from `last_segment + 1`
4. Updates job progress in Redis

This prevents re-processing already completed segments.

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run infrastructure tests (requires Redis & MinIO)
go test ./internal/infrastructure/...

# Run with coverage
go test -cover ./...

# Run specific test
go test -v ./internal/infrastructure/minio -run TestMinIOClient
```

## 🐳 Docker Volumes

- **luca_s3_data**: Shared between MinIO and worker (persistent)
- **/tmp/ffmpeg**: Temporary processing directory (ephemeral)

The shared volume allows the worker to access MinIO files directly without network transfer.

## 📊 Monitoring

### Job States

- **pending**: Job created, waiting to start
- **processing**: Currently transcoding
- **done**: Successfully completed
- **failed**: Error occurred

### Redis Keys

Jobs are stored in Redis with TTL of 7 days:
```
transcode:job:{job_id}
```

## 🔒 Security Considerations

1. **No Authentication**: Implement auth in your API gateway
2. **MinIO Access**: Use IAM policies to restrict bucket access
3. **Network**: Run worker in private network, expose only gRPC port
4. **Secrets**: Use Docker secrets or vault for credentials

## 🚧 Limitations

- No built-in rate limiting (implement in your API gateway)
- No job prioritization (FIFO processing)
- No distributed processing (single worker instance)
- No webhook notifications (implement in your application)

## 📝 License

MIT License - See LICENSE file for details

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Submit a pull request

## 📞 Support

For issues and questions:
- Open an issue on GitHub
- Check existing issues for solutions

---

**Built with ❤️ using Go, FFmpeg, MinIO, and Redis**
