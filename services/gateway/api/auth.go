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
	var request AuthenticateRequestSchema
	if err := json.Unmarshal(b, &request); err != nil {
		return err
	}

	if !request.Login.HasValue && !request.PhoneNumber.HasValue && !request.Email.HasValue {
		return ErrorNoUserId{}
	}

	// checking that exactly one user id has been passed
	if request.Login.HasValue && (request.Email.HasValue || request.PhoneNumber.HasValue) {
		return ErrorTooMuchUserId{}
	}
	if request.Email.HasValue && request.PhoneNumber.HasValue {
		return ErrorTooMuchUserId{}
	}

	r.AuthenticateRequestSchema = request

	return nil
}

func (r *CreateApiTokenRequest) UnmarshalJSON(b []byte) error {
	var request CreateApiTokenRequestSchema
	if err := json.Unmarshal(b, &request); err != nil {
		return err
	}

	if request.Ttl.Nanoseconds() <= 0 {
		return errors.New("ttl must be > 0")
	}

	r.CreateApiTokenRequestSchema = request
	return nil
}
