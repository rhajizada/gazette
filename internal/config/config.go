package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

// RedisConfig holds Redis settings loaded from environment variables.
type RedisConfig struct {
	Addr     string `env:"GAZETTE_REDIS_ADDR,notEmpty"`
	Username string `env:"GAZETTE_REDIS_USERNAME"`
	Password string `env:"GAZETTE_REDIS_PASSWORD"`
	DB       int    `env:"GAZETTE_REDIS_DB" envDefault:"0"`
}

// PostgresConfig holds Postgres settings loaded from environment variables.
type PostgresConfig struct {
	Host     string `env:"GAZETTE_POSTGRES_HOST,notEmpty"`
	Port     int    `env:"GAZETTE_POSTGRES_PORT" envDefault:"5432"`
	User     string `env:"GAZETTE_POSTGRES_USER,notEmpty"`
	Password string `env:"GAZETTE_POSTGRES_PASSWORD,notEmpty"`
	DBName   string `env:"GAZETTE_POSTGRES_DBNAME,notEmpty"`
	SSLMode  string `env:"GAZETTE_POSTGRES_SSLMODE" envDefault:"disable"`
}

// ServerConfig holds server-related settings.
type ServerConfig struct {
	Port      int    `env:"GAZETTE_PORT" envDefault:"8080"`
	SecretKey string `env:"GAZETTE_SECRET_KEY,notEmpty"`
	Database  PostgresConfig
	Redis     RedisConfig
	OAuth     OAuthConfig
}

// OAuthConfig holds OAuth provider settings.
type OAuthConfig struct {
	ClientID     string `env:"GAZETTE_OAUTH_CLIENT_ID,notEmpty"`
	ClientSecret string `env:"GAZETTE_OAUTH_CLIENT_SECRET,notEmpty"`
	IssuerURL    string `env:"GAZETTE_OAUTH_ISSUER_URL,notEmpty"`
	RedirectURL  string `env:"GAZETTE_OAUTH_REDIRECT_URL,notEmpty"`
}

// OllamaConfig holds ollama settings.
type OllamaConfig struct {
	BaseUrl         string `env:"GAZETTE_OLLAMA_URL,notEmpty"`
	EmbeddingsModel string `env:"GAZETTE_OLLAMA_EMBEDDINGS_MODEL,notEmpty"`
}

// QueuesConfig holds worker queue settings.
type QueuesConfig struct {
	Critical int `env:"GAZETTE_CRITICAL_QUEUES_COUNT" envDefault:"4"`
	Default  int `env:"GAZETTE_DEFAULT_QUEUES_COUNT" envDefault:"2"`
	Low      int `env:"GAZETTE_LOW_QUEUES_COUNT" envDefault:"2"`
}

// WorkerConfig holds worker-related settings.
type WorkerConfig struct {
	Database PostgresConfig
	Redis    RedisConfig
	Ollama   OllamaConfig
	Queues   QueuesConfig
}

// SchedulerConfig holds scheduler-related settings.
type SchedulerConfig struct {
	Redis             RedisConfig
	HeartbeatInterval time.Duration  `env:"GAZETTE_HEARTBEAT_INTERVAL" envDefault:"30s"`
	Location          *time.Location `env:"GAZETTE_LOCATION" envDefault:"UTC"`
}

// LoadServer populates ServerConfig from environment variables prefixed with GAZETTE_.
func LoadServer() (*ServerConfig, error) {
	var cfg ServerConfig
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// LoadWorker populates WorkerConfig from environment variables prefixed with GAZETTE_.
func LoadWorker() (*WorkerConfig, error) {
	var cfg WorkerConfig
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// LoadScheduler populates SchedulerConfig from environment variables prefixed with GAZETTE_.
func LoadScheduler() (*SchedulerConfig, error) {
	var cfg SchedulerConfig
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
