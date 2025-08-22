package soatoken

import (
	"encoding/base64"
	"encoding/binary"
	"errors"

	"github.com/google/uuid"
)

type Token struct {
	ProfileID uuid.UUID
	AccountID int
	Salt      [28]byte
}

const TOKEN_LENGTH = 64

func Parse(tokenStr string) (Token, error) {
	if len(tokenStr) != TOKEN_LENGTH {
		return Token{}, errors.New("not a soa token")
	}

	var rawToken [48]byte
	_, err := base64.StdEncoding.Decode(rawToken[:], []byte(tokenStr))
	if err != nil {
		return Token{}, err
	}

	profileId, err := uuid.FromBytes(rawToken[:16])
	if err != nil {
		return Token{}, err
	}

	return Token{
		ProfileID: profileId,
		AccountID: int(binary.LittleEndian.Uint32(rawToken[16:20])),
		Salt:      [28]byte(rawToken[20:]),
	}, nil
}
