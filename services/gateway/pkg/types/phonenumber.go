package types

import (
	"encoding/json"
	"fmt"
	"regexp"
)

type PhoneNumber string

const MIN_PHONE_NUMBER_LENGTH = 8
const MAX_PHONE_NUMBER_LENGTH = 16

var phone_number_regexp = regexp.MustCompile(`^\+\d{7,15}$`)
var error_phone_number_bad_length = fmt.Errorf("len(phone_number) must be in [%d;%d]", MIN_PHONE_NUMBER_LENGTH, MAX_PHONE_NUMBER_LENGTH)
var error_phone_number_invalid = fmt.Errorf("phone_number must match %s", phone_number_regexp)

func (p *PhoneNumber) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	if !(MIN_PHONE_NUMBER_LENGTH <= len(s) && len(s) <= MAX_PHONE_NUMBER_LENGTH) {
		return error_phone_number_bad_length
	}

	if !phone_number_regexp.MatchString(s) {
		return error_phone_number_invalid
	}

	*p = PhoneNumber(s)
	return nil
}

func (p PhoneNumber) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(p))
}
