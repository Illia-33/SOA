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
	service    *service.StatsService
}

func New(cfg service.Config) (Server, error) {
	service, err := service.New(cfg)
	if err != nil {
		return Server{}, err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterStatsServiceServer(grpcServer, &service)

	return Server{
		grpcServer: grpcServer,
		service:    &service,
	}, nil
}

func (s *Server) Run(port int) error {
	s.service.Start()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	return s.grpcServer.Serve(lis)
}
