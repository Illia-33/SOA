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

func (c *GatewayService) createAccountsServiceStub(qp *query.Params) (pb.AccountsServiceClient, error) {
	client, err := createGrpcClient("accounts-service:50051", qp)
	if err != nil {
		return nil, err
	}

	return pb.NewAccountsServiceClient(client), nil
}

func (c *GatewayService) RegisterProfile(qp *query.Params, req *api.RegisterProfileRequest) (api.RegisterProfileResponse, httperr.Err) {
	stub, err := c.createAccountsServiceStub(qp)
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
	stub, err := c.createAccountsServiceStub(qp)
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
	stub, err := c.createAccountsServiceStub(qp)
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
	stub, err := c.createAccountsServiceStub(qp)
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
	stub, err := c.createAccountsServiceStub(qp)
	if err != nil {
		return api.AuthenticateResponse{}, httperr.New(http.StatusInternalServerError, err)
	}

	var grpcRequest pb.AuthenticateRequest
	{
		if len(req.Login) > 0 {
			grpcRequest.UserId = &pb.AuthenticateRequest_Login{
				Login: req.Login,
			}
		} else if len(req.Email) > 0 {
			grpcRequest.UserId = &pb.AuthenticateRequest_Email{
				Email: req.Email,
			}
		} else if len(req.PhoneNumber) > 0 {
			grpcRequest.UserId = &pb.AuthenticateRequest_PhoneNumber{
				PhoneNumber: req.PhoneNumber,
			}
		} else {
			panic("shouldn't reach here")
		}
	}
	grpcRequest.Password = req.Password

	resp, err := stub.Authenticate(context.Background(), &grpcRequest)

	if err != nil {
		return api.AuthenticateResponse{}, httperr.New(http.StatusInternalServerError, err)
	}

	return api.AuthenticateResponse{
		Token: resp.Token,
	}, httperr.Ok()
}
