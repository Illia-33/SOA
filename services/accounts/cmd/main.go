package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"errors"
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

func extractJwtPrivateKey() (ed25519.PrivateKey, error) {
	envPrivateKey, err := extractEnv("JWT_ED25519_PRIVATE_KEY")
	if err != nil {
		log.Println("cannot find JWT_ED25519_PRIVATE_KEY environment variable, generating jwt private key...")
		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return ed25519.PrivateKey{}, err
		}

		log.Printf("jwt public ed25519 key: %s", hex.EncodeToString(pub))
		return priv, nil
	}

	seed := make([]byte, ed25519.SeedSize)

	_, err = hex.Decode(seed, []byte(envPrivateKey))
	if err != nil {
		return ed25519.PrivateKey{}, errors.New("cannot decode hex encoded jwt private key seed")
	}

	return ed25519.NewKeyFromSeed(seed), nil
}

func extractServiceConfig() (server.AccountsServiceConfig, error) {
	dbHost, err := extractEnv("DB_HOST")
	if err != nil {
		return server.AccountsServiceConfig{}, err
	}

	dbUser, err := extractEnv("DB_USER")
	if err != nil {
		return server.AccountsServiceConfig{}, err
	}

	dbPassword, err := extractEnv("DB_PASSWORD")
	if err != nil {
		return server.AccountsServiceConfig{}, err
	}

	dbPoolSizeStr, err := extractEnv("DB_POOL_SIZE")
	if err != nil {
		return server.AccountsServiceConfig{}, err
	}
	dbPoolSize, err := strconv.Atoi(dbPoolSizeStr)
	if err != nil {
		return server.AccountsServiceConfig{}, err
	}

	jwtPrivateKey, err := extractJwtPrivateKey()
	if err != nil {
		return server.AccountsServiceConfig{}, err
	}

	return server.AccountsServiceConfig{
		DbHost:        dbHost,
		DbUser:        dbUser,
		DbPassword:    dbPassword,
		DbPoolSize:    dbPoolSize,
		JwtPrivateKey: jwtPrivateKey,
	}, nil
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
