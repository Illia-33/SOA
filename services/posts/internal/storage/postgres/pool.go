package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type connectionPool struct {
	*pgxpool.Pool
}

func newPool(cfg ConnectionConfig) (connectionPool, error) {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=5432 dbname=posts-postgres sslmode=disable pool_max_conns=%d", cfg.User, cfg.Password, cfg.Host, cfg.PoolSize)
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return connectionPool{}, nil
	}

	return connectionPool{
		Pool: pool,
	}, nil
}
