package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	reader    *kafka.Reader
	closeChan chan struct{}
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
		// Logger:      log.Default(),
		// ErrorLogger: log.Default(),
	}

	reader := kafka.NewReader(cfg)
	closeChan := make(chan struct{})
	{
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-c
			closeChan <- struct{}{}
		}()
	}

	consumer := Consumer[MsgType]{
		reader:    reader,
		closeChan: closeChan,
	}

	go func() {
		<-closeChan
		err := consumer.close()
		if err != nil {
			log.Printf("error occurred while closing kafka consumer: %v", err)
		}
	}()

	return consumer, nil
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

func (c *Consumer[MsgType]) Close() {
	if c.closeChan != nil {
		c.closeChan <- struct{}{}
		c.closeChan = nil
	}
}

func (c *Consumer[MsgType]) close() error {
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
