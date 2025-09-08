package postgres

type ConnectionConfig struct {
	Host     string
	User     string
	Password string
	PoolSize int
}
