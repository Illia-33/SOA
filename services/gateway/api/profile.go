package api

import "soa-socialnetwork/services/gateway/pkg/types"

type GetProfileResponse struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Birthday string `json:"birthday"`
	Bio      string `json:"bio"`
}

type EditProfileRequest struct {
	Name     types.Optional[types.Name]    `json:"name"`
	Surname  types.Optional[types.Surname] `json:"surname"`
	Birthday types.Optional[types.Date]    `json:"birthday"`
	Bio      types.Optional[types.Bio]     `json:"bio"`
}
