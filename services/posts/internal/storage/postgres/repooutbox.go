package postgres

import (
	"context"
	"soa-socialnetwork/services/posts/internal/models"
	"soa-socialnetwork/services/posts/internal/repos"
	"time"
)

type outboxRepo struct {
	ctx   context.Context
	scope pgxScope
}

func (r outboxRepo) Put(event models.OutboxEvent) error {
	sql := `
	INSERT INTO outbox(event_type, payload)
	VALUES ($1, $2::jsonb);
	`
	_, err := r.scope.Exec(r.ctx, sql, event.Type, event.Payload)
	return err
}

func (r outboxRepo) Fetch(params repos.OutboxFetchParams) ([]models.OutboxEvent, error) {
	sql := `
	WITH cte AS (
		SELECT id, event_type, payload, created_at
		FROM outbox
		WHERE 
			created_at > $1 
			AND is_processed = FALSE
		ORDER BY created_at ASC
		LIMIT $2
		FOR UPDATE SKIP LOCKED
	)
	UPDATE outbox AS o
	SET is_processed = TRUE
	FROM cte
	WHERE o.id = cte.id
	RETURNING cte.event_type, cte.payload, cte.created_at;
	`

	rows, err := r.scope.Query(r.ctx, sql, params.LastCreatedAt, params.Limit)
	if err != nil {
		return nil, err
	}

	events := make([]models.OutboxEvent, 0, params.Limit)
	for {
		if !rows.Next() {
			if err := rows.Err(); err != nil {
				return nil, err
			}
			break
		}

		var (
			eventType   string
			jsonPayload string
			createdAt   time.Time
		)

		err := rows.Scan(&eventType, &jsonPayload, &createdAt)
		if err != nil {
			return nil, err
		}

		events = append(events, models.OutboxEvent{
			Type:      eventType,
			Payload:   models.OutboxEventPayload(jsonPayload),
			CreatedAt: createdAt,
		})
	}

	return events, nil
}
