# Quick Start Guide

## 🚀 Getting Started in 5 Minutes

### 1. Clone and Setup

```bash
git clone https://github.com/devlucas-java/luca-s3.git
cd luca-s3
```

### 2. Start Services

```bash
# Start Redis, MinIO, and Worker
docker-compose up -d

# Check if services are running
docker-compose ps
```

### 3. Upload a Video to MinIO

```bash
# Install MinIO client (if not installed)
# Linux/Mac: brew install minio/stable/mc
# Or download from: https://min.io/docs/minio/linux/reference/minio-mc.html

# Configure MinIO client
mc alias set local http://localhost:9000 username password

# Create bucket (if not exists)
mc mb local/videos

# Upload a video
mc cp your-video.mp4 local/videos/test-video.mp4
```

### 4. Test Transcoding

Create a simple test client:

```go
// test-client.go
package main

import (
    "context"
    "log"
    "time"

    pb "github.com/devlucas-java/luca-s3/internal/delivery/grpc/pb"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main() {
    // Connect to worker
    conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    client := pb.NewVideoServiceClient(conn)
    ctx := context.Background()

    // Start transcoding
    resp, err := client.TranscodeVideo(ctx, &pb.TranscodeVideoRequest{
        VideoId:      "test-video",
        OriginalPath: "videos/test-video.mp4",
        Resolutions: []pb.Resolution{
            pb.Resolution_RESOLUTION_720P,
            pb.Resolution_RESOLUTION_480P,
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Job started: %s", resp.JobId)

    // Check status
    for {
        status, err := client.GetJobStatus(ctx, &pb.JobStatusRequest{
            JobId: resp.JobId,
        })
        if err != nil {
            log.Fatal(err)
        }

        log.Printf("Status: %s", status.Status)
        
        if status.Status == pb.JobStatus_JOB_STATUS_DONE {
            log.Println("✅ Transcoding completed!")
            break
        }
        
        if status.Status == pb.JobStatus_JOB_STATUS_FAILED {
            log.Println("❌ Transcoding failed!")
            break
        }

        time.Sleep(5 * time.Second)
    }

    // Get HLS manifest URL
    manifest, err := client.GetHLSManifest(ctx, &pb.GetHLSManifestRequest{
        VideoId:   "test-video",
        ExpiresIn: 3600,
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("HLS Manifest: %s", manifest.Url)
}
```

Run the test:

```bash
go run test-client.go
```

### 5. Play the Video

Use any HLS player to play the transcoded video:

**HTML5 with hls.js:**

```html
<!DOCTYPE html>
<html>
<head>
    <script src="https://cdn.jsdelivr.net/npm/hls.js@latest"></script>
</head>
<body>
    <video id="video" controls width="640"></video>
    <script>
        const video = document.getElementById('video');
        const hls = new Hls();
        
        // Replace with your manifest URL
        hls.loadSource('http://localhost:9000/videos/test-video/hls/master.m3u8');
        hls.attachMedia(video);
    </script>
</body>
</html>
```

### 6. Check Logs

```bash
# Worker logs
docker-compose logs -f worker

# Redis logs
docker-compose logs -f redis

# MinIO logs
docker-compose logs -f minio
```

### 7. Stop Services

```bash
docker-compose down
```

## 🔧 Troubleshooting

### Proto files not generating?

```bash
# Make sure protoc is installed
protoc --version

# Regenerate proto files
./scripts/generate-proto.sh
```

### Build errors?

```bash
# Clean and rebuild
make clean
go mod tidy
make build
```

### Services not starting?

```bash
# Check Docker logs
docker-compose logs

# Restart services
docker-compose restart
```

### Can't connect to MinIO?

```bash
# Check MinIO is running
curl http://localhost:9000/minio/health/live

# Access MinIO console
open http://localhost:9001
# Login: username / password
```

### Redis connection issues?

```bash
# Test Redis connection
redis-cli -h localhost -p 6379 ping
# Should return: PONG
```

## 📚 Next Steps

- Read the full [README.md](README.md)
- Check [API Reference](README.md#-api-reference)
- See [Architecture](README.md#-architecture)
- Run [Tests](README.md#-testing)

## 🆘 Need Help?

- Open an issue on GitHub
- Check existing issues for solutions
- Read the documentation
