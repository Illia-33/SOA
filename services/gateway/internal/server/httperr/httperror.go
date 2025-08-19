package httperr

import (
	"net/http"
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

func Ok() Err {
	return Err{
		StatusCode: http.StatusOK,
		Err:        nil,
	}
}

func (e *Err) IsOk() bool {
	return e.StatusCode == http.StatusOK && e.Err == nil
}
