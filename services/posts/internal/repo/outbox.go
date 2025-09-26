package repo

import (
	"soa-socialnetwork/services/posts/internal/models"
	"time"
)

type OutboxRepository interface {
	Put(models.OutboxEvent) error
	Fetch(OutboxFetchParams) ([]models.OutboxEvent, error)
}

type OutboxFetchParams struct {
	Limit         int
	LastCreatedAt time.Time
}
