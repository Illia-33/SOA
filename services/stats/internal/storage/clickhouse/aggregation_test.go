package clickhouse

import (
	"fmt"
	"soa-socialnetwork/services/stats/pkg/models"
	"time"
)

func (s *testSuite) TestAggregationUsers() {
	for i := 1; i <= 100; i++ {
		err := s.db.Registrations().Put(models.RegistrationEvent{
			AccountId: models.AccountId(i),
			ProfileId: fmt.Sprintf("profile-id-%d", i),
			Timestamp: time.Now(),
		})
		s.Require().NoError(err, "cannot put registration")
	}

	viewerId := models.AccountId(101)
	{
		err := s.db.Registrations().Put(models.RegistrationEvent{
			AccountId: viewerId,
			ProfileId: "viewer",
			Timestamp: time.Now(),
		})
		s.Require().NoError(err, "cannot put viewer registration")
	}

	for i := 1; i <= 100; i++ {
		err := s.db.Posts().Put(models.PostEvent{
			PostId:    models.PostId(i),
			AuthorId:  models.AccountId(i),
			Timestamp: time.Now(),
		})
		s.Require().NoError(err, "cannot put post")
	}

	for i := 1; i <= 100; i++ {
		views := make([]models.PostViewEvent, i)
		for j := range i {
			views[j] = models.PostViewEvent{
				PostId:          models.PostId(i),
				ViewerAccountId: viewerId,
				Timestamp:       time.Now(),
			}
		}

		err := s.db.PostsViews().Put(views...)
		s.Require().NoError(err, "cannot put views")
	}

	top10Users, err := s.db.Aggregation().GetTop10UsersByMetric(models.METRIC_VIEW_COUNT)
	s.Require().NoError(err, "cannot get top 10 users")

	s.Require().Equal(10, len(top10Users))
	for i, userStats := range top10Users {
		s.Assert().EqualValues(100-i, userStats.AccountId)
		s.Assert().EqualValues(100-i, userStats.MetricValue)
	}
}

func (s *testSuite) TestAggregationPosts() {
	posterId := models.AccountId(100)
	{
		err := s.db.Registrations().Put(models.RegistrationEvent{
			AccountId: posterId,
			ProfileId: "poster",
			Timestamp: time.Now(),
		})
		s.Require().NoError(err, "cannot put poster registration")
	}

	viewerId := models.AccountId(101)
	{
		err := s.db.Registrations().Put(models.RegistrationEvent{
			AccountId: viewerId,
			ProfileId: "viewer",
			Timestamp: time.Now(),
		})
		s.Require().NoError(err, "cannot put viewer registration")
	}

	for i := 1; i <= 100; i++ {
		err := s.db.Posts().Put(models.PostEvent{
			PostId:    models.PostId(i),
			AuthorId:  posterId,
			Timestamp: time.Now(),
		})
		s.Require().NoError(err, "cannot put post")
	}

	for i := 1; i <= 100; i++ {
		views := make([]models.PostViewEvent, i)
		for j := range i {
			views[j] = models.PostViewEvent{
				PostId:          models.PostId(i),
				ViewerAccountId: viewerId,
				Timestamp:       time.Now(),
			}
		}

		err := s.db.PostsViews().Put(views...)
		s.Require().NoError(err, "cannot put views")
	}

	top10Posts, err := s.db.Aggregation().GetTop10PostsByMetric(models.METRIC_VIEW_COUNT)
	s.Require().NoError(err, "cannot get top 10 users")

	s.Require().Equal(10, len(top10Posts))
	for i, postStats := range top10Posts {
		s.Assert().EqualValues(100-i, postStats.PostId)
		s.Assert().EqualValues(100-i, postStats.MetricValue)
	}
}
