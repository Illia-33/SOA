package server

import (
	"errors"
	"regexp"
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

// Parse a birthday from the string s in format DD-MM-YYYY
func parseBirthday(s string) (b birthday, err error) {
	birthdayRegexp := regexp.MustCompile(`^\d{2}-\d{2}-\d{4}$`)
	if !birthdayRegexp.MatchString(s) {
		err = errors.New("birthday must be in format DD-MM-YYYY")
		return
	}

	b.day = stringToInt(s[0:2])
	b.month = stringToInt(s[3:5])
	b.year = stringToInt(s[6:10])
	return
}

// Checks if birthday represents a valid date
func (b birthday) isValid() bool {
	if b.year < 0 {
		return false
	}

	if !(1 <= b.month && b.month <= 12) {
		return false
	}

	return 1 <= b.day && b.day <= daysPerMonth(b.month, b.year)
}
