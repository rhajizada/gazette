package database

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rhajizada/gazette/internal/config"
)

// CreatePool creates a *pgxpool.Pool instance using pgx/v5.
func CreatePool(cfg *config.PostgresConfig) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}
	return pool, nil
}

// CreateRedisClient create Redis Client connection for asynq
func CreateRedisClient(cfg *config.QueueConfig) *asynq.RedisClientOpt {
	conn := asynq.RedisClientOpt{
		Addr:     cfg.Addr,
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.DB,
	}
	return &conn
}
