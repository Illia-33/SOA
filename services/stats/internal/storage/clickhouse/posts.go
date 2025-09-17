package clickhouse

import (
	"context"
	"soa-socialnetwork/services/stats/pkg/models"

	chDriver "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type postsRepo struct {
	ctx  context.Context
	conn chDriver.Conn
}

func (r *postsRepo) Put(events ...models.PostEvent) error {
	sql := `
	INSERT INTO posts(post_id, author_id, post_time)
	`

	batch, err := r.conn.PrepareBatch(r.ctx, sql)
	if err != nil {
		return err
	}
	defer batch.Close()

	for _, event := range events {
		err := batch.Append(event.PostId, event.AuthorId, event.Timestamp)
		if err != nil {
			return err
		}
	}

	return batch.Send()
}
