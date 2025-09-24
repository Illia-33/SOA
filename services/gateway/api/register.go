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
	ProfileId string `json:"profile_id"`
}
