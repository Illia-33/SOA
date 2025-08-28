package dbclient

import (
	"context"
	req "soa-socialnetwork/services/posts/internal/server/dbclient/requests"
)

type DatabaseClient interface {
	GetPageData(context.Context, req.GetPageDataRequest) (req.GetPageDataResponse, error)
	EditPageSettings(context.Context, req.EditPageSettingsRequest) error
	NewPost(context.Context, req.NewPostRequest) (req.NewPostResponse, error)
	NewComment(context.Context, req.NewCommentRequest) (req.NewCommentResponse, error)
}
