package postgres

import (
	"context"
	"soa-socialnetwork/services/posts/internal/repo"

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

func (o *poolDatabase) OpenConnection(ctx context.Context) (repo.Connection, error) {
	conn, err := o.connPool.Pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	return poolConnection{
		ctx:  ctx,
		conn: conn,
	}, nil
}

func (o *poolDatabase) BeginTransaction(ctx context.Context) (repo.Transaction, error) {
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

func (c poolConnection) Pages() repo.PagesRepository {
	return pagesRepo{
		ctx:   c.ctx,
		scope: c.conn,
	}
}
func (c poolConnection) Posts() repo.PostsRepository {
	return postsRepo{
		ctx:   c.ctx,
		scope: c.conn,
	}
}
func (c poolConnection) Comments() repo.CommentsRepository {
	return commentsRepo{
		ctx:   c.ctx,
		scope: c.conn,
	}
}
func (c poolConnection) Metrics() repo.MetricsRepository {
	return metricsRepo{
		ctx:   c.ctx,
		scope: c.conn,
	}
}
func (c poolConnection) Outbox() repo.OutboxRepository {
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

func (c transaction) Pages() repo.PagesRepository {
	return pagesRepo{
		ctx:   c.ctx,
		scope: c.tx,
	}
}
func (c transaction) Posts() repo.PostsRepository {
	return postsRepo{
		ctx:   c.ctx,
		scope: c.tx,
	}
}
func (c transaction) Comments() repo.CommentsRepository {
	return commentsRepo{
		ctx:   c.ctx,
		scope: c.tx,
	}
}
func (c transaction) Metrics() repo.MetricsRepository {
	return metricsRepo{
		ctx:   c.ctx,
		scope: c.tx,
	}
}
func (c transaction) Outbox() repo.OutboxRepository {
	return outboxRepo{
		ctx:   c.ctx,
		scope: c.tx,
	}
}
