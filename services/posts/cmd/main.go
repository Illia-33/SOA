package main

import (
	"log"
	"soa-socialnetwork/services/common/envvar"
	"soa-socialnetwork/services/posts/internal/server"
	"soa-socialnetwork/services/posts/internal/service"
)

func main() {
	port := envvar.MustIntFromEnv("POSTS_SERVICE_PORT")

	s, err := server.Create(service.PostsServiceConfig{
		DbHost:       envvar.MustStringFromEnv("DB_HOST"),
		DbUser:       envvar.MustStringFromEnv("DB_USER"),
		DbPassword:   envvar.MustStringFromEnv("DB_PASSWORD"),
		DbPoolSize:   envvar.MustIntFromEnv("DB_POOL_SIZE"),
		JwtPublicKey: envvar.MustEd25519PubKeyFromEnv("JWT_ED25519_PUBLIC_KEY"),
	})
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
