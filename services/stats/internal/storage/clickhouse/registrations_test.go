package clickhouse

import (
	"soa-socialnetwork/services/stats/pkg/models"
	"time"
)

func (s *testSuite) TestRegistrationsSimple() {
	err := s.db.Registrations().Put(models.RegistrationEvent{
		AccountId: 111,
		ProfileId: "some_profile_id",
		Timestamp: time.Now(),
	})
	s.Require().NoError(err)
}
