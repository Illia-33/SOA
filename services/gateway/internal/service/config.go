package service

import (
	"crypto/ed25519"
)

// Set of parameters of gateway service
type Config struct {
	// ed25519 public key used for verifying jwt tokens
	JwtPublicKey ed25519.PublicKey

	// Hostname of accounts service
	AccountsServiceHost string

	// Port of accounts service
	AccountsServicePort int

	// Hostname of posts service
	PostsServiceHost string

	// Port of posts service
	PostsServicePort int

	// Hostname of stats service
	StatsServiceHost string

	// Port of stats service
	StatsServicePort int
}
