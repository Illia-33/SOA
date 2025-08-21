package jsonextractor

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"soa-socialnetwork/services/gateway/api"
	"soa-socialnetwork/services/gateway/internal/server/birthday"
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
	case *api.RegisterProfileRequest, *api.EditProfileRequest, *api.AuthenticateRequest:
		return ctx.BindJSON(v)

	case *EmptyRequest:
		return nil

	default:
		return errors.New("unsupported request type")
	}
}

func (j *JsonExtractor) validateRequest(r any) error {
	switch v := r.(type) {
	case *api.RegisterProfileRequest:
		return j.validateRegisterProfileRequest(v)

	case *api.EditProfileRequest:
		return j.validateEditProfileRequest(v)

	case *api.AuthenticateRequest:
		return j.validateAuthenticateRequest(v)

	case *EmptyRequest:
		return nil

	default:
		panic("shouldn't reach here")
	}
}

func (j *JsonExtractor) validateRegisterProfileRequest(req *api.RegisterProfileRequest) error {
	if !(1 <= len(req.Login) && len(req.Login) <= 32) {
		return errors.New("login: length must be in [1; 32]")
	}

	if !(6 <= len(req.Password) && len(req.Password) <= 32) {
		return errors.New("password: length must be in [6; 32]")
	}

	if !j.validateEmail(req.Email) {
		return errors.New("email: invalid")
	}

	if !j.validatePhoneNumber(req.PhoneNumber) {
		return errors.New("phone number: must be in format +0123456789")
	}

	if !(1 <= len(req.Name) && len(req.Name) <= 32) {
		return errors.New("name: length must be in [1; 32]")
	}

	if !(1 <= len(req.Surname) && len(req.Surname) <= 32) {
		return errors.New("surname: length must be in [1; 32]")
	}

	return nil
}

func (j *JsonExtractor) validateEditProfileRequest(req *api.EditProfileRequest) error {
	if !(0 <= len(req.Name) && len(req.Name) <= 32) {
		return errors.New("name: length must be in [1; 32]")
	}

	if !(0 <= len(req.Surname) && len(req.Surname) <= 32) {
		return errors.New("surname: length must be in [1; 32]")
	}

	if len(req.Birthday) > 0 && !j.validateBirthday(req.Birthday) {
		return errors.New("birthday: must be in format YYYY-MM-DD")
	}

	if len(req.Bio) > 256 {
		return errors.New("bio: length must be <= 256")
	}

	if len(req.PhoneNumber) > 0 && !j.validatePhoneNumber(req.PhoneNumber) {
		return errors.New("phone number: must be in format +0123456789")
	}

	if len(req.Email) > 0 && !j.validateEmail(req.Email) {
		return errors.New("email: invalid")
	}

	return nil
}

func (j *JsonExtractor) validateAuthenticateRequest(req *api.AuthenticateRequest) error {
	if len(req.Login) == 0 && len(req.Email) == 0 && len(req.PhoneNumber) == 0 {
		return errors.New("no userid")
	}

	// checking that exactly one user id has been passed
	if len(req.Login) > 0 && (len(req.Email) > 0 || len(req.PhoneNumber) > 0) {
		return errors.New("too much userid")
	}
	if len(req.Email) > 0 && len(req.PhoneNumber) > 0 {
		return errors.New("too much userid")
	}

	// check validity of passed user id
	if !(0 <= len(req.Login) && len(req.Login) <= 32) {
		return errors.New("login: length must be in [1; 32]")
	}
	if len(req.PhoneNumber) > 0 && !j.validatePhoneNumber(req.PhoneNumber) {
		return errors.New("phone number: must be in format +0123456789")
	}
	if len(req.Email) > 0 && !j.validateEmail(req.Email) {
		return errors.New("email: invalid")
	}

	if !(6 <= len(req.Password) && len(req.Password) <= 32) {
		return errors.New("password: length must be in [6; 32]")
	}

	return nil
}

func (j *JsonExtractor) validateEmail(s string) bool {
	emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegexp.MatchString(s)
}

func (j *JsonExtractor) validatePhoneNumber(s string) bool {
	phoneNumberRegexp := regexp.MustCompile(`^\+\d{7,15}$`)
	return phoneNumberRegexp.MatchString(s)
}

func (j *JsonExtractor) validateBirthday(s string) bool {
	b, err := birthday.Parse(s)
	if err != nil {
		return false
	}

	return b.IsValid()
}
