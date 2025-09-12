package main

import (
	"fmt"
	"log"
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
	log.Println("Statistics service")
	port := MustIntFromEnv("STATS_SERVICE_PORT")
	cfg := service.Config{
		DbHost:     MustStringFromEnv("DB_HOST"),
		DbPort:     MustIntFromEnv("DB_PORT"),
		DbUser:     MustStringFromEnv("DB_USER"),
		DbPassword: MustStringFromEnv("DB_PASSWORD"),
		KafkaHost:  MustStringFromEnv("KAFKA_HOST"),
		KafkaPort:  MustIntFromEnv("KAFKA_PORT"),
	}

	log.Printf("creating server with cfg: %+v", cfg)
	s, err := server.New(port, cfg)
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
