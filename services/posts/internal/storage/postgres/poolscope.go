package postgres

import (
	"context"
	"soa-socialnetwork/services/posts/internal/repos"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPoolDatabase(cfg ConnectionConfig) (poolDatabase, error) {
	pool, err := newPool(cfg)
	if err != nil {
		return poolDatabase{}, err
	}

	return poolDatabase{
		connPool: pool,
	}, nil
}

type poolDatabase struct {
	connPool connectionPool
}

func (o *poolDatabase) OpenConnection(ctx context.Context) (repos.Connection, error) {
	conn, err := o.connPool.Pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	return poolConnection{
		ctx:  ctx,
		conn: conn,
	}, nil
}

func (o *poolDatabase) BeginTransaction(ctx context.Context) (repos.Transaction, error) {
	conn, err := o.connPool.Pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		conn.Release()
		return nil, err
	}

	return transaction{
		ctx:  ctx,
		conn: conn,
		tx:   tx,
	}, nil
}

type poolConnection struct {
	ctx  context.Context
	conn *pgxpool.Conn
}

func (c poolConnection) Close() error {
	c.conn.Release()
	return nil
}

func (c poolConnection) Pages() repos.PagesRepository {
	return pagesRepo{
		ctx:   c.ctx,
		scope: c.conn,
	}
}
func (c poolConnection) Posts() repos.PostsRepository {
	return postsRepo{
		ctx:   c.ctx,
		scope: c.conn,
	}
}
func (c poolConnection) Comments() repos.CommentsRepository {
	return commentsRepo{
		ctx:   c.ctx,
		scope: c.conn,
	}
}
func (c poolConnection) Metrics() repos.MetricsRepository {
	return metricsRepo{
		ctx:   c.ctx,
		scope: c.conn,
	}
}
func (c poolConnection) Outbox() repos.OutboxRepository {
	return outboxRepo{
		ctx:   c.ctx,
		scope: c.conn,
	}
}

type transaction struct {
	ctx  context.Context
	conn *pgxpool.Conn
	tx   pgx.Tx
}

func (c transaction) Commit() error {
	return c.tx.Commit(c.ctx)
}

func (c transaction) Rollback() error {
	return c.tx.Rollback(c.ctx)
}

func (c transaction) Close() error {
	c.tx.Rollback(c.ctx)
	c.conn.Release()
	return nil
}

func (c transaction) Pages() repos.PagesRepository {
	return pagesRepo{
		ctx:   c.ctx,
		scope: c.tx,
	}
}
func (c transaction) Posts() repos.PostsRepository {
	return postsRepo{
		ctx:   c.ctx,
		scope: c.tx,
	}
}
func (c transaction) Comments() repos.CommentsRepository {
	return commentsRepo{
		ctx:   c.ctx,
		scope: c.tx,
	}
}
func (c transaction) Metrics() repos.MetricsRepository {
	return metricsRepo{
		ctx:   c.ctx,
		scope: c.tx,
	}
}
func (c transaction) Outbox() repos.OutboxRepository {
	return outboxRepo{
		ctx:   c.ctx,
		scope: c.tx,
	}
}
