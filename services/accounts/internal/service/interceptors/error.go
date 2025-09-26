package interceptors

import (
	"context"
	"soa-socialnetwork/services/accounts/internal/storage/postgres"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func convertToGrpcCode(err error) codes.Code {
	switch err.(type) {
	case postgres.ErrorTokenNotFound, postgres.ErrorAccountNotFound, postgres.ErrorProfileNotFound, postgres.ErrorUserIdNotFound:
		return codes.NotFound

	case postgres.ErrorPasswordsDoNotMatch:
		return codes.PermissionDenied

	case ErrorAccessDenied, ErrorInvalidToken, ErrorNoAuth, ErrorNoReadAccess, ErrorNoWriteAccess:
		return codes.PermissionDenied

	case ErrorUnknownAuthKind:
		return codes.InvalidArgument

	case ErrorNoMetadata:
		return codes.Internal

	default:
		return codes.Internal
	}
}

func convertToGrpcError(err error) error {
	return status.Error(convertToGrpcCode(err), err.Error())
}

func ConvertErrors() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		return nil, convertToGrpcError(err)
	}
}
