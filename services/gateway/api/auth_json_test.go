package api

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func unmarshal[T any](s string) (T, error) {
	var unmarshalled T
	err := json.Unmarshal([]byte(s), &unmarshalled)
	return unmarshalled, err
}

func TestAuthSimple(t *testing.T) {
	rawJson := `
{
	"login": "some_login",
	"password": "some_passwd"
}
	`

	req, err := unmarshal[AuthenticateRequest](rawJson)

	require.NoError(t, err, "valid json unmarshalling failed")

	assert.True(t, req.Login.HasValue)
	assert.Equal(t, "some_login", string(req.Login.Value))

	assert.Equal(t, "some_passwd", string(req.Password))
}

func TestAuthTooMuchUserId(t *testing.T) {
	rawJson := `
{
	"login": "some_login",
	"email": "some_email@yahoo.com",
	"password": "some_passwd"
}
	`

	_, err := unmarshal[AuthenticateRequest](rawJson)
	require.ErrorAs(t, err, &ErrorTooMuchUserId{})
}

func TestAuthNoUserId(t *testing.T) {
	rawJson := `
{
	"password": "some_passwd"
}
	`

	_, err := unmarshal[AuthenticateRequest](rawJson)
	require.ErrorAs(t, err, &ErrorNoUserId{})
}
