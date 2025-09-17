package models

type PostId int32
type AccountId int32
type CommentId int32

type Metric int

const (
	METRIC_VIEW_COUNT Metric = iota
	METRIC_LIKE_COUNT
	METRIC_COMMENT_COUNT
)
