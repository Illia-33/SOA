package dbclient

import (
	"context"
	req "soa-socialnetwork/services/posts/internal/server/dbclient/requests"
)

type DatabaseClient interface {
	GetPageData(context.Context, req.GetPageDataRequest) (req.GetPageDataResponse, error)
	EditPageSettings(context.Context, req.EditPageSettingsRequest) error

	NewPost(context.Context, req.NewPostRequest) (req.NewPostResponse, error)
	GetPost(context.Context, req.GetPostRequest) (req.GetPostResponse, error)
	GetPosts(context.Context, req.GetPostsRequest) (req.GetPostsResponse, error)
	EditPost(context.Context, req.EditPostRequest) error
	DeletePost(context.Context, req.DeletePostRequest) error

	NewComment(context.Context, req.NewCommentRequest) (req.NewCommentResponse, error)
	GetComments(context.Context, req.GetCommentsRequest) (req.GetCommentsResponse, error)

	NewView(context.Context, req.NewViewRequest) error
	NewLike(context.Context, req.NewLikeRequest) error
}
