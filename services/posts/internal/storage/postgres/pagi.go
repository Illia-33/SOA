package postgres

import (
	"encoding/base64"
	"encoding/json"

	"soa-socialnetwork/services/posts/internal/repos"
)

func decodePagiToken[DecodedToken any](encoded repos.PagiToken) (decoded DecodedToken, err error) {
	raw, err := base64.RawURLEncoding.DecodeString(string(encoded))
	if err != nil {
		return
	}
	err = json.Unmarshal(raw, &decoded)
	return
}

func encodePagiToken[EncodedToken any](token EncodedToken) (repos.PagiToken, error) {
	raw, err := json.Marshal(&token)
	if err != nil {
		return "", err
	}

	encoded := base64.RawURLEncoding.EncodeToString(raw)
	return repos.PagiToken(encoded), nil
}
