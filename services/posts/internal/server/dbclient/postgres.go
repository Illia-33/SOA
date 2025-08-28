package dbclient

import (
	"context"
	"errors"
	"fmt"
	"log"
	dbreq "soa-socialnetwork/services/posts/internal/server/dbclient/requests"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresConfig struct {
	Host     string
	User     string
	Password string
	PoolSize int
}

type PostgresDbClient struct {
	connPool *pgxpool.Pool
}

func NewPostgresDbClient(cfg PostgresConfig) (*PostgresDbClient, error) {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=5432 dbname=posts-postgres sslmode=disable pool_max_conns=%d", cfg.User, cfg.Password, cfg.Host, cfg.PoolSize)
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	return &PostgresDbClient{
		connPool: pool,
	}, nil
}

func (p *PostgresDbClient) GetPageData(ctx context.Context, req dbreq.GetPageDataRequest) (resp dbreq.GetPageDataResponse, err error) {
	switch id := req.EntityId.(type) {
	case dbreq.AccountId:
		{
			accountId := id
			sql := `
			SELECT id, visible_for_unauthorized, comments_enabled, anyone_can_post
			FROM pages
			WHERE account_id = $1;
			`

			row := p.connPool.QueryRow(ctx, sql, accountId)
			err = row.Scan(&resp.Id, &resp.AnyoneCanPost, &resp.CommentsEnabled, &resp.VisibleForUnauthorized)
			if err == nil {
				// ok, returning
				return
			}

			// otherwise, creating new page for user
			sql = `
			INSERT INTO pages(account_id)
			VALUES ($1)
			RETURNING id, visible_for_unauthorized, comments_enabled, anyone_can_post;
			`

			row = p.connPool.QueryRow(ctx, sql, accountId)
			err = row.Scan(&resp.Id, &resp.AnyoneCanPost, &resp.CommentsEnabled, &resp.VisibleForUnauthorized)
			return
		}

	case dbreq.PageId:
		{
			pageId := id
			sql := `
			SELECT id, visible_for_unauthorized, comments_enabled, anyone_can_post
			FROM pages
			WHERE id = $1;
			`

			row := p.connPool.QueryRow(ctx, sql, pageId)
			err = row.Scan(&resp.Id, &resp.AnyoneCanPost, &resp.CommentsEnabled, &resp.VisibleForUnauthorized)
			return
		}

	case dbreq.PostId:
		{
			postId := id
			sql := `
			SELECT id, visible_for_unauthorized, comments_enabled, anyone_can_post
			FROM pages
			WHERE id IN (
				SELECT page_id
				FROM posts
				WHERE id = $1
			);
			`

			row := p.connPool.QueryRow(ctx, sql, postId)
			err = row.Scan(&resp.Id, &resp.AnyoneCanPost, &resp.CommentsEnabled, &resp.VisibleForUnauthorized)
			return
		}
	}

	err = errors.New("unknown PageEntityId type")
	return
}

func (p *PostgresDbClient) EditPageSettings(ctx context.Context, req dbreq.EditPageSettingsRequest) error {
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

	pgVisibleForUnauthorized := pgtype.Bool{Bool: req.VisibleForUnauthorized.Value, Valid: req.VisibleForUnauthorized.HasValue}
	pgCommentsEnabled := pgtype.Bool{Bool: req.CommentsEnabled.Value, Valid: req.CommentsEnabled.HasValue}
	pgAnyoneCanPost := pgtype.Bool{Bool: req.AnyoneCanPost.Value, Valid: req.AnyoneCanPost.HasValue}

	row := p.connPool.QueryRow(ctx, sql, pgVisibleForUnauthorized, pgCommentsEnabled, pgAnyoneCanPost, req.PageId)
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

func (p *PostgresDbClient) NewPost(ctx context.Context, req dbreq.NewPostRequest) (resp dbreq.NewPostResponse, err error) {
	sql := `
	INSERT INTO posts(page_id, author_account_id, text_content, source_post_id)
	VALUES ($1, $2, $3, $4)
	RETURNING id;
	`

	pgSourcePostId := pgtype.Int4{Int32: int32(req.Content.SourcePostId.Value), Valid: req.Content.SourcePostId.HasValue}

	row := p.connPool.QueryRow(ctx, sql, req.PageId, req.AuthorId, req.Content.TextContent, pgSourcePostId)
	err = row.Scan(&resp.Id)
	return
}

func (p *PostgresDbClient) NewComment(ctx context.Context, req dbreq.NewCommentRequest) (resp dbreq.NewCommentResponse, err error) {
	sql := `
	INSERT INTO comments(post_id, author_account_id, content, reply_comment_id)
	VALUES ($1, $2, $3, $4)
	RETURNING id;
	`

	pgReplyCommentId := pgtype.Int4{Int32: int32(req.ReplyCommentId.Value), Valid: req.ReplyCommentId.HasValue}

	row := p.connPool.QueryRow(ctx, sql, req.PostId, req.AuthorId, req.Content, pgReplyCommentId)
	err = row.Scan(&resp.Id)
	return
}
