package soatoken

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"

	"github.com/google/uuid"
)

type Payload struct {
	AccountId int32
	ProfileId uuid.UUID
}

func NewSoaToken(payload Payload) string {
	const RAW_TOKEN_LENGTH = 48
	var rawToken [RAW_TOKEN_LENGTH]byte

	n := copy(rawToken[:16], payload.ProfileId[:])
	if n != 16 {
		panic("must copy exact 16 bytes")
	}

	binary.LittleEndian.PutUint32(rawToken[16:20], uint32(payload.AccountId))

	_, err := rand.Read(rawToken[20:])
	if err != nil {
		panic("unexpected error from rand.Read")
	}

	return base64.RawStdEncoding.EncodeToString(rawToken[:])

}
