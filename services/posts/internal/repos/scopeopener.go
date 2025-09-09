package repos

import "context"

// RepositoryProvider groups accessors for all repositories that operate on
// the same underlying data source (e.g., a transaction or a shared connection).
type RepositoryProvider interface {
	Pages() PagesRepository
	Posts() PostsRepository
	Comments() CommentsRepository
	Metrics() MetricsRepository
	Outbox() OutboxRepository
}

// Transaction is a unit of work that scopes repositories to a single
// database transaction. Use Commit to persist changes or Rollback to cancel them.
// Either commit or rollback must be called to finish transaction properly.
// Close releases associated resources when the work is done.
type Transaction interface {
	RepositoryProvider

	Commit() error
	Rollback() error
	Close() error
}

// Connection scopes repositories to a shared, non-transactional connection.
// Close releases associated resources when the work is done.
type Connection interface {
	RepositoryProvider

	Close() error
}

// RepoScopeOpener constructs units of work that provide scoped repository access.
// BeginTransaction starts a transactional unit; OpenConnection opens a non-transactional scope.
// The provided context carries deadlines and cancellation to the underlying driver.
type RepoScopeOpener interface {
	OpenConnection(context.Context) (Connection, error)
	BeginTransaction(context.Context) (Transaction, error)
}
