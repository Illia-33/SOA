package soagrpc

import (
	"soa-socialnetwork/services/gateway/internal/server/query"
)

type StubCreator[TStub any] interface {
	New(target string, qp *query.Params) (TStub, error)
}
