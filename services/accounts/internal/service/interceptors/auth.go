package interceptors

import (
	"context"
	"soa-socialnetwork/services/accounts/internal/service/interceptors/errs"
	"soa-socialnetwork/services/accounts/pkg/soajwt"
	"soa-socialnetwork/services/accounts/pkg/soatoken"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type TokenVerifiers struct {
	Jwt soajwt.Verifier
	Soa soatoken.Verifier
}

func Auth(verifiers TokenVerifiers) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		authReqs := getAuthRequirements(info.FullMethod)
		if !authReqs.needAuth {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errs.NoMetadata{}
		}

		parsedToken, err := parseTokenFromMetadata(md)
		if err != nil {
			return nil, err
		}

		authInfo, err := func() (AuthInfo, error) {
			switch parsedToken.kind {
			case auth_kind_jwt:
				return verifyJwtToken(parsedToken.value, verifiers.Jwt)

			case auth_kind_soa:
				return verifySoaToken(parsedToken.value, verifiers.Soa, authReqs)

			default:
				panic("unknown auth token kind")
			}
		}()

		if err != nil {
			return nil, err
		}

		return handler(context.WithValue(ctx, AuthInfoKey, authInfo), req)
	}
}

type AuthInfo struct {
	ProfileId string
	AccountId int32
}

type AuthInfoKeyType struct{}

var AuthInfoKey AuthInfoKeyType

type authKind int

const (
	auth_kind_jwt = iota
	auth_kind_soa
)

type parsedToken struct {
	kind  authKind
	value string
}

func parseTokenFromMetadata(md metadata.MD) (parsedToken, error) {
	auth, ok := md["authorization"]
	if !ok {
		return parsedToken{}, errs.NoAuth{}
	}

	split := strings.Split(auth[0], " ")
	if len(split) != 2 {
		return parsedToken{}, errs.InvalidToken{}
	}

	authKind, authToken := split[0], split[1]

	switch authKind {
	case "Bearer":
		return parsedToken{
			kind:  auth_kind_jwt,
			value: authToken,
		}, nil

	case "SoaToken":
		return parsedToken{
			kind:  auth_kind_soa,
			value: authToken,
		}, nil

	default:
		return parsedToken{}, errs.UnknownAuthKind{}
	}
}

func verifyJwtToken(token string, verifier soajwt.Verifier) (AuthInfo, error) {
	parsedToken, err := verifier.Verify(token)
	if err != nil {
		return AuthInfo{}, err
	}

	return AuthInfo{
		ProfileId: parsedToken.Subject,
		AccountId: int32(parsedToken.AccountId),
	}, nil
}

func verifySoaToken(token string, verifier soatoken.Verifier, reqs authRequirements) (AuthInfo, error) {
	parsedToken, err := soatoken.Parse(token)
	if err != nil {
		return AuthInfo{}, err
	}

	err = verifier.Verify(token, soatoken.RightsRequirements{
		Read:  reqs.needReadAccess,
		Write: reqs.needWriteAccess,
	})
	if err != nil {
		return AuthInfo{}, err
	}

	return AuthInfo{
		ProfileId: parsedToken.ProfileId.String(),
		AccountId: parsedToken.AccountId,
	}, nil
}
