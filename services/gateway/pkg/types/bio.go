package types

import (
	"encoding/json"
	"fmt"
)

type Bio string

const BIO_MIN_LENGTH = 0
const BIO_MAX_LENGTH = 256

var error_bio_bad_length = fmt.Errorf("len(bio) must be in [%d;%d]", BIO_MIN_LENGTH, BIO_MAX_LENGTH)

func (n *Bio) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	if !(BIO_MIN_LENGTH <= len(s) && len(s) <= BIO_MAX_LENGTH) {
		return error_bio_bad_length
	}

	*n = Bio(s)
	return nil
}

func (n Bio) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(n))
}
