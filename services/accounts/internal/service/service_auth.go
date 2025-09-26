package service

import (
	"context"
	"log"
	"soa-socialnetwork/services/accounts/internal/models"
	"soa-socialnetwork/services/accounts/internal/repo"
	"soa-socialnetwork/services/accounts/internal/soajwtissuer"
	"time"

	pb "soa-socialnetwork/services/accounts/proto"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const JWT_DEFAULT_TTL = 30 * time.Second

func (s *AccountsService) Authenticate(ctx context.Context, req *pb.AuthByPassword) (*pb.AuthResponse, error) {
	conn, err := s.Db.OpenConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	accountParams, err := s.checkPassword(conn, req)
	if err != nil {
		return nil, err
	}

	profileId, err := conn.Profiles().ResolveAccountId(accountParams.Id)
	if err != nil {
		return nil, err
	}

	token, err := s.jwtIssuer.Issue(soajwtissuer.PersonalData{
		AccountId: int(accountParams.Id),
		ProfileId: string(profileId),
	}, JWT_DEFAULT_TTL)
	if err != nil {
		log.Printf("cannot create jwt token: %v", err)
		return nil, err
	}

	return &pb.AuthResponse{
		Token: token,
	}, nil
}

func (s *AccountsService) CreateApiToken(ctx context.Context, req *pb.CreateApiTokenRequest) (*pb.CreateApiTokenResponse, error) {
	if req.Params == nil {
		return nil, status.Error(codes.InvalidArgument, "no api token params")
	}

	if req.Auth == nil {
		return nil, status.Error(codes.InvalidArgument, "no auth params")
	}

	conn, err := s.Db.OpenConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	accountParams, err := s.checkPassword(conn, req.Auth)
	if err != nil {
		return nil, err
	}

	profileId, err := conn.Profiles().ResolveAccountId(accountParams.Id)
	if err != nil {
		return nil, err
	}

	token, err := buildSoaApiToken(accountData{
		accountId: int(accountParams.Id),
		profileId: uuid.MustParse(string(profileId)),
	})
	if err != nil {
		return nil, err
	}

	tokenBase64 := token.toBase64()

	validUntil, err := conn.ApiTokens().Put(models.ApiToken(tokenBase64), repo.ApiTokenParams{
		AccountId:   int(accountParams.Id),
		ReadAccess:  req.Params.ReadAccess,
		WriteAccess: req.Params.WriteAccess,
		Ttl:         req.Params.Ttl.AsDuration(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateApiTokenResponse{
		Token:      tokenBase64,
		ValidUntil: timestamppb.New(validUntil),
	}, nil
}

func (s *AccountsService) ValidateApiToken(ctx context.Context, req *pb.ApiToken) (*pb.ApiTokenValidity, error) {
	conn, err := s.Db.OpenConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	tokenData, err := conn.ApiTokens().Get(models.ApiToken(req.Token))
	if err != nil {
		return nil, err
	}

	now := time.Now()
	if now.After(tokenData.ValidUntil) {
		return &pb.ApiTokenValidity{
			Result: &pb.ApiTokenValidity_Invalid_{
				Invalid: &pb.ApiTokenValidity_Invalid{},
			},
		}, nil
	}

	return &pb.ApiTokenValidity{
		Result: &pb.ApiTokenValidity_Valid_{
			Valid: &pb.ApiTokenValidity_Valid{
				ReadAccess:  tokenData.ReadAccess,
				WriteAccess: tokenData.WriteAccess,
				ValidUntil:  timestamppb.New(tokenData.ValidUntil),
				CreatedAt:   timestamppb.New(tokenData.CreatedAt),
			},
		},
	}, nil
}

func (s *AccountsService) checkPassword(conn repo.Connection, authData *pb.AuthByPassword) (models.AccountParams, error) {
	switch userId := authData.UserId.(type) {
	case *pb.AuthByPassword_Login:
		return conn.Accounts().CheckPasswordByLogin(userId.Login, authData.Password)

	case *pb.AuthByPassword_Email:
		return conn.Accounts().CheckPasswordByEmail(userId.Email, authData.Password)

	case *pb.AuthByPassword_PhoneNumber:
		return conn.Accounts().CheckPasswordByPhoneNumber(userId.PhoneNumber, authData.Password)

	default:
		panic("unknown user id")
	}
}
