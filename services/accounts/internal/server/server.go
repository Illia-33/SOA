package server

import (
	"fmt"
	"net"

	pb "soa-socialnetwork/services/accounts/proto"

	"google.golang.org/grpc"
)

type Server struct {
	listener   net.Listener
	grpcServer *grpc.Server
}

func Create(port int, cfg AccountsServiceConfig) (*Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer()
	service, err := createAccountsService(cfg)
	if err != nil {
		return nil, err
	}

	pb.RegisterAccountsServiceServer(grpcServer, service)

	return &Server{
		listener:   lis,
		grpcServer: grpcServer,
	}, nil
}

func (s *Server) Run() error {
	return s.grpcServer.Serve(s.listener)
}
