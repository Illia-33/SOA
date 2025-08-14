package server

import (
	"soa-socialnetwork/services/gateway/api"
)

type gatewayServiceContext struct {
}

func createContext() gatewayServiceContext {
	return gatewayServiceContext{}
}

func (c *gatewayServiceContext) RegisterProfile(req *api.RegisterProfileRequest) (api.GetProfileResponse, httpError) {
	return api.GetProfileResponse{}, httpOK() // TODO
}

func (c *gatewayServiceContext) GetProfileInfo(profileId string) (api.GetProfileResponse, httpError) {
	return api.GetProfileResponse{}, httpOK() // TODO
}

func (c *gatewayServiceContext) EditProfileInfo(profileId string, req *api.EditProfileRequest) httpError {
	return httpOK() // TODO
}

func (c *gatewayServiceContext) DeleteProfile(profileId string) httpError {
	return httpOK() // TODO
}
