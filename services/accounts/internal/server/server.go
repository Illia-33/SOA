package server

import (
	"fmt"
	"net"

	"soa-socialnetwork/services/accounts/internal/server/interceptors"
	pb "soa-socialnetwork/services/accounts/proto"

	"google.golang.org/grpc"
)

type Server struct {
	listener   net.Listener
	grpcServer *grpc.Server
	service    *AccountsService
}

func Create(port int, cfg AccountsServiceConfig) (*Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	service, err := createAccountsService(cfg)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptors.Auth(service.jwtVerifier, &service.soaVerifier)),
	)

	if err != nil {
		return nil, err
	}

	pb.RegisterAccountsServiceServer(grpcServer, service)

	return &Server{
		listener:   lis,
		grpcServer: grpcServer,
		service:    service,
	}, nil
}

func (s *Server) Run() error {
	s.service.Start()
	return s.grpcServer.Serve(s.listener)
}
