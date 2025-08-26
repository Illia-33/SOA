package server

import (
	"crypto/ed25519"
	"soa-socialnetwork/services/gateway/internal/server/soagrpc"
)

type GatewayServiceConfig struct {
	JwtPublicKey               ed25519.PublicKey
	AccountsServiceHost        string
	AccountsServicePort        int
	AccountsServiceStubFactory soagrpc.AccountsStubFactory
}
