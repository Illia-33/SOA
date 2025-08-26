package httperr

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Err struct {
	StatusCode int
	Err        error
}

func New(statusCode int, err error) Err {
	return Err{
		StatusCode: statusCode,
		Err:        err,
	}
}

var grpcToHttp = map[codes.Code]int{
	codes.OK:                 http.StatusOK,
	codes.Canceled:           http.StatusInternalServerError,
	codes.Unknown:            http.StatusInternalServerError,
	codes.InvalidArgument:    http.StatusBadRequest,
	codes.DeadlineExceeded:   http.StatusGatewayTimeout,
	codes.NotFound:           http.StatusNotFound,
	codes.AlreadyExists:      http.StatusConflict,
	codes.PermissionDenied:   http.StatusForbidden,
	codes.ResourceExhausted:  http.StatusTooManyRequests,
	codes.FailedPrecondition: http.StatusBadRequest,
	codes.Aborted:            http.StatusConflict,
	codes.OutOfRange:         http.StatusBadRequest,
	codes.Unimplemented:      http.StatusNotImplemented,
	codes.Internal:           http.StatusInternalServerError,
	codes.Unavailable:        http.StatusServiceUnavailable,
	codes.DataLoss:           http.StatusInternalServerError,
	codes.Unauthenticated:    http.StatusUnauthorized,
}

func FromGrpcError(err error) Err {
	s := status.Convert(err)
	if s != nil {
		httpCode, ok := grpcToHttp[s.Code()]
		if !ok {
			httpCode = http.StatusInternalServerError
		}

		return Err{
			StatusCode: httpCode,
			Err:        err,
		}
	}

	return Err{
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}

func Ok() Err {
	return Err{
		StatusCode: http.StatusOK,
		Err:        nil,
	}
}

func (e *Err) IsOk() bool {
	return e.StatusCode == http.StatusOK && e.Err == nil
}
