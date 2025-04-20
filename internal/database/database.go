package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rhajizada/gazette/internal/config"
)

// CreatePool creates a *pgxpool.Pool instance using pgx/v5.
func CreatePool(c *config.PostgresConfig) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}
	return pool, nil
}
