package domain

import (
	opt "soa-socialnetwork/services/common/option"
	"time"
)

type CommentId int32

type Comment struct {
	Id        CommentId
	PostId    PostId
	AuthorId  AccountId
	Content   string
	ReplyId   opt.Option[CommentId]
	CreatedAt time.Time
}
