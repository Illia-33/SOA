package service

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"

	"github.com/google/uuid"
)

const SOA_API_TOKEN_LEN = 48

type accountData struct {
	accountId int
	profileId uuid.UUID
}

type soaApiToken [SOA_API_TOKEN_LEN]byte

func (t soaApiToken) toBase64() string {
	return base64.RawStdEncoding.EncodeToString(t[:])
}

func buildSoaApiToken(udata accountData) (token soaApiToken, err error) {
	n := copy(token[:16], udata.profileId[:])
	if n != 16 {
		panic("must copy exact 16 bytes")
	}

	binary.LittleEndian.PutUint32(token[16:20], uint32(udata.accountId))

	_, err = rand.Read(token[20:])
	if err != nil {
		return
	}

	return
}
