package postgres

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite

	db testPgContainerDatabase
}

func (s *testSuite) SetupSuite() {
	s.db = newTestPostgresDatabase(s.T())
}

func (s *testSuite) AfterTest(suiteName, testName string) {
	s.db.cleanup(s.T())
}

func TestRunSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}
