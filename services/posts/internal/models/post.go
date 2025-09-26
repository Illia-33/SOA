package models

import (
	opt "soa-socialnetwork/services/common/option"
	"time"
)

type PostId int32

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
