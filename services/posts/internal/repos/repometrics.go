package repos

import (
	dom "soa-socialnetwork/services/posts/internal/domain"
)

type MetricsRepository interface {
	NewView(dom.AccountId, dom.PostId) error
	NewLike(dom.AccountId, dom.PostId) error
}
