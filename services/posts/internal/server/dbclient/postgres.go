package dbclient

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	dbReq "soa-socialnetwork/services/posts/internal/server/dbclient/requests"
	dbTypes "soa-socialnetwork/services/posts/internal/server/dbclient/types"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (p *PostgresDbClient) GetPageData(ctx context.Context, req dbReq.GetPageDataRequest) (resp dbReq.GetPageDataResponse, err error) {
	switch id := req.EntityId.(type) {
	case dbReq.AccountId:
		{
			accountId := id
			sql := `
			SELECT id, visible_for_unauthorized, comments_enabled, anyone_can_post
			FROM pages
			WHERE account_id = $1;
			`

			row := p.connPool.QueryRow(ctx, sql, accountId)
			err = row.Scan(&resp.Id, &resp.VisibleForUnauthorized, &resp.CommentsEnabled, &resp.AnyoneCanPost)
			if err == nil { // ok, returning
				resp.AccountId = dbTypes.AccountId(accountId)
				return
			}

			// otherwise, creating new page for user
			sql = `
			INSERT INTO pages(account_id)
			VALUES ($1)
			RETURNING id, visible_for_unauthorized, comments_enabled, anyone_can_post;
			`

			row = p.connPool.QueryRow(ctx, sql, accountId)
			err = row.Scan(&resp.Id, &resp.VisibleForUnauthorized, &resp.CommentsEnabled, &resp.AnyoneCanPost)
			if err == nil {
				resp.AccountId = dbTypes.AccountId(accountId)
			}
			return
		}

	case dbReq.PageId:
		{
			pageId := id
			sql := `
			SELECT account_id, visible_for_unauthorized, comments_enabled, anyone_can_post
			FROM pages
			WHERE id = $1;
			`

			row := p.connPool.QueryRow(ctx, sql, pageId)
			err = row.Scan(&resp.AccountId, &resp.VisibleForUnauthorized, &resp.CommentsEnabled, &resp.AnyoneCanPost)
			if err != nil {
				err = status.Error(codes.NotFound, "page not found")
				return
			}

			resp.Id = dbTypes.PageId(pageId)
			return
		}

	case dbReq.PostId:
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

			row := p.connPool.QueryRow(ctx, sql, postId)
			err = row.Scan(&resp.Id, &resp.AccountId, &resp.VisibleForUnauthorized, &resp.CommentsEnabled, &resp.AnyoneCanPost)
			if err != nil {
				err = status.Error(codes.NotFound, "page not found")
			}
			return
		}
	}

	err = errors.New("unknown PageEntityId type")
	return
}

