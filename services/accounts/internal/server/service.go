package server

import (
	"context"
	"crypto/ed25519"
	"fmt"

	"soa-socialnetwork/services/accounts/internal/server/soajwtissuer"
	"soa-socialnetwork/services/accounts/pkg/soajwt"
	pb "soa-socialnetwork/services/accounts/proto"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountsServiceConfig struct {
	DbHost        string
	DbUser        string
	DbPassword    string
	DbPoolSize    int
	JwtPrivateKey ed25519.PrivateKey
}

type AccountsService struct {
	pb.UnimplementedAccountsServiceServer

	dbPool      *pgxpool.Pool
	jwtIssuer   soajwtissuer.Issuer
	jwtVerifier soajwt.Verifier
	soaVerifier serviceSoaTokenVerifier
}

func createAccountsService(cfg AccountsServiceConfig) (*AccountsService, error) {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=5432 dbname=accounts-postgres sslmode=disable pool_max_conns=%d", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPoolSize)
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	jwtIssuer := soajwtissuer.New(cfg.JwtPrivateKey)

	pubkey := cfg.JwtPrivateKey.Public().(ed25519.PublicKey)
	service := &AccountsService{
		dbPool:      pool,
		jwtIssuer:   jwtIssuer,
		jwtVerifier: soajwt.NewVerifier(pubkey),
	}
	service.soaVerifier = serviceSoaTokenVerifier{service: service}

	return service, nil
}
