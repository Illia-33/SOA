package jsonextractor

import (
	"errors"
	"fmt"
	"net/http"
	"soa-socialnetwork/services/gateway/api"
	"soa-socialnetwork/services/gateway/internal/server/httperr"

	"github.com/gin-gonic/gin"
)

type EmptyRequest struct{}

type JsonExtractor struct {
}

func New() JsonExtractor {
	return JsonExtractor{}
}

func (j *JsonExtractor) Extract(r any, ctx *gin.Context) httperr.Err {
	err := j.bindJSON(r, ctx)
	if err != nil {
		return httperr.New(http.StatusBadRequest, fmt.Errorf("cannot bind json: %v", err))
	}

	err = j.validateRequest(r)
	if err != nil {
		return httperr.New(http.StatusBadRequest, fmt.Errorf("bad request: %v", err))
	}

	return httperr.Ok()
}

func (j *JsonExtractor) bindJSON(r any, ctx *gin.Context) error {
	switch v := r.(type) {
	case *api.RegisterProfileRequest, *api.EditProfileRequest, *api.AuthenticateRequest, *api.CreateApiTokenRequest:
		return ctx.BindJSON(v)

	case *api.EditPageSettingsRequest, *api.NewPostRequest, *api.NewCommentRequest:
		return ctx.BindJSON(v)

	case *EmptyRequest:
		return nil

	default:
		return errors.New("unsupported request type")
	}
}

func (j *JsonExtractor) validateRequest(r any) error {
	switch v := r.(type) {
	case *api.AuthenticateRequest:
		return j.validateAuthenticateRequest(v)

	case *api.CreateApiTokenRequest:
		return j.validateCreateApiTokenRequest(v)
	}

	return nil
}

func (j *JsonExtractor) validateAuthenticateRequest(req *api.AuthenticateRequest) error {
	if !req.Login.HasValue && !req.PhoneNumber.HasValue && !req.Email.HasValue {
		return errors.New("no userid")
	}

	// checking that exactly one user id has been passed
	if req.Login.HasValue && (req.Email.HasValue || req.PhoneNumber.HasValue) {
		return errors.New("too much userid")
	}
	if req.Email.HasValue && req.PhoneNumber.HasValue {
		return errors.New("too much userid")
	}

	return nil
}

func (j *JsonExtractor) validateCreateApiTokenRequest(req *api.CreateApiTokenRequest) error {
	if err := j.validateAuthenticateRequest(&req.Auth); err != nil {
		return err
	}

	if req.Ttl.Nanoseconds() <= 0 {
		return errors.New("ttl must be > 0")
	}

	return nil
}
