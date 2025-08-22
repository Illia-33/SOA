package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"soa-socialnetwork/services/accounts/internal/server/soajwtissuer"
	"time"

	pb "soa-socialnetwork/services/accounts/proto"
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

type userData struct {
	accountId int
	profileId string
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

func (s *AccountsService) fetchUserData(ctx context.Context, uid userId) (userData, error) {
	sql := fmt.Sprintf(`
	SELECT a.id, p.profile_id, a.password
	FROM accounts AS a
	JOIN profiles AS p ON a.id = p.account_id
	WHERE a.%s = $1;
	`, uid.kind.asDbColumn())

	row := s.dbPool.QueryRow(ctx, sql, uid.value)
	var data userData

	if err := row.Scan(&data.accountId, &data.profileId, &data.password); err != nil {
		return userData{}, errors.New("user not found")
	}

	return data, nil
}

func (s *AccountsService) Authenticate(ctx context.Context, req *pb.AuthByPassword) (*pb.AuthResponse, error) {
	userData, err := s.fetchUserData(ctx, extractUserIdFromProto(req.UserId))
	if err != nil {
		return nil, err
	}

	if userData.password != req.Password {
		return nil, errors.New("passwords doesn't match")
	}

	token, err := s.jwtIssuer.Issue(soajwtissuer.PersonalData{
		AccountId: userData.accountId,
		ProfileId: userData.profileId,
	}, JWT_DEFAULT_TTL)
	if err != nil {
		log.Printf("cannot create jwt token: %v", err)
		return nil, err
	}

	return &pb.AuthResponse{
		Token: token,
	}, nil
}
