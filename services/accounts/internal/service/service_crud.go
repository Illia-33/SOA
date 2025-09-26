package service

import (
	"context"
	"encoding/json"
	"time"

	"soa-socialnetwork/services/accounts/internal/models"
	"soa-socialnetwork/services/accounts/internal/repo"
	"soa-socialnetwork/services/accounts/internal/service/errs"
	"soa-socialnetwork/services/accounts/internal/service/interceptors"
	pb "soa-socialnetwork/services/accounts/proto"
	"soa-socialnetwork/services/common/option"
	statsModels "soa-socialnetwork/services/stats/pkg/models"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func getAuthInfo(ctx context.Context) interceptors.AuthInfo {
	authInfo := ctx.Value(interceptors.AuthInfoKey).(interceptors.AuthInfo)
	return authInfo
}

func (s *AccountsService) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	tx, err := s.Db.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Close()

	profileId := uuid.New().String()
	registrationData := models.RegistrationData{
		Login:       req.Login,
		Password:    req.Password,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
		Surname:     req.Surname,
	}

	accountId, err := tx.Accounts().New(registrationData)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Profiles().New(models.ProfileId(profileId), accountId, registrationData)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	payload, err := json.Marshal(statsModels.RegistrationEvent{
		AccountId: statsModels.AccountId(accountId),
		ProfileId: profileId,
		Timestamp: time.Now(),
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Outbox().Put(models.OutboxEvent{
		Type:      "registration",
		Payload:   payload,
		CreatedAt: time.Now(),
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &pb.RegisterUserResponse{
		ProfileId: profileId,
	}, nil
}

func (s *AccountsService) UnregisterUser(ctx context.Context, req *pb.UnregisterUserRequest) (*pb.Empty, error) {
	authInfo := getAuthInfo(ctx)
	if authInfo.ProfileId != req.ProfileId {
		return nil, errs.AccessDenied{}
	}

	profileId := models.ProfileId(req.ProfileId)
	accountId, err := func() (models.AccountId, error) {
		conn, err := s.Db.OpenConnection(ctx)
		if err != nil {
			return -1, err
		}
		defer conn.Close()

		return conn.Profiles().ResolveProfileId(profileId)
	}()
	if err != nil {
		return nil, err
	}

	tx, err := s.Db.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Close()

	err = tx.Accounts().Delete(accountId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Profiles().Delete(profileId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *AccountsService) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.Profile, error) {
	conn, err := s.Db.OpenConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	profileId := models.ProfileId(req.ProfileId)

	data, err := conn.Profiles().GetByProfileId(profileId)
	if err != nil {
		return nil, err
	}

	return &pb.Profile{
		Name:      data.Name,
		Surname:   data.Surname,
		ProfileId: req.ProfileId,
		Birthday:  timestamppb.New(data.Birthday),
		Bio:       data.Bio,
	}, nil
}

func (s *AccountsService) EditProfile(ctx context.Context, req *pb.EditProfileRequest) (*pb.Empty, error) {
	authInfo := getAuthInfo(ctx)
	if authInfo.ProfileId != req.ProfileId {
		return nil, errs.AccessDenied{}
	}

	if req.EditedProfileData == nil {
		return nil, status.Error(codes.InvalidArgument, "empty edit data")
	}

	pbEditedData := req.EditedProfileData

	edited := func() repo.EditedProfileData {
		var edited repo.EditedProfileData

		if pbEditedData.Name != "" {
			edited.Name = option.Some(pbEditedData.Name)
		}

		if pbEditedData.Surname != "" {
			edited.Surname = option.Some(pbEditedData.Surname)
		}

		if pbEditedData.Birthday != nil {
			edited.Birthday = option.Some(pbEditedData.Birthday.AsTime())
		}

		if pbEditedData.Bio != "" {
			edited.Bio = option.Some(pbEditedData.Bio)
		}

		return edited
	}()

	conn, err := s.Db.OpenConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	err = conn.Profiles().Edit(models.ProfileId(req.ProfileId), edited)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *AccountsService) ResolveProfileId(ctx context.Context, req *pb.ResolveProfileIdRequest) (*pb.ResolveProfileIdResponse, error) {
	profileId := models.ProfileId(req.ProfileId)

	conn, err := s.Db.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	accountId, err := conn.Profiles().ResolveProfileId(profileId)
	if err != nil {
		return nil, err
	}

	return &pb.ResolveProfileIdResponse{
		AccountId: int32(accountId),
	}, nil
}

func (s *AccountsService) ResolveAccountId(ctx context.Context, req *pb.ResolveAccountIdRequest) (*pb.ResolveAccountIdResponse, error) {
	accountId := models.AccountId(req.AccountId)
	conn, err := s.Db.OpenConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	profileId, err := conn.Profiles().ResolveAccountId(accountId)
	if err != nil {
		return nil, err
	}

	return &pb.ResolveAccountIdResponse{
		ProfileId: string(profileId),
	}, nil
}
