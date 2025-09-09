package server

import (
	"fmt"
	"net"
	"soa-socialnetwork/services/posts/internal/server/interceptors"
	pb "soa-socialnetwork/services/posts/proto"

	"google.golang.org/grpc"
)

type PostsServer struct {
	grpcServer *grpc.Server
	service    *PostsService
	port       int
}

func (s *PostsServer) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}

	s.service.Start()
	return s.grpcServer.Serve(lis)
}

func Create(cfg PostsServiceConfig) (PostsServer, error) {
	service, err := newPostsService(cfg)
	if err != nil {
		return PostsServer{}, err
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.WithAuth(&service.JwtVerifier),
		),
	)

	pb.RegisterPostsServiceServer(grpcServer, &service)
	return PostsServer{
		grpcServer: grpcServer,
		service:    &service,
		port:       cfg.Port,
	}, nil
}
