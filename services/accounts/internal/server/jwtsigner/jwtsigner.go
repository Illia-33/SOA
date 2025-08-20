package jwtsigner

import (
	"crypto/ed25519"
	"encoding/json"
	"time"

	"soa-socialnetwork/internal/soajwt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Signer struct {
	privateKey ed25519.PrivateKey
}

type PersonalData struct {
	AccountId int
	ProfileId string
}

func New(privateKey ed25519.PrivateKey) (Signer, error) {
	return Signer{
		privateKey: privateKey,
	}, nil
}

func (j *Signer) generateClaims(data PersonalData, ttl time.Duration) (jwt.Claims, error) {
	jwtUuid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	token := soajwt.Token{
		ExpiresAt: now.Add(ttl),
		IssuedAt:  now,
		JwtId:     jwtUuid.String(),
		AccountId: data.AccountId,
		ProfileId: data.ProfileId,
	}

	// TODO optimize soajwt-to-map conversion

	asJson, err := json.Marshal(token)
	if err != nil {
		return nil, err
	}

	var claims map[string]any
	err = json.Unmarshal(asJson, &claims)
	if err != nil {
		return nil, err
	}

	return jwt.MapClaims(claims), nil
}

func (j *Signer) Sign(data PersonalData, ttl time.Duration) (string, error) {
	claims, err := j.generateClaims(data, ttl)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	tokenStr, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
