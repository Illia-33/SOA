package service

import (
	"context"
	"crypto/ed25519"
	"time"

	"soa-socialnetwork/services/accounts/internal/repo"
	"soa-socialnetwork/services/accounts/internal/soajwtissuer"
	"soa-socialnetwork/services/accounts/internal/storage/postgres"
	"soa-socialnetwork/services/accounts/pkg/soajwt"
	"soa-socialnetwork/services/accounts/pkg/soatoken"
	pb "soa-socialnetwork/services/accounts/proto"
	"soa-socialnetwork/services/common/backjob"
)

type AccountsService struct {
	pb.UnimplementedAccountsServiceServer

	Db          repo.Database
	JwtVerifier soajwt.Verifier
	SoaVerifier soatoken.Verifier

	outboxJob backjob.TickerJob
	jwtIssuer soajwtissuer.Issuer
}

func NewAccountsService(cfg Config) (*AccountsService, error) {
	ctx := context.Background()

	db, err := postgres.NewDatabase(ctx, postgres.Config{
		Host:     cfg.DbHost,
		User:     cfg.DbUser,
		Password: cfg.DbPassword,
		PoolSize: cfg.DbPoolSize,
	})
	if err != nil {
		return nil, err
	}

	pubkey := cfg.JwtPrivateKey.Public().(ed25519.PublicKey)

	jwtIssuer := soajwtissuer.New(cfg.JwtPrivateKey)
	jwtVerifier := soajwt.NewEd25519Verifier(pubkey)
	soaVerifier := soaVerifier{db: &db}

	service := &AccountsService{
		Db:          &db,
		JwtVerifier: &jwtVerifier,
		SoaVerifier: &soaVerifier,

		outboxJob: backjob.NewTickerJob(3*time.Second, checkOutboxJob(&db)),
		jwtIssuer: jwtIssuer,
	}

	return service, nil
}

func (s *AccountsService) Start() {
	s.outboxJob.Run()
}
