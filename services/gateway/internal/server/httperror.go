package server

import (
	"net/http"
)

type httpError struct {
	StatusCode int
	Err        error
}

func newHttpError(statusCode int, err error) httpError {
	return httpError{
		StatusCode: statusCode,
		Err:        err,
	}
}

func httpOK() httpError {
	return httpError{
		StatusCode: http.StatusOK,
		Err:        nil,
	}
}

func (e *httpError) IsOK() bool {
	return e.StatusCode == http.StatusOK && e.Err == nil
}
