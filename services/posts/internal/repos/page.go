package repos

import (
	"context"
	opt "soa-socialnetwork/services/common/option"
	dom "soa-socialnetwork/services/posts/internal/domain"
)

type PagesRepository interface {
	GetByAccountId(context.Context, dom.AccountId) (dom.Page, error)
	GetByPageId(context.Context, dom.PageId) (dom.Page, error)
	GetByPostId(context.Context, dom.PostId) (dom.Page, error)
	Edit(context.Context, dom.PageId, EditedPageSettings) error
}

type EditedPageSettings struct {
	VisibleForUnauthorized opt.Option[bool]
	CommentsEnabled        opt.Option[bool]
	AnyoneCanPost          opt.Option[bool]
}
