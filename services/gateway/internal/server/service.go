package server

import (
	"context"
	"fmt"
	"net/http"
	"soa-socialnetwork/services/accounts/pkg/soajwt"
	accountsPb "soa-socialnetwork/services/accounts/proto"
	"soa-socialnetwork/services/gateway/api"
	"soa-socialnetwork/services/gateway/internal/server/httperr"
	"soa-socialnetwork/services/gateway/internal/server/query"
	"soa-socialnetwork/services/gateway/internal/server/soagrpc"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GatewayService struct {
	jwtVerifier         soajwt.Verifier
	accountsGrpcTarget  string
	accountsStubFactory soagrpc.AccountsStubFactory
}

func newGatewayService(cfg GatewayServiceConfig) GatewayService {
	if cfg.AccountsServiceStubFactory == nil {
		cfg.AccountsServiceStubFactory = defaultAccountsStubFactory{}
	}
	return GatewayService{
		jwtVerifier:         soajwt.NewVerifier(cfg.JwtPublicKey),
		accountsGrpcTarget:  fmt.Sprintf("%s:%d", cfg.AccountsServiceHost, cfg.AccountsServicePort),
		accountsStubFactory: cfg.AccountsServiceStubFactory,
	}
}

func (s *GatewayService) createAccountsStub(qp *query.Params) (accountsPb.AccountsServiceClient, error) {
	return s.accountsStubFactory.New(s.accountsGrpcTarget, qp)
}

func (s *GatewayService) RegisterProfile(qp *query.Params, req *api.RegisterProfileRequest) (api.RegisterProfileResponse, httperr.Err) {
	stub, err := s.createAccountsStub(qp)
	if err != nil {
		return api.RegisterProfileResponse{}, httperr.New(http.StatusInternalServerError, err)
	}

	resp, err := stub.RegisterUser(context.Background(), &accountsPb.RegisterUserRequest{
		Login:       string(req.Login),
		Password:    string(req.Password),
		Email:       string(req.Email),
		PhoneNumber: string(req.PhoneNumber),
		Name:        string(req.Name),
		Surname:     string(req.Surname),
	})
	if err != nil {
		return api.RegisterProfileResponse{}, httperr.New(http.StatusInternalServerError, err)
	}

	return api.RegisterProfileResponse{
		ProfileId: resp.ProfileId,
	}, httperr.Ok()
}

func (s *GatewayService) GetProfileInfo(qp *query.Params) (api.GetProfileResponse, httperr.Err) {
	stub, err := s.createAccountsStub(qp)
	if err != nil {
		return api.GetProfileResponse{}, httperr.New(http.StatusInternalServerError, err)
	}

	resp, err := stub.GetProfile(context.Background(), &accountsPb.GetProfileRequest{
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

func (s *GatewayService) EditProfileInfo(qp *query.Params, req *api.EditProfileRequest) httperr.Err {
	stub, err := s.createAccountsStub(qp)
	if err != nil {
		return httperr.New(http.StatusInternalServerError, err)
	}

	var pbBirthday *timestamppb.Timestamp = nil
	if req.Birthday.HasValue {
		pbBirthday = timestamppb.New(req.Birthday.Value.Time)
	}

	_, err = stub.EditProfile(context.Background(), &accountsPb.EditProfileRequest{
		ProfileId: qp.ProfileId,
		EditedProfileData: &accountsPb.Profile{
			Name:     string(req.Name.Value),
			Surname:  string(req.Surname.Value),
			Birthday: pbBirthday,
			Bio:      string(req.Bio.Value),
		},
	})
	if err != nil {
		return httperr.New(http.StatusInternalServerError, err)
	}

	return httperr.Ok()
}

func (s *GatewayService) DeleteProfile(qp *query.Params) httperr.Err {
	stub, err := s.createAccountsStub(qp)
	if err != nil {
		return httperr.New(http.StatusInternalServerError, err)
	}

	_, err = stub.UnregisterUser(context.Background(), &accountsPb.UnregisterUserRequest{
		ProfileId: qp.ProfileId,
	})
	if err != nil {
		return httperr.New(http.StatusInternalServerError, err)
	}

	return httperr.Ok()
}

func (s *GatewayService) buildAuthByPassword(req *api.AuthenticateRequest) (proto accountsPb.AuthByPassword) {
	if req.Login.HasValue {
		proto.UserId = &accountsPb.AuthByPassword_Login{
			Login: string(req.Login.Value),
		}
	} else if req.Email.HasValue {
		proto.UserId = &accountsPb.AuthByPassword_Email{
			Email: string(req.Email.Value),
		}
	} else if req.PhoneNumber.HasValue {
		proto.UserId = &accountsPb.AuthByPassword_PhoneNumber{
			PhoneNumber: string(req.PhoneNumber.Value),
		}
	} else {
		panic("at least one user id must be provided")
	}

	proto.Password = string(req.Password)
	return
}

func (s *GatewayService) Authenticate(qp *query.Params, req *api.AuthenticateRequest) (api.AuthenticateResponse, httperr.Err) {
	stub, err := s.createAccountsStub(qp)
	if err != nil {
		return api.AuthenticateResponse{}, httperr.New(http.StatusInternalServerError, err)
	}

	protoRequest := s.buildAuthByPassword(req)
	resp, err := stub.Authenticate(context.Background(), &protoRequest)

	if err != nil {
		return api.AuthenticateResponse{}, httperr.New(http.StatusInternalServerError, err)
	}

	return api.AuthenticateResponse{
		Token: resp.Token,
	}, httperr.Ok()
}

func (s *GatewayService) CreateApiToken(qp *query.Params, req *api.CreateApiTokenRequest) (api.CreateApiTokenResponse, httperr.Err) {
	stub, err := s.createAccountsStub(qp)
	if err != nil {
		return api.CreateApiTokenResponse{}, httperr.New(http.StatusInternalServerError, err)
	}

	protoAuthByPassword := s.buildAuthByPassword(&req.Auth)

	resp, err := stub.CreateApiToken(context.Background(), &accountsPb.CreateApiTokenRequest{
		Auth: &protoAuthByPassword,
		Params: &accountsPb.AuthTokenParams{
			ReadAccess:  req.ReadAccess,
			WriteAccess: req.WriteAccess,
			Ttl:         durationpb.New(req.Ttl.Duration),
		},
	})

	if err != nil {
		return api.CreateApiTokenResponse{}, httperr.New(http.StatusInternalServerError, err)
	}

	return api.CreateApiTokenResponse{
		Token: resp.Token,
	}, httperr.Ok()
}
