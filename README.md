# Luca-S3 Transcode Worker

[🇧🇷 Versão em Português](README.pt-BR.md) | [🇪🇸 Versión en Español](README.es-ES.md)

* [📊🧩 DRAW.IO ARCHITECTURE DIAGRAM 🖍️🌐☁️](https://viewer.diagrams.net/?tags=%7B%7D&lightbox=1&highlight=0000ff&edit=_blank&layers=1&nav=1&title=architecture.drawio&dark=1#R%3Cmxfile%3E%3Cdiagram%20name%3D%22Alto%20Nivel%22%20id%3D%22alto%22%3E7Zhtb5swEMc%2FDdL2YpXBAZKXbdq1nVq1Wqp1e%2BnABdw6dmScJt2nnw02gaBpW5p0ShQpQtzfh4%2Fc%2FfyEh4fT5aUks%2FxWpMC8AKVLD597QRDGgb4a4bUSer1BJWSSppXkr4QR%2FQlWRFad0xSKlqMSgik6a4uJ4BwS1dKIlGLRdpsI1o46Ixl0hFFCWFd9pKnKrepHg1XDFdAst6H7QVw1TIlztv%2BkyEkqFg0JX3hB5AV4STx85hntTz88lEKoDR5sdzJdDoGZKrkKVG%2F4efsd19mTwNWOY1nQXgib29qdzhhNiDfE3mlcXrHQDtdcgeSkiq4xw2cfCpjrhjFJnoGnH2251KtjQIo5T8GE8bX3IqcKRjOSmNaFhl5ruZoy2zyhjA0FE1LbXHDtdFYoKZ7BiTpqlPRhPDHOgquRDeQ7uxoDPi7zsp4%2Bm9EXkAqWDelt6VxP7SWIKSj5qm0bpG8ptsMYI2svGoPC%2BeSN8RBbjdhxmNU9b4UC3YcFYes84Q5PN%2FOEfBppHT0K%2BQyygVD29X5Y9nMaIhT65jY2wtXDw32l95FLzy7J6gdjHEV7RVYvWENr0EWrnkQPA61eB61byq%2FvGjzdjZ%2FMahagkRLSLFElQwPkht0uGUqjcRTuF0Nx1GaonoqaDLncHQZDYYehr5DSosHQFzEuASLK4hPheLB7fMb9sBei%2FcInbOPTczgdLj5RB597Rl7NmoaubkbNzdFE6rq9z75oEMWY7NfMs7Yv8oMuOQe2L3InswY6%2Fom25zMmSGoa9IlNdGCBVJ%2BrrCmkykUmOGEXK7UJCyyp%2Bt64%2F6Hv0UloLK6TZJqQM1Ztv6WiEHOZQOuYoIjMQLWWY%2FOG70GOBEYUfWkfMPeQA7%2FDQWA4eJCEF4lI4VvFATL%2FQ2%2BZiUryckmqp5YbWiitFLXPNS9mettjH6z9LkHpOemWcDqBQtXe58BAuShvoW1TbvCRm0246R7UseEmFQtuZxAvMJFMns13lMKVvVyYdl5ofJwgtlTo7gm6Zwr9VO5Li2pf%2Bj%2FKGR7LuUk5u6fW8KSsngQypTxzIxSlVIJ6jzk5Og7Vv%2F5yrJ%2F%2Bh4%2FHNlrjmz6%2B%2BAU%3D%3C%2Fdiagram%3E%3Cdiagram%20name%3D%22Medio%20Nivel%22%20id%3D%22medio%22%3E1Zldb5swFIZ%2FDdJ2UQlsMHCZZm3XVVWrZtq0SwcMWHNsZJwm6a%2BfSUz4CJPSDIZ608bHx8bvc8zxiWPB%2BWp7J3GePYqYMAvY8daCXywAHORB%2Fa%2B07A6WwAkPhlTS2DjVhgV9I8ZoG%2BuaxqRoOSohmKJ52xgJzkmkWjYspdi03RLB2k%2FNcUpODIsIs1PrTxqr7KjLrTu%2BEppm1aMdZASucOVtpBQZjsWmYYI3FpxLIdTh02o7J6ykV4E5jLv9S%2B9xZZJwdc4AcDrAzFGoXaVXijWPSTnCseD1JqOKLHIclb0bHWFty9SKme6EMjYXTMj9WBhjEiSRthdKit%2Bk0YOigCyTcoTgamGe5lTtQ9AdDeT6sJ5XzNZmPbOc0Qhbc2jN%2FP1fKLTDPVdEcmy8iVRk2xBk5N8RsSJK7rRL1ogQMuHYNMNpbNu2i9mzMDBtbDZTepy5xq0%2FGOL99OHo9D0SxG4f%2FQAsIUId2k4f7fTlWQuxF0RqqBZAerB%2BEbUjnHm27TmDAXfbwEHYJg68AYh7YxNPkgREvfs9RkvknUX8u8S8iHTaLJlT%2FdyRtrTrtgmjAQCjjwD4nhe5PhhGxuvBNl7XGYCvPzrfICL9fJeB53r2OXyfaU4Y5Zqs%2FSkVejVKNz43ckeS5FIsy37rBlgBssL9QZyscpJ2jCpbr5YcUzZUmIDdDlNgd%2FLMEJk9GDtMxNG53e8LU4h8iM96DR4pv3%2FSpjmj5SpHyuOO3XkPhuAbfgS%2BLySmhTZ9E8sXkouCKmFUjoHZ6VQonj0A56pOnT6fv6dANBu7WaqEdoVjBPSwc5CCQdD3FPPTpPr3oDd7vokeQT%2F8b%2BjdIdBXCYvEJ18ET2Mh1jIyXuAkPOUMFT8hVSZSwTG7qa2GoMIyJWZieAK1H5UkDCv62l7hP8kGl8mGg8j2JpPtTSkbTSYbXSbbG0S2P5ls%2FzLZ%2FiCyg8lkB1NGOxxa9n7oTEq8azjkgnJVNGZ%2BLg2NwwJ0KtGqZKoZHqasiR7Xdh7k8DLI6GPvLffC4zIYRHZVok6g27lMdziMbmcq3W8%2FwCJ6C4j9sHx%2B%2BBVewafb6yv3rHJVl3%2BqXZP2laztepQLTjrFqzFhRlOumxEp7361oSwwaYTZzHSsaBzv8dVlst17o1zsolI%2BsNPyKrn0Vvv7iwuqVthTtYJ21ep38tAFF2%2B6Wf9KcEhV9Y8t8OYP%3C%2Fdiagram%3E%3Cdiagram%20name%3D%22Baixo%20Nivel%22%20id%3D%22baixo%22%3E1Zpbc9o6EMc%2FjWd6Hg6DLV8fE5o0OdPOyQSmfWSELWw1tuWRRQj99GeN5Sv0FBxxyQtjdleS9%2F%2BzpLVAQ5Pk7QvHWfSNBSTWjHHwpqHPmmHYHoLPwrApDbrnuqUl5DSQtsYwpb%2BINI6ldUUDkncCBWOxoFnX6LM0Jb7o2DDnbN0NW7K4O2qGQ7JjmPo43rX%2BoIGIpNW2zMbxQGgYVUPrtld6ElxFy1TyCAds3TKhOw1NOGOivEreJiQu1KuEKdvd%2F8Zb3xknqTikgbHbQPaRi02VL2erNCBFC11Dt%2BuICjLNsF9410AYbJFIYule0jiesJjxbVsUWMQNTLDngrMX0vK4xgLZdtGCpUJCBsLothz%2FFccrOX74%2FDQBywNOg5hwzbCh9TYyKi2jkMk2hAvy1kpDJv2FsIQIvoGQqMXFkhDWLYimtMleqq%2FVo1pxw%2FIRCuueG5HhQuq8X3Pz1Jovl0vD9%2FdpHtgL2zpI8xnHae7D1J0S%2Fkph3EZ1UbnmeelTqb%2Fd1b%2BOkQCQAv2tj6D%2FY5pnsHbtqk9Lxzm0d3va6wrEt08uvuuT%2FeIvXMu0xoeIj1Mcb36R77DNsE9%2FtcR%2FLSxz6S7Xndq3XGacLcggGvYeGm6Xhu10aXgKYDgfAUa93MzYw9dpB0cU5%2FPavYMjyUhYpELCBNLL4dLIz0WnnkvvweN%2BBDwhSQnHgsyiVbJIMY07hKT4c1F55zKe%2FQbXefAo2Ue8U%2BMhOlRPzj48nu0gfNA%2B8o2mj%2F%2BCaRLT4i4bwROaUgBw728dXRhTqNDHN0%2BPynD0qiqvV1ZZKnaWqo%2Fr5vEPWzyTjOUUZsCmJTknAYXV6f4nW8x5HdDFUi91GrqBOPjUnFtI3%2Fl8KlCu3QX1txJQe95SLlOAydH0IyZSDWOx8l8ITKjx9mWtsVdIRklm7jGDOLBrwedoNBpEzf0zNV1HveXOVIHt5O%2BKh25Hx2B7LudVzWE2%2BwoBTnH8QHHbkZE0oGmo3Rmaa2ueAcWcT%2FK8MLUoVt6ApVDp3S9hR4NkT8axt046KjCia8H4f9igS0H4dxzT4JnkLF4JytL84Fo8Jwn0t8ry7aHNmYoKNLYU4Dn96YACPD6MD%2B%2BmZWF9cEmev9BsW5BrxdlcksVEDHt5HQLHVTB3qmWVBDsHhLu02Ir7MsrYAVj0UK1jjIuIhQwe5LvGKjUXmIdEdmzuYNivFCcxFvS1e4fvSlu%2FZNrWxdI2h6VtKknbVp32tukN53jTCsgY3b4a1z0%2FFYZmGller1K32gr%2BMVzvHj%2FDRXkDjf51JochsS6JxLkOJM5xSPonE55iJPYlkbjXgcQ6DkkvvNrPlCFxLomkKjAvzQQdx6QXjpCllol7USbKa4dhTPTjmPTCTV3xPPGGMXGUMPFUIzn498eBVax1nWkPehJR72DLMtU%2BWubAklnRdB9f7NkyhuXtqclbv1jeaFjeFan3Jm6cKXH42vw5pZwXzX980N1%2F%3C%2Fdiagram%3E%3C%2Fmxfile%3E#%7B%22pageId%22%3A%22baixo%22%7D)



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
