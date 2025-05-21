package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type DatabaseConfig struct {
	Host     string
	Port     uint16
	User     string
	Password string
	Database string
}

type PostgreSQL struct {
	conn *pgx.Conn
}

// func New(ctx context.Context, config DatabaseConfig) (*PostgreSQL, error) {
func New(ctx context.Context, config DatabaseConfig) (*PostgreSQL, error) {
	const connectionStringTemplate = "postgres://%s:%s@%s:%d/%s?sslmode=disable"
	connectionString := fmt.Sprintf(connectionStringTemplate, config.User, config.Password, config.Host, config.Port, config.Database)

	connConfig, err := pgx.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	conn, err := pgx.ConnectConfig(ctx, connConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &PostgreSQL{
		conn: conn,
	}, nil
}
