package api

import (
	"encoding/json"
	"errors"
	"soa-socialnetwork/services/gateway/pkg/types"
)

// Exactly one of {login,email,phone_number} must have value
type AuthenticateRequestSchema struct {
	Login       types.Optional[types.Login]       `json:"login"`
	Email       types.Optional[types.Email]       `json:"email"`
	PhoneNumber types.Optional[types.PhoneNumber] `json:"phone_number"`
	Password    types.Password                    `json:"password"`
}

type AuthenticateRequest struct {
	AuthenticateRequestSchema
}

type AuthenticateResponse struct {
	Token string `json:"token"`
}

// Ttl (time to live) must be positive
type CreateApiTokenRequestSchema struct {
	Auth        AuthenticateRequest `json:"auth"`
	ReadAccess  bool                `json:"read_access"`
	WriteAccess bool                `json:"write_access"`
	Ttl         types.Duration      `json:"ttl"`
}

type CreateApiTokenRequest struct {
	CreateApiTokenRequestSchema
}

type CreateApiTokenResponse struct {
	Token string `json:"token"`
}

func (r *AuthenticateRequest) UnmarshalJSON(b []byte) error {
	var schema AuthenticateRequestSchema
	if err := json.Unmarshal(b, &schema); err != nil {
		return err
	}

	if !r.Login.HasValue && !r.PhoneNumber.HasValue && !r.Email.HasValue {
		return errors.New("no userid")
	}

	// checking that exactly one user id has been passed
	if r.Login.HasValue && (r.Email.HasValue || r.PhoneNumber.HasValue) {
		return errors.New("too much userid")
	}
	if r.Email.HasValue && r.PhoneNumber.HasValue {
		return errors.New("too much userid")
	}

	r.AuthenticateRequestSchema = schema

	return nil
}

func (r *CreateApiTokenRequest) UnmarshalJSON(b []byte) error {
	var schema CreateApiTokenRequestSchema
	if err := json.Unmarshal(b, &schema); err != nil {
		return err
	}

	if r.Ttl.Nanoseconds() <= 0 {
		return errors.New("ttl must be > 0")
	}

	r.CreateApiTokenRequestSchema = schema
	return nil
}