func (p *PostgresDbClient) EditPageSettings(ctx context.Context, req dbReq.EditPageSettingsRequest) error {
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

func (p *PostgresDbClient) NewPost(ctx context.Context, req dbReq.NewPostRequest) (resp dbReq.NewPostResponse, err error) {
	sql := `
	INSERT INTO posts(page_id, author_account_id, text_content, source_post_id)
	VALUES ($1, $2, $3, $4)
	RETURNING id;
	`

	pgSourcePostId := pgtype.Int4{Int32: int32(req.Content.SourcePostId.Value), Valid: req.Content.SourcePostId.HasValue}

	row := p.connPool.QueryRow(ctx, sql, req.PageId, req.AuthorId, req.Content.Text, pgSourcePostId)
	err = row.Scan(&resp.Id)
	return
}

func (p *PostgresDbClient) GetPost(ctx context.Context, req dbReq.GetPostRequest) (resp dbReq.GetPostResponse, err error) {
	sql := `
	SELECT page_id, author_account_id, text_content, source_post_id, pinned, views_count, created_at
	FROM posts
	WHERE id = $1;
	`

	var post dbTypes.Post
	var pgSourcePostId pgtype.Int4

	row := p.connPool.QueryRow(ctx, sql, req.PostId)
	err = row.Scan(&post.PageId, &post.AuthorAccountId, &post.Content.Text, &pgSourcePostId, &post.Pinned, &post.ViewsCount, &post.CreatedAt)
	if err != nil {
		err = status.Error(codes.NotFound, "post not found")
		return
	}

	post.Id = req.PostId
	post.Content.SourcePostId = dbTypes.Option[dbTypes.PostId]{
		Value:    dbTypes.PostId(pgSourcePostId.Int32),
		HasValue: pgSourcePostId.Valid,
	}
	resp.Post = post
	return
}

const POSTGRES_POSTS_PAGE_SIZE = 10

func (p *PostgresDbClient) GetPosts(ctx context.Context, req dbReq.GetPostsRequest) (dbReq.GetPostsResponse, error) {
	sql := fmt.Sprintf(`
	SELECT id, author_account_id, text_content, source_post_id, pinned, views_count, created_at
	FROM posts
	WHERE page_id = $1 AND created_at < $2
	ORDER BY created_at DESC
	LIMIT %d;
	`, POSTGRES_POSTS_PAGE_SIZE)

	page, err := decodePgPostsPagiToken(req.PageToken)
	if err != nil {
		return dbReq.GetPostsResponse{}, err
	}

	rows, err := p.connPool.Query(ctx, sql, req.PageId, page.LastCreatedAt)
	if err != nil {
		return dbReq.GetPostsResponse{}, err
	}

	posts := make([]dbTypes.Post, 0, POSTGRES_POSTS_PAGE_SIZE)

	for {
		if !rows.Next() {
			err := rows.Err()
			if err != nil {
				return dbReq.GetPostsResponse{}, err
			}

			break
		}

		var post dbTypes.Post
		var pgSourcePostId pgtype.Int4
		err := rows.Scan(&post.Id, &post.AuthorAccountId, &post.Content.Text, &pgSourcePostId, &post.Pinned, &post.ViewsCount, &post.CreatedAt)
		if err != nil {
			return dbReq.GetPostsResponse{}, err
		}

		post.PageId = req.PageId
		post.Content.SourcePostId = dbTypes.Option[dbTypes.PostId]{Value: dbTypes.PostId(pgSourcePostId.Int32), HasValue: pgSourcePostId.Valid}
		posts = append(posts, post)
	}

	var nextPageToken dbReq.PaginationToken
	if len(posts) > 0 {
		token := pgPostsPagiToken{
			LastCreatedAt: posts[len(posts)-1].CreatedAt,
		}
		encodedToken, err := encodePgPostsPagiToken(token)

		if err != nil {
			log.Printf("warning: cannot encode paginating token (%v): %v", token, err)
		} else {
			nextPageToken = encodedToken
		}
	}

	return dbReq.GetPostsResponse{
		Posts:         posts,
		NextPageToken: dbReq.PaginationToken(nextPageToken),
	}, nil
}

func (p *PostgresDbClient) NewComment(ctx context.Context, req dbReq.NewCommentRequest) (resp dbReq.NewCommentResponse, err error) {
	sql := `
	INSERT INTO comments(post_id, author_account_id, text_content, reply_comment_id)
	VALUES ($1, $2, $3, $4)
	RETURNING id;
	`

	pgReplyCommentId := pgtype.Int4{Int32: int32(req.ReplyCommentId.Value), Valid: req.ReplyCommentId.HasValue}

	row := p.connPool.QueryRow(ctx, sql, req.PostId, req.AuthorId, req.Content, pgReplyCommentId)
	err = row.Scan(&resp.Id)
	return
}

func (p *PostgresDbClient) EditPost(ctx context.Context, req dbReq.EditPostRequest) error {
	sql := `
	WITH affected_rows AS (
		UPDATE posts
		SET
			text_content = COALESCE($1, text_content),
			pinned = COALESCE($2, pinned)
		WHERE id = $3
		RETURNING 1
	)
	SELECT count(*) FROM affected_rows;
	`
	pgTextContent := pgtype.Text{
		String: string(req.Text.Value),
		Valid:  req.Text.HasValue,
	}
	pgPinned := pgtype.Bool{
		Bool:  req.Pinned.Value,
		Valid: req.Pinned.HasValue,
	}

	row := p.connPool.QueryRow(ctx, sql, pgTextContent, pgPinned, req.PostId)
	var countAffected int
	if err := row.Scan(&countAffected); err != nil {
		return err
	}

	if countAffected == 0 {
		return status.Error(codes.NotFound, "post not found")
	}

	if countAffected != 1 {
		log.Printf("warning: more than 1 post with id = %d", req.PostId)
	}

	return nil
}

func (p *PostgresDbClient) DeletePost(ctx context.Context, req dbReq.DeletePostRequest) error {
	sql := `
	WITH affected_rows AS (
		DELETE FROM posts
		WHERE id = $1
		RETURNING 1
	)
	SELECT count(*) FROM affected_rows;
	`

	row := p.connPool.QueryRow(ctx, sql, req.PostId)
	var countAffected int
	if err := row.Scan(&countAffected); err != nil {
		return err
	}

	if countAffected == 0 {
		return status.Error(codes.NotFound, "post not found")
	}

	if countAffected != 1 {
		log.Printf("warning: more than 1 post with id = %d", req.PostId)
	}

	return nil
}

func (p *PostgresDbClient) NewView(ctx context.Context, req dbReq.NewViewRequest) error {
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

	row := p.connPool.QueryRow(ctx, sql, req.PostId)
	var countAffected int
	if err := row.Scan(&countAffected); err != nil {
		return err
	}

	if countAffected == 0 {
		return status.Error(codes.NotFound, "post not found")
	}

	if countAffected != 1 {
		log.Println("warning: more than 1 post with id %d in posts table", req.PostId)
	}

	return nil
}

type pgPostsPagiToken struct {
	LastCreatedAt time.Time `json:"lcr"`
}

func decodePgPostsPagiToken(token dbReq.PaginationToken) (pgPostsPagiToken, error) {
	if token == "" {
		return pgPostsPagiToken{
			LastCreatedAt: time.Date(9999, time.December, 31, 23, 59, 59, 0, time.UTC),
		}, nil
	}

	raw, err := base64.RawURLEncoding.DecodeString(string(token))
	if err != nil {
		return pgPostsPagiToken{}, err
	}

	var decoded pgPostsPagiToken
	err = json.Unmarshal(raw, &decoded)
	if err != nil {
		return pgPostsPagiToken{}, err
	}

	return decoded, nil
}

func encodePgPostsPagiToken(token pgPostsPagiToken) (dbReq.PaginationToken, error) {
	raw, err := json.Marshal(&token)
	if err != nil {
		return "", err
	}

	encoded := base64.RawURLEncoding.EncodeToString(raw)
	return dbReq.PaginationToken(encoded), nil
}
