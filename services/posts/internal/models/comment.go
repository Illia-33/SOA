package models

import (
	opt "soa-socialnetwork/services/common/option"
	"time"
)

type CommentId int32

type Comment struct {
	Id        CommentId
	PostId    PostId
	AuthorId  AccountId
	Content   Text
	ReplyId   opt.Option[CommentId]
	CreatedAt time.Time
}
