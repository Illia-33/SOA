package clickhouse

import (
	"context"
	"soa-socialnetwork/services/stats/internal/repo"
	"soa-socialnetwork/services/stats/pkg/models"

	chDriver "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type postsCommentsRepo struct {
	ctx  context.Context
	conn chDriver.Conn
}

func (r *postsCommentsRepo) GetCountForPost(id models.PostId) (uint64, error) {
	sql := `
	SELECT count(*)
	FROM posts_comments
	WHERE post_id = ?;
	`
	row := r.conn.QueryRow(r.ctx, sql, id)
	var count uint64
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *postsCommentsRepo) GetDynamicsForPost(id models.PostId) (repo.CommentDynamics, error) {
	sql := `
	SELECT day, count(*)
	FROM posts_comments
	WHERE post_id = ?
	GROUP BY toYYYYMMDD(view_time) AS day;
	`
	rows, err := r.conn.Query(r.ctx, sql, id)
	if err != nil {
		return nil, err
	}

	var dynamics repo.CommentDynamics

	for rows.Next() {
		var dailyStat repo.DailyCommentsStat
		err := rows.Scan(&dailyStat.Date, &dailyStat.Count)
		if err != nil {
			return nil, err
		}

		dynamics = append(dynamics, dailyStat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dynamics, nil
}

func (r *postsCommentsRepo) Put(events ...models.PostCommentEvent) error {
	sql := `
	INSERT INTO posts_comments(post_id, author_account_id, comment_id, post_time)
	`

	batch, err := r.conn.PrepareBatch(r.ctx, sql)
	if err != nil {
		return err
	}
	defer batch.Close()

	for _, event := range events {
		err := batch.Append(event.PostId, event.AuthorAccountId, event.CommentId, event.Timestamp)
		if err != nil {
			return err
		}
	}

	return batch.Send()
}
