package postgres

import (
	"context"
	"errors"
	"log"
	"soa-socialnetwork/services/posts/internal/models"
	"soa-socialnetwork/services/posts/internal/repo"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type pagesRepo struct {
	ctx   context.Context
	scope pgxScope
}

func (r pagesRepo) GetByAccountId(accountId models.AccountId) (models.Page, error) {
	sql := `
	INSERT INTO pages(account_id)
	VALUES ($1)
	ON CONFLICT(account_id) DO UPDATE
	SET account_id = pages.account_id
	RETURNING id, visible_for_unauthorized, comments_enabled, anyone_can_post;
	`
	var page models.Page

	row := r.scope.QueryRow(r.ctx, sql, accountId)
	err := row.Scan(&page.Id, &page.VisibleForUnauthorized, &page.CommentsEnabled, &page.AnyoneCanPost)
	if err != nil {
		return models.Page{}, err
	}

	page.AccountId = models.AccountId(accountId)
	return page, nil
}

func (r pagesRepo) GetByPageId(pageId models.PageId) (models.Page, error) {
	sql := `
	SELECT account_id, visible_for_unauthorized, comments_enabled, anyone_can_post
	FROM pages
	WHERE id = $1;
	`
	var page models.Page

	row := r.scope.QueryRow(r.ctx, sql, pageId)
	err := row.Scan(&page.AccountId, &page.VisibleForUnauthorized, &page.CommentsEnabled, &page.AnyoneCanPost)
	if err != nil {
		return models.Page{}, status.Error(codes.NotFound, "page not found")
	}

	page.Id = models.PageId(pageId)
	return page, nil
}

func (r pagesRepo) GetByPostId(postId models.PostId) (models.Page, error) {
	sql := `
	SELECT id, account_id, visible_for_unauthorized, comments_enabled, anyone_can_post
	FROM pages
	WHERE id IN (
		SELECT page_id
		FROM posts
		WHERE id = $1
	);
	`
	var page models.Page

	row := r.scope.QueryRow(r.ctx, sql, postId)
	err := row.Scan(&page.Id, &page.AccountId, &page.VisibleForUnauthorized, &page.CommentsEnabled, &page.AnyoneCanPost)
	if err != nil {
		return models.Page{}, status.Error(codes.NotFound, "page not found")
	}

	return page, nil
}

func (r pagesRepo) Edit(pageId models.PageId, edited repo.EditedPageSettings) error {
	sql := `
	WITH affected_rows AS (
		UPDATE pages
		SET
			visible_for_unauthorized = COALESCE($1, visible_for_unauthorized),
			comments_enabled = COALESCE($2, comments_enabled),
			anyone_can_post = COALESCE($3, anyone_can_post)
		WHERE id = $4
		RETURNING 1
	)
	SELECT count(*) FROM affected_rows;
	`

	pgVisibleForUnauthorized := pgtype.Bool{Bool: edited.VisibleForUnauthorized.Value, Valid: edited.VisibleForUnauthorized.HasValue}
	pgCommentsEnabled := pgtype.Bool{Bool: edited.CommentsEnabled.Value, Valid: edited.CommentsEnabled.HasValue}
	pgAnyoneCanPost := pgtype.Bool{Bool: edited.AnyoneCanPost.Value, Valid: edited.AnyoneCanPost.HasValue}

	row := r.scope.QueryRow(r.ctx, sql, pgVisibleForUnauthorized, pgCommentsEnabled, pgAnyoneCanPost, pageId)
	var count int
	row.Scan(&count)
	if count == 0 {
		return errors.New("page not found")
	}

	if count != 1 {
		log.Println("warning: more than 1 profile has been edited while editting page settings")
	}

	return nil
}
