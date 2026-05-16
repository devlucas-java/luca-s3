package grpc

import (
	"fmt"
	"net"

	"github.com/devlucas-java/luca-s3/internal/delivery/grpc/pb"
	"github.com/devlucas-java/luca-s3/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	handler *Handler
}

func NewServer(handler *Handler) *Server {
	return &Server{handler: handler}
}

func (s *Server) Start(port string) error {
	log := logger.Instance()

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(100*1024*1024), // 100MB
		grpc.MaxSendMsgSize(100*1024*1024),
	)

	pb.RegisterVideoServiceServer(grpcServer, s.handler)
	reflection.Register(grpcServer)

	log.Infof("grpc: listening on :%s", port)
	return grpcServer.Serve(lis)
}
