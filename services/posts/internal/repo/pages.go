package repo

import (
	opt "soa-socialnetwork/services/common/option"
	"soa-socialnetwork/services/posts/internal/models"
)

type PagesRepository interface {
	GetByAccountId(models.AccountId) (models.Page, error)
	GetByPageId(models.PageId) (models.Page, error)
	GetByPostId(models.PostId) (models.Page, error)
	Edit(models.PageId, EditedPageSettings) error
}

type EditedPageSettings struct {
	VisibleForUnauthorized opt.Option[bool]
	CommentsEnabled        opt.Option[bool]
	AnyoneCanPost          opt.Option[bool]
}
