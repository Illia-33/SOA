package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	pb "soa-socialnetwork/services/accounts/proto"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AccountsServiceConfig struct {
	DbHost     string
	DbUser     string
	DbPassword string
	DbPoolSize int
}

type AccountsService struct {
	pb.UnimplementedAccountsServiceServer

	dbpool *pgxpool.Pool
}

func (s *AccountsService) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	profileUUID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	sql := `
	WITH acc_id AS (
		INSERT INTO accounts(login, password, email, phone_number)
		VALUES ('$1', '$2', '$3', '$4')
		RETURNING id;
	)
	INSERT INTO profiles(account_id, name, surname, profile_id)
	VALUES (acc_id, '$5', '$6', '$7');
	`

	_, err = s.dbpool.Exec(ctx, sql, req.Login, req.Password, req.Email, req.PhoneNumber, req.Name, req.Surname, profileUUID)
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
		DELETE FROM profiles WHERE profile_id = '$1'
		RETURNING id;
	)
	DELETE FROM accounts WHERE id = acc_id;
	`

	_, err := s.dbpool.Exec(ctx, sql, req.ProfileId)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *AccountsService) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.Profile, error) {
	sql := `
	SELECT name, surname, birthday, bio
	FROM profiles
	WHERE profile_id = '$1'
	LIMIT 1;
	`
	row := s.dbpool.QueryRow(ctx, sql, req.ProfileId)
	var (
		name     string
		surname  string
		birthday time.Time
		bio      string
	)

	if err := row.Scan(&name, &surname, &birthday, &bio); err != nil {
		return nil, errors.New("profile not found")
	}

	return &pb.Profile{
		Name:      name,
		Surname:   surname,
		ProfileId: req.ProfileId,
		Birthday:  timestamppb.New(birthday),
		Bio:       bio,
	}, nil
}
func (s *AccountsService) EditProfile(ctx context.Context, req *pb.EditProfileRequest) (*pb.Empty, error) {
	if req.EditedProfileData == nil {
		return nil, errors.New("empty edit data")
	}

	p := req.EditedProfileData

	toPostgresOptArg := func(obj any) string {
		switch v := obj.(type) {
		case string:
			if len(v) > 0 {
				return fmt.Sprintf("'%s'", v)
			}
		case *timestamppb.Timestamp:
			if v != nil {
				year, month, day := v.AsTime().Date()
				return fmt.Sprintf("%04d-%02d-%02d", year, int(month), day)
			}
		}

		return "NULL"
	}

	sql := `
	UPDATE profiles
	SET
		name = COALESCE($1, name),
		surname = COALESCE($2, surname),
		birthday = COALESCE($3, birthday),
		bio = COALESCE($4, bio)
	WHERE profile_id = '$5'
	`

	_, err := s.dbpool.Exec(ctx, sql, toPostgresOptArg(p.Name), toPostgresOptArg(p.Surname), toPostgresOptArg(p.Birthday), toPostgresOptArg(p.Bio), req.ProfileId)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func createAccountsService(cfg AccountsServiceConfig) (*AccountsService, error) {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=5432 dbname=accounts-postgres sslmode=disable pool_max_conns=%d", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPoolSize)
	pool, err := pgxpool.New(context.Background(), connStr)

	if err != nil {
		return nil, err
	}

	return &AccountsService{
		dbpool: pool,
	}, nil
}
