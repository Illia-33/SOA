package postgres

import (
	"context"
	"fmt"
	"soa-socialnetwork/services/common/envvar"
	"soa-socialnetwork/services/posts/internal/repo"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const test_postgres_user = "test_postgres_user"
const test_postgres_password = "test_postgres_password"

type testPgContainerDatabase struct {
	cont       *testcontainers.DockerContainer
	globalConn *pgx.Conn

	hostname string
	port     int
}

func newTestPostgresDatabase(t *testing.T) testPgContainerDatabase {
	pathToMigrations := envvar.MustStringFromEnv("DB_MIGRATIONS")
	ctx := context.Background()

	cont, err := testcontainers.Run(
		ctx, "postgres:17-alpine",
		testcontainers.WithExposedPorts("5432"),
		testcontainers.WithEnv(map[string]string{
			"POSTGRES_USER":     test_postgres_user,
			"POSTGRES_PASSWORD": test_postgres_password,
		}),
		testcontainers.WithMounts(
			testcontainers.BindMount(pathToMigrations, testcontainers.ContainerMountTarget("/docker-entrypoint-initdb.d")),
		),
		testcontainers.WithAdditionalWaitStrategy(
			wait.ForExec([]string{"pg_isready"}).
				WithExitCode(0).
				WithStartupTimeout(3*time.Second),
		),
	)
	require.NoError(t, err, "cannot run test postgres container")

	hostname, err := cont.Host(ctx)
	require.NoError(t, err, "cannot get postgres container hostname")

	port := func() int {
		natPort, err := cont.MappedPort(ctx, "5432/tcp")
		require.NoError(t, err, "cannot get container exposed port")

		portStr := string(natPort)[:strings.Index(string(natPort), "/")]
		port, err := strconv.Atoi(portStr)
		require.NoError(t, err, "cannot convert port")

		return port
	}()

	time.Sleep(3 * time.Second)

	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%d sslmode=disable", test_postgres_user, test_postgres_password, hostname, port)
	conn, err := pgx.Connect(ctx, connStr)
	require.NoError(t, err, "cannot create connection to database")

	return testPgContainerDatabase{
		cont:       cont,
		globalConn: conn,
		hostname:   hostname,
		port:       port,
	}
}

func (d *testPgContainerDatabase) newConn() (*pgx.Conn, error) {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%d sslmode=disable", test_postgres_user, test_postgres_password, d.hostname, d.port)
	conn, err := pgx.Connect(context.Background(), connStr)
	return conn, err
}

func (d *testPgContainerDatabase) cleanup(t *testing.T) {
	conn := d.globalConn
	_, err := conn.Exec(context.Background(), `
		TRUNCATE TABLE pages;
		TRUNCATE TABLE posts;
		TRUNCATE TABLE comments;
		TRUNCATE TABLE likes;
		TRUNCATE TABLE outbox;
	`)

	require.NoError(t, err, "database cleanup failed")
}

func (d *testPgContainerDatabase) OpenConnection(context.Context) (repo.Connection, error) {
	conn, err := d.newConn()
	if err != nil {
		return nil, err
	}
	return &testConnection{
		testRepoProvider: testRepoProvider{
			scope: conn,
		},
		conn: conn,
	}, nil
}

func (d *testPgContainerDatabase) BeginTransaction(context.Context) (repo.Transaction, error) {
	conn, err := d.newConn()
	if err != nil {
		return nil, err
	}

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return nil, err
	}

	return &testTransaction{
		testRepoProvider: testRepoProvider{
			scope: tx,
		},
		conn: conn,
		tx:   tx,
	}, nil
}

type testRepoProvider struct {
	scope pgxScope
}

func (p *testRepoProvider) Pages() repo.PagesRepository {
	return pagesRepo{
		ctx:   context.Background(),
		scope: p.scope,
	}
}

func (p *testRepoProvider) Posts() repo.PostsRepository {
	return postsRepo{
		ctx:   context.Background(),
		scope: p.scope,
	}
}

func (p *testRepoProvider) Comments() repo.CommentsRepository {
	return commentsRepo{
		ctx:   context.Background(),
		scope: p.scope,
	}
}

func (p *testRepoProvider) Metrics() repo.MetricsRepository {
	return metricsRepo{
		ctx:   context.Background(),
		scope: p.scope,
	}
}

func (p *testRepoProvider) Outbox() repo.OutboxRepository {
	return outboxRepo{
		ctx:   context.Background(),
		scope: p.scope,
	}
}

type testConnection struct {
	testRepoProvider

	conn *pgx.Conn
}

func (c *testConnection) Close() error {
	return c.conn.Close(context.Background())
}

type testTransaction struct {
	testRepoProvider

	conn *pgx.Conn
	tx   pgx.Tx
}

func (c *testTransaction) Close() error {
	c.tx.Rollback(context.Background())
	return c.conn.Close(context.Background())
}

func (c *testTransaction) Commit() error {
	return c.tx.Commit(context.Background())
}

func (c *testTransaction) Rollback() error {
	return c.tx.Rollback(context.Background())
}
