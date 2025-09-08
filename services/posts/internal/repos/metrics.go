package repos

import (
	"context"
	dom "soa-socialnetwork/services/posts/internal/domain"
)

type MetricsRepository interface {
	NewView(context.Context, dom.AccountId, dom.PostId) error
	NewLike(context.Context, dom.AccountId, dom.PostId) error
}
