package main

import (
	"github.com/devlucas-java/luca-s3/configs"
	"github.com/devlucas-java/luca-s3/internal/application/service"
	grpcserver "github.com/devlucas-java/luca-s3/internal/delivery/grpc"
	minioclient "github.com/devlucas-java/luca-s3/internal/infrastructure/minio"
	redisclient "github.com/devlucas-java/luca-s3/internal/infrastructure/redis"
	"github.com/devlucas-java/luca-s3/pkg/logger"
)

func main() {
	log := logger.Instance()
	log.Info("starting luca-s3 transcode worker")

	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	if err := cfg.Validate(); err != nil {
		log.Fatalf("config validate: %v", err)
	}

	// Redis for job state
	redis, err := redisclient.New(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("redis: %v", err)
	}
	defer redis.Close()

	jobRepo := redisclient.NewJobRepository(redis)

	// MinIO client
	minio, err := minioclient.New(cfg.MinIOEndpoint, cfg.MinIOAccessKey, cfg.MinIOSecretKey, cfg.MinIOUseSSL)
	if err != nil {
		log.Fatalf("minio: %v", err)
	}

	transcodeSvc := service.NewTranscodeService(minio, jobRepo)
	handler := grpcserver.NewHandler(transcodeSvc, minio)
	grpcSrv := grpcserver.NewServer(handler)

	log.Info("worker ready - MinIO is the source of truth, Redis for job state")
	if err := grpcSrv.Start(cfg.ServerPort); err != nil {
		log.Fatalf("grpc server: %v", err)
	}
}
