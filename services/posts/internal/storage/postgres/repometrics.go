package postgres

import (
	"context"
	"log"
	dom "soa-socialnetwork/services/posts/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type metricsRepo struct {
	ctx   context.Context
	scope pgxScope
}

func (r metricsRepo) NewView(accountId dom.AccountId, postId dom.PostId) error {
	sql := `
	WITH affected_rows AS (
		UPDATE posts
		SET
			views_count = views_count + 1
		WHERE id = $1
		RETURNING 1
	)
	SELECT count(*) FROM affected_rows;
	`

	row := r.scope.QueryRow(r.ctx, sql, postId)
	var countAffected int
	if err := row.Scan(&countAffected); err != nil {
		return err
	}

	if countAffected == 0 {
		return status.Error(codes.NotFound, "post not found")
	}

	if countAffected != 1 {
		log.Printf("warning: more than 1 post with id %d in posts table", postId)
	}

	return nil
}

func (r metricsRepo) NewLike(accountId dom.AccountId, postId dom.PostId) error {
	sql := `
	INSERT INTO likes(post_id, author_account_id)
	VALUES ($1, $2)
	`

	_, err := r.scope.Exec(r.ctx, sql, postId, accountId)
	if err != nil {
		return err
	}

	return nil
}
