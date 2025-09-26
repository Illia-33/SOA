package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type connection struct {
	repoProvider

	conn *pgxpool.Conn
}

func (c *connection) Close() error {
	if c.conn != nil {
		c.conn.Release()
		c.conn = nil
	}
	return nil
}
