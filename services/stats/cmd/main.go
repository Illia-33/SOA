package main

import (
	"log"
	"soa-socialnetwork/services/common/envvar"
	"soa-socialnetwork/services/stats/internal/server"
	"soa-socialnetwork/services/stats/internal/service"
)

func main() {
	log.Println("Statistics service")
	port := envvar.MustIntFromEnv("STATS_SERVICE_PORT")
	cfg := service.Config{
		DbHost:     envvar.MustStringFromEnv("DB_HOST"),
		DbPort:     envvar.MustIntFromEnv("DB_PORT"),
		DbUser:     envvar.MustStringFromEnv("DB_USER"),
		DbPassword: envvar.MustStringFromEnv("DB_PASSWORD"),
		KafkaHost:  envvar.MustStringFromEnv("KAFKA_HOST"),
		KafkaPort:  envvar.MustIntFromEnv("KAFKA_PORT"),
	}

	log.Printf("creating server with cfg: %+v", cfg)
	s, err := server.New(cfg)
	if err != nil {
		panic(err)
	}
	log.Println("server created successfully")

	log.Println("running...")
	err = s.Run(port)
	if err != nil {
		panic(err)
	}
}
