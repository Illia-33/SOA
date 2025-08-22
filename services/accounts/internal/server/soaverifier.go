package server

import (
	"context"
	"errors"
	"soa-socialnetwork/services/accounts/pkg/soatoken"
	pb "soa-socialnetwork/services/accounts/proto"
)

type serviceSoaTokenVerifier struct {
	service *AccountsService
}

func (v serviceSoaTokenVerifier) Verify(token string, req soatoken.RightsRequirements) error {
	response, err := v.service.ValidateApiToken(context.Background(), &pb.ApiToken{
		Token: token,
	})

	if err != nil {
		return err
	}

	if response.GetInvalid() != nil {
		return errors.New("access denied")
	}

	valid := response.GetValid()
	if valid == nil {
		panic("shouldn't reach here")
	}

	if req.Read && !valid.ReadAccess {
		return errors.New("need read access")
	}

	if req.Write && !valid.WriteAccess {
		return errors.New("need write access")
	}

	return nil
}
