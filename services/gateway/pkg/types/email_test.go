package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailValid(t *testing.T) {
	validEmails := []string{
		`simple@example.com`,
		`very.common@example.com`,
		`disposable.style.email.with+symbol@example.com`,
		`other.email-with-hyphen@example.com`,
		`fully-qualified-domain@example.co.uk`,
		`user.name+tag+sorting@example.com`,
		`x@example.com`,
		// `"much.more unusual"@example.com`, // TODO
		// `"very.unusual.@.unusual.com"@example.com`, // TODO
		`admin@mailserver1`,
		`example@s.example`,
		// `user@[192.168.1.1]`, // TODO
	}
	for _, inputEmail := range validEmails {
		email, err := unmarshalFromString[Email](inputEmail)
		assert.NoError(t, err, `cannot unmarshall "%s"`, inputEmail)
		if err == nil {
			assert.Equal(t, inputEmail, string(email))
		}
	}
}

func TestEmailInvalid(t *testing.T) {
	validEmails := []string{
		`plainaddress`,
		`@missing-local.org`,
		`username@`,
		`username@.com`,
		`username@.com.`,
		// `.username@yahoo.com`, // TODO fix it
		`username@yahoo.com.`,
		`username@yahoo..com`,
		// `username@yahoo.c`, // TODO fix short TLD
		`username@yahoo,com`,
		`username@-example.com`,
		`username@example..com`,
	}
	for _, inputEmail := range validEmails {
		_, err := unmarshalFromString[Email](inputEmail)
		assert.Error(t, err, `unmarshalling "%s" is successful`, inputEmail)
	}
}
