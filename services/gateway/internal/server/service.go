package server

import (
	"context"
	"errors"
	"net/http"
	"soa-socialnetwork/internal/soajwt"
	pb "soa-socialnetwork/services/accounts/proto"
	"soa-socialnetwork/services/gateway/api"
	"soa-socialnetwork/services/gateway/internal/server/birthday"
	"soa-socialnetwork/services/gateway/internal/server/httperr"
	"soa-socialnetwork/services/gateway/internal/server/query"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GatewayService struct {
	JwtVerifier soajwt.Verifier
}

func initService(cfg GatewayServiceConfig) GatewayService {
	return GatewayService{
		JwtVerifier: soajwt.NewVerifier(cfg.JwtPublicKey),
	}
}

func (c *GatewayService) createAccountsServiceStub() (pb.AccountsServiceClient, error) {
	conn, err := grpc.NewClient("accounts-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return pb.NewAccountsServiceClient(conn), nil
}

func (c *GatewayService) RegisterProfile(qp *query.Params, req *api.RegisterProfileRequest) (api.RegisterProfileResponse, httperr.Err) {
	stub, err := c.createAccountsServiceStub()
	if err != nil {
		return api.RegisterProfileResponse{}, httperr.New(http.StatusInternalServerError, err)
	}

	resp, err := stub.RegisterUser(context.Background(), &pb.RegisterUserRequest{
		Login:       req.Login,
		Password:    req.Password,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
		Surname:     req.Surname,
	})
	if err != nil {
		return api.RegisterProfileResponse{}, httperr.New(http.StatusInternalServerError, err)
	}

	return api.RegisterProfileResponse{
		ProfileID: resp.ProfileId,
	}, httperr.Ok()
}

func (c *GatewayService) GetProfileInfo(qp *query.Params) (api.GetProfileResponse, httperr.Err) {
	stub, err := c.createAccountsServiceStub()
	if err != nil {
		return api.GetProfileResponse{}, httperr.New(http.StatusInternalServerError, err)
	}

	resp, err := stub.GetProfile(context.Background(), &pb.GetProfileRequest{
		ProfileId: qp.ProfileId,
	})
	if err != nil {
		return api.GetProfileResponse{}, httperr.New(http.StatusInternalServerError, err)
	}

	return api.GetProfileResponse{
		Name:     resp.Name,
		Surname:  resp.Surname,
		Birthday: resp.Birthday.AsTime().Format("2006-01-02"),
		Bio:      resp.Bio,
	}, httperr.Ok()
}

func (c *GatewayService) EditProfileInfo(qp *query.Params, req *api.EditProfileRequest) httperr.Err {
	stub, err := c.createAccountsServiceStub()
	if err != nil {
		return httperr.New(http.StatusInternalServerError, err)
	}

	var pbBirthday *timestamppb.Timestamp = nil
	if len(req.Birthday) > 0 {
		bday, err := birthday.Parse(req.Birthday)
		if err != nil {
			return httperr.New(http.StatusInternalServerError, errors.New("bad birthday, check validity"))
		}
		pbBirthday = timestamppb.New(bday.AsTime())
	}

	_, err = stub.EditProfile(context.Background(), &pb.EditProfileRequest{
		ProfileId: qp.ProfileId,
		EditedProfileData: &pb.Profile{
			Name:     req.Name,
			Surname:  req.Surname,
			Birthday: pbBirthday,
			Bio:      req.Bio,
		},
	})
	if err != nil {
		return httperr.New(http.StatusInternalServerError, err)
	}

	return httperr.Ok()
}

func (c *GatewayService) DeleteProfile(qp *query.Params) httperr.Err {
	stub, err := c.createAccountsServiceStub()
	if err != nil {
		return httperr.New(http.StatusInternalServerError, err)
	}

	_, err = stub.UnregisterUser(context.Background(), &pb.UnregisterUserRequest{
		ProfileId: qp.ProfileId,
	})
	if err != nil {
		return httperr.New(http.StatusInternalServerError, err)
	}

	return httperr.Ok()
}

func (c *GatewayService) Authenticate(qp *query.Params, req *api.AuthenticateRequest) (api.AuthenticateResponse, httperr.Err) {
	stub, err := c.createAccountsServiceStub()
	if err != nil {
		return api.AuthenticateResponse{}, httperr.New(http.StatusInternalServerError, err)
	}

	resp, err := stub.Authenticate(context.Background(), &pb.AuthenticateRequest{
		Login:    req.Login,
		Password: req.Password,
	})

	if err != nil {
		return api.AuthenticateResponse{}, httperr.New(http.StatusInternalServerError, err)
	}

	return api.AuthenticateResponse{
		Token: resp.Token,
	}, httperr.Ok()
}
