package types

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBioValid(t *testing.T) {
	validBioStr := "this is valid bio"
	bio, err := unmarshalFromString[Bio](validBioStr)
	require.NoError(t, err, "cannot unmarshal valid bio")
	assert.Equal(t, validBioStr, string(bio))
}

func TestBioEmpty(t *testing.T) {
	bio, err := unmarshalFromString[Bio]("")
	require.NoError(t, err, "cannot unmarshal empty bio")
	assert.Equal(t, "", string(bio))
}

func TestBioTooLong(t *testing.T) {
	bioStr := strings.Repeat("x", BIO_MAX_LENGTH+5)
	_, err := unmarshalFromString[Bio](bioStr)
	require.Error(t, err, "bio is too long, but unmarshable")
}
