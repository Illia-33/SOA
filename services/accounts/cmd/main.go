package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"soa-socialnetwork/services/accounts/internal/server"
)

func extractEnv(key string) (string, error) {
	val, exists := os.LookupEnv(key)
	if !exists {
		return "", fmt.Errorf("environment variable %s doesn't exist", key)
	}

	return val, nil
}

func extractPort() (int, error) {
	port, exists := os.LookupEnv("ACCOUNTS_SERVICE_PORT")
	if !exists {
		return 50501, nil
	}

	num, err := strconv.Atoi(port)
	if err != nil {
		return 0, fmt.Errorf("cannot convert ACCOUNTS_SERVICE_PORT environment variable to number (%v)", err)
	}
	return num, nil
}

func extractServiceConfig() (cfg server.AccountsServiceConfig, err error) {
	cfg.DbHost, err = extractEnv("DB_HOST")
	if err != nil {
		return
	}

	cfg.DbUser, err = extractEnv("DB_USER")
	if err != nil {
		return
	}

	cfg.DbPassword, err = extractEnv("DB_PASSWORD")
	if err != nil {
		return
	}

	dbPoolSizeStr, err := extractEnv("DB_POOL_SIZE")
	if err != nil {
		return
	}
	dbPoolSize, err := strconv.Atoi(dbPoolSizeStr)
	if err != nil {
		return
	}
	cfg.DbPoolSize = dbPoolSize

	return
}

func main() {
	log.Println("Accounts Service")

	port, err := extractPort()
	if err != nil {
		log.Fatalf("failed to get port: %v", err)
	}

	cfg, err := extractServiceConfig()
	if err != nil {
		log.Fatalf("cannot extract service config: %v", err)
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
