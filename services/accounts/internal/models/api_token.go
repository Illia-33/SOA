package models

import "time"

type ApiToken string

type ApiTokenData struct {
	Token       ApiToken
	AccountId   int
	ReadAccess  bool
	WriteAccess bool
	CreatedAt   time.Time
	ValidUntil  time.Time
}
