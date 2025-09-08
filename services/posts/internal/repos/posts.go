package repos

import (
	"context"
	opt "soa-socialnetwork/services/common/option"
	dom "soa-socialnetwork/services/posts/internal/domain"
)

type PostsRepository interface {
	New(context.Context, dom.PageId, NewPostData) (dom.PostId, error)
	List(context.Context, dom.PageId, PagiToken) (PostsList, error)
	Get(context.Context, dom.PostId) (dom.Post, error)
	Edit(context.Context, dom.PostId, EditedPostData) error
	Delete(context.Context, dom.PostId) error
}

type NewPostData struct {
	AuthorId dom.AccountId
	Content  dom.PostContent
}

type EditedPostData struct {
	Text   opt.Option[dom.Text]
	Pinned opt.Option[bool]
}

type PostsList struct {
	Posts         []dom.Post
	NextPagiToken PagiToken
}
