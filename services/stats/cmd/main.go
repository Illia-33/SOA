package main

import (
	"fmt"
	"os"
	"soa-socialnetwork/services/stats/internal/server"
	"soa-socialnetwork/services/stats/internal/service"
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

func main() {
	port := MustIntFromEnv("STATS_SERVICE_PORT")
	cfg := service.Config{
		DbHost:     MustStringFromEnv("DB_HOST"),
		DbPort:     MustIntFromEnv("DB_PORT"),
		DbUser:     MustStringFromEnv("DB_USER"),
		DbPassword: MustStringFromEnv("DB_PASSWORD"),
	}

	s, err := server.New(port, cfg)
	if err != nil {
		panic(err)
	}

	err = s.Run()
	if err != nil {
		panic(err)
	}
}
