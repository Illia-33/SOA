package clickhouse

import (
	"context"
	"soa-socialnetwork/services/stats/pkg/models"

	chDriver "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type registrationsRepo struct {
	ctx  context.Context
	conn chDriver.Conn
}

func (r *registrationsRepo) Put(events ...models.RegistrationEvent) error {
	sql := `
	INSERT INTO registrations(account_id, profile_id, register_time)
	`

	batch, err := r.conn.PrepareBatch(r.ctx, sql)
	if err != nil {
		return err
	}

	for _, event := range events {
		err := batch.Append(event.AccountId, event.ProfileId, event.Timestamp)
		if err != nil {
			return err
		}
	}

	return batch.Send()
}
