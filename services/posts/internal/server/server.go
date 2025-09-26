package server

import (
	"fmt"
	"net"
	"soa-socialnetwork/services/posts/internal/service"
	"soa-socialnetwork/services/posts/internal/service/interceptors"
	pb "soa-socialnetwork/services/posts/proto"

	"google.golang.org/grpc"
)

type PostsServer struct {
	grpcServer *grpc.Server
	service    *service.PostsService
}

func Create(cfg service.PostsServiceConfig) (PostsServer, error) {
	service, err := service.New(cfg)
	if err != nil {
		return PostsServer{}, err
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.WithAuth(service.JwtVerifier),
		),
	)

	pb.RegisterPostsServiceServer(grpcServer, &service)
	return PostsServer{
		grpcServer: grpcServer,
		service:    &service,
	}, nil
}

func (s *PostsServer) Run(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	s.service.Start()
	return s.grpcServer.Serve(lis)
}
