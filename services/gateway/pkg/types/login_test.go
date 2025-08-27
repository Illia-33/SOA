package types

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginValid(t *testing.T) {
	validLoginStr := "valid_login"
	login, err := unmarshalFromString[Login](validLoginStr)
	require.NoError(t, err, `cannot unmarshal valid login: "%s"`, validLoginStr)
	assert.Equal(t, validLoginStr, string(login))
}

func TestLoginEmpty(t *testing.T) {
	_, err := unmarshalFromString[Login]("")
	require.Error(t, err, "empty login is invalid")
}

func TestLoginTooLong(t *testing.T) {
	bioStr := strings.Repeat("x", LOGIN_MAX_LENGTH+5)
	_, err := unmarshalFromString[Login](bioStr)
	require.Error(t, err, "bio is too long, but unmarshable")
}
