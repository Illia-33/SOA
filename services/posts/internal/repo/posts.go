package repo

import (
	opt "soa-socialnetwork/services/common/option"
	"soa-socialnetwork/services/posts/internal/models"
)

type PostsRepository interface {
	New(models.PageId, NewPostData) (models.PostId, error)
	List(models.PageId, PagiToken) (PostsList, error)
	Get(models.PostId) (models.Post, error)
	Edit(models.PostId, EditedPostData) error
	Delete(models.PostId) error
}

type NewPostData struct {
	AuthorId models.AccountId
	Content  models.PostContent
}

type EditedPostData struct {
	Text   opt.Option[models.Text]
	Pinned opt.Option[bool]
}

type PostsList struct {
	Posts         []models.Post
	NextPagiToken PagiToken
}
