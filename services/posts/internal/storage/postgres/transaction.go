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

func (t *transaction) Commit() error {
	return t.tx.Commit(t.ctx)
}

func (t *transaction) Rollback() error {
	return t.tx.Rollback(t.ctx)
}

func (t *transaction) Close() error {
	if t.conn != nil {
		t.tx.Rollback(t.ctx)
		t.conn.Release()
		t.conn = nil
	}

	return nil
}
