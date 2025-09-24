package envvar

import (
	"crypto/ed25519"
	"encoding/hex"
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
