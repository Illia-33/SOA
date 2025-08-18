package server

import (
	"context"
	"errors"
	"net/http"
	pb "soa-socialnetwork/services/accounts/proto"
	"soa-socialnetwork/services/gateway/api"
	"soa-socialnetwork/services/gateway/internal/birthday"
	"soa-socialnetwork/services/gateway/internal/httperr"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type gatewayService struct {
}

func initService() gatewayService {
	return gatewayService{}
}

func (c *gatewayService) createAccountsServiceStub() (pb.AccountsServiceClient, error) {
	conn, err := grpc.NewClient("accounts-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return pb.NewAccountsServiceClient(conn), nil
}

func (c *gatewayService) RegisterProfile(req *api.RegisterProfileRequest) (response api.RegisterProfileResponse, httpErr httperr.Err) {
	stub, err := c.createAccountsServiceStub()
	if err != nil {
		httpErr = httperr.New(http.StatusInternalServerError, err)
		return
	}

	r, err := stub.RegisterUser(context.Background(), &pb.RegisterUserRequest{
		Login:       req.Login,
		Password:    req.Password,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
		Surname:     req.Surname,
	})
	if err != nil {
		httpErr = httperr.New(http.StatusInternalServerError, err)
		return
	}

	response.ProfileID = r.ProfileId
	httpErr = httperr.OK()
	return
}

func (c *gatewayService) GetProfileInfo(profileId string) (response api.GetProfileResponse, httpErr httperr.Err) {
	stub, err := c.createAccountsServiceStub()
	if err != nil {
		httpErr = httperr.New(http.StatusInternalServerError, err)
		return
	}

	r, err := stub.GetProfile(context.Background(), &pb.GetProfileRequest{
		ProfileId: profileId,
	})
	if err != nil {
		httpErr = httperr.New(http.StatusInternalServerError, err)
		return
	}

	response = api.GetProfileResponse{
		Name:     r.Name,
		Surname:  r.Surname,
		Birthday: r.Birthday.AsTime().Format("2006-01-02"),
		Bio:      r.Bio,
	}
	httpErr = httperr.OK()
	return
}

func (c *gatewayService) EditProfileInfo(profileId string, req *api.EditProfileRequest) httperr.Err {
	stub, err := c.createAccountsServiceStub()
	if err != nil {
		return httperr.New(http.StatusInternalServerError, err)
	}

	birthday, err := birthday.Parse(req.Birthday)
	if err != nil {
		return httperr.New(http.StatusInternalServerError, errors.New("bad birthday, check validity"))
	}

	_, err = stub.EditProfile(context.Background(), &pb.EditProfileRequest{
		ProfileId: profileId,
		EditedProfileData: &pb.Profile{
			Name:     req.Name,
			Surname:  req.Surname,
			Birthday: timestamppb.New(birthday.AsTime()),
			Bio:      req.Bio,
		},
	})
	if err != nil {
		return httperr.New(http.StatusInternalServerError, err)
	}

	return httperr.OK()
}

func (c *gatewayService) DeleteProfile(profileId string) httperr.Err {
	stub, err := c.createAccountsServiceStub()
	if err != nil {
		return httperr.New(http.StatusInternalServerError, err)
	}

	_, err = stub.UnregisterUser(context.Background(), &pb.UnregisterUserRequest{
		ProfileId: profileId,
	})
	if err != nil {
		return httperr.New(http.StatusInternalServerError, err)
	}

	return httperr.OK()
}
