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
	messages     chan kafka.Message[MsgType]
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
	const MSG_CAPACITY = 100
	r.messages = make(chan kafka.Message[MsgType], MSG_CAPACITY)
	r.startConsume(ctx)
	r.startWork(ctx)
}

func (r *topicWorker[MsgType]) startConsume(ctx context.Context) {
	go func() {
		for {
			msg, err := r.consumer.FetchMessage(ctx)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}

				log.Printf("error occured while fetching %T message from kafka: %v", msg, err)
				time.Sleep(time.Second)
				continue
			}

			r.messages <- msg
		}

		close(r.messages)
		r.close()
	}()
}

func (r *topicWorker[MsgType]) startWork(ctx context.Context) {
	go func() {
		const BATCH_CAPACITY = 10
		batch := make(messageBatch[MsgType], 0, BATCH_CAPACITY)

		processBatch := func() error {
			if len(batch) == 0 {
				return nil
			}
			err := r.processBatch(ctx, batch)
			if err != nil {
				return err
			}

			batch = batch[:0]
			return nil
		}

		processBatchWithLogging := func() error {
			err := processBatch()
			if err != nil {
				log.Printf("error occured while processing %T batch: %v", batch, err)
			}

			return err
		}

		done := false
		for !done {
			for len(batch) >= BATCH_CAPACITY {
				err := processBatchWithLogging()
				if err != nil {
					time.Sleep(500 * time.Millisecond)
				}
			}

			select {
			case msg, ok := <-r.messages:
				{
					if !ok {
						done = true
						break
					}

					batch = append(batch, msg)
				}

			case <-time.NewTimer(500 * time.Millisecond).C:
				{
					processBatchWithLogging()
				}
			}
		}
	}()
}

func (r *topicWorker[MsgType]) close() error {
	return r.consumer.Close()
}
