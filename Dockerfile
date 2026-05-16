FROM golang:1.25-alpine AS builder

# Install FFmpeg and build dependencies
RUN apk add --no-cache \
    ffmpeg \
    git \
    ca-certificates

WORKDIR /app

# Copy go mod files and config
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/worker ./cmd/server

# Final stage
FROM alpine:latest

# Install FFmpeg and runtime dependencies
RUN apk add --no-cache \
    ffmpeg \
    ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/worker .

# Copy config file
COPY .default.env .

# Create temp directory for FFmpeg work
RUN mkdir -p /tmp/ffmpeg && chmod 777 /tmp/ffmpeg

EXPOSE 50051

CMD ["./worker"]
