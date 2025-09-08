package repos

import (
	"context"
	opt "soa-socialnetwork/services/common/option"
	dom "soa-socialnetwork/services/posts/internal/domain"
)

type PagesRepository interface {
	Get(context.Context, PageEntityId) (dom.Page, error)
	Edit(context.Context, dom.PageId, EditedPageSettings) error
}

type PageEntityId interface {
	isPageEntityId()
}

type AccountId dom.AccountId
type PageId dom.PageId
type PostId dom.PostId

func (AccountId) isPageEntityId() {}
func (PageId) isPageEntityId()    {}
func (PostId) isPageEntityId()    {}

type EditedPageSettings struct {
	VisibleForUnauthorized opt.Option[bool]
	CommentsEnabled        opt.Option[bool]
	AnyoneCanPost          opt.Option[bool]
}
