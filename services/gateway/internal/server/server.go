package server

import (
	"fmt"
)

type GatewayServer struct {
	context GatewayService
	router  httpRouter
	port    int
}

func (s *GatewayServer) Run() error {
	return s.router.Run(fmt.Sprintf(":%d", s.port))
}

func Create(port int, cfg GatewayServiceConfig) (GatewayServer, error) {
	service := initService(cfg)
	router := createRouter(&service)
	return GatewayServer{
		context: service,
		router:  router,
		port:    port,
	}, nil
}
