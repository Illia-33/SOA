package server

import (
	"context"
	"crypto/ed25519"
	"errors"
	"fmt"
	"log"
	"time"

	"soa-socialnetwork/services/accounts/internal/server/jwtsigner"
	pb "soa-socialnetwork/services/accounts/proto"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AccountsServiceConfig struct {
	DbHost        string
	DbUser        string
	DbPassword    string
	DbPoolSize    int
	JwtPrivateKey ed25519.PrivateKey
}

type AccountsService struct {
	pb.UnimplementedAccountsServiceServer

	dbpool    *pgxpool.Pool
	jwtSigner jwtsigner.Signer
}

func createAccountsService(cfg AccountsServiceConfig) (*AccountsService, error) {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=5432 dbname=accounts-postgres sslmode=disable pool_max_conns=%d", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPoolSize)
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	signer, err := jwtsigner.New(cfg.JwtPrivateKey)
	if err != nil {
		return nil, err
	}

	return &AccountsService{
		dbpool:    pool,
		jwtSigner: signer,
	}, nil
}

func (s *AccountsService) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	profileUUID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	sql := `
	WITH acc_id AS (
		INSERT INTO accounts(login, password, email, phone_number)
		VALUES ($1, $2, $3, $4)
		RETURNING id as account_id
	)
	INSERT INTO profiles(account_id, name, surname, profile_id)
	SELECT account_id, name, surname, profile_id
	FROM (VALUES ($5, $6, $7::uuid)) as t (name, surname, profile_id)
	CROSS JOIN acc_id;
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
		DELETE FROM profiles 
		WHERE profile_id = $1
		RETURNING account_id
	)
	DELETE FROM accounts
	WHERE id IN (SELECT account_id FROM acc_id);
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
	WHERE profile_id = $1
	LIMIT 1;
	`
	row := s.dbpool.QueryRow(ctx, sql, req.ProfileId)
	var (
		name       pgtype.Text
		surname    pgtype.Text
		pgBirthday pgtype.Date
		bio        pgtype.Text
	)

	if err := row.Scan(&name, &surname, &pgBirthday, &bio); err != nil {
		log.Printf("error while getting %s profile: %v", req.ProfileId, err)
		return nil, errors.New("profile not found")
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

	row := s.dbpool.QueryRow(ctx, sql, pgName, pgSurname, pgBirthday, pgBio, req.ProfileId)
	var cnt int
	if err := row.Scan(&cnt); err != nil {
		return nil, err
	}

	if cnt == 0 {
		return nil, errors.New("user not found")
	}

	if cnt != 1 {
		log.Printf("warning: there more than one users with profile id %s", req.ProfileId)
	}

	return &pb.Empty{}, nil
}

func (s *AccountsService) Authenticate(ctx context.Context, req *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	sql := `
	SELECT a.id, p.profile_id
	FROM accounts AS a
	JOIN profiles AS p ON a.id = p.account_id
	WHERE a.login = $1 AND a.password = $2;
	`

	row := s.dbpool.QueryRow(ctx, sql, req.Login, req.Password)
	var data jwtsigner.PersonalData
	if err := row.Scan(&data.AccountId, &data.ProfileId); err != nil {
		return nil, errors.New("user not found")
	}

	token, err := s.jwtSigner.Sign(data, 30*time.Second)
	if err != nil {
		log.Printf("cannot create jwt token: %v", err)
		return nil, err
	}

	return &pb.AuthenticateResponse{
		Token: token,
	}, nil
}
