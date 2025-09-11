package server

import (
	"fmt"
	"net"
	"soa-socialnetwork/services/stats/internal/service"
	pb "soa-socialnetwork/services/stats/proto"

	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	service    service.StatsService
	port       int
}

func New(port int, cfg service.Config) (Server, error) {
	service, err := service.New(cfg)
	if err != nil {
		return Server{}, err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterStatsServiceServer(grpcServer, &service)

	return Server{
		grpcServer: grpcServer,
		service:    service,
		port:       port,
	}, nil
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}

	return s.grpcServer.Serve(lis)
}
