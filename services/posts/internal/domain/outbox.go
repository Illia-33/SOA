package domain

import "time"

type OutboxEventType string
type OutboxEventPayload []byte

type OutboxEvent struct {
	Type      string
	Payload   OutboxEventPayload
	CreatedAt time.Time
}
