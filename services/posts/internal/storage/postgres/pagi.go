package postgres

import (
	"encoding/base64"
	"encoding/json"

	"soa-socialnetwork/services/posts/internal/repo"
)

func decodePagiToken[DecodedToken any](encoded repo.PagiToken) (decoded DecodedToken, err error) {
	raw, err := base64.RawURLEncoding.DecodeString(string(encoded))
	if err != nil {
		return
	}
	err = json.Unmarshal(raw, &decoded)
	return
}

func encodePagiToken[EncodedToken any](token EncodedToken) (repo.PagiToken, error) {
	raw, err := json.Marshal(&token)
	if err != nil {
		return "", err
	}

	encoded := base64.RawURLEncoding.EncodeToString(raw)
	return repo.PagiToken(encoded), nil
}
