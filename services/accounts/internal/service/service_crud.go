package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	pb "soa-socialnetwork/services/accounts/proto"
	statsModels "soa-socialnetwork/services/stats/pkg/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *AccountsService) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	profileUUID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	tx, err := s.dbPool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	sql := `
	INSERT INTO accounts (login, password, email, phone_number)
	VALUES ($1, $2, $3, $4)
	RETURNING id;
	`
	row := tx.QueryRow(ctx, sql, req.Login, req.Password, req.Email, req.PhoneNumber)
	var accountId int
	err = row.Scan(&accountId)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	sql = `
	INSERT INTO profiles(account_id, profile_id, name, surname)
	VALUES ($1, $2::uuid, $3, $4);
	`
	_, err = tx.Exec(ctx, sql, accountId, profileUUID, req.Name, req.Surname)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	payload, err := json.Marshal(statsModels.RegistrationEvent{
		AccountId: statsModels.AccountId(accountId),
		ProfileId: profileUUID.String(),
		Timestamp: time.Now(),
	})
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	sql = `
	INSERT INTO outbox(event_type, payload)
	VALUES ('registration', $1::jsonb);
	`
	_, err = tx.Exec(ctx, sql, payload)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterUserResponse{
		ProfileId: profileUUID.String(),
	}, nil
}

func (s *AccountsService) UnregisterUser(ctx context.Context, req *pb.UnregisterUserRequest) (*pb.Empty, error) {
	sql := `
	WITH acc_id AS (
		DELETE FROM profiles 
		WHERE profile_id = $1
		RETURNING account_id
	)
	DELETE FROM accounts
	WHERE id IN (SELECT account_id FROM acc_id);
	`

	_, err := s.dbPool.Exec(ctx, sql, req.ProfileId)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *AccountsService) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.Profile, error) {
	sql := `
	SELECT name, surname, birthday, bio
	FROM profiles
	WHERE profile_id = $1;
	`
	row := s.dbPool.QueryRow(ctx, sql, req.ProfileId)
	var (
		name       pgtype.Text
		surname    pgtype.Text
		pgBirthday pgtype.Date
		bio        pgtype.Text
	)

	if err := row.Scan(&name, &surname, &pgBirthday, &bio); err != nil {
		return nil, status.Error(codes.NotFound, "profile not found")
	}

	var birthday *timestamppb.Timestamp = nil
	if pgBirthday.Valid {
		birthday = timestamppb.New(pgBirthday.Time)
	}

	return &pb.Profile{
		Name:      name.String,
		Surname:   surname.String,
		ProfileId: req.ProfileId,
		Birthday:  birthday,
		Bio:       bio.String,
	}, nil
}

func (s *AccountsService) EditProfile(ctx context.Context, req *pb.EditProfileRequest) (*pb.Empty, error) {
	if req.EditedProfileData == nil {
		return nil, errors.New("empty edit data")
	}

	p := req.EditedProfileData

	stringToPgText := func(s string) pgtype.Text {
		if len(s) == 0 {
			return pgtype.Text{}
		}

		return pgtype.Text{
			String: s,
			Valid:  true,
		}
	}

	pbTimestampToPgDate := func(ts *timestamppb.Timestamp) pgtype.Date {
		if ts == nil {
			return pgtype.Date{}
		}

		return pgtype.Date{
			Time:             ts.AsTime(),
			InfinityModifier: pgtype.Finite,
			Valid:            true,
		}
	}

	pgName := stringToPgText(p.Name)
	pgSurname := stringToPgText(p.Surname)
	pgBirthday := pbTimestampToPgDate(p.Birthday)
	pgBio := stringToPgText(p.Bio)

	sql := `
	WITH affected_rows AS (
		UPDATE profiles
		SET
			name = COALESCE($1, name),
			surname = COALESCE($2, surname),
			birthday = COALESCE($3, birthday),
			bio = COALESCE($4, bio)
		WHERE profile_id = $5
		RETURNING 1
	)
	SELECT count(*) FROM affected_rows;
	`

	row := s.dbPool.QueryRow(ctx, sql, pgName, pgSurname, pgBirthday, pgBio, req.ProfileId)
	var cnt int
	if err := row.Scan(&cnt); err != nil {
		return nil, err
	}

	if cnt == 0 {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	if cnt != 1 {
		log.Printf("warning: there more than one users with profile id %s", req.ProfileId)
	}

	return &pb.Empty{}, nil
}

func (s *AccountsService) ResolveProfileId(ctx context.Context, req *pb.ResolveProfileIdRequest) (*pb.ResolveProfileIdResponse, error) {
	profileId := req.ProfileId
	sql := `
	SELECT account_id
	FROM profiles
	WHERE profile_id = $1;
	`

	row := s.dbPool.QueryRow(ctx, sql, profileId)
	var accountId int
	if err := row.Scan(&accountId); err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &pb.ResolveProfileIdResponse{
		AccountId: int32(accountId),
	}, nil
}

func (s *AccountsService) ResolveAccountId(ctx context.Context, req *pb.ResolveAccountIdRequest) (*pb.ResolveAccountIdResponse, error) {
	accountId := req.AccountId
	sql := `
	SELECT profile_id
	FROM profiles
	WHERE account_id = $1;
	`

	row := s.dbPool.QueryRow(ctx, sql, accountId)
	var profileId string
	if err := row.Scan(&profileId); err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &pb.ResolveAccountIdResponse{
		ProfileId: profileId,
	}, nil
}
