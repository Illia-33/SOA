package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"soa-socialnetwork/services/posts/internal/server"
	"soa-socialnetwork/services/posts/internal/service"
	"strconv"
)

func TryStringFromEnv(key string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("env variable %s not found", key)
	}

	return val, nil
}

func MustStringFromEnv(key string) string {
	val, err := TryStringFromEnv(key)
	if err != nil {
		panic(err)
	}
	return val
}

func TryIntFromEnv(key string) (int, error) {
	s, err := TryStringFromEnv(key)
	if err != nil {
		return 0, err
	}

	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("cannot unparse %s env variable: %v", key, err)
	}

	return val, nil
}

func MustIntFromEnv(key string) int {
	val, err := TryIntFromEnv(key)
	if err != nil {
		panic(err)
	}
	return val
}

func TryEd25519PubKeyFromEnv(key string) (ed25519.PublicKey, error) {
	pubkeyStr, err := TryStringFromEnv(key)
	if err != nil {
		return ed25519.PublicKey{}, err
	}

	pubkey := make(ed25519.PublicKey, ed25519.PublicKeySize)

	_, err = hex.Decode([]byte(pubkey), []byte(pubkeyStr))
	if err != nil {
		return nil, fmt.Errorf("hex decoding error while parsing ed25519 public key from %s (%s): %v", key, pubkeyStr, err)
	}

	return pubkey, nil

}

func MustEd25519PubKeyFromEnv(key string) ed25519.PublicKey {
	val, err := TryEd25519PubKeyFromEnv(key)
	if err != nil {
		panic(err)
	}
	return val
}

func main() {
	s, err := server.Create(service.PostsServiceConfig{
		Port:         MustIntFromEnv("POSTS_SERVICE_PORT"),
		DbHost:       MustStringFromEnv("DB_HOST"),
		DbUser:       MustStringFromEnv("DB_USER"),
		DbPassword:   MustStringFromEnv("DB_PASSWORD"),
		DbPoolSize:   MustIntFromEnv("DB_POOL_SIZE"),
		JwtPublicKey: MustEd25519PubKeyFromEnv("JWT_ED25519_PUBLIC_KEY"),
	})
	if err != nil {
		panic(err)
	}
	log.Println("server created successfully")

	log.Println("running...")
	err = s.Run()
	if err != nil {
		panic(err)
	}
}
