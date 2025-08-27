package types

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSurnameValid(t *testing.T) {
	validSurnameStr := "Smirnov"
	surname, err := unmarshalFromString[Surname](validSurnameStr)
	require.NoError(t, err, `cannot unmarshal valid name: "%s": %v`, validSurnameStr, err)
	assert.Equal(t, validSurnameStr, string(surname))
}

func TestSurnameEmpty(t *testing.T) {
	_, err := unmarshalFromString[Surname]("")
	require.Error(t, err, "empty name is invalid")
}

func TestSurnameTooLong(t *testing.T) {
	surnameStr := strings.Repeat("x", SURNAME_MAX_LENGTH+5)
	_, err := unmarshalFromString[Surname](surnameStr)
	require.Error(t, err, "name is too long, but unmarshable")
}
