package repos

import (
	opt "soa-socialnetwork/services/common/option"
	dom "soa-socialnetwork/services/posts/internal/domain"
)

type PagesRepository interface {
	GetByAccountId(dom.AccountId) (dom.Page, error)
	GetByPageId(dom.PageId) (dom.Page, error)
	GetByPostId(dom.PostId) (dom.Page, error)
	Edit(dom.PageId, EditedPageSettings) error
}

type EditedPageSettings struct {
	VisibleForUnauthorized opt.Option[bool]
	CommentsEnabled        opt.Option[bool]
	AnyoneCanPost          opt.Option[bool]
}
