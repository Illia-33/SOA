package models

import "time"

type OutboxEvent struct {
	Type      string
	Payload   []byte
	CreatedAt time.Time
}
