package api

import "soa-socialnetwork/services/gateway/pkg/types"

type NewCommentRequest struct {
	Content        string                `json:"content"`
	ReplyCommentId types.Optional[int32] `json:"reply_comment_id"`
}

type NewCommentResponse struct {
	CommentId int32 `json:"comment_id"`
}

type Comment struct {
	Id             int32                 `json:"id"`
	AuthorId       int32                 `json:"author_id"`
	Content        string                `json:"content"`
	ReplyCommentId types.Optional[int32] `json:"reply_comment_id"`
}

type GetCommentsRequest struct {
	PageToken string `json:"page_token"`
}

type GetCommentsResponse struct {
	Comments      []Comment `json:"comments"`
	NextPageToken string    `json:"next_page_token"`
}
