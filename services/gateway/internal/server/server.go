package server

import (
	"fmt"
	"soa-socialnetwork/services/gateway/internal/service"
)

type GatewayServer struct {
	service service.GatewayService
	router  httpRouter
}

func (s *GatewayServer) Run(port int) error {
	return s.router.Run(fmt.Sprintf(":%d", port))
}

func Create(cfg service.Config) (GatewayServer, error) {
	service := service.NewGatewayService(cfg)
	router := newHttpRouter(&service)
	return GatewayServer{
		service: service,
		router:  router,
	}, nil
}
