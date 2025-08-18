package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"soa-socialnetwork/services/gateway/internal/server"
)

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

func main() {
	log.Println("Gateway Service")

	port, err := extractPort()
	if err != nil {
		log.Fatalf("failed to get port: %v", err)
	}

	log.Printf("creating server with port %d...", port)
	s, err := server.Create(port)
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
