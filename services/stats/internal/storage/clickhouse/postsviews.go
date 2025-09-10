package clickhouse

import (
	"context"
	"soa-socialnetwork/services/stats/internal/repo"
	"soa-socialnetwork/services/stats/pkg/models"

	chDriver "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type postsViewsRepo struct {
	ctx  context.Context
	conn chDriver.Conn
}

func (r *postsViewsRepo) GetCountForPost(id models.PostId) (int64, error) {
	sql := `
	SELECT count(*)
	FROM posts_views
	WHERE post_id = ?;
	`
	row := r.conn.QueryRow(r.ctx, sql, id)
	var count int64
	if err := row.Scan(&count); err != nil {
		return -1, err
	}

	return count, nil
}

func (r *postsViewsRepo) GetDynamicsForPost(id models.PostId) (repo.ViewDynamics, error) {
	sql := `
	SELECT day, count(*)
	FROM posts_views
	WHERE post_id = ?
	GROUP BY toYYYYMMDD(view_time) AS day;
	`
	rows, err := r.conn.Query(r.ctx, sql, id)
	if err != nil {
		return nil, err
	}

	var dynamics repo.ViewDynamics

	for rows.Next() {
		var dailyStat repo.DailyViewsStat
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

func (r *postsViewsRepo) Put(events []models.PostViewEvent) error {
	sql := `
	INSERT INTO posts_views(post_id, viewer_account_id, view_time)
	`

	batch, err := r.conn.PrepareBatch(r.ctx, sql)
	if err != nil {
		return err
	}

	for _, event := range events {
		err := batch.Append(event.PostId, event.ViewerAccountId, event.Timestamp)
		if err != nil {
			return err
		}
	}

	return batch.Send()
}
