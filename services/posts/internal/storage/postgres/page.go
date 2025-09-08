package postgres

import (
	"context"
	"errors"
	"log"
	dom "soa-socialnetwork/services/posts/internal/domain"
	"soa-socialnetwork/services/posts/internal/repos"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PageRepo struct {
	ConnPool connectionPool
}

func (r *PageRepo) Get(ctx context.Context, entityId repos.PageEntityId) (page dom.Page, err error) {
	switch id := entityId.(type) {
	case repos.AccountId:
		{
			accountId := id
			sql := `
			SELECT id, visible_for_unauthorized, comments_enabled, anyone_can_post
			FROM pages
			WHERE account_id = $1;
			`

			row := r.ConnPool.QueryRow(ctx, sql, accountId)
			err = row.Scan(&page.Id, &page.VisibleForUnauthorized, &page.CommentsEnabled, &page.AnyoneCanPost)
			if err == nil { // ok, returning
				page.AccountId = dom.AccountId(accountId)
				return
			}

			// otherwise, creating new page for user
			sql = `
			INSERT INTO pages(account_id)
			VALUES ($1)
			RETURNING id, visible_for_unauthorized, comments_enabled, anyone_can_post;
			`

			row = r.ConnPool.QueryRow(ctx, sql, accountId)
			err = row.Scan(&page.Id, &page.VisibleForUnauthorized, &page.CommentsEnabled, &page.AnyoneCanPost)
			if err == nil {
				page.AccountId = dom.AccountId(accountId)
			}
			return
		}

	case repos.PageId:
		{
			pageId := id
			sql := `
			SELECT account_id, visible_for_unauthorized, comments_enabled, anyone_can_post
			FROM pages
			WHERE id = $1;
			`

			row := r.ConnPool.QueryRow(ctx, sql, pageId)
			err = row.Scan(&page.AccountId, &page.VisibleForUnauthorized, &page.CommentsEnabled, &page.AnyoneCanPost)
			if err != nil {
				err = status.Error(codes.NotFound, "page not found")
				return
			}

			page.Id = dom.PageId(pageId)
			return
		}

	case repos.PostId:
		{
			postId := id
			sql := `
			SELECT id, account_id, visible_for_unauthorized, comments_enabled, anyone_can_post
			FROM pages
			WHERE id IN (
				SELECT page_id
				FROM posts
				WHERE id = $1
			);
			`

			row := r.ConnPool.QueryRow(ctx, sql, postId)
			err = row.Scan(&page.Id, &page.AccountId, &page.VisibleForUnauthorized, &page.CommentsEnabled, &page.AnyoneCanPost)
			if err != nil {
				err = status.Error(codes.NotFound, "page not found")
			}
			return
		}
	}

	err = errors.New("unknown PageEntityId type")
	return

}

func (r *PageRepo) Edit(ctx context.Context, pageId dom.PageId, edited repos.EditedPageSettings) error {
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

	row := r.ConnPool.QueryRow(ctx, sql, pgVisibleForUnauthorized, pgCommentsEnabled, pgAnyoneCanPost, pageId)
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
