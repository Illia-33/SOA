package repo

import (
	"soa-socialnetwork/services/accounts/internal/models"
	"time"
)

type OutboxRepo interface {
	Put(models.OutboxEvent) error
	Fetch(OutboxFetchParams) ([]models.OutboxEvent, error)
}

type OutboxFetchParams struct {
	Limit         int
	LastCreatedAt time.Time
}
