package postgres

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type transaction struct {
	repoProvider

	conn *pgxpool.Conn
	tx   pgx.Tx
}

func (c *transaction) Commit() error {
	return c.tx.Commit(c.ctx)
}

func (c *transaction) Rollback() error {
	return c.tx.Rollback(c.ctx)
}

func (c *transaction) Close() error {
	if c.conn != nil {
		c.tx.Rollback(c.ctx)
		c.conn.Release()
		c.conn = nil
		c.tx = nil
	}
	return nil
}
