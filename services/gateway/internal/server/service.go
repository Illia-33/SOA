package server

import (
	"soa-socialnetwork/services/gateway/api"
)

type gatewayService struct {
}

func initService() gatewayService {
	return gatewayService{}
}

func (c *gatewayService) RegisterProfile(req *api.RegisterProfileRequest) (api.GetProfileResponse, httpError) {
	return api.GetProfileResponse{}, httpOK() // TODO
}

func (c *gatewayService) GetProfileInfo(profileId string) (api.GetProfileResponse, httpError) {
	return api.GetProfileResponse{}, httpOK() // TODO
}

func (c *gatewayService) EditProfileInfo(profileId string, req *api.EditProfileRequest) httpError {
	return httpOK() // TODO
}

func (c *gatewayService) DeleteProfile(profileId string) httpError {
	return httpOK() // TODO
}
