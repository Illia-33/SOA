package server

import (
	"log"
	"strconv"
)

// Converts string to int.
// Panics if an error is occurred, so use it carefully.
func stringToInt(s string) int {
	number, err := strconv.Atoi(s)
	if err != nil {
		log.Panic(err)
	}

	return number
}
