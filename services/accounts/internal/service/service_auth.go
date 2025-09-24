package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"soa-socialnetwork/services/accounts/internal/soajwtissuer"
	"time"

	pb "soa-socialnetwork/services/accounts/proto"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type userIdKind int

const (
	user_id_login userIdKind = iota
	user_id_email
	user_id_phone_number
)

const (
	login_column        = "login"
	email_column        = "email"
	phone_number_column = "phone_number"
)

func (id userIdKind) asDbColumn() string {
	switch id {
	case user_id_login:
		{
			return login_column
		}

	case user_id_email:
		{
			return email_column
		}

	case user_id_phone_number:
		{
			return phone_number_column
		}
	}

	return "UNKNOWN"
}

type userId struct {
	kind  userIdKind
	value string
}

type accountData struct {
	accountId int
	profileId uuid.UUID
	password  string
}

const JWT_DEFAULT_TTL = 30 * time.Second

func extractUserIdFromProto(protoUserId any) userId {
	switch v := protoUserId.(type) {
	case *pb.AuthByPassword_Login:
		{
			return userId{
				kind:  user_id_login,
				value: v.Login,
			}
		}
	case *pb.AuthByPassword_Email:
		{
			return userId{
				kind:  user_id_email,
				value: v.Email,
			}
		}
	case *pb.AuthByPassword_PhoneNumber:
		{
			return userId{
				kind:  user_id_phone_number,
				value: v.PhoneNumber,
			}
		}
	}

	panic("shouldn't reach here")
}

func (s *AccountsService) fetchAccountData(ctx context.Context, uid userId) (accountData, error) {
	sql := fmt.Sprintf(`
	SELECT a.id, p.profile_id, a.password
	FROM accounts AS a
	JOIN profiles AS p ON a.id = p.account_id
	WHERE a.%s = $1;
	`, uid.kind.asDbColumn())

	row := s.dbPool.QueryRow(ctx, sql, uid.value)
	var data accountData

	if err := row.Scan(&data.accountId, &data.profileId, &data.password); err != nil {
		return accountData{}, errors.New("account not found")
	}

	return data, nil
}

func (s *AccountsService) Authenticate(ctx context.Context, req *pb.AuthByPassword) (*pb.AuthResponse, error) {
	accData, err := s.fetchAccountData(ctx, extractUserIdFromProto(req.UserId))
	if err != nil {
		return nil, err
	}

	if accData.password != req.Password {
		return nil, errors.New("passwords doesn't match")
	}

	token, err := s.jwtIssuer.Issue(soajwtissuer.PersonalData{
		AccountId: accData.accountId,
		ProfileId: accData.profileId.String(),
	}, JWT_DEFAULT_TTL)
	if err != nil {
		log.Printf("cannot create jwt token: %v", err)
		return nil, err
	}

	return &pb.AuthResponse{
		Token: token,
	}, nil
}

const SOA_API_TOKEN_LEN = 48

type soaApiToken [SOA_API_TOKEN_LEN]byte

func (t soaApiToken) toBase64() string {
	return base64.RawStdEncoding.EncodeToString(t[:])
}

func buildSoaApiToken(udata accountData) (token soaApiToken, err error) {
	n := copy(token[:16], udata.profileId[:])
	if n != 16 {
		panic("must copy exact 16 bytes")
	}

	binary.LittleEndian.PutUint32(token[16:20], uint32(udata.accountId))

	_, err = rand.Read(token[20:])
	if err != nil {
		return
	}

	return
}

func (s *AccountsService) CreateApiToken(ctx context.Context, req *pb.CreateApiTokenRequest) (*pb.CreateApiTokenResponse, error) {
	if req.Params == nil {
		return nil, errors.New("no api token params")
	}

	if req.Auth == nil {
		return nil, errors.New("no auth params")
	}

	accData, err := s.fetchAccountData(ctx, extractUserIdFromProto(req.Auth.UserId))
	if err != nil {
		return nil, err
	}

	if accData.password != req.Auth.Password {
		return nil, errors.New("passwords doesn't match")
	}

	token, err := buildSoaApiToken(accData)
	if err != nil {
		return nil, err
	}

	tokenBase64 := token.toBase64()
	sql := `
	INSERT INTO api_tokens(account_id, token, valid_until, read_access, write_access)
	VALUES ($1, $2, NOW() + $3, $4, $5)
	RETURNING valid_until;
	`

	row := s.dbPool.QueryRow(ctx, sql, accData.accountId, tokenBase64, req.Params.Ttl.AsDuration(), req.Params.ReadAccess, req.Params.WriteAccess)
	var validUntil time.Time
	if err := row.Scan(&validUntil); err != nil {
		return nil, err
	}

	return &pb.CreateApiTokenResponse{
		Token:      tokenBase64,
		ValidUntil: timestamppb.New(validUntil),
	}, nil
}

func (s *AccountsService) ValidateApiToken(ctx context.Context, req *pb.ApiToken) (*pb.ApiTokenValidity, error) {
	sql := `
		SELECT read_access, write_access, valid_until, created_at
		FROM api_tokens
		WHERE token = $1
	`

	row := s.dbPool.QueryRow(ctx, sql, req.Token)
	var (
		readAccess  bool
		writeAccess bool
		validUntil  time.Time
		createdAt   time.Time
	)

	if err := row.Scan(&readAccess, &writeAccess, &validUntil, &createdAt); err != nil {
		return nil, err
	}

	now := time.Now()
	if now.After(validUntil) {
		return &pb.ApiTokenValidity{
			Result: &pb.ApiTokenValidity_Invalid_{
				Invalid: &pb.ApiTokenValidity_Invalid{},
			},
		}, nil
	}

	return &pb.ApiTokenValidity{
		Result: &pb.ApiTokenValidity_Valid_{
			Valid: &pb.ApiTokenValidity_Valid{
				ReadAccess:  readAccess,
				WriteAccess: writeAccess,
				ValidUntil:  timestamppb.New(validUntil),
				CreatedAt:   timestamppb.New(createdAt),
			},
		},
	}, nil
}
