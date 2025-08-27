package types

import (
	"encoding/json"
	"errors"
	"time"
)

type Duration struct {
	time.Duration
}

var error_duration_negative = errors.New("duration must be >= 0")
var error_duration_invalid = errors.New("invalid duration")

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	dur, err := time.ParseDuration(s)
	if err != nil {
		return error_duration_invalid
	}

	if dur.Nanoseconds() < 0 {
		return error_duration_negative
	}

	*d = Duration{
		Duration: dur,
	}
	return nil
}
