package dbclient

import (
	"encoding/base64"
	"encoding/json"
	req "soa-socialnetwork/services/posts/internal/server/dbclient/requests"
	types "soa-socialnetwork/services/posts/internal/server/dbclient/types"
	"time"
)

func decodePagiToken[PagiToken any](encoded req.PagiToken) (decoded PagiToken, err error) {
	raw, err := base64.RawURLEncoding.DecodeString(string(encoded))
	if err != nil {
		return
	}
	err = json.Unmarshal(raw, &decoded)
	return
}

func encodePagiToken[PagiToken any](token PagiToken) (req.PagiToken, error) {
	raw, err := json.Marshal(&token)
	if err != nil {
		return "", err
	}

	encoded := base64.RawURLEncoding.EncodeToString(raw)
	return req.PagiToken(encoded), nil
}

type pgPostsPagiToken struct {
	LastCreatedAt time.Time `json:"lcr"`
}

func decodePgPostsPagiToken(token req.PagiToken) (pgPostsPagiToken, error) {
	if token == "" {
		return pgPostsPagiToken{
			LastCreatedAt: time.Date(9999, time.December, 31, 23, 59, 59, 0, time.UTC),
		}, nil
	}

	return decodePagiToken[pgPostsPagiToken](token)
}

func encodePgPostsPagiToken(token pgPostsPagiToken) (req.PagiToken, error) {
	return encodePagiToken(token)
}

type pgCommentsPagiToken struct {
	LastId types.CommentId `json:"lid"`
}

func decodePgCommentsPagiToken(token req.PagiToken) (pgCommentsPagiToken, error) {
	if token == "" {
		return pgCommentsPagiToken{}, nil
	}
	return decodePagiToken[pgCommentsPagiToken](token)
}

func encodePgCommentsPagiToken(token pgCommentsPagiToken) (req.PagiToken, error) {
	return encodePagiToken(token)
}
