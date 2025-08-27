package types

import (
	"encoding/json"
	"fmt"
)

type Login string

const LOGIN_MIN_LENGTH = 1
const LOGIN_MAX_LENGTH = 32

var error_login_bad_length = fmt.Errorf("len(login) must be in [%d;%d]", LOGIN_MIN_LENGTH, LOGIN_MAX_LENGTH)

func (l *Login) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	if !(LOGIN_MIN_LENGTH <= len(s) && len(s) <= LOGIN_MAX_LENGTH) {
		return error_login_bad_length
	}

	*l = Login(s)
	return nil
}

func (l Login) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(l))
}
