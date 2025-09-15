package repo

import "context"

type Database interface {
	PostsViews(context.Context) PostsViewsRepo
	PostsLikes(context.Context) PostsLikesRepo
	PostsComments(context.Context) PostsCommentsRepo
	Registrations(context.Context) RegistrationsRepo
	Posts(context.Context) PostsRepo
}
