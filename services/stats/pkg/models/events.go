package models

import "time"

type PostViewEvent struct {
	PostId          PostId    `json:"post_id"`
	ViewerAccountId AccountId `json:"viewer_account_id"`
	Timestamp       time.Time `json:"timestamp"`
}

type PostLikeEvent struct {
	PostId         PostId    `json:"post_id"`
	LikerAccountId AccountId `json:"liker_account_id"`
	Timestamp      time.Time `json:"timestamp"`
}

type PostCommentEvent struct {
	CommentId       CommentId `json:"comment_id"`
	PostId          PostId    `json:"post_id"`
	AuthorAccountId AccountId `json:"author_account_id"`
	Timestamp       time.Time `json:"timestamp"`
}

type RegistrationEvent struct {
	AccountId AccountId `json:"account_Id"`
	ProfileId string    `json:"profile_id"`
	Timestamp time.Time `json:"timestamp"`
}

type PostEvent struct {
	PostId    PostId
	AuthorId  AccountId
	Timestamp time.Time
}
