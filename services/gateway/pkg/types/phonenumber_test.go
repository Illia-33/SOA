package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPhoneNumberValid(t *testing.T) {
	validPhoneNumbers := []string{
		"+14155552671",
		"+442071838750",
		"+4915123456789",
		"+8613800138000",
		"+918527001234",
		"+81312345678",
		"+33123456789",
		"+5511998765432",
		"+61298765432",
		"+380501234567",
		"+79991234567",
	}
	for _, inputPhoneNumber := range validPhoneNumbers {
		phoneNumber, err := unmarshalFromString[PhoneNumber](inputPhoneNumber)
		assert.NoError(t, err, `cannot unmarshall "%s": %v`, inputPhoneNumber, err)
		if err == nil {
			assert.Equal(t, inputPhoneNumber, string(phoneNumber))
		}
	}
}

func TestPhoneNumberInvalid(t *testing.T) {
	validPhoneNumbers := []string{
		"12345",
		"+1",
		// "+999123456789", // TODO fix non-existent country codes
		"0044123456789",
		"+44-20-7183-8750",
		"+44 2071 838 750",
		"phone123456",
		"+14155552671#123",
		"+14155552671x123",
		"+",
		"++14155552671",
		"+14155552671abc",
	}
	for _, inputPhoneNumber := range validPhoneNumbers {
		_, err := unmarshalFromString[PhoneNumber](inputPhoneNumber)
		assert.Error(t, err, `unmarshalling "%s" is successful`, inputPhoneNumber)
	}
}
