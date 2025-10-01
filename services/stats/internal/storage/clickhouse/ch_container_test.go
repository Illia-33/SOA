package clickhouse

import (
	"context"
	"soa-socialnetwork/services/common/envvar"
	"soa-socialnetwork/services/stats/internal/repo"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

type clickhouseTestDb struct {
	cont         *testcontainers.DockerContainer
	underlyingDb *Database
}

func newClickhouseTestDb(t *testing.T) clickhouseTestDb {
	ctx := context.Background()
	migrations := envvar.MustStringFromEnv("DB_MIGRATIONS")

	cont, err := testcontainers.Run(
		ctx,
		"clickhouse/clickhouse-server:25.6",
		testcontainers.WithEnv(
			map[string]string{
				"CLICKHOUSE_USER":     "user",
				"CLICKHOUSE_PASSWORD": "password",
			},
		),
		testcontainers.WithExposedPorts("9000"),
		testcontainers.WithMounts(
			testcontainers.BindMount(
				migrations,
				"/docker-entrypoint-initdb.d",
			),
		),
	)
	require.NoError(t, err, "cannot create clickhouse container")

	hostname, err := cont.Host(ctx)
	require.NoError(t, err, "cannot get container hostname")

	port := func() int {
		natPort, err := cont.MappedPort(ctx, "9000/tcp")
		require.NoError(t, err, "cannot get container exposed port")

		portStr := string(natPort)[:strings.Index(string(natPort), "/")]
		port, err := strconv.Atoi(portStr)
		require.NoError(t, err, "cannot convert port")

		return port
	}()

	time.Sleep(5 * time.Second)

	db, err := NewDB(Config{
		Hostname:     hostname,
		Port:         port,
		Database:     "default",
		Username:     "user",
		Password:     "password",
		MaxIdleConns: 10,
		MaxOpenConns: 5,
	})
	require.NoError(t, err, "cannot init database")

	return clickhouseTestDb{
		cont:         cont,
		underlyingDb: &db,
	}
}

func (d *clickhouseTestDb) clean(t *testing.T) {
	err := d.underlyingDb.connection.Exec(
		context.Background(),
		`
	TRUNCATE ALL TABLES FROM default;
	`,
	)
	require.NoError(t, err, "cannot perform cleanup")
}

func (d *clickhouseTestDb) PostsViews() repo.PostsViewsRepo {
	return d.underlyingDb.PostsViews(context.Background())
}

func (d *clickhouseTestDb) PostsLikes() repo.PostsLikesRepo {
	return d.underlyingDb.PostsLikes(context.Background())
}

func (d *clickhouseTestDb) PostsComments() repo.PostsCommentsRepo {
	return d.underlyingDb.PostsComments(context.Background())
}

func (d *clickhouseTestDb) Registrations() repo.RegistrationsRepo {
	return d.underlyingDb.Registrations(context.Background())
}

func (d *clickhouseTestDb) Posts() repo.PostsRepo {
	return d.underlyingDb.Posts(context.Background())
}

func (d *clickhouseTestDb) Aggregation() repo.AggregationRepo {
	return d.underlyingDb.Aggregation(context.Background())
}
