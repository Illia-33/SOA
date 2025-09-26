package domain

type PageId int32

type Page struct {
	Id                     PageId
	AccountId              AccountId
	VisibleForUnauthorized bool
	CommentsEnabled        bool
	AnyoneCanPost          bool
}
