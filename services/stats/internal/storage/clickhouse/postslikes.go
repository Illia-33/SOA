package clickhouse

import (
	"context"
	"soa-socialnetwork/services/stats/internal/repo"
	"soa-socialnetwork/services/stats/pkg/models"

	chDriver "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type postsLikesRepo struct {
	ctx  context.Context
	conn chDriver.Conn
}

func (r *postsLikesRepo) GetCountForPost(id models.PostId) (uint64, error) {
	sql := `
	SELECT count(*)
	FROM posts_likes
	WHERE post_id = ?;
	`
	row := r.conn.QueryRow(r.ctx, sql, id)
	var count uint64
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *postsLikesRepo) GetDynamicsForPost(id models.PostId) (repo.LikeDynamics, error) {
	sql := `
	SELECT day, count(*)
	FROM posts_likes
	WHERE post_id = ?
	GROUP BY toYYYYMMDD(like_time) AS day;
	`
	rows, err := r.conn.Query(r.ctx, sql, id)
	if err != nil {
		return nil, err
	}

	var dynamics repo.LikeDynamics

	for rows.Next() {
		var dailyStat repo.DailyLikeStat
		var ymd uint32
		err := rows.Scan(&ymd, &dailyStat.Count)
		if err != nil {
			return nil, err
		}

		dailyStat.Date = dateFromYYYYMMDD(ymd)

		dynamics = append(dynamics, dailyStat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dynamics, nil
}

func (r *postsLikesRepo) Put(events ...models.PostLikeEvent) error {
	sql := `
	INSERT INTO posts_likes(post_id, liker_account_id, like_time)
	`

	batch, err := r.conn.PrepareBatch(r.ctx, sql)
	if err != nil {
		return err
	}
	defer batch.Close()

	for _, event := range events {
		err := batch.Append(event.PostId, event.LikerAccountId, event.Timestamp)
		if err != nil {
			return err
		}
	}

	return batch.Send()
}
