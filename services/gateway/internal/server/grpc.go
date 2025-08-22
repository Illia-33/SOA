package server

import (
	"context"
	"fmt"
	"soa-socialnetwork/services/gateway/internal/server/query"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func createGrpcClient(target string, qp *query.Params) (*grpc.ClientConn, error) {
	return grpc.NewClient(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(newCreds(qp)),
	)
}

type tokenCredentials struct {
	token string
	kind  query.AuthTokenKind
}

func newCreds(qp *query.Params) tokenCredentials {
	if qp.AuthKind == query.AUTH_TOKEN_JWT {
		return jwtCreds(qp.AuthToken)
	}

	if qp.AuthKind == query.AUTH_TOKEN_SOA {
		return soaCreds(qp.AuthToken)
	}

	return tokenCredentials{}
}

func jwtCreds(rawToken string) tokenCredentials {
	return tokenCredentials{
		token: rawToken,
		kind:  query.AUTH_TOKEN_JWT,
	}
}

func soaCreds(rawToken string) tokenCredentials {
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
