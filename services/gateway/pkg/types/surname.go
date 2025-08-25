package types

import (
	"encoding/json"
	"fmt"
)

type Surname string

const SURNAME_MIN_LENGTH = 1
const SURNAME_MAX_LENGTH = 32

var error_surname_bad_length = fmt.Errorf("len(surname) must be in [%d;%d]", SURNAME_MIN_LENGTH, SURNAME_MAX_LENGTH)

func (n *Surname) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	if !(SURNAME_MIN_LENGTH <= len(s) && len(s) <= SURNAME_MAX_LENGTH) {
		return error_surname_bad_length
	}

	*n = Surname(s)
	return nil
}

func (n Surname) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(n))
}
