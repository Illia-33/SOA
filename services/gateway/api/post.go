package api

import "soa-socialnetwork/services/gateway/pkg/types"

type NewPostRequest struct {
	Text   string                `json:"text"`
	Repost types.Optional[int32] `json:"repost"`
}

type NewPostResponse struct {
	PostId int32 `json:"post_id"`
}

type Post struct {
	Id           int32                 `json:"id"`
	AuthorId     int32                 `json:"author_id"`
	Text         string                `json:"text"`
	SourcePostId types.Optional[int32] `json:"source_post_id"`
	Pinned       bool                  `json:"pinned"`
	ViewsCount   int32                 `json:"views_count"`
}

type GetPostsRequest struct {
	PageToken string `json:"page_token"`
}

type GetPostsResponse struct {
	Posts         []Post `json:"posts"`
	NextPageToken string `json:"next_page_token"`
}

type EditPostRequest struct {
	Text   types.Optional[string] `json:"text"`
	Pinned types.Optional[bool]   `json:"pinned"`
}
