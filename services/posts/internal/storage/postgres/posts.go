package postgres

import (
	"context"
	"fmt"
	"log"
	opt "soa-socialnetwork/services/common/option"
	dom "soa-socialnetwork/services/posts/internal/domain"
	"soa-socialnetwork/services/posts/internal/repos"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostsRepo struct {
	ConnPool connectionPool
}

func (r *PostsRepo) New(ctx context.Context, pageId dom.PageId, data repos.NewPostData) (dom.PostId, error) {
	sql := `
	INSERT INTO posts(page_id, author_account_id, text_content, source_post_id)
	VALUES ($1, $2, $3, $4)
	RETURNING id;
	`

	pgSourcePostId := pgtype.Int4{Int32: int32(data.Content.SourcePostId.Value), Valid: data.Content.SourcePostId.HasValue}

	row := r.ConnPool.QueryRow(ctx, sql, pageId, data.AuthorId, data.Content.Text, pgSourcePostId)
	var postId dom.PostId
	err := row.Scan(&postId)
	if err != nil {
		return 0, err
	}

	return postId, nil
}

const POSTS_PAGE_SIZE = 10

type postsPagiToken struct {
	LastCreatedAt time.Time `json:"lcr"`
}

func decodePostsPagiToken(token repos.PagiToken) (postsPagiToken, error) {
	if token == "" {
		return postsPagiToken{
			LastCreatedAt: time.Date(9999, time.December, 31, 23, 59, 59, 0, time.UTC),
		}, nil
	}
	return decodePagiToken[postsPagiToken](token)
}

func encodePostsPagiToken(token postsPagiToken) (repos.PagiToken, error) {
	return encodePagiToken(token)
}

func (r *PostsRepo) List(ctx context.Context, pageId dom.PageId, encodedPagiToken repos.PagiToken) (repos.PostsList, error) {
	sql := fmt.Sprintf(`
	SELECT id, author_account_id, text_content, source_post_id, pinned, views_count, created_at
	FROM posts
	WHERE page_id = $1 AND created_at < $2
	ORDER BY created_at DESC
	LIMIT %d;
	`, POSTS_PAGE_SIZE)

	pagiToken, err := decodePostsPagiToken(encodedPagiToken)
	if err != nil {
		return repos.PostsList{}, err
	}

	rows, err := r.ConnPool.Query(ctx, sql, pageId, pagiToken.LastCreatedAt)
	if err != nil {
		return repos.PostsList{}, err
	}

	posts := make([]dom.Post, 0, POSTS_PAGE_SIZE)

	for {
		if !rows.Next() {
			err := rows.Err()
			if err != nil {
				return repos.PostsList{}, err
			}

			break
		}

		var post dom.Post
		var pgSourcePostId pgtype.Int4
		err := rows.Scan(&post.Id, &post.AuthorAccountId, &post.Content.Text, &pgSourcePostId, &post.Pinned, &post.ViewsCount, &post.CreatedAt)
		if err != nil {
			return repos.PostsList{}, err
		}

		post.PageId = pageId
		post.Content.SourcePostId = opt.Option[dom.PostId]{Value: dom.PostId(pgSourcePostId.Int32), HasValue: pgSourcePostId.Valid}
		posts = append(posts, post)
	}

	var nextPagiToken repos.PagiToken
	if len(posts) > 0 {
		token := postsPagiToken{
			LastCreatedAt: posts[len(posts)-1].CreatedAt,
		}
		encodedToken, err := encodePostsPagiToken(token)

		if err != nil {
			log.Printf("warning: cannot encode paginating token (%v): %v", token, err)
		} else {
			nextPagiToken = encodedToken
		}
	}

	return repos.PostsList{
		Posts:         posts,
		NextPagiToken: nextPagiToken,
	}, nil
}

func (r *PostsRepo) Get(ctx context.Context, postId dom.PostId) (dom.Post, error) {
	sql := `
	SELECT page_id, author_account_id, text_content, source_post_id, pinned, views_count, created_at
	FROM posts
	WHERE id = $1;
	`

	var post dom.Post
	var pgSourcePostId pgtype.Int4

	row := r.ConnPool.QueryRow(ctx, sql, postId)
	err := row.Scan(&post.PageId, &post.AuthorAccountId, &post.Content.Text, &pgSourcePostId, &post.Pinned, &post.ViewsCount, &post.CreatedAt)
	if err != nil {
		return dom.Post{}, err
	}

	post.Id = postId
	post.Content.SourcePostId = opt.Option[dom.PostId]{
		Value:    dom.PostId(pgSourcePostId.Int32),
		HasValue: pgSourcePostId.Valid,
	}

	return post, nil
}

func (r *PostsRepo) Edit(ctx context.Context, postId dom.PostId, edited repos.EditedPostData) error {
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
		String: string(edited.Text.Value),
		Valid:  edited.Text.HasValue,
	}
	pgPinned := pgtype.Bool{
		Bool:  edited.Pinned.Value,
		Valid: edited.Pinned.HasValue,
	}

	row := r.ConnPool.QueryRow(ctx, sql, pgTextContent, pgPinned, postId)
	var countAffected int
	if err := row.Scan(&countAffected); err != nil {
		return err
	}

	if countAffected == 0 {
		return status.Error(codes.NotFound, "post not found")
	}

	if countAffected != 1 {
		log.Printf("warning: more than 1 post with id = %d", postId)
	}

	return nil
}

func (r *PostsRepo) Delete(ctx context.Context, postId dom.PostId) error {
	sql := `
	WITH affected_rows AS (
		DELETE FROM posts
		WHERE id = $1
		RETURNING 1
	)
	SELECT count(*) FROM affected_rows;
	`

	row := r.ConnPool.QueryRow(ctx, sql, postId)
	var countAffected int
	if err := row.Scan(&countAffected); err != nil {
		return err
	}

	if countAffected == 0 {
		return status.Error(codes.NotFound, "post not found")
	}

	if countAffected != 1 {
		log.Printf("warning: more than 1 post with id = %d", postId)
	}

	return nil
}
