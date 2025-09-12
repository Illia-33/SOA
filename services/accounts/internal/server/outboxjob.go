package server

import (
	"context"
	"log"
	"soa-socialnetwork/services/common/backjob"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
)

func checkOutboxJob(dbPool *pgxpool.Pool) backjob.JobCallback {
	lastCreatedAt := time.Time{}
	kafkaWriter := kafka.Writer{
		Addr:                   kafka.TCP("stats-kafka:9092"),
		RequiredAcks:           kafka.RequireAll,
		AllowAutoTopicCreation: true,
	}
	return func(ctx context.Context) error {
		tx, err := dbPool.Begin(ctx)
		if err != nil {
			return err
		}

		sql := `
		WITH cte AS (
			SELECT id, event_type, payload, created_at
			FROM outbox
			WHERE 
				created_at > $1 
				AND is_processed = FALSE
			ORDER BY created_at ASC
			LIMIT 100
			FOR UPDATE SKIP LOCKED
		)
		UPDATE outbox AS o
		SET is_processed = TRUE
		FROM cte
		WHERE o.id = cte.id
		RETURNING cte.event_type, cte.payload, cte.created_at;
		`

		rows, err := tx.Query(ctx, sql, lastCreatedAt)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}

		currentLastCreatedAt := time.Time{}
		var messages []kafka.Message
		for {
			if !rows.Next() {
				if err := rows.Err(); err != nil {
					return err
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
				return err
			}

			currentLastCreatedAt = createdAt
			log.Printf("got event of type '%s' (%v) in outbox, payload = '%s'", eventType, createdAt, jsonPayload)
			messages = append(messages, kafka.Message{
				Topic: eventType,
				Value: []byte(jsonPayload),
				Time:  createdAt,
			})
		}

		if len(messages) == 0 {
			tx.Rollback(ctx)
			return nil
		}

		err = kafkaWriter.WriteMessages(ctx, messages...)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}

		err = tx.Commit(ctx)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}

		lastCreatedAt = currentLastCreatedAt

		return nil
	}
}
