package server

import (
	accountsPb "soa-socialnetwork/services/accounts/proto"
	"soa-socialnetwork/services/gateway/internal/server/query"
	"soa-socialnetwork/services/gateway/internal/server/soagrpc"
	postsPb "soa-socialnetwork/services/posts/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func defaultGrpcClient(target string, qp *query.Params) (*grpc.ClientConn, error) {
	return grpc.NewClient(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(soagrpc.NewCreds(qp)),
	)
}

type defaultAccountsStubFactory struct {
}

func (f defaultAccountsStubFactory) New(target string, qp *query.Params) (accountsPb.AccountsServiceClient, error) {
	client, err := defaultGrpcClient(target, qp)
	if err != nil {
		return nil, err
	}
	return accountsPb.NewAccountsServiceClient(client), nil
}

type defaultPostsStubFactory struct {
}

func (f defaultPostsStubFactory) New(target string, qp *query.Params) (postsPb.PostsServiceClient, error) {
	client, err := defaultGrpcClient(target, qp)
	if err != nil {
		return nil, err
	}
	return postsPb.NewPostsServiceClient(client), nil
}
