package clickhouse

import (
	"soa-socialnetwork/services/stats/pkg/models"
	"sync"
	"time"
)

func (s *testSuite) TestPostsViewsSimple() {
	for i := range 10 {
		err := s.db.PostsViews().Put(models.PostViewEvent{
			PostId:          1,
			ViewerAccountId: models.AccountId(i),
			Timestamp:       time.Now(),
		})

		s.Require().NoError(err, "cannot put post view")
	}

	cnt, err := s.db.PostsViews().GetCountForPost(1)
	s.Require().NoError(err, "cannot get count for post")
	s.Assert().EqualValues(10, cnt)
}

func (s *testSuite) TestPostsViewsConcurrent() {
	errs := make([]error, 100)
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := range 10 {
		go func(i int) {
			for j := range 10 {
				idx := 10*i + j
				errs[idx] = s.db.PostsViews().Put(models.PostViewEvent{
					PostId:          1,
					ViewerAccountId: models.AccountId(idx),
					Timestamp:       time.Now(),
				})
			}
			wg.Done()
		}(i)
	}

	wg.Wait()

	for _, err := range errs {
		s.Require().NoError(err)
	}

	cnt, err := s.db.PostsViews().GetCountForPost(1)
	s.Require().NoError(err)
	s.Assert().EqualValues(100, cnt)
}

func (s *testSuite) TestPostsViewsDayDynamics() {
	now := time.Now()
	for i := range 10 {
		err := s.db.PostsViews().Put(models.PostViewEvent{
			PostId:          1,
			ViewerAccountId: models.AccountId(i),
			Timestamp:       now.Add(-time.Duration(24*i) * time.Hour),
		})
		s.Require().NoError(err, "failed to put event")
	}

	dynamics, err := s.db.PostsViews().GetDynamicsForPost(1)
	s.Require().NoError(err, "cannot get post dynamics")
	s.Assert().Equal(10, len(dynamics))
	for _, dayStat := range dynamics {
		s.Assert().EqualValues(1, dayStat.Count)
	}
}
