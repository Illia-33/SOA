package clickhouse

import (
	"context"
	"fmt"
	"soa-socialnetwork/services/stats/internal/repo"

	ch "github.com/ClickHouse/clickhouse-go/v2"
	chDriver "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type Database struct {
	connection chDriver.Conn
}

type Config struct {
	Hostname     string
	Port         int
	Database     string
	Username     string
	Password     string
	MaxIdleConns int
	MaxOpenConns int
}

func NewDB(cfg Config) (Database, error) {
	conn, err := ch.Open(&ch.Options{
		Addr: []string{fmt.Sprintf("%s:%d", cfg.Hostname, cfg.Port)},
		Auth: ch.Auth{
			Database: cfg.Database,
			Username: cfg.Username,
			Password: cfg.Password,
		},
		MaxIdleConns: cfg.MaxIdleConns,
		MaxOpenConns: cfg.MaxOpenConns,
	})

	if err != nil {
		return Database{}, err
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return Database{}, nil
	}

	return Database{
		connection: conn,
	}, nil
}

func (d *Database) PostsViews(ctx context.Context) repo.PostsViewsRepo {
	return &postsViewsRepo{
		ctx:  ctx,
		conn: d.connection,
	}
}

func (d *Database) PostsLikes(ctx context.Context) repo.PostsLikesRepo {
	return &postsLikesRepo{
		ctx:  ctx,
		conn: d.connection,
	}
}

func (d *Database) PostsComments(ctx context.Context) repo.PostsCommentsRepo {
	return &postsCommentsRepo{
		ctx:  ctx,
		conn: d.connection,
	}
}

func (d *Database) Registrations(ctx context.Context) repo.RegistrationsRepo {
	return &registrationsRepo{
		ctx:  ctx,
		conn: d.connection,
	}
}

func (d *Database) Posts(ctx context.Context) repo.PostsRepo {
	return &postsRepo{
		ctx:  ctx,
		conn: d.connection,
	}
}

func (d *Database) Aggregation(ctx context.Context) repo.AggregationRepo {
	return &aggregationRepo{
		ctx:  ctx,
		conn: d.connection,
	}
}
