package kafkajobs

import (
	"context"
	"errors"
	"io"
	"log"
	"soa-socialnetwork/services/stats/internal/kafka"
	"time"
)

type messageBatch[MsgType any] []kafka.Message[MsgType]
type batchProcessor[MsgType any] func(context.Context, messageBatch[MsgType]) error

type topicWorker[MsgType any] struct {
	consumer     kafka.Consumer[MsgType]
	processBatch batchProcessor[MsgType]
}

func newTopicProcessor[MsgType any](
	connCfg kafka.ConnectionConfig,
	readerCfg kafka.ConsumerConfig,
	processBatch batchProcessor[MsgType],
) (topicWorker[MsgType], error) {
	r, err := kafka.NewConsumer[MsgType](connCfg, readerCfg)
	if err != nil {
		return topicWorker[MsgType]{}, err
	}

	return topicWorker[MsgType]{
		consumer:     r,
		processBatch: processBatch,
	}, nil
}

func (r *topicWorker[MsgType]) start(ctx context.Context) {
	go func() {
		const BATCH_CAPACITY = 10
		batch := make(messageBatch[MsgType], 0, BATCH_CAPACITY)
		batchProcessed := false
		messages := make(chan kafka.Message[MsgType], BATCH_CAPACITY)

		go func() {
			for {
				msg, err := r.consumer.FetchMessage(ctx)
				if err != nil {
					if errors.Is(err, io.EOF) {
						break
					}

					log.Printf("error occured while fetching %T message from kafka: %v", msg, err)
					time.Sleep(500 * time.Millisecond)
					continue
				}
				messages <- msg
			}

			close(messages)
		}()

		tryProcessBatch := func() error {
			if len(batch) == 0 {
				return nil
			}

			if !batchProcessed {
				err := r.processBatch(ctx, batch)
				if err != nil {
					log.Printf("error occured while processing %T batch: %v", batch, err)
					return err
				}
				batchProcessed = true
			}

			err := r.consumer.CommitMessages(ctx, batch...)
			if err != nil {
				log.Printf("error occured while committing %T batch: %v", batch, err)
				return err
			}

			batch = batch[:0]
			batchProcessed = false
			return nil
		}

		done := false
		for !done {
			if len(batch) >= BATCH_CAPACITY {
				err := tryProcessBatch()
				if err != nil {
					time.Sleep(500 * time.Millisecond)
				}
				continue
			}

			select {
			case msg, ok := <-messages:
				{
					if !ok {
						done = true
						break
					}
					batch = append(batch, msg)
				}

			case <-time.NewTimer(500 * time.Millisecond).C:
				{
					tryProcessBatch()
				}
			}

		}

		r.close()
	}()
}

func (r *topicWorker[MsgType]) close() error {
	return r.consumer.Close()
}
