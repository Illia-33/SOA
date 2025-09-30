package postgres

import (
	"context"
	"fmt"
	"soa-socialnetwork/services/accounts/internal/models"
	"sync"

	"github.com/google/uuid"
)

func (s *testSuite) TestProfilesSimple() {
	ctx := context.Background()
	conn, err := s.db.OpenConnection(ctx)
	s.Require().NoError(err)

	registrationData := models.RegistrationData{
		Login:       "login",
		Password:    "password",
		Email:       "email@mail.com",
		PhoneNumber: "+333333333333",
		Name:        "name",
		Surname:     "surname",
	}

	accountId, err := conn.Accounts().New(registrationData)
	s.Require().NoError(err)

	profileId := models.ProfileId(uuid.NewString())
	err = conn.Profiles().New(profileId, accountId, registrationData)
	s.Require().NoError(err)

	{
		profileData, err := conn.Profiles().GetByProfileId(profileId)
		s.Require().NoError(err)

		s.Assert().Equal(accountId, profileData.AccountId)
		s.Assert().Equal(profileId, profileData.ProfileId)
		s.Assert().Equal(registrationData.Name, profileData.Name)
		s.Assert().Equal(registrationData.Surname, profileData.Surname)
	}

	{
		profileData, err := conn.Profiles().GetByAccountId(accountId)
		s.Require().NoError(err)

		s.Assert().Equal(accountId, profileData.AccountId)
		s.Assert().Equal(profileId, profileData.ProfileId)
		s.Assert().Equal(registrationData.Name, profileData.Name)
		s.Assert().Equal(registrationData.Surname, profileData.Surname)
	}
}

func (s *testSuite) TestProfilesResolving() {
	ctx := context.Background()
	conn, err := s.db.OpenConnection(ctx)
	s.Require().NoError(err)

	registrationData := models.RegistrationData{
		Login:       "login",
		Password:    "password",
		Email:       "email@mail.com",
		PhoneNumber: "+333333333333",
		Name:        "name",
		Surname:     "surname",
	}

	accountId, err := conn.Accounts().New(registrationData)
	s.Require().NoError(err)

	profileId := models.ProfileId(uuid.NewString())
	err = conn.Profiles().New(profileId, accountId, registrationData)
	s.Require().NoError(err)

	{
		resolvedAccountId, err := conn.Profiles().ResolveProfileId(profileId)
		s.Require().NoError(err)
		s.Require().Equal(accountId, resolvedAccountId)
	}

	{
		resolvedProfileId, err := conn.Profiles().ResolveAccountId(accountId)
		s.Require().NoError(err)
		s.Require().Equal(profileId, resolvedProfileId)
	}
}

func (s *testSuite) TestProfilesConcurrent() {
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
		accountId models.AccountId
		profileId models.ProfileId
		err       error
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
				accountId, err := conn.Accounts().New(registrations[idx])
				if err != nil {
					newResults[idx].err = err
					continue
				}

				profileId := models.ProfileId(uuid.NewString())
				err = conn.Profiles().New(profileId, accountId, registrations[idx])
				if err != nil {
					newResults[idx].err = err
					continue
				}

				newResults[idx] = newResult{
					accountId: accountId,
					profileId: profileId,
					err:       nil,
				}
			}
			wg.Done()
		}(i)
	}

	wg.Wait()

	for i := range registrations {
		s.Require().NoError(newResults[i].err, "failed to create user with registration data: %+v", registrations[i])
	}

	conn, err := s.db.OpenConnection(ctx)
	s.Require().NoError(err)

	for i, result := range newResults {
		data, err := conn.Profiles().GetByProfileId(result.profileId)
		s.Require().NoError(err)

		s.Assert().Equal(result.accountId, data.AccountId)
		s.Assert().Equal(result.profileId, data.ProfileId)
		s.Assert().Equal(registrations[i].Name, data.Name)
		s.Assert().Equal(registrations[i].Surname, data.Surname)
	}
}
