package interceptors

import (
	"context"
	"soa-socialnetwork/services/accounts/pkg/soajwt"
	"soa-socialnetwork/services/accounts/pkg/soatoken"
	"soa-socialnetwork/services/posts/internal/models"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AccountIdCtxKey string

const AUTHOR_ACCOUNT_ID_CTX_KEY AccountIdCtxKey = "account_id"

func WithAuth(verifier soajwt.Verifier) grpc.UnaryServerInterceptor {
	validateToken := func(t authToken) validationInfo {
		switch t.kind {
		case auth_kind_unknown:
			{
				return validationInfo{}
			}

		case auth_kind_jwt:
			{
				token, err := verifier.Verify(string(t.value))
				if err != nil {
					return validationInfo{}
				}

				return validationInfo{
					valid:     true,
					accountId: models.AccountId(token.AccountId),
				}
			}

		case auth_kind_soa:
			{
				token, err := soatoken.Parse(string(t.value))
				if err != nil {
					return validationInfo{}
				}

				return validationInfo{
					valid:     true,
					accountId: models.AccountId(token.AccountId),
				}
			}
		}

		panic("unknown auth kind")
	}

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Internal, "cannot get metadata from context")
		}

		authToken := fetchTokenFromMetadata(md)
		tokenInfo := validateToken(authToken)

		if !tokenInfo.valid {
			return handler(ctx, req)
		}

		return handler(context.WithValue(ctx, AUTHOR_ACCOUNT_ID_CTX_KEY, tokenInfo.accountId), req)
	}
}

type authKind int

const (
	auth_kind_unknown authKind = iota
	auth_kind_jwt
	auth_kind_soa
)

func authKindFromString(s string) authKind {
	if s == "Bearer" {
		return auth_kind_jwt
	}

	if s == "SoaToken" {
		return auth_kind_soa
	}

	return auth_kind_unknown
}

type authTokenValue string

type authToken struct {
	kind  authKind
	value authTokenValue
}

type validationInfo struct {
	valid     bool
	accountId models.AccountId
}

func fetchTokenFromMetadata(md metadata.MD) authToken {
	auth, ok := md["authorization"]
	if !ok || len(auth) != 1 {
		return authToken{}
	}

	splitAuth := strings.Split(auth[0], " ")
	if len(splitAuth) != 2 {
		return authToken{}
	}

	kindRaw, val := splitAuth[0], splitAuth[1]
	kind := authKindFromString(kindRaw)

	if kind == auth_kind_unknown {
		return authToken{}
	}

	return authToken{
		kind:  kind,
		value: authTokenValue(val),
	}
}
