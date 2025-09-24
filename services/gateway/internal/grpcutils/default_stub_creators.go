package grpcutils

import (
	"soa-socialnetwork/services/gateway/internal/query"

	accountsPb "soa-socialnetwork/services/accounts/proto"
	postsPb "soa-socialnetwork/services/posts/proto"
	statsPb "soa-socialnetwork/services/stats/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func defaultGrpcClient(target string, qp *query.Params) (*grpc.ClientConn, error) {
	return grpc.NewClient(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(NewCreds(qp)),
	)
}

type DefaultAccountsStubCreator struct {
}

func (f DefaultAccountsStubCreator) New(target string, qp *query.Params) (accountsPb.AccountsServiceClient, error) {
	client, err := defaultGrpcClient(target, qp)
	if err != nil {
		return nil, err
	}
	return accountsPb.NewAccountsServiceClient(client), nil
}

type DefaultPostsStubCreator struct {
}

func (f DefaultPostsStubCreator) New(target string, qp *query.Params) (postsPb.PostsServiceClient, error) {
	client, err := defaultGrpcClient(target, qp)
	if err != nil {
		return nil, err
	}
	return postsPb.NewPostsServiceClient(client), nil
}

type DefaultStatsStubCreator struct {
}

func (f DefaultStatsStubCreator) New(target string, qp *query.Params) (statsPb.StatsServiceClient, error) {
	client, err := defaultGrpcClient(target, qp)
	if err != nil {
		return nil, err
	}
	return statsPb.NewStatsServiceClient(client), nil
}
