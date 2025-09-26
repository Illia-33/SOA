package models

import "time"

type ProfileId string

type ProfileData struct {
	AccountId AccountId
	ProfileId ProfileId
	Name      string
	Surname   string
	Birthday  time.Time
	Bio       string
}
