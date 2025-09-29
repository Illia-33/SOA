package postgres

import (
	"context"
	"fmt"
	"soa-socialnetwork/services/accounts/internal/models"
	"soa-socialnetwork/services/accounts/internal/repo"
	"sync"
	"time"
)

func (s *testSuite) TestApiTokensSimple() {
	ctx := context.Background()
	conn, err := s.db.OpenConnection(ctx)
	s.Require().NoError(err)

	token := models.ApiToken("some_api_token")
	tokenParams := repo.ApiTokenParams{
		AccountId:   111,
		ReadAccess:  true,
		WriteAccess: false,
		Ttl:         time.Hour,
	}

	validUntil, err := conn.ApiTokens().Put(token, tokenParams)
	s.Require().NoError(err)

	tokenData, err := conn.ApiTokens().Get(token)
	s.Require().NoError(err)

	s.Assert().Equal(tokenParams.AccountId, tokenData.AccountId)
	s.Assert().Equal(tokenParams.ReadAccess, tokenData.ReadAccess)
	s.Assert().Equal(tokenParams.WriteAccess, tokenData.WriteAccess)
	s.Assert().Equal(validUntil, tokenData.ValidUntil)
}

func (s *testSuite) TestApiTokensConcurrent() {
	ctx := context.Background()

	type result struct {
		validUntil time.Time
		err        error
	}

	results := make([]result, 100)
	tokensParams := make([]repo.ApiTokenParams, 100)

	for i := range 100 {
		tokensParams[i] = repo.ApiTokenParams{
			AccountId:   i,
			ReadAccess:  true,
			WriteAccess: false,
			Ttl:         time.Hour,
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := range 10 {
		conn, err := s.db.OpenConnection(ctx)
		s.Require().NoError(err)
		go func(i int) {
			defer conn.Close()
			for j := range 10 {
				idx := 10*i + j
				token := models.ApiToken(fmt.Sprintf("some_api_token_%d", idx))
				validUntil, err := conn.ApiTokens().Put(token, tokensParams[idx])
				results[idx] = result{
					validUntil: validUntil,
					err:        err,
				}
			}

			wg.Done()
		}(i)
	}

	wg.Wait()

	for _, result := range results {
		s.Require().NoError(result.err)
	}

	conn, err := s.db.OpenConnection(ctx)
	s.Require().NoError(err)

	for i, result := range results {
		tokenData, err := conn.ApiTokens().Get(models.ApiToken(fmt.Sprintf("some_api_token_%d", i)))
		s.Require().NoError(err)
		s.Assert().Equal(tokensParams[i].AccountId, tokenData.AccountId)
		s.Assert().Equal(tokensParams[i].ReadAccess, tokenData.ReadAccess)
		s.Assert().Equal(tokensParams[i].WriteAccess, tokenData.WriteAccess)
		s.Assert().Equal(result.validUntil, tokenData.ValidUntil)
	}
}
