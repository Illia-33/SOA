package server

import (
	"crypto/ed25519"
	"soa-socialnetwork/services/gateway/internal/server/soagrpc"
)

// Set of parameters of gateway service
type GatewayServiceConfig struct {
	// ed25519 public key used for verifying jwt tokens
	JwtPublicKey ed25519.PublicKey

	// Hostname of accounts service
	AccountsServiceHost string

	// Port of accounts service
	AccountsServicePort int

	// Factory used for creating stubs for communication with accounts service.
	// Could be nil, then default gRPC communication will be used.
	AccountsServiceStubFactory soagrpc.AccountsStubFactory

	// Hostname of posts service
	PostsServiceHost string

	// Port of posts service
	PostsServicePort int

	// Factory used for creating stubs for communication with posts service.
	// Could be nil, then default gRPC communication will be used.
	PostsServiceStubFactory soagrpc.PostsStubFactory
}
