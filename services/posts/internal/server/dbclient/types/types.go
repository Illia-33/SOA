package types

import "time"

type AccountId int32
type PageId int32
type PostId int32
type CommentId int32

type PageData struct {
	Id                     PageId
	AccountId              AccountId
	VisibleForUnauthorized bool
	CommentsEnabled        bool
	AnyoneCanPost          bool
}

type Text string

type PostContent struct {
	TextContent  Text
	SourcePostId Option[PostId]
}

type PostData struct {
	Id              PostId
	PageId          PageId
	AuthorAccountId Option[AccountId]
	Content         PostContent
	Pinned          bool
	CreatedAt       time.Time
}

type CommentData struct {
	Id        CommentId
	PostId    PostId
	AuthorId  AccountId
	Content   Text
	ReplyId   Option[CommentId]
	CreatedAt time.Time
}

type Option[T any] struct {
	Value    T
	HasValue bool
}

func Some[T any](v T) Option[T] {
	return Option[T]{
		Value:    v,
		HasValue: true,
	}
}

func None[T any]() Option[T] {
	return Option[T]{
		HasValue: false,
	}
}

func OptionFromPtr[T any](p *T) Option[T] {
	if p == nil {
		return None[T]()
	}

	return Some(*p)
}
