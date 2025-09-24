package api

import "soa-socialnetwork/services/gateway/pkg/types"

// Exactly one of {login,email,phone_number} must have value
type AuthenticateRequest struct {
	Login       types.Optional[types.Login]       `json:"login"`
	Email       types.Optional[types.Email]       `json:"email"`
	PhoneNumber types.Optional[types.PhoneNumber] `json:"phone_number"`
	Password    types.Password                    `json:"password"`
}

type AuthenticateResponse struct {
	Token string `json:"token"`
}

// Ttl (time to live) must be positive
type CreateApiTokenRequest struct {
	Auth        AuthenticateRequest `json:"auth"`
	ReadAccess  bool                `json:"read_access"`
	WriteAccess bool                `json:"write_access"`
	Ttl         types.Duration      `json:"ttl"`
}

type CreateApiTokenResponse struct {
	Token string `json:"token"`
}
