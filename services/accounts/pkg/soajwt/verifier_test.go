package soajwt

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildValidJwt(token Token, priv ed25519.PrivateKey) string {
	algoJson := `{"alg":"EdDSA","typ":"JWT"}`
	algo := base64.RawURLEncoding.EncodeToString([]byte(algoJson))

	payloadBytes, err := json.Marshal(token)
	if err != nil {
		panic(err)
	}
	payload := base64.RawURLEncoding.EncodeToString(payloadBytes)

	verified := algo + "." + payload
	signRaw, err := priv.Sign(nil, []byte(verified), crypto.Hash(0))
	if err != nil {
		panic(err)
	}
	sign := base64.RawURLEncoding.EncodeToString(signRaw)
	return algo + "." + payload + "." + sign
}

func TestVerifierValid(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	now := time.Now()
	token := Token{
		Issuer:    "test-issuer",
		Subject:   "test-subject",
		Audience:  []string{"test-audience"},
		ExpiresAt: now.Add(time.Hour),
		NotBefore: now,
		IssuedAt:  now,
		JwtId:     "test-id",
		AccountId: 228,
	}

	jwt := buildValidJwt(token, priv)

	verifier := NewEd25519Verifier(pub)
	verifiedToken, err := verifier.Verify(jwt)
	require.NoError(t, err, `verification of valid token failed`)
	assert.Equal(t, token.Issuer, verifiedToken.Issuer)
	assert.Equal(t, token.Subject, verifiedToken.Subject)
	assert.Equal(t, token.Audience, verifiedToken.Audience)
	// assert.Equal(t, token.ExpiresAt, verifiedToken.ExpiresAt) // TODO fix time comparison
	// assert.Equal(t, token.NotBefore, verifiedToken.NotBefore)
	// assert.Equal(t, token.IssuedAt, verifiedToken.IssuedAt)
	assert.Equal(t, token.JwtId, verifiedToken.JwtId)
	assert.Equal(t, token.AccountId, verifiedToken.AccountId)
}

func TestVerifierTimeout(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	now := time.Now()
	token := Token{
		Issuer:    "test-issuer",
		Subject:   "test-subject",
		Audience:  []string{"test-audience"},
		ExpiresAt: now.Add(5 * time.Second),
		NotBefore: now,
		IssuedAt:  now,
		JwtId:     "test-id",
		AccountId: 229,
	}

	jwt := buildValidJwt(token, priv)
	verifier := NewEd25519Verifier(pub)

	time.Sleep(6 * time.Second)

	_, err = verifier.Verify(jwt)
	require.Error(t, err, `token is expired, but verifiable`)
}

func TestVerifierBrokenPayload(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	now := time.Now()
	token := Token{
		Issuer:    "test-issuer",
		Subject:   "test-subject",
		Audience:  []string{"test-audience"},
		ExpiresAt: now.Add(time.Hour),
		NotBefore: now,
		IssuedAt:  now,
		JwtId:     "test-id",
		AccountId: 2210,
	}

	jwt := func() string {
		valid := buildValidJwt(token, priv)
		validBytes := []byte(valid)
		rand.Reader.Read(validBytes[:4])
		return string(validBytes)
	}()

	verifier := NewEd25519Verifier(pub)
	_, err = verifier.Verify(jwt)
	require.Error(t, err, `broken jwt token is verified`)
}

func TestVerifierBrokenSign(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	now := time.Now()
	token := Token{
		Issuer:    "test-issuer",
		Subject:   "test-subject",
		Audience:  []string{"test-audience"},
		ExpiresAt: now.Add(time.Hour),
		NotBefore: now,
		IssuedAt:  now,
		JwtId:     "test-id",
		AccountId: 2211,
	}

	jwt := func() string {
		valid := buildValidJwt(token, priv)
		validBytes := []byte(valid)
		rand.Reader.Read(validBytes[len(validBytes)-4:])
		return string(validBytes)
	}()

	verifier := NewEd25519Verifier(pub)
	_, err = verifier.Verify(jwt)
	require.Error(t, err, `broken jwt token is verified`)
}
