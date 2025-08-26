package server

import "crypto/ed25519"

type GatewayServiceConfig struct {
	JwtPublicKey        ed25519.PublicKey
	AccountsServiceHost string
	AccountsServicePort int
}
