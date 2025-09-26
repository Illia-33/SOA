package postgres

import (
	"context"
	"fmt"
	"soa-socialnetwork/services/accounts/internal/repo"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	connPool *pgxpool.Pool
}

func NewDatabase(ctx context.Context, cfg Config) (Database, error) {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=5432 dbname=accounts-postgres sslmode=disable pool_max_conns=%d", cfg.User, cfg.Password, cfg.Host, cfg.PoolSize)
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return Database{}, err
	}

	return Database{
		connPool: pool,
	}, nil
}

func (d *Database) OpenConnection(ctx context.Context) (repo.Connection, error) {
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

func (d *Database) BeginTransaction(ctx context.Context) (repo.Transaction, error) {
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
			scope: tx,
		},
		conn: conn,
		tx:   tx,
	}, nil
}
