package clickhouse

import (
	"soa-socialnetwork/services/stats/pkg/models"
	"time"
)

func (s *testSuite) TestPostsSimple() {
	err := s.db.Posts().Put(models.PostEvent{
		PostId:    111,
		AuthorId:  1,
		Timestamp: time.Now(),
	})
	s.Require().NoError(err)
}
