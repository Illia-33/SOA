package soagrpc

import (
	accountPb "soa-socialnetwork/services/accounts/proto"
	"soa-socialnetwork/services/gateway/internal/server/query"
)

type AccountsStubFactory interface {
	New(target string, qp *query.Params) (accountPb.AccountsServiceClient, error)
}
