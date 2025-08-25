package types

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNameValid(t *testing.T) {
	validNameStr := "Spiridon"
	name, err := unmarshalFromString[Name](validNameStr)
	require.NoError(t, err, `cannot unmarshal valid name: "%s"`, validNameStr)
	assert.Equal(t, validNameStr, string(name))
}

func TestNameEmpty(t *testing.T) {
	_, err := unmarshalFromString[Name]("")
	require.Error(t, err, "empty name is invalid")
}

func TestNameTooLong(t *testing.T) {
	nameStr := strings.Repeat("x", NAME_MAX_LENGTH+5)
	_, err := unmarshalFromString[Name](nameStr)
	require.Error(t, err, "name is too long, but unmarshable")
}
