package postgres

import (
	"context"
	"soa-socialnetwork/services/posts/internal/models"
	"sync"
)

func (s *testSuite) TestPagesConcurrent() {
	ctx := context.Background()
	accountId := models.AccountId(111)

	errs := make([]error, 10)

	barrier := sync.WaitGroup{}
	wg := sync.WaitGroup{}
	barrier.Add(1)
	wg.Add(10)
	for i := range 10 {
		conn, err := s.db.OpenConnection(ctx)
		s.Require().NoError(err)

		go func(i int) {
			barrier.Wait()

			_, err := conn.Pages().GetByAccountId(accountId)
			errs[i] = err

			wg.Done()
		}(i)
	}

	barrier.Done()
	wg.Wait()

	for _, err := range errs {
		s.Require().NoError(err)
	}
}
