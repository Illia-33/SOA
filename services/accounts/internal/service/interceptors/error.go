package interceptors

import (
	"context"
	serviceErrs "soa-socialnetwork/services/accounts/internal/service/errs"
	"soa-socialnetwork/services/accounts/internal/service/interceptors/errs"
	pgErrs "soa-socialnetwork/services/accounts/internal/storage/postgres/errs"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func convertToGrpcCode(err error) (code codes.Code, knownError bool) {
	switch err.(type) {
	case errs.InvalidToken, errs.NoAuth:
		return codes.PermissionDenied, true

	case errs.UnknownAuthKind:
		return codes.InvalidArgument, true

	case errs.NoMetadata:
		return codes.Internal, true

	case pgErrs.TokenNotFound, pgErrs.AccountNotFound, pgErrs.ProfileNotFound, pgErrs.UserIdNotFound:
		return codes.NotFound, true

	case pgErrs.PasswordsDoNotMatch:
		return codes.PermissionDenied, true

	case serviceErrs.NoReadAccess, serviceErrs.NoWriteAccess, serviceErrs.TokenExpired, serviceErrs.AccessDenied:
		return codes.PermissionDenied, true

	default:
		return codes.Internal, false
	}
}

func convertToGrpcError(err error) error {
	code, knownError := convertToGrpcCode(err)
	if !knownError {
		return err
	}

	return status.Error(code, err.Error())
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
