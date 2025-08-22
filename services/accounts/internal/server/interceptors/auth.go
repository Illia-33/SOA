package interceptors

import (
	"context"
	"errors"
	"soa-socialnetwork/services/accounts/pkg/soajwt"
	"soa-socialnetwork/services/accounts/pkg/soatoken"
	"strings"

	pb "soa-socialnetwork/services/accounts/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// TODO refactor it

func Auth(jwtVerifier *soajwt.Verifier, soaVerifier soatoken.Verifier) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if !(info.FullMethod == pb.AccountsService_EditProfile_FullMethodName || info.FullMethod == pb.AccountsService_UnregisterUser_FullMethodName) {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("metatada is not found")
		}

		auth := md["authorization"]
		if len(auth) == 0 {
			return nil, errors.New("auth token is not found")
		}

		split := strings.Split(auth[0], " ")
		if len(split) != 2 {
			return nil, errors.New("bad token")
		}

		authKind, authToken := split[0], split[1]

		if !(authKind == "Bearer" || authKind == "SoaToken") {
			return nil, errors.New("unknown auth kind")
		}

		profileId := func() string {
			switch v := req.(type) {
			case *pb.EditProfileRequest:
				{
					return v.ProfileId
				}

			case *pb.UnregisterUserRequest:
				{
					return v.ProfileId
				}
			}

			panic("shouldn't reach here")
		}()

		if authKind == "Bearer" {
			token, err := jwtVerifier.Verify(authToken)
			if err != nil {
				return nil, err
			}

			if token.Subject != profileId {
				return nil, errors.New("access denied")
			}

			return handler(ctx, req)
		}

		if authKind == "SoaToken" {
			token, err := soatoken.Parse(authToken)
			if err != nil {
				return nil, err
			}

			if token.ProfileID.String() != profileId {
				return nil, errors.New("access denied")
			}

			err = soaVerifier.Verify(authToken, soatoken.RightsRequirements{
				Read:  false,
				Write: true,
			})

			if err != nil {
				return nil, err
			}

			return handler(ctx, req)
		}

		panic("shouldn't reach here")
	}
}
