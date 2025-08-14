package server

import (
	"fmt"
	"net"

	pb "soa-socialnetwork/services/accounts/proto"

	"google.golang.org/grpc"
)

type AccountsServiceServer struct {
	pb.UnimplementedAccountsServiceServer

	listener net.Listener
	server   *grpc.Server
}

func Create(port int) (*AccountsServiceServer, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := grpc.NewServer()
	pb.RegisterAccountsServiceServer(s, &AccountsServiceServer{})

	return &AccountsServiceServer{
		listener: lis,
		server:   s,
	}, nil
}

func (s *AccountsServiceServer) Run() error {
	return s.server.Serve(s.listener)
}
