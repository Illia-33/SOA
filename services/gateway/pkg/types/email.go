package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
)

type Email string

const EMAIL_MIN_LENGTH = 4
const EMAIL_MAX_LENGTH = 320

var error_email_bad_length = fmt.Errorf("len(email) must be in [%d;%d]", EMAIL_MIN_LENGTH, EMAIL_MAX_LENGTH)
var error_email_invalid = errors.New("invalid email")

var email_regexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func (e *Email) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	if !(EMAIL_MIN_LENGTH <= len(s) && len(s) <= EMAIL_MAX_LENGTH) {
		return error_email_bad_length
	}

	if !email_regexp.MatchString(s) {
		return error_email_invalid
	}

	*e = Email(s)
	return nil
}

func (e Email) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(e))
}
