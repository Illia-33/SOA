package types

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPasswordValid(t *testing.T) {
	validPasswordStr := "passwd123"
	password, err := unmarshalFromString[Password](validPasswordStr)
	require.NoError(t, err, `cannot unmarshal valid password: "%s"`, validPasswordStr)
	assert.Equal(t, validPasswordStr, string(password))
}

func TestPasswordEmpty(t *testing.T) {
	_, err := unmarshalFromString[Password]("")
	require.Error(t, err, "empty password is invalid")
}

func TestPasswordTooLong(t *testing.T) {
	passwordStr := strings.Repeat("x", PASSWORD_MAX_LENGTH+5)
	_, err := unmarshalFromString[Password](passwordStr)
	require.Error(t, err, "password is too long, but unmarshable")
}

func TestPasswordTooShort(t *testing.T) {
	passwordStr := strings.Repeat("x", PASSWORD_MIN_LENGTH-1)
	_, err := unmarshalFromString[Password](passwordStr)
	require.Error(t, err, "password is too short, but unmarshable")
}
