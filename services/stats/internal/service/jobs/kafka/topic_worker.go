package kafkajobs

import (
	"context"
	"errors"
	"io"
	"log"
	"soa-socialnetwork/services/stats/internal/kafka"
	"time"
)

type iTopicWorker interface {
	start(ctx context.Context)
	close() error
}

type messageBatch[TMsg any] []kafka.Message[TMsg]
type batchProcessor[TMsg any] func(context.Context, messageBatch[TMsg]) error

type topicWorker[TMsg any] struct {
	consumer     kafka.Consumer[TMsg]
	processBatch batchProcessor[TMsg]
}

func newTopicWorker[TMsg any](
	connCfg kafka.ConnectionConfig,
	readerCfg kafka.ConsumerConfig,
	processBatch batchProcessor[TMsg],
) (topicWorker[TMsg], error) {
	r, err := kafka.NewConsumer[TMsg](connCfg, readerCfg)
	if err != nil {
		return topicWorker[TMsg]{}, err
	}

	return topicWorker[TMsg]{
		consumer:     r,
		processBatch: processBatch,
	}, nil
}

func (w *topicWorker[TMsg]) start(ctx context.Context) {
	const DEFAULT_BATCH_CAPACITY = 10
	const DEFAULT_CHAN_CAPACITY = 30

	messagesChan := make(chan kafka.Message[TMsg], DEFAULT_CHAN_CAPACITY)
	batchContext := topicProcessorContext[TMsg]{
		bufferSize: DEFAULT_BATCH_CAPACITY,
		processor:  w.processBatch,
		consumer:   &w.consumer,
	}
	batchContext.init()

	w.startConsumerRoutine(ctx, messagesChan, &w.consumer)
	w.startProcessorRoutine(ctx, messagesChan, batchContext)
}

func (w *topicWorker[TMsg]) startConsumerRoutine(ctx context.Context, c chan<- kafka.Message[TMsg], consumer *kafka.Consumer[TMsg]) {
	go func() {
		for {
			msg, err := consumer.FetchMessage(ctx)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}

				log.Printf("error occured while fetching %T message from kafka: %v", msg, err)
				time.Sleep(500 * time.Millisecond)
				continue
			}
			c <- msg
		}

		close(c)
	}()
}

func (w *topicWorker[TMsg]) startProcessorRoutine(ctx context.Context, c <-chan kafka.Message[TMsg], processorCtx topicProcessorContext[TMsg]) {
	go func() {
		done := false
		for !done {
			if processorCtx.isBatchBufferFull() {
				err := processorCtx.tryProcessBatch(ctx)
				if err != nil {
					time.Sleep(500 * time.Millisecond)
				}
				continue
			}

			select {
			case msg, ok := <-c:
				{
					if !ok {
						done = true
						break
					}
					processorCtx.putMessage(msg)
				}

			case <-time.NewTimer(500 * time.Millisecond).C:
				{
					processorCtx.tryProcessBatch(ctx)
				}
			}

		}

		w.close()
	}()
}

func (w *topicWorker[TMsg]) close() error {
	return w.consumer.Close()
}

type topicProcessorContext[TMsg any] struct {
	bufferSize int
	processor  batchProcessor[TMsg]
	consumer   *kafka.Consumer[TMsg]

	batch            messageBatch[TMsg]
	isBatchProcessed bool
}

func (c *topicProcessorContext[TMsg]) init() {
	c.batch = make(messageBatch[TMsg], 0, c.bufferSize)
	c.isBatchProcessed = false
}

func (c *topicProcessorContext[TMsg]) putMessage(msg kafka.Message[TMsg]) {
	c.batch = append(c.batch, msg)
}

func (c *topicProcessorContext[TMsg]) isBatchBufferFull() bool {
	return len(c.batch) >= c.bufferSize
}

func (c *topicProcessorContext[TMsg]) tryProcessBatch(ctx context.Context) error {
	if len(c.batch) == 0 {
		return nil
	}

	if !c.isBatchProcessed {
		err := c.processor(ctx, c.batch)
		if err != nil {
			log.Printf("error occured while processing %T batch: %v", c.batch, err)
			return err
		}

		c.isBatchProcessed = true
	}

	err := c.consumer.CommitMessages(ctx, c.batch...)
	if err != nil {
		log.Printf("error occured while committing %T batch: %v", c.batch, err)
		return err
	}

	c.batch = c.batch[:0]
	c.isBatchProcessed = false
	return nil
}
