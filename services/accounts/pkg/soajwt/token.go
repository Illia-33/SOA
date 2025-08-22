package soajwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	Issuer    string    `json:"iss"`
	Subject   string    `json:"sub"`
	Audience  []string  `json:"aud"`
	ExpiresAt time.Time `json:"exp"`
	NotBefore time.Time `json:"nbf"`
	IssuedAt  time.Time `json:"iat"`
	JwtId     string    `json:"jti"`
	AccountId int       `json:"accid"`
}

func (t *Token) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(t.ExpiresAt), nil
}

func (t *Token) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(t.IssuedAt), nil
}

func (t *Token) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(t.NotBefore), nil
}

func (t *Token) GetIssuer() (string, error) {
	return t.Issuer, nil
}

func (t *Token) GetSubject() (string, error) {
	return t.Subject, nil
}

func (t *Token) GetAudience() (jwt.ClaimStrings, error) {
	return t.Audience, nil
}
