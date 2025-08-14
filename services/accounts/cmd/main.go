package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"soa-socialnetwork/services/accounts/internal/server"
)

func extractPort() (int, error) {
	port, exists := os.LookupEnv("SOA_ACCOUNTS_SERVICE_PORT")
	if !exists {
		return 50501, nil
	}

	num, err := strconv.Atoi(port)
	if err != nil {
		return 0, fmt.Errorf("cannot convert SOA_ACCOUNTS_SERVICE_PORT environment variable to number (%v)", err)
	}
	return num, nil
}

func main() {
	log.Println("Accounts Service")

	port, err := extractPort()
	if err != nil {
		log.Fatalf("failed to get port: %v", err)
	}

	log.Printf("creating server with port %d...", port)
	s, err := server.Create(port)
	if err != nil {
		log.Fatalf("initializing server failed: %v", err)
	}
	log.Print("success")

	log.Print("running server...")
	err = s.Run()
	if err != nil {
		log.Fatalf("server error: %v", err)
	}
}
