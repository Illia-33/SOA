package requests

import "soa-socialnetwork/services/posts/internal/server/dbclient/types"

type EditPageSettingsRequest struct {
	PageId                 types.PageId
	VisibleForUnauthorized types.Option[bool]
	CommentsEnabled        types.Option[bool]
	AnyoneCanPost          types.Option[bool]
}

type PageEntityId interface {
	isPageEntityId()
}

type AccountId types.AccountId

func (AccountId) isPageEntityId() {}

type PageId types.PageId

func (PageId) isPageEntityId() {}

type PostId types.PostId

func (PostId) isPageEntityId() {}

type PaginationToken string

type GetPageDataRequest struct {
	EntityId PageEntityId
}

type GetPageDataResponse struct {
	Id                     types.PageId
	VisibleForUnauthorized bool
	CommentsEnabled        bool
	AnyoneCanPost          bool
}

type NewPostRequest struct {
	PageId   types.PageId
	AuthorId types.AccountId
	Content  types.PostContent
}

type NewPostResponse struct {
	Id types.PostId
}

type GetPostRequest struct {
	PostId types.PostId
}

type GetPostResponse struct {
	Post types.Post
}

type GetPostsRequest struct {
	PageId    types.PageId
	PageToken PaginationToken
}

type GetPostsResponse struct {
	Posts         []types.Post
	NextPageToken PaginationToken
}

type NewCommentRequest struct {
	PostId         types.PostId
	AuthorId       types.AccountId
	Content        types.Text
	ReplyCommentId types.Option[types.CommentId]
}

type NewCommentResponse struct {
	Id types.CommentId
}
