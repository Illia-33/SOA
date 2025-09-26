package server

import (
	"fmt"
	"net"

	"soa-socialnetwork/services/accounts/internal/service"
	"soa-socialnetwork/services/accounts/internal/service/interceptors"
	pb "soa-socialnetwork/services/accounts/proto"

	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	service    *service.AccountsService
}

func Create(cfg service.Config) (*Server, error) {
	service, err := service.NewAccountsService(cfg)
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.ConvertErrors(),
			interceptors.Auth(service),
		),
	)
	pb.RegisterAccountsServiceServer(grpcServer, service)

	return &Server{
		grpcServer: grpcServer,
		service:    service,
	}, nil
}

func (s *Server) Run(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	s.service.Start()
	return s.grpcServer.Serve(lis)
}
