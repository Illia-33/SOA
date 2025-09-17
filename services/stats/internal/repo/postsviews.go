package repo

import (
	"soa-socialnetwork/services/stats/pkg/models"
	"time"
)

type DailyViewsStat struct {
	Date  time.Time
	Count uint64
}

type ViewDynamics []DailyViewsStat

type PostsViewsRepo interface {
	GetCountForPost(models.PostId) (uint64, error)
	GetDynamicsForPost(models.PostId) (ViewDynamics, error)

	Put(...models.PostViewEvent) error
}
