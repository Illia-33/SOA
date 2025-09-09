package server

import (
	"context"
	"log"
	"soa-socialnetwork/services/common/backjob"
	"soa-socialnetwork/services/posts/internal/repos"
	"time"

	"github.com/segmentio/kafka-go"
)

func newCheckOutboxCallback(db repos.RepoScopeOpener, eventsPerCall int) backjob.JobCallback {
	kafkaWriter := kafka.Writer{
		Addr:                   kafka.TCP("kafka:9092"),
		RequiredAcks:           kafka.RequireAll,
		AllowAutoTopicCreation: true,
	}

	lastCreatedAt := time.Time{}

	return func(ctx context.Context) error {
		tx, err := db.BeginTransaction(ctx)
		if err != nil {
			return err
		}
		defer tx.Close()

		events, err := tx.Outbox().Fetch(repos.OutboxFetchParams{
			Limit:         eventsPerCall,
			LastCreatedAt: lastCreatedAt,
		})

		if err != nil {
			tx.Rollback()
			return err
		}

		if len(events) == 0 {
			tx.Rollback()
			return nil
		}

		log.Println("got events")
		for i, event := range events {
			log.Printf("event #%d: %+v", i, event)
		}

		messages := make([]kafka.Message, len(events))
		for i, event := range events {
			messages[i] = kafka.Message{
				Topic: event.Type,
				Value: []byte(event.Payload),
			}
		}

		err = kafkaWriter.WriteMessages(ctx)
		if err != nil {
			tx.Rollback()
			return err
		}

		err = tx.Commit()
		if err != nil {
			return err
		}

		lastCreatedAt = events[len(events)-1].CreatedAt
		return nil
	}
}
