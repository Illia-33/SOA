package domain

import (
	opt "soa-socialnetwork/services/common/option"
	"time"
)

type AccountId int32
type PageId int32
type PostId int32
type CommentId int32
type Text string
type OutboxEventType string
type OutboxEventPayload []byte

type Page struct {
	Id                     PageId
	AccountId              AccountId
	VisibleForUnauthorized bool
	CommentsEnabled        bool
	AnyoneCanPost          bool
}

type PostContent struct {
	Text         Text
	SourcePostId opt.Option[PostId]
}

type Post struct {
	Id              PostId
	PageId          PageId
	AuthorAccountId AccountId
	Content         PostContent
	Pinned          bool
	ViewsCount      int32
	CreatedAt       time.Time
}

type Comment struct {
	Id        CommentId
	PostId    PostId
	AuthorId  AccountId
	Content   Text
	ReplyId   opt.Option[CommentId]
	CreatedAt time.Time
}

type OutboxEvent struct {
	Type      string
	Payload   OutboxEventPayload
	CreatedAt time.Time
}
