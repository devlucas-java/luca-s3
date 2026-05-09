# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -o server ./cmd/server

# Development stage
FROM golang:1.26-alpine

WORKDIR /app

# Copy built binary from builder
COPY --from=builder /build/server .

# Copy source code for development
COPY . .

# Copy .env file for development
COPY .default.env .default.env

# Expose port
EXPOSE 8080

# Run the application
CMD ["./server"]
