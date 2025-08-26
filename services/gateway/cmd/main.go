package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"

	"soa-socialnetwork/services/gateway/internal/server"
)

func extractEnv(key string) (string, error) {
	val, exists := os.LookupEnv(key)
	if !exists {
		return "", fmt.Errorf("environment variable %s doesn't exist", key)
	}

	return val, nil
}

func extractPort() (int, error) {
	port, exists := os.LookupEnv("GATEWAY_SERVICE_PORT")
	if !exists {
		return 8080, nil
	}

	num, err := strconv.Atoi(port)
	if err != nil {
		return 0, fmt.Errorf("cannot convert GATEWAY_SERVICE_PORT environment variable to number (%v)", err)
	}

	return num, nil
}

func extractJwtPublicKey() (ed25519.PublicKey, error) {
	envPublicKey, err := extractEnv("JWT_ED25519_PUBLIC_KEY")
	if err != nil {
		return nil, err
	}

	key := make(ed25519.PublicKey, ed25519.PublicKeySize)

	_, err = hex.Decode(key, []byte(envPublicKey))
	if err != nil {
		return nil, fmt.Errorf("cannot decode hex encoded jwt public key: %v", err)
	}

	return key, nil
}

func extractConfig() (server.GatewayServiceConfig, error) {
	jwtPubKey, err := extractJwtPublicKey()
	if err != nil {
		return server.GatewayServiceConfig{}, err
	}

	accountsServiceHost, err := extractEnv("ACCOUNTS_SERVICE_HOST")
	if err != nil {
		return server.GatewayServiceConfig{}, err
	}

	accountsServicePortStr, err := extractEnv("ACCOUNTS_SERVICE_PORT")
	if err != nil {
		return server.GatewayServiceConfig{}, err
	}

	accountsServicePort, err := strconv.Atoi(accountsServicePortStr)
	if err != nil {
		return server.GatewayServiceConfig{}, fmt.Errorf("bad ACCOUNTS_SERVICE PORT: %s", accountsServicePortStr)
	}

	return server.GatewayServiceConfig{
		JwtPublicKey:        jwtPubKey,
		AccountsServiceHost: accountsServiceHost,
		AccountsServicePort: accountsServicePort,
	}, nil
}

func main() {
	log.Println("Gateway Service")

	port, err := extractPort()
	if err != nil {
		log.Fatalf("failed to get port: %v", err)
	}

	cfg, err := extractConfig()
	if err != nil {
		log.Fatalf("failed to extract config: %v", err)
	}

	log.Printf("creating server with port %d...", port)
	s, err := server.Create(port, cfg)
	if err != nil {
		log.Fatalf("initializing server failed: %v", err)
	}
	log.Println("success")

	log.Println("running server...")
	err = s.Run()
	if err != nil {
		log.Fatalf("server error: %v", err)
	}
}
