package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"soa-socialnetwork/services/accounts/internal/models"
	"soa-socialnetwork/services/accounts/internal/repo"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func marshalJson[T any](t *testing.T, obj T) []byte {
	data, err := json.Marshal(obj)
	require.NoError(t, err)
	return data
}

func unmarshalJson[T any](t *testing.T, data []byte) T {
	var obj T
	err := json.Unmarshal(data, &obj)
	require.NoError(t, err)
	return obj
}

type testPayload struct {
	Desc string `json:"desc"`
}

func (s *testSuite) TestOutboxSimple() {
	ctx := context.Background()
	conn, err := s.db.OpenConnection(ctx)
	s.Require().NoError(err)

	inputPayload := testPayload{
		Desc: "payload description",
	}

	inputEvent := models.OutboxEvent{
		Type:      "test_type",
		Payload:   marshalJson(s.T(), inputPayload),
		CreatedAt: time.Now(),
	}

	err = conn.Outbox().Put(inputEvent)
	s.Require().NoError(err)

	events, err := conn.Outbox().Fetch(repo.OutboxFetchParams{
		Limit: 100,
	})
	s.Require().NoError(err)
	s.Require().Equal(1, len(events))

	outputEvent := events[0]

	s.Assert().Equal(inputEvent.Type, outputEvent.Type)
	s.Assert().Equal(inputPayload.Desc, unmarshalJson[testPayload](s.T(), outputEvent.Payload).Desc)
}

func (s *testSuite) TestOutboxConcurrent() {
	ctx := context.Background()

	errs := make([]error, 100)
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := range 10 {
		conn, err := s.db.OpenConnection(ctx)
		s.Require().NoError(err)

		go func(i int) {
			for j := range 10 {
				idx := 10*i + j
				inputPayload := testPayload{
					Desc: fmt.Sprintf("%d", idx),
				}

				inputEvent := models.OutboxEvent{
					Type:      "test_type",
					Payload:   marshalJson(s.T(), inputPayload),
					CreatedAt: time.Now(),
				}

				errs[idx] = conn.Outbox().Put(inputEvent)
			}

			wg.Done()
		}(i)
	}

	wg.Wait()

	for _, err := range errs {
		s.Require().NoError(err)
	}

	conn, err := s.db.OpenConnection(ctx)
	s.Require().NoError(err)

	events, err := conn.Outbox().Fetch(repo.OutboxFetchParams{
		Limit: 100,
	})
	s.Require().NoError(err)
	s.Require().Equal(100, len(events))

	idxFound := make([]bool, 100)

	for _, event := range events {
		s.Assert().Equal("test_type", event.Type)
		desc := unmarshalJson[testPayload](s.T(), event.Payload).Desc
		idx, err := strconv.Atoi(desc)
		s.Require().NoError(err)
		s.Require().True(0 <= idx && idx < 100)
		s.Require().False(idxFound[idx])
		idxFound[idx] = true
	}
}
