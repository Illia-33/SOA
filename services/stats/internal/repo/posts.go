package repo

import "soa-socialnetwork/services/stats/pkg/models"

type PostsRepo interface {
	Put(...models.PostEvent) error
}
