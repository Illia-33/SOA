package grpcutils

import (
	"soa-socialnetwork/services/gateway/internal/query"
)

type StubCreator[TStub any] interface {
	New(target string, qp *query.Params) (TStub, error)
}
