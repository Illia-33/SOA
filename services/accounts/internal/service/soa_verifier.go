package service

import (
	"context"
	"soa-socialnetwork/services/accounts/internal/models"
	"soa-socialnetwork/services/accounts/internal/repo"
	"soa-socialnetwork/services/accounts/internal/service/errs"
	"soa-socialnetwork/services/accounts/pkg/soatoken"
	"time"
)

type soaVerifier struct {
	db repo.Database
}

func (v *soaVerifier) Verify(token string, reqs soatoken.RightsRequirements) error {
	conn, err := v.db.OpenConnection(context.Background())
	if err != nil {
		return err
	}
	defer conn.Close()

	tokenData, err := conn.ApiTokens().Get(models.ApiToken(token))
	if err != nil {
		return err
	}

	if reqs.Read && !tokenData.ReadAccess {
		return errs.NoReadAccess{}
	}

	if reqs.Write && !tokenData.WriteAccess {
		return errs.NoWriteAccess{}
	}

	now := time.Now()
	if now.After(tokenData.ValidUntil) {
		return errs.TokenExpired{}
	}

	return nil
}
