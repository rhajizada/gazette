package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// RedisConfig holds Redis settings loaded from environment variables.
type RedisConfig struct {
	Addr     string `env:"GAZETTE_REDIS_ADDR" env-required`
	Username string `env:"GAZETTE_REDIS_USERNAME"`
	Password string `env:"GAZETTE_REDIS_PASSWORD"`
	DB       int    `env:"GAZETTE_REDIS_DB" env-default:"0"`
}

// PostgresConfig holds Postgres settings loaded from environment variables.
type PostgresConfig struct {
	Host     string `env:"GAZETTE_POSTGRES_HOST" env-required`
	Port     int    `env:"GAZETTE_POSTGRES_PORT" env-default:"5432"`
	User     string `env:"GAZETTE_POSTGRES_USER" env-required`
	Password string `env:"GAZETTE_POSTGRES_PASSWORD" env-required`
	DBName   string `env:"GAZETTE_POSTGRES_DBNAME" env-required`
	SSLMode  string `env:"GAZETTE_POSTGRES_SSLMODE" env-default:"disable"`
}

// ServerConfig holds server-related settings.
type ServerConfig struct {
	Port      int    `env:"GAZETTE_PORT" env-default:"8080"`
	SecretKey string `env:"GAZETTE_SECRET_KEY" env-required`
	Database  PostgresConfig
	Redis     RedisConfig
	OAuth     OAuthConfig
}

type OAuthConfig struct {
	ClientID     string `env:"GAZETTE_OAUTH_CLIENT_ID" env-required`
	ClientSecret string `env:"GAZETTE_OAUTH_CLIENT_SECRET" env-required`
	IssuerURL    string `env:"GAZETTE_OAUTH_ISSUER_URL" env-required`
	RedirectURL  string `env:"GAZETTE_OAUTH_REDIRECT_URL" env-required`
}

// WorkerConfig holds worker-related settings.
type WorkerConfig struct {
	Database PostgresConfig
	Redis    RedisConfig
}

// SchedulerConfig holds scheduler-related settings.
type SchedulerConfig struct {
	Redis             RedisConfig
	HeartbeatInterval time.Duration  `env:"GAZETTE_HEARTBEAT_INTERVAL" env-default:"30s"`
	Location          *time.Location `env:"GAZETTE_LOCATION" env-default:"UTC"`
}

// LoadServer populates ServerConfig from environment variables prefixed with GAZETTE_.
func LoadServer() (*ServerConfig, error) {
	var cfg ServerConfig
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// LoadWorker populates WorkerConfig from environment variables prefixed with GAZETTE_.
func LoadWorker() (*WorkerConfig, error) {
	var cfg WorkerConfig
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// LoadScheduler populates SchedulerConfig from environment variables prefixed with GAZETTE_.
func LoadScheduler() (*SchedulerConfig, error) {
	var cfg SchedulerConfig
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
