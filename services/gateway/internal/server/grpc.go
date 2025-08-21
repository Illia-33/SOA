package server

import (
	"context"
	"soa-socialnetwork/services/gateway/internal/server/query"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func createGrpcClient(target string, qp *query.Params) (*grpc.ClientConn, error) {
	return grpc.NewClient(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(jwtCreds(qp.RawJwtToken)),
	)
}

type jwtCredentials struct {
	token string
}

func jwtCreds(rawToken string) jwtCredentials {
	return jwtCredentials{
		token: rawToken,
	}
}

func (j jwtCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	if len(j.token) == 0 {
		return nil, nil
	}

	return map[string]string{
		"authorization": j.token,
	}, nil
}

func (j jwtCredentials) RequireTransportSecurity() bool {
	return false
}
