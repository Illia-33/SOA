package types

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testDuration struct {
	hour        int
	minute      int
	second      int
	millisecond int
	microsecond int
	nanosecond  int
}

func (d testDuration) String() string {
	result := ""
	if d.hour != 0 {
		result += strconv.Itoa(d.hour) + "h"
	}
	if d.minute != 0 {
		result += strconv.Itoa(d.minute) + "m"
	}
	if d.second != 0 {
		result += strconv.Itoa(d.second) + "s"
	}
	if d.millisecond != 0 {
		result += strconv.Itoa(d.millisecond) + "ms"
	}
	if d.microsecond != 0 {
		result += strconv.Itoa(d.microsecond) + "us"
	}
	if d.nanosecond != 0 {
		result += strconv.Itoa(d.nanosecond) + "ns"
	}

	if len(result) == 0 {
		return "0ns"
	}

	return result
}

func (d testDuration) ToStdDuration() time.Duration {
	dur, err := time.ParseDuration(d.String())
	if err != nil {
		panic(err)
	}
	return dur
}

func TestDurationValid(t *testing.T) {
	validDurations := []testDuration{
		{hour: 43800},
		{hour: 25},
		{minute: 34},
		{minute: 3600},
		{second: 23},
		{millisecond: 24814},
		{nanosecond: 19345232},
		{hour: 2, minute: 30},
		{hour: 14, minute: 25, second: 30, millisecond: 25, microsecond: 21, nanosecond: 9},
		{nanosecond: 0},
	}

	for _, inputDur := range validDurations {
		dur, err := unmarshalFromString[Duration](inputDur.String())
		assert.NoError(t, err)
		if err == nil {
			assert.Equal(t, inputDur.ToStdDuration(), dur.Duration)
		}
	}
}

func TestDurationNegative(t *testing.T) {
	negativeDurations := []testDuration{
		{nanosecond: -30},
		{hour: -5},
		{minute: -34},
	}

	for _, inputDur := range negativeDurations {
		_, err := unmarshalFromString[Duration](inputDur.String())
		assert.Error(t, err, `successfully parsed negative duration from "%s"`, inputDur)
	}
}

func TestDurationInvalidInput(t *testing.T) {
	inputs := []string{
		"1 hour",
		"3600",
		"invalid duration",
		"14d",
		"1y14d",
		"3y",
	}

	for _, input := range inputs {
		_, err := unmarshalFromString[Duration](input)
		assert.Error(t, err, `successfully parsed negative duration from "%s"`, input)
	}
}
