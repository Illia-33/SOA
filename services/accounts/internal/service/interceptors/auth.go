package interceptors

import (
	"context"
	"soa-socialnetwork/services/accounts/internal/service"
	"soa-socialnetwork/services/accounts/pkg/soajwt"
	"soa-socialnetwork/services/accounts/pkg/soatoken"
	"strings"

	pb "soa-socialnetwork/services/accounts/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func needAuth(methodName string) bool {
	return methodName == pb.AccountsService_EditProfile_FullMethodName || methodName == pb.AccountsService_UnregisterUser_FullMethodName
}

type parsedToken struct {
	kind  string
	value string
}

func parseTokenFromMetadata(md metadata.MD) (parsedToken, error) {
	auth, ok := md["authorization"]
	if !ok {
		return parsedToken{}, ErrorNoAuth{}
	}

	split := strings.Split(auth[0], " ")
	if len(split) != 2 {
		return parsedToken{}, ErrorInvalidToken{}
	}

	authKind, authToken := split[0], split[1]
	if !(authKind == "Bearer" || authKind == "SoaToken") {
		return parsedToken{}, ErrorUnknownAuthKind{}
	}

	return parsedToken{
		kind:  authKind,
		value: authToken,
	}, nil
}

func fetchProfileIdFromRequest(request any) string {
	switch v := request.(type) {
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
}

func verifyJwtToken(tokenStr string, verifier soajwt.Verifier, requestProfileId string) error {
	token, err := verifier.Verify(tokenStr)
	if err != nil {
		return err
	}

	if token.Subject != requestProfileId {
		return ErrorAccessDenied{}
	}

	return nil
}

func verifySoaToken(tokenStr string, s *service.AccountsService, requestProfileId string, req soatoken.RightsRequirements) error {
	token, err := soatoken.Parse(tokenStr)
	if err != nil {
		return err
	}

	if token.ProfileId.String() != requestProfileId {
		return ErrorAccessDenied{}
	}

	response, err := s.ValidateApiToken(context.Background(), &pb.ApiToken{
		Token: tokenStr,
	})

	if err != nil {
		return err
	}

	if response.GetInvalid() != nil {
		return ErrorInvalidToken{}
	}

	valid := response.GetValid()
	if valid == nil {
		panic("shouldn't reach here")
	}

	if req.Read && !valid.ReadAccess {
		return ErrorNoReadAccess{}
	}

	if req.Write && !valid.WriteAccess {
		return ErrorNoWriteAccess{}
	}

	return nil
}

func getSoaTokenRightsRequirements(methodName string) soatoken.RightsRequirements {
	switch methodName {
	case pb.AccountsService_EditProfile_FullMethodName, pb.AccountsService_UnregisterUser_FullMethodName:
		return soatoken.RightsRequirements{
			Read:  false,
			Write: true,
		}

	default:
		return soatoken.RightsRequirements{
			Read:  false,
			Write: false,
		}
	}
}

func Auth(s *service.AccountsService) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if !needAuth(info.FullMethod) {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, ErrorNoMetadata{}
		}

		parsedToken, err := parseTokenFromMetadata(md)
		if err != nil {
			return nil, err
		}

		profileId := fetchProfileIdFromRequest(req)

		if parsedToken.kind == "Bearer" {
			err := verifyJwtToken(parsedToken.value, s.JwtVerifier, profileId)
			if err != nil {
				return nil, err
			}
			return handler(ctx, req)
		}

		if parsedToken.kind == "SoaToken" {
			err := verifySoaToken(parsedToken.value, s, profileId, getSoaTokenRightsRequirements(info.FullMethod))
			if err != nil {
				return nil, err
			}
			return handler(ctx, req)
		}

		panic("shouldn't reach here")
	}
}

type ErrorUnknownAuthKind struct{}
type ErrorAccessDenied struct{}
type ErrorNoAuth struct{}
type ErrorInvalidToken struct{}
type ErrorNoWriteAccess struct{}
type ErrorNoReadAccess struct{}
type ErrorNoMetadata struct{}

func (ErrorUnknownAuthKind) Error() string {
	return "unknown auth kind"
}

func (ErrorAccessDenied) Error() string {
	return "access denied"
}

func (ErrorNoAuth) Error() string {
	return "auth required"
}

func (ErrorInvalidToken) Error() string {
	return "token is invalid"
}

func (ErrorNoWriteAccess) Error() string {
	return "no write access"
}

func (ErrorNoReadAccess) Error() string {
	return "no read access"
}

func (ErrorNoMetadata) Error() string {
	return "metadata not found"
}
