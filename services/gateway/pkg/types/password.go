package types

import (
	"encoding/json"
	"fmt"
)

type Password string

const PASSWORD_MIN_LENGTH = 6
const PASSWORD_MAX_LENGTH = 32

var error_password_bad_length = fmt.Errorf("len(password) must be in [%d;%d]", PASSWORD_MIN_LENGTH, PASSWORD_MAX_LENGTH)

func (p *Password) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	if !(PASSWORD_MIN_LENGTH <= len(s) && len(s) <= PASSWORD_MAX_LENGTH) {
		return error_password_bad_length
	}

	*p = Password(s)
	return nil
}

func (p Password) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(p))
}
