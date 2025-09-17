package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
)

type ConnectionConfig struct {
	Host string
	Port int
}

type ConsumerConfig struct {
	Topic   string
	GroupId string
}

type Consumer[MsgType any] struct {
	reader *kafka.Reader
}

type Message[MsgType any] struct {
	Value MsgType

	kafkaMsg kafka.Message
}

func NewConsumer[MsgType any](connCfg ConnectionConfig, readerCfg ConsumerConfig) (Consumer[MsgType], error) {
	cfg := kafka.ReaderConfig{
		Brokers: []string{fmt.Sprintf("%s:%d", connCfg.Host, connCfg.Port)},
		GroupID: readerCfg.GroupId,
		Topic:   readerCfg.Topic,
	}

	reader := kafka.NewReader(cfg)
	{
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-c
			reader.Close()
		}()
	}

	return Consumer[MsgType]{
		reader: reader,
	}, nil
}

func (c *Consumer[MsgType]) FetchMessage(ctx context.Context) (Message[MsgType], error) {
	kafkaMsg, err := c.reader.FetchMessage(ctx)
	if err != nil {
		return Message[MsgType]{}, err
	}

	var msg MsgType
	json.Unmarshal(kafkaMsg.Value, &msg)

	return Message[MsgType]{
		Value:    msg,
		kafkaMsg: kafkaMsg,
	}, nil
}

func (c *Consumer[MsgType]) CommitMessages(ctx context.Context, messages ...Message[MsgType]) error {
	if len(messages) == 0 {
		return nil
	}

	kafkaMsgs := make([]kafka.Message, len(messages))
	for i := range messages {
		kafkaMsgs[i] = messages[i].kafkaMsg
	}

	return c.reader.CommitMessages(ctx, kafkaMsgs...)
}

func (c *Consumer[MsgType]) Close() error {
	if c.reader == nil {
		return nil
	}

	err := c.reader.Close()
	if err != nil {
		return err
	}

	c.reader = nil
	return nil
}
