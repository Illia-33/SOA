package birthday

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

const (
	JANUARY   int = 1
	FEBRUARY  int = 2
	MARCH     int = 3
	APRIL     int = 4
	MAY       int = 5
	JUNE      int = 6
	JULY      int = 7
	AUGUST    int = 8
	SEPTEMBER int = 9
	OCTOBER   int = 10
	NOVEMBER  int = 11
	DECEMBER  int = 12
)

func isLeapYear(year int) bool {
	if year%400 == 0 {
		return true
	}

	if year%100 == 0 {
		return false
	}

	return year%4 == 0
}

func daysPerMonth(month int, year int) int {
	switch month {
	case JANUARY, MARCH, MAY, JULY, AUGUST, OCTOBER, DECEMBER:
		{
			return 31
		}

	case APRIL, JUNE, SEPTEMBER, NOVEMBER:
		{
			return 30
		}

	case FEBRUARY:
		{
			if isLeapYear(year) {
				return 29
			}

			return 28
		}
	}

	panic("shouldn't reach here")
}

type birthday struct {
	day   int
	month int
	year  int
}

// Parse a birthday from the string s in format YYYY-MM-DD
func Parse(s string) (b birthday, err error) {
	birthdayRegexp := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	if !birthdayRegexp.MatchString(s) {
		err = errors.New("birthday must be in format YYYY-MM-DD")
		return
	}

	b.day = stringToInt(s[0:2])
	b.month = stringToInt(s[3:5])
	b.year = stringToInt(s[6:10])
	return
}

// Checks if birthday represents a valid date
func (b birthday) IsValid() bool {
	if b.year < 0 {
		return false
	}

	if !(1 <= b.month && b.month <= 12) {
		return false
	}

	return 1 <= b.day && b.day <= daysPerMonth(b.month, b.year)
}

func (b birthday) YYYY_MM_DD() string {
	return fmt.Sprintf("%04d-%02d-%02d", b.year, b.month, b.day)
}

func (b birthday) AsTime() time.Time {
	return time.Date(b.year, time.Month(b.month), b.day, 0, 0, 0, 0, nil)
}
