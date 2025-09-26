package postgres

import (
	"context"
	"fmt"
	"log"
	opt "soa-socialnetwork/services/common/option"
	"soa-socialnetwork/services/posts/internal/models"
	"soa-socialnetwork/services/posts/internal/repo"

	"github.com/jackc/pgx/v5/pgtype"
)

type commentsRepo struct {
	ctx   context.Context
	scope pgxScope
}

func (r commentsRepo) New(postId models.PostId, data repo.NewCommentData) (models.CommentId, error) {
	sql := `
	INSERT INTO comments(post_id, author_account_id, text_content, reply_comment_id)
	VALUES ($1, $2, $3, $4)
	RETURNING id;
	`

	pgReplyCommentId := pgtype.Int4{Int32: int32(data.ReplyCommentId.Value), Valid: data.ReplyCommentId.HasValue}
	row := r.scope.QueryRow(r.ctx, sql, postId, data.AuthorId, data.Content, pgReplyCommentId)

	var id models.CommentId
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

const COMMENTS_PAGE_SIZE = 10

type commentsPagiToken struct {
	LastId models.CommentId `json:"lid"`
}

func decodeCommentsPagiToken(token repo.PagiToken) (commentsPagiToken, error) {
	if token == "" {
		return commentsPagiToken{}, nil
	}
	return decodePagiToken[commentsPagiToken](token)
}

func encodeCommentsPagiToken(token commentsPagiToken) (repo.PagiToken, error) {
	return encodePagiToken(token)
}

func (r commentsRepo) List(postId models.PostId, encodedPagiToken repo.PagiToken) (repo.CommentsList, error) {
	pagiToken, err := decodeCommentsPagiToken(encodedPagiToken)
	if err != nil {
		return repo.CommentsList{}, err
	}

	sql := fmt.Sprintf(`
	SELECT id, author_account_id, text_content, reply_comment_id, created_at
	FROM comments
	WHERE post_id = $1 AND id > $2
	ORDER BY id
	LIMIT %d;
	`, COMMENTS_PAGE_SIZE)

	rows, err := r.scope.Query(r.ctx, sql, postId, pagiToken.LastId)
	if err != nil {
		return repo.CommentsList{}, err
	}

	comments := make([]models.Comment, 0, COMMENTS_PAGE_SIZE)
	for {
		if !rows.Next() {
			err := rows.Err()
			if err != nil {
				return repo.CommentsList{}, err
			}
			break
		}

		var pgReplyCommentId pgtype.Int4
		var comment models.Comment

		err := rows.Scan(&comment.Id, &comment.AuthorId, &comment.Content, &pgReplyCommentId, &comment.CreatedAt)
		if err != nil {
			return repo.CommentsList{}, err
		}

		comment.PostId = postId
		if pgReplyCommentId.Valid {
			comment.ReplyId = opt.Some(models.CommentId(pgReplyCommentId.Int32))
		}

		comments = append(comments, comment)
	}

	var nextPagiToken repo.PagiToken
	if len(comments) > 0 {
		token := commentsPagiToken{
			LastId: comments[len(comments)-1].Id,
		}
		encodedToken, err := encodeCommentsPagiToken(token)
		if err != nil {
			log.Printf("warning: cannot encode comments paginating token (%v): %v", token, err)
		} else {
			nextPagiToken = encodedToken
		}
	}

	return repo.CommentsList{
		Comments:      comments,
		NextPagiToken: nextPagiToken,
	}, nil
}
