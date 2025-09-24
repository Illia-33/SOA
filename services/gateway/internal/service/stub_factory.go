package service

import (
	accountsPb "soa-socialnetwork/services/accounts/proto"
	"soa-socialnetwork/services/gateway/internal/grpcutils"
	"soa-socialnetwork/services/gateway/internal/query"
	postsPb "soa-socialnetwork/services/posts/proto"
	statsPb "soa-socialnetwork/services/stats/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func defaultGrpcClient(target string, qp *query.Params) (*grpc.ClientConn, error) {
	return grpc.NewClient(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(grpcutils.NewCreds(qp)),
	)
}

type defaultAccountsStubCreator struct {
}

func (f defaultAccountsStubCreator) New(target string, qp *query.Params) (accountsPb.AccountsServiceClient, error) {
	client, err := defaultGrpcClient(target, qp)
	if err != nil {
		return nil, err
	}
	return accountsPb.NewAccountsServiceClient(client), nil
}

type defaultPostsStubCreator struct {
}

func (f defaultPostsStubCreator) New(target string, qp *query.Params) (postsPb.PostsServiceClient, error) {
	client, err := defaultGrpcClient(target, qp)
	if err != nil {
		return nil, err
	}
	return postsPb.NewPostsServiceClient(client), nil
}

type defaultStatsStubCreator struct {
}

func (f defaultStatsStubCreator) New(target string, qp *query.Params) (statsPb.StatsServiceClient, error) {
	client, err := defaultGrpcClient(target, qp)
	if err != nil {
		return nil, err
	}
	return statsPb.NewStatsServiceClient(client), nil
}
