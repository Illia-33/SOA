package repo

import (
	"soa-socialnetwork/services/accounts/internal/models"
	"time"
)

type ApiTokensRepo interface {
	Put(models.ApiToken, ApiTokenParams) (validUntil time.Time, err error)
	Get(models.ApiToken) (models.ApiTokenData, error)
}

type ApiTokenParams struct {
	AccountId   int
	ReadAccess  bool
	WriteAccess bool
	Ttl         time.Duration
}
