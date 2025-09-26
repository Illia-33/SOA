package service

import "crypto/ed25519"

type Config struct {
	DbHost        string
	DbUser        string
	DbPassword    string
	DbPoolSize    int
	JwtPrivateKey ed25519.PrivateKey
}
