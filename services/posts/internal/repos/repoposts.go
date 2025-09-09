package repos

import (
	opt "soa-socialnetwork/services/common/option"
	dom "soa-socialnetwork/services/posts/internal/domain"
)

type PostsRepository interface {
	New(dom.PageId, NewPostData) (dom.PostId, error)
	List(dom.PageId, PagiToken) (PostsList, error)
	Get(dom.PostId) (dom.Post, error)
	Edit(dom.PostId, EditedPostData) error
	Delete(dom.PostId) error
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
