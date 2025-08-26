package server

import (
	"context"
	"soa-socialnetwork/services/accounts/pkg/soatoken"
	pb "soa-socialnetwork/services/accounts/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return status.Error(codes.PermissionDenied, "token invalid")
	}

	valid := response.GetValid()
	if valid == nil {
		panic("shouldn't reach here")
	}

	if req.Read && !valid.ReadAccess {
		return status.Error(codes.PermissionDenied, "need read access")
	}

	if req.Write && !valid.WriteAccess {
		return status.Error(codes.PermissionDenied, "need write access")
	}

	return nil
}
