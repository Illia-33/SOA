package repos

import (
	dom "soa-socialnetwork/services/posts/internal/domain"
	"time"
)

type OutboxRepository interface {
	Put(dom.OutboxEvent) error
	Fetch(OutboxFetchParams) ([]dom.OutboxEvent, error)
}

type OutboxFetchParams struct {
	Limit         int
	LastCreatedAt time.Time
}
