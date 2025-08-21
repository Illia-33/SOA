package interceptors

import (
	"context"
	"errors"
	"soa-socialnetwork/internal/soajwt"
	"strings"

	pb "soa-socialnetwork/services/accounts/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func Auth(verifier *soajwt.Verifier) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if !(strings.HasSuffix(info.FullMethod, "EditProfile") || strings.HasSuffix(info.FullMethod, "UnregisterUser")) {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("metatada is not found")
		}

		auth := md["authorization"]
		if len(auth) == 0 || len(auth[0]) == 0 {
			return nil, errors.New("auth token is not found")
		}

		rawToken := auth[0]

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

		token, err := verifier.Verify(rawToken)
		if err != nil {
			return nil, err
		}

		if token.Subject != profileId {
			return nil, errors.New("access denied")
		}

		return handler(ctx, req)
	}
}
