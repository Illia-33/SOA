package types

import (
	"encoding/json"
	"fmt"
)

type Name string

const NAME_MIN_LENGTH = 1
const NAME_MAX_LENGTH = 32

var error_name_bad_length = fmt.Errorf("len(name) must be in [%d;%d]", NAME_MIN_LENGTH, NAME_MAX_LENGTH)

func (n *Name) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	if !(NAME_MIN_LENGTH <= len(s) && len(s) <= NAME_MAX_LENGTH) {
		return error_name_bad_length
	}

	*n = Name(s)
	return nil
}

func (n Name) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(n))
}
