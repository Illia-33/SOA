package service

import "crypto/ed25519"

type PostsServiceConfig struct {
	DbHost       string
	DbUser       string
	DbPassword   string
	DbPoolSize   int
	JwtPublicKey ed25519.PublicKey
}
