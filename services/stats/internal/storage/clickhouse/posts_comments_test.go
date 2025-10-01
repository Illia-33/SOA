package clickhouse

import (
	"soa-socialnetwork/services/stats/pkg/models"
	"sync"
	"time"
)

func (s *testSuite) TestPostsCommentsSimple() {
	for i := range 10 {
		err := s.db.PostsComments().Put(models.PostCommentEvent{
			PostId:          1,
			AuthorAccountId: models.AccountId(i),
			Timestamp:       time.Now(),
		})

		s.Require().NoError(err, "cannot put event")
	}

	cnt, err := s.db.PostsComments().GetCountForPost(1)
	s.Require().NoError(err, "cannot get count for post")
	s.Assert().EqualValues(10, cnt)
}

func (s *testSuite) TestPostsCommentsConcurrent() {
	errs := make([]error, 100)
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := range 10 {
		go func(i int) {
			for j := range 10 {
				idx := 10*i + j
				errs[idx] = s.db.PostsComments().Put(models.PostCommentEvent{
					PostId:          1,
					AuthorAccountId: models.AccountId(idx),
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

	cnt, err := s.db.PostsComments().GetCountForPost(1)
	s.Require().NoError(err)
	s.Assert().EqualValues(100, cnt)
}

func (s *testSuite) TestPostsCommentsDayDynamics() {
	now := time.Now()
	for i := range 10 {
		err := s.db.PostsComments().Put(models.PostCommentEvent{
			PostId:          1,
			AuthorAccountId: models.AccountId(i),
			Timestamp:       now.Add(-time.Duration(24*i) * time.Hour),
		})
		s.Require().NoError(err, "failed to put event")
	}

	dynamics, err := s.db.PostsComments().GetDynamicsForPost(1)
	s.Require().NoError(err, "cannot get post dynamics")
	s.Assert().Equal(10, len(dynamics))
	for _, dayStat := range dynamics {
		s.Assert().EqualValues(1, dayStat.Count)
	}
}
