package kafkajobs

import (
	"context"
	"encoding/json"
	"soa-socialnetwork/services/stats/internal/kafka"
	"sync"
	"testing"
	"time"

	kafkago "github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
)

func marshalJson[T any](t *testing.T, obj T) []byte {
	data, err := json.Marshal(obj)
	require.NoError(t, err)
	return data
}

type testKafkaMessage struct {
	Value int
}

func (s *kafkaTestSuite) TestConsumeTopic() {
	ctx := context.Background()

	func() {
		writer := s.kafkaWriter()
		defer writer.Close()

		msgs := make([]kafkago.Message, 100)
		for i := range msgs {
			msgs[i] = kafkago.Message{
				Value: marshalJson(s.T(), testKafkaMessage{Value: i}),
			}
		}

		err := writer.WriteMessages(ctx, msgs...)
		s.Require().NoError(err, "cannot write test messages")
	}()

	consumer, err := kafka.NewConsumer[testKafkaMessage](
		kafka.ConnectionConfig{
			Host: s.hostname,
			Port: s.port,
		},
		kafka.ConsumerConfig{
			Topic:   kafka_test_topic,
			GroupId: "test-consume-topic",
		},
	)
	s.Require().NoError(err, "cannot create consumer")

	mutex := sync.Mutex{}
	msgs := make([]testKafkaMessage, 0, 100)

	callback := func(_ context.Context, batch messageBatch[testKafkaMessage]) error {
		mutex.Lock()
		defer mutex.Unlock()
		for _, msg := range batch {
			msgs = append(msgs, msg.Value)
		}
		return nil
	}

	worker := topicWorker[testKafkaMessage]{
		consumer:             consumer,
		processBatchCallback: callback,

		batchCapacity:         10,
		messagesChanCapacity:  30,
		processorTick:         50 * time.Millisecond,
		consumerErrorWaitTime: 50 * time.Millisecond,
	}

	worker.start(ctx)

	time.Sleep(10 * time.Second)

	mutex.Lock()
	defer mutex.Unlock()
	s.Require().Equal(100, len(msgs))
	for i, msg := range msgs {
		s.Assert().EqualValues(i, msg.Value)
	}
}
