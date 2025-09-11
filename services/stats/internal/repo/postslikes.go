package repo

import (
	"soa-socialnetwork/services/stats/pkg/models"
	"time"
)

type DailyLikeStat struct {
	Date  time.Time
	Count int64
}

type LikeDynamics []DailyLikeStat

type PostsLikesRepo interface {
	GetCountForPost(models.PostId) (int64, error)
	GetDynamicsForPost(models.PostId) (LikeDynamics, error)

	Put(...models.PostLikeEvent) error
}
