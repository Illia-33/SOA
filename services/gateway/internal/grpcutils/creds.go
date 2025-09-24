package grpcutils

import (
	"context"
	"fmt"
	"soa-socialnetwork/services/gateway/internal/query"
)

type tokenCredentials struct {
	token string
	kind  query.AuthTokenKind
}

func NewCreds(qp *query.Params) tokenCredentials {
	if qp.AuthKind == query.AUTH_TOKEN_JWT {
		return JwtCreds(qp.AuthToken)
	}

	if qp.AuthKind == query.AUTH_TOKEN_SOA {
		return SoaCreds(qp.AuthToken)
	}

	return tokenCredentials{}
}

func JwtCreds(rawToken string) tokenCredentials {
	return tokenCredentials{
		token: rawToken,
		kind:  query.AUTH_TOKEN_JWT,
	}
}

func SoaCreds(rawToken string) tokenCredentials {
	return tokenCredentials{
		token: rawToken,
		kind:  query.AUTH_TOKEN_SOA,
	}
}

func (j tokenCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	switch j.kind {
	case query.AUTH_TOKEN_JWT:
		{
			return map[string]string{
				"authorization": fmt.Sprintf("Bearer %s", j.token),
			}, nil
		}

	case query.AUTH_TOKEN_SOA:
		{
			return map[string]string{
				"authorization": fmt.Sprintf("SoaToken %s", j.token),
			}, nil
		}

	default:
		return nil, nil
	}
}

func (j tokenCredentials) RequireTransportSecurity() bool {
	return false
}
