package postgres

import (
	"context"
	"fmt"
	"soa-socialnetwork/services/posts/internal/repo"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDatabase(ctx context.Context, cfg Config) (database, error) {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=5432 dbname=posts-postgres sslmode=disable pool_max_conns=%d", cfg.User, cfg.Password, cfg.Host, cfg.PoolSize)
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return database{}, nil
	}

	if err != nil {
		return database{}, err
	}

	return database{
		connPool: pool,
	}, nil
}

type database struct {
	connPool *pgxpool.Pool
}

func (d *database) OpenConnection(ctx context.Context) (repo.Connection, error) {
	conn, err := d.connPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	return &connection{
		repoProvider: repoProvider{
			ctx:   ctx,
			scope: conn,
		},
		conn: conn,
	}, nil
}

func (d *database) BeginTransaction(ctx context.Context) (repo.Transaction, error) {
	conn, err := d.connPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		conn.Release()
		return nil, err
	}

	return &transaction{
		repoProvider: repoProvider{
			ctx:   ctx,
			scope: conn,
		},
		conn: conn,
		tx:   tx,
	}, nil
}
