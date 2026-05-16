#!/bin/bash

# Create pb directory if it doesn't exist
mkdir -p internal/delivery/grpc/pb

# Generate protobuf files
protoc --go_out=internal/delivery/grpc/pb \
       --go-grpc_out=internal/delivery/grpc/pb \
       --go_opt=paths=source_relative \
       --go-grpc_opt=paths=source_relative \
       proto/video.proto

# Move files if they were generated in proto subdirectory
if [ -d "internal/delivery/grpc/pb/proto" ]; then
    mv internal/delivery/grpc/pb/proto/*.go internal/delivery/grpc/pb/ 2>/dev/null
    rmdir internal/delivery/grpc/pb/proto 2>/dev/null
fi

echo "✅ Proto files generated successfully"
ls -la internal/delivery/grpc/pb/*.pb.go
