package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testDate struct {
	year  int
	month int
	day   int
}

func (d testDate) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.year, d.month, d.day)
}

func makeDate(y int, m int, d int) testDate {
	return testDate{
		year:  y,
		month: m,
		day:   d,
	}
}

func TestDateValid(t *testing.T) {
	validDates := []testDate{makeDate(2003, 6, 12), makeDate(1995, 10, 23), makeDate(1960, 5, 5), makeDate(2024, 2, 29)}
	for _, inputDate := range validDates {
		dateStr := inputDate.String()
		date, err := unmarshalFromString[Date](dateStr)
		require.NoError(t, err, `cannot unmarshal "%s"`, inputDate)
		assert.Equal(t, inputDate.year, date.Year())
		assert.Equal(t, inputDate.month, int(date.Month()))
		assert.Equal(t, inputDate.day, date.Day())
	}
}

func TestDateInvalid(t *testing.T) {
	invalidDates := []testDate{makeDate(2005, 2, 29), makeDate(1990, 13, 10), makeDate(1985, 6, 31)}
	for _, inputDate := range invalidDates {
		dateStr := inputDate.String()
		_, err := unmarshalFromString[Date](dateStr)
		require.Error(t, err, `invalid date "%s" has been unmarshalled successfully`, inputDate)
	}
}

func TestDateInvalidInput(t *testing.T) {
	inputs := []string{"123-45-67", "19900301", "date", "999-01-01"}
	for _, input := range inputs {
		_, err := unmarshalFromString[Date](input)
		require.Error(t, err, `invalid input "%s" has been unmarshalled successfully`, input)
	}
}
