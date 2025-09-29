package postgres

import (
	"context"
	"fmt"
	"soa-socialnetwork/services/accounts/internal/models"
	"sync"
)

func (s *testSuite) TestAccountsSimple() {
	ctx := context.Background()
	conn, err := s.db.OpenConnection(ctx)
	s.Require().NoError(err, "cannot create db connection")
	defer conn.Close()

	registrationData := models.RegistrationData{
		Login:       "login",
		Password:    "password",
		Email:       "email@mail.com",
		PhoneNumber: "+333333333333",
		Name:        "name",
		Surname:     "surname",
	}

	id, err := conn.Accounts().New(registrationData)
	s.Require().NoError(err)

	verifyCheckOk := func(params models.AccountParams, err error) {
		s.Assert().NoError(err)
		if err != nil {
			s.Assert().Equal(id, params.Id)
		}
	}

	verifyCheckOk(conn.Accounts().CheckPasswordByLogin(registrationData.Login, registrationData.Password))
	verifyCheckOk(conn.Accounts().CheckPasswordByLogin(registrationData.Login, registrationData.Password))
	verifyCheckOk(conn.Accounts().CheckPasswordByLogin(registrationData.Login, registrationData.Password))
}

func (s *testSuite) TestAccountNew() {
	ctx := context.Background()
	conn, err := s.db.OpenConnection(ctx)
	s.Require().NoError(err, "cannot create db connection")
	defer conn.Close()

	registrations := make([]models.RegistrationData, 100)
	for i := range registrations {
		registrations[i] = models.RegistrationData{
			Login:       fmt.Sprintf("login_%d", i),
			Password:    "password",
			Email:       fmt.Sprintf("email_%d@mail.com", i),
			PhoneNumber: fmt.Sprintf("+333333333%d", i),
			Name:        fmt.Sprintf("name_%d", i),
			Surname:     fmt.Sprintf("surname_%d", i),
		}
	}

	ids := make([]models.AccountId, 0, 100)
	for i := range registrations {
		id, err := conn.Accounts().New(registrations[i])
		s.Require().NoError(err)
		ids = append(ids, id)
	}

	for i, data := range registrations {
		sql := `
		WITH cte AS (
			SELECT id
			FROM accounts
			WHERE id = $1 AND login = $2 AND password = $3 AND email = $4 AND phone_number = $5
		)
		SELECT count(*) FROM cte;
		`

		globalConn := s.db.globalConn
		row := globalConn.QueryRow(ctx, sql, ids[i], data.Login, data.Password, data.Email, data.PhoneNumber)

		var cnt int
		err := row.Scan(&cnt)
		s.Require().NoError(err)
		s.Assert().Equal(1, cnt, "required row not found or there are duplicated")
	}
}

func (s *testSuite) TestAccountNewConcurrent() {
	ctx := context.Background()

	registrations := make([]models.RegistrationData, 100)
	for i := range registrations {
		registrations[i] = models.RegistrationData{
			Login:       fmt.Sprintf("login_%d", i),
			Password:    "password",
			Email:       fmt.Sprintf("email_%d@mail.com", i),
			PhoneNumber: fmt.Sprintf("+333333333%d", i),
			Name:        fmt.Sprintf("name_%d", i),
			Surname:     fmt.Sprintf("surname_%d", i),
		}
	}

	type newResult struct {
		id  models.AccountId
		err error
	}

	newResults := make([](newResult), 100)

	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := range 10 {
		conn, err := s.db.OpenConnection(ctx)
		s.Require().NoError(err, "cannot create db connection")
		go func(i int) {
			defer conn.Close()
			for j := range 10 {
				idx := 10*i + j
				id, err := conn.Accounts().New(registrations[idx])
				newResults[idx] = newResult{
					id:  id,
					err: err,
				}
			}
			wg.Done()
		}(i)
	}

	wg.Wait()

	for i := range registrations {
		s.Require().NoError(newResults[i].err, "failed to create user with registration data: %+v", registrations[i])
	}

	for i, data := range registrations {
		sql := `
		WITH cte AS (
			SELECT id
			FROM accounts
			WHERE id = $1 AND login = $2 AND password = $3 AND email = $4 AND phone_number = $5
		)
		SELECT count(*) FROM cte;
		`

		globalConn := s.db.globalConn
		row := globalConn.QueryRow(ctx, sql, newResults[i].id, data.Login, data.Password, data.Email, data.PhoneNumber)

		var cnt int
		err := row.Scan(&cnt)
		s.Require().NoError(err)
		s.Assert().Equal(1, cnt, "required row not found or there are duplicated")
	}
}

func (s *testSuite) TestAccountsDelete() {
	ctx := context.Background()
	conn, err := s.db.OpenConnection(ctx)
	s.Require().NoError(err, "cannot create db connection")
	defer conn.Close()

	registrationData := models.RegistrationData{
		Login:       "login",
		Password:    "password",
		Email:       "email@mail.com",
		PhoneNumber: "+333333333333",
		Name:        "name",
		Surname:     "surname",
	}

	id, err := conn.Accounts().New(registrationData)
	s.Require().NoError(err)

	{
		params, err := conn.Accounts().CheckPasswordByLogin(registrationData.Login, registrationData.Password)
		s.Assert().NoError(err)
		if err != nil {
			s.Assert().Equal(id, params.Id)
		}
	}

	err = conn.Accounts().Delete(id)
	s.Require().NoError(err)

	{
		_, err := conn.Accounts().CheckPasswordByLogin(registrationData.Login, registrationData.Password)
		s.Require().Error(err)
	}
}
