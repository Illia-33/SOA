package server

import (
	accountPb "soa-socialnetwork/services/accounts/proto"
	"soa-socialnetwork/services/gateway/internal/server/query"
	"soa-socialnetwork/services/gateway/internal/server/soagrpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type defaultAccountsStubFactory struct {
}

func (f defaultAccountsStubFactory) New(target string, qp *query.Params) (accountPb.AccountsServiceClient, error) {
	client, err := grpc.NewClient(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(soagrpc.NewCreds(qp)),
	)

	if err != nil {
		return nil, err
	}

	return accountPb.NewAccountsServiceClient(client), nil
}
