package soajwt

import (
	"crypto/ed25519"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type Verifier interface {
	Verify(jwt string) (Token, error)
}

type ed25519Verifier struct {
	publicKey ed25519.PublicKey
}

func NewEd25519Verifier(pubkey ed25519.PublicKey) Verifier {
	return &ed25519Verifier{
		publicKey: pubkey,
	}
}

func (v *ed25519Verifier) Verify(tokenString string) (Token, error) {
	var token Token
	parsedToken, err := jwt.ParseWithClaims(tokenString, &token, func(t *jwt.Token) (any, error) {
		return v.publicKey, nil
	})
	if err != nil {
		return Token{}, err
	}

	if !parsedToken.Valid {
		return Token{}, errors.New("invalid token")
	}

	return token, nil
}
