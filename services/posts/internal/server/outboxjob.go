package server

import (
	"context"
	"log"
	"soa-socialnetwork/services/common/backjob"
	"soa-socialnetwork/services/posts/internal/repos"

	"github.com/segmentio/kafka-go"
)

func newCheckOutboxCallback(outbox repos.OutboxRepository, eventsPerCall int) backjob.JobCallback {
	kafkaWriter := kafka.Writer{
		Addr:                   kafka.TCP("kafka:9092"),
		RequiredAcks:           kafka.RequireAll,
		AllowAutoTopicCreation: true,
	}
	return func(ctx context.Context) error {
		events, tx, err := outbox.Fetch(ctx, repos.OutboxFetchParams{
			Limit: eventsPerCall,
		})
		if err != nil {
			return err
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
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				log.Printf("warning: cannot rollback outbox transaction: %v", rollbackErr)
			}
			return err
		}

		return tx.Commit(ctx)
	}
}
