package types

import (
	"encoding/json"
	"errors"
	"time"
)

type Date struct {
	time.Time
}

var error_date_bad_format = errors.New("date invalid, must be in format YYYY-MM-DD")

func (d *Date) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	tm, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return error_date_bad_format
	}

	*d = Date{
		Time: tm,
	}
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Format(time.DateOnly))
}
