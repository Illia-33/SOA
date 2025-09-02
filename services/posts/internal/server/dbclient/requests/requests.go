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

type PagiToken string

// GetPage

type GetPageDataRequest struct {
	EntityId PageEntityId
}

type GetPageDataResponse struct {
	Id                     types.PageId
	AccountId              types.AccountId
	VisibleForUnauthorized bool
	CommentsEnabled        bool
	AnyoneCanPost          bool
}

// NewPost

type NewPostRequest struct {
	PageId   types.PageId
	AuthorId types.AccountId
	Content  types.PostContent
}

type NewPostResponse struct {
	Id types.PostId
}

// GetPost

type GetPostRequest struct {
	PostId types.PostId
}

type GetPostResponse struct {
	Post types.Post
}

// GetPosts

type GetPostsRequest struct {
	PageId    types.PageId
	PageToken PagiToken
}

type GetPostsResponse struct {
	Posts         []types.Post
	NextPageToken PagiToken
}

// NewComment

type NewCommentRequest struct {
	PostId         types.PostId
	AuthorId       types.AccountId
	Content        types.Text
	ReplyCommentId types.Option[types.CommentId]
}

type NewCommentResponse struct {
	Id types.CommentId
}

// GetComments

type GetCommentsRequest struct {
	PostId    types.PostId
	PageToken string
}

type GetCommentsResponse struct {
	Comments      []types.Comment
	NextPageToken PagiToken
}

// EditPost

type EditPostRequest struct {
	PostId types.PostId
	Text   types.Option[types.Text]
	Pinned types.Option[bool]
}

// DeletePost

type DeletePostRequest struct {
	PostId types.PostId
}

// NewView

type NewViewRequest struct {
	AccountId types.AccountId
	PostId    types.PostId
}

// NewLike

type NewLikeRequest struct {
	AccountId types.AccountId
	PostId    types.PostId
}
