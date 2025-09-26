package repo

import "context"

type RepoProvider interface {
	Accounts() AccountsRepo
	Profiles() ProfilesRepo
	ApiTokens() ApiTokensRepo
	Outbox() OutboxRepo
}

type Transaction interface {
	RepoProvider

	Commit() error
	Rollback() error
	Close() error
}

type Connection interface {
	RepoProvider

	Close() error
}

type Database interface {
	OpenConnection(context.Context) (Connection, error)
	BeginTransaction(context.Context) (Transaction, error)
}
