package repos

import (
	opt "soa-socialnetwork/services/common/option"
	"soa-socialnetwork/services/posts/internal/models"
)

type CommentsRepository interface {
	New(models.PostId, NewCommentData) (models.CommentId, error)
	List(models.PostId, PagiToken) (CommentsList, error)
}

type NewCommentData struct {
	AuthorId       models.AccountId
	Content        models.Text
	ReplyCommentId opt.Option[models.CommentId]
}

type CommentsList struct {
	Comments      []models.Comment
	NextPagiToken PagiToken
}
