package repo

import (
	"soa-socialnetwork/services/stats/pkg/models"
	"time"
)

type DailyCommentsStat struct {
	Date  time.Time
	Count int64
}

type CommentDynamics []DailyCommentsStat

type PostsCommentsRepo interface {
	GetCountForPost(models.PostId) (int64, error)
	GetDynamicsForPost(models.PostId) (CommentDynamics, error)

	Put([]models.PostCommentEvent) error
}
