package server

import (
	"fmt"
)

type GatewayServer struct {
	context gatewayServiceContext
	router  httpRouter
	port    int
}

func (s *GatewayServer) Run() error {
	return s.router.Run(fmt.Sprintf(":%d", s.port))
}

func Create(port int) (GatewayServer, error) {
	context := createContext()
	router := createRouter(&context)
	return GatewayServer{
		context: context,
		router:  router,
		port:    port,
	}, nil
}
