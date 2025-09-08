package repos

import (
	"context"
	opt "soa-socialnetwork/services/common/option"
	dom "soa-socialnetwork/services/posts/internal/domain"
)

type CommentsRepository interface {
	New(context.Context, dom.PostId, NewCommentData) (dom.CommentId, error)
	List(context.Context, dom.PostId, PagiToken) (CommentsList, error)
}

type NewCommentData struct {
	AuthorId       dom.AccountId
	Content        dom.Text
	ReplyCommentId opt.Option[dom.CommentId]
}

type CommentsList struct {
	Comments      []dom.Comment
	NextPagiToken PagiToken
}
