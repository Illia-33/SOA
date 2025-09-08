package repos

import (
	"context"
	dom "soa-socialnetwork/services/posts/internal/domain"
)

type OutboxRepository interface {
	Put(context.Context, dom.OutboxEvent) error
	Fetch(context.Context, OutboxFetchParams) ([]dom.OutboxEvent, Transaction, error)
}

type OutboxFetchParams struct {
	Limit int
}
