package service

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"time"

	"soa-socialnetwork/services/accounts/internal/soajwtissuer"
	"soa-socialnetwork/services/accounts/pkg/soajwt"
	pb "soa-socialnetwork/services/accounts/proto"
	"soa-socialnetwork/services/common/backjob"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountsService struct {
	pb.UnimplementedAccountsServiceServer

	JwtVerifier soajwt.Verifier

	dbPool    *pgxpool.Pool
	outboxJob backjob.TickerJob
	jwtIssuer soajwtissuer.Issuer
}

func NewAccountsService(cfg Config) (*AccountsService, error) {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=5432 dbname=accounts-postgres sslmode=disable pool_max_conns=%d", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPoolSize)
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	jwtIssuer := soajwtissuer.New(cfg.JwtPrivateKey)
	pubkey := cfg.JwtPrivateKey.Public().(ed25519.PublicKey)
	service := &AccountsService{
		dbPool:      pool,
		outboxJob:   backjob.NewTickerJob(5*time.Second, checkOutboxJob(pool)),
		jwtIssuer:   jwtIssuer,
		JwtVerifier: soajwt.NewEd25519Verifier(pubkey),
	}

	return service, nil
}

func (s *AccountsService) Start() {
	s.outboxJob.Run()
}
