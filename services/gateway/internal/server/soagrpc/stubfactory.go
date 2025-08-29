package soagrpc

import (
	accountsPb "soa-socialnetwork/services/accounts/proto"
	"soa-socialnetwork/services/gateway/internal/server/query"
	postsPb "soa-socialnetwork/services/posts/proto"
)

type AccountsStubFactory interface {
	New(target string, qp *query.Params) (accountsPb.AccountsServiceClient, error)
}

type PostsStubFactory interface {
	New(target string, qp *query.Params) (postsPb.PostsServiceClient, error)
}
