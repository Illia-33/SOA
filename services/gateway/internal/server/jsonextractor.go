package server

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"soa-socialnetwork/services/gateway/api"
)

type emptyRequest struct{}

type jsonExtractor struct {
}

func (j *jsonExtractor) extract(r any, ctx httpContext) httpError {
	err := j.bindJSON(r, ctx)
	if err != nil {
		return newHttpError(http.StatusBadRequest, fmt.Errorf("cannot bind json: %v", err))
	}

	err = j.validateRequest(r)
	if err != nil {
		return newHttpError(http.StatusBadRequest, fmt.Errorf("bad request: %v", err))
	}

	return httpOK()
}

func (j *jsonExtractor) bindJSON(r any, ctx httpContext) error {
	switch v := r.(type) {
	case *api.RegisterProfileRequest, *api.EditProfileRequest:
		return ctx.BindJSON(v)

	case emptyRequest:
		return nil

	default:
		return errors.New("unsupported request type")
	}
}

func (j *jsonExtractor) validateRequest(r any) error {
	switch v := r.(type) {
	case *api.RegisterProfileRequest:
		return j.validateRegisterProfileRequest(v)

	case *api.EditProfileRequest:
		return j.validateEditProfileRequest(v)

	case emptyRequest:
		return nil

	default:
		panic("shouldn't reach here")
	}
}

func (j *jsonExtractor) validateRegisterProfileRequest(req *api.RegisterProfileRequest) error {
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

func (j *jsonExtractor) validateEditProfileRequest(req *api.EditProfileRequest) error {
	if !(1 <= len(req.Name) && len(req.Name) <= 32) {
		return errors.New("name: length must be in [1; 32]")
	}

	if !(1 <= len(req.Surname) && len(req.Surname) <= 32) {
		return errors.New("surname: length must be in [1; 32]")
	}

	if !j.validateBirthday(req.Birthday) {
		return errors.New("birthday: must be in format DD-MM-YYYY")
	}

	if len(req.Bio) > 256 {
		return errors.New("bio: length must be <= 256")
	}

	if !j.validatePhoneNumber(req.PhoneNumber) {
		return errors.New("phone number: must be in format +0123456789")
	}

	if !j.validateEmail(req.Email) {
		return errors.New("email: invalid")
	}

	return nil
}

func (j *jsonExtractor) validateEmail(s string) bool {
	emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegexp.MatchString(s)
}

func (j *jsonExtractor) validatePhoneNumber(s string) bool {
	phoneNumberRegexp := regexp.MustCompile(`^\+\d{7,15}$`)
	return phoneNumberRegexp.MatchString(s)
}

func (j *jsonExtractor) validateBirthday(s string) bool {
	b, err := parseBirthday(s)
	if err != nil {
		return false
	}

	return b.isValid()
}
