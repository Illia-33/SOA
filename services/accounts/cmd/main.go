package main

import (
	"log"

	"soa-socialnetwork/services/accounts/internal/server"
	"soa-socialnetwork/services/accounts/internal/service"
	"soa-socialnetwork/services/common/envvar"
)

func extractServiceConfig() service.Config {
	return service.Config{
		DbHost:        envvar.MustStringFromEnv("DB_HOST"),
		DbUser:        envvar.MustStringFromEnv("DB_USER"),
		DbPassword:    envvar.MustStringFromEnv("DB_PASSWORD"),
		DbPoolSize:    envvar.MustIntFromEnv("DB_POOL_SIZE"),
		JwtPrivateKey: envvar.MustEd25519PrivKeyFromEnv("JWT_ED25519_PRIVATE_KEY"),
	}
}

func main() {
	log.Println("Accounts Service")

	port := envvar.MustIntFromEnv("ACCOUNTS_SERVICE_PORT")
	cfg := extractServiceConfig()

	log.Printf("creating server with port %d...", port)
	s, err := server.Create(cfg)
	if err != nil {
		log.Fatalf("initializing server failed: %v", err)
	}
	log.Println("success")

	log.Println("running server...")
	err = s.Run(port)
	if err != nil {
		log.Fatalf("server error: %v", err)
	}
}
