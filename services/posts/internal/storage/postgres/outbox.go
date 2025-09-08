package postgres

import (
	"context"
	"fmt"
	dom "soa-socialnetwork/services/posts/internal/domain"
	"soa-socialnetwork/services/posts/internal/repos"
	"time"

	"github.com/jackc/pgx/v5"
)

type OutboxRepo struct {
	ConnPool connectionPool

	lastCreatedAt time.Time
}

func (r *OutboxRepo) Put(ctx context.Context, event dom.OutboxEvent) error {
	sql := `
	INSERT INTO outbox(event_type, payload)
	VALUES ($1, $2::jsonb);
	`
	_, err := r.ConnPool.Exec(ctx, sql, event.Type, event.Payload)
	return err
}

func (r *OutboxRepo) Fetch(ctx context.Context, params repos.OutboxFetchParams) ([]dom.OutboxEvent, repos.Transaction, error) {
	tx, err := r.ConnPool.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}

	sql := fmt.Sprintf(`
	WITH cte AS (
		SELECT id, event_type, payload, created_at
		FROM outbox
		WHERE 
			created_at > $1 
			AND is_processed = FALSE
		ORDER BY created_at ASC
		LIMIT %d
		FOR UPDATE SKIP LOCKED
	)
	UPDATE outbox AS o
	SET is_processed = TRUE
	FROM cte
	WHERE o.id = cte.id
	RETURNING cte.event_type, cte.payload, cte.created_at;
	`, params.Limit)

	rows, err := tx.Query(ctx, sql, r.lastCreatedAt)
	if err != nil {
		tx.Rollback(ctx)
		return nil, nil, err
	}

	currentLastCreatedAt := time.Time{}
	events := make([]dom.OutboxEvent, 0, params.Limit)
	for {
		if !rows.Next() {
			if err := rows.Err(); err != nil {
				return nil, nil, err
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
			tx.Rollback(ctx)
			return nil, nil, err
		}

		currentLastCreatedAt = createdAt
		events = append(events, dom.OutboxEvent{
			Type:    eventType,
			Payload: dom.OutboxEventPayload(jsonPayload),
		})
	}

	return events, outboxRepoTransaction{repo: r, tx: tx, txLastCreatedAt: currentLastCreatedAt}, nil
}

type outboxRepoTransaction struct {
	repo            *OutboxRepo
	tx              pgx.Tx
	txLastCreatedAt time.Time
}

func (t outboxRepoTransaction) Commit(ctx context.Context) error {
	err := t.tx.Commit(ctx)
	if err != nil {
		return err
	}

	t.repo.lastCreatedAt = t.txLastCreatedAt
	return nil
}

func (t outboxRepoTransaction) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}
