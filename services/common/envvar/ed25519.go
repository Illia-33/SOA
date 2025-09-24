package envvar

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
)

func TryEd25519PubKeyFromEnv(key string) (ed25519.PublicKey, error) {
	pubkeyStr, err := TryStringFromEnv(key)
	if err != nil {
		return ed25519.PublicKey{}, err
	}

	pubkey := make(ed25519.PublicKey, ed25519.PublicKeySize)

	_, err = hex.Decode([]byte(pubkey), []byte(pubkeyStr))
	if err != nil {
		return nil, fmt.Errorf("hex decoding error while parsing ed25519 public key from %s (%s): %v", key, pubkeyStr, err)
	}

	return pubkey, nil

}

func MustEd25519PubKeyFromEnv(key string) ed25519.PublicKey {
	val, err := TryEd25519PubKeyFromEnv(key)
	if err != nil {
		panic(err)
	}
	return val
}

func TryEd25519PrivKeyFromEnv(key string) (ed25519.PrivateKey, error) {
	envPrivateKey, err := TryStringFromEnv(key)
	if err != nil {
		return ed25519.PrivateKey{}, err
	}

	seed := make([]byte, ed25519.SeedSize)

	_, err = hex.Decode(seed, []byte(envPrivateKey))
	if err != nil {
		return ed25519.PrivateKey{}, errors.New("cannot decode hex encoded jwt private key seed")
	}

	return ed25519.NewKeyFromSeed(seed), nil
}

func MustEd25519PrivKeyFromEnv(key string) ed25519.PrivateKey {
	val, err := TryEd25519PrivKeyFromEnv(key)
	if err != nil {
		panic(err)
	}

	return val
}
