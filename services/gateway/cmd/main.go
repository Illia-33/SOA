package main

import (
	"log"

	"soa-socialnetwork/services/common/envvar"
	"soa-socialnetwork/services/gateway/internal/server"
	"soa-socialnetwork/services/gateway/internal/service"
)

func extractConfig() service.Config {
	return service.Config{
		JwtPublicKey:        envvar.MustEd25519PubKeyFromEnv("JWT_ED25519_PUBLIC_KEY"),
		AccountsServiceHost: envvar.MustStringFromEnv("ACCOUNTS_SERVICE_HOST"),
		AccountsServicePort: envvar.MustIntFromEnv("ACCOUNTS_SERVICE_PORT"),
		PostsServiceHost:    envvar.MustStringFromEnv("POSTS_SERVICE_HOST"),
		PostsServicePort:    envvar.MustIntFromEnv("POSTS_SERVICE_PORT"),
		StatsServiceHost:    envvar.MustStringFromEnv("STATS_SERVICE_HOST"),
		StatsServicePort:    envvar.MustIntFromEnv("STATS_SERVICE_PORT"),
	}
}

func main() {
	log.Println("Gateway Service")

	port := envvar.MustIntFromEnv("GATEWAY_SERVICE_PORT")
	cfg := extractConfig()

	log.Printf("creating server...")
	s, err := server.Create(cfg)
	if err != nil {
		log.Fatalf("initializing server failed: %v", err)
	}
	log.Println("success")

	log.Printf("running server at port %d...", port)
	err = s.Run(port)
	if err != nil {
		log.Fatalf("server error: %v", err)
	}
}
