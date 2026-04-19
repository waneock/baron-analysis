package grpc

import (
	"fmt"
	"net"
	pb "skinbaron-analyzer/proto/parsing/v1"

	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	addr       string
}

func NewServer(addr string, handler pb.ParsingServiceServer, opts ...grpc.ServerOption) *Server {
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterParsingServiceServer(grpcServer, handler)

	return &Server{
		grpcServer: grpcServer,
		addr:       addr,
	}
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("listen tcp: %w", err)
	}

	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("serve grpc error: %w", err)
	}

	return nil
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}
