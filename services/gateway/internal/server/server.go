package server

import (
	"fmt"
	"soa-socialnetwork/services/gateway/internal/service"
)

type GatewayServer struct {
	context service.GatewayService
	router  httpRouter
	port    int
}

func (s *GatewayServer) Run() error {
	return s.router.Run(fmt.Sprintf(":%d", s.port))
}

func Create(port int, cfg service.Config) (GatewayServer, error) {
	service := service.NewGatewayService(cfg)
	router := newHttpRouter(&service)
	return GatewayServer{
		context: service,
		router:  router,
		port:    port,
	}, nil
}
