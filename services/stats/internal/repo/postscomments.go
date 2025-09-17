package repo

import (
	"soa-socialnetwork/services/stats/pkg/models"
	"time"
)

type DailyCommentsStat struct {
	Date  time.Time
	Count uint64
}

type CommentDynamics []DailyCommentsStat

type PostsCommentsRepo interface {
	GetCountForPost(models.PostId) (uint64, error)
	GetDynamicsForPost(models.PostId) (CommentDynamics, error)

	Put(...models.PostCommentEvent) error
}
