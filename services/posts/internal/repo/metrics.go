package repo

import (
	"soa-socialnetwork/services/posts/internal/models"
)

type MetricsRepository interface {
	NewView(models.AccountId, models.PostId) error
	NewLike(models.AccountId, models.PostId) error
}
