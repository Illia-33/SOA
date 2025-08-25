package api

import "soa-socialnetwork/services/gateway/pkg/types"

type RegisterProfileRequest struct {
	Login       types.Login       `json:"login"`
	Password    types.Password    `json:"password"`
	Email       types.Email       `json:"email"`
	PhoneNumber types.PhoneNumber `json:"phone_number"`
	Name        types.Name        `json:"name"`
	Surname     types.Surname     `json:"surname"`
}

type RegisterProfileResponse struct {
	ProfileID string `json:"profile_id"`
}

type GetProfileResponse struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Birthday string `json:"birthday"`
	Bio      string `json:"bio"`
}

type EditProfileRequest struct {
	Name        types.Optional[types.Name]        `json:"name"`
	Surname     types.Optional[types.Surname]     `json:"surname"`
	Birthday    types.Optional[types.Date]        `json:"birthday"`
	Bio         types.Optional[types.Bio]         `json:"bio"`
	PhoneNumber types.Optional[types.PhoneNumber] `json:"phone_number"`
	Email       types.Optional[types.Email]       `json:"email"`
}

type AuthenticateRequest struct {
	Login       types.Optional[types.Login]       `json:"login"`
	Email       types.Optional[types.Email]       `json:"email"`
	PhoneNumber types.Optional[types.PhoneNumber] `json:"phone_number"`
	Password    types.Password                    `json:"password"`
}

type AuthenticateResponse struct {
	Token string `json:"token"`
}

type CreateApiTokenRequest struct {
	Auth        AuthenticateRequest `json:"auth"`
	ReadAccess  bool                `json:"read_access"`
	WriteAccess bool                `json:"write_access"`
	Ttl         types.Duration      `json:"ttl"`
}

type CreateApiTokenResponse struct {
	Token string `json:"token"`
}
