package soajwtissuer

import (
	"crypto/ed25519"
	"time"

	"soa-socialnetwork/services/accounts/pkg/soajwt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Issuer struct {
	privateKey ed25519.PrivateKey
}

type PersonalData struct {
	AccountId int
	ProfileId string
}

func New(privateKey ed25519.PrivateKey) Issuer {
	return Issuer{
		privateKey: privateKey,
	}
}

func (j *Issuer) Issue(data PersonalData, ttl time.Duration) (string, error) {
	jwtUuid, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	now := time.Now()

	token := soajwt.Token{
		Issuer:    "accounts-service",
		Subject:   data.ProfileId,
		Audience:  []string{"user"},
		ExpiresAt: now.Add(ttl),
		NotBefore: now,
		IssuedAt:  now,
		JwtId:     jwtUuid.String(),
		AccountId: data.AccountId,
	}

	tokenStr, err := jwt.NewWithClaims(jwt.SigningMethodEdDSA, &token).SignedString(j.privateKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
