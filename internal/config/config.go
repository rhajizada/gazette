package config

import (
	"github.com/syntaqx/env"
)

type ServerConfig struct {
	Port      int            `env:"GAZETTE_PORT,default=8080"`
	SecretKey string         `env:"GAZETTE_SECRET_KEY,required"`
	Database  PostgresConfig `env:"GAZETTE_POSTGRES"`
	Redis     RedisConfig    `env:"GAZETTE_REDIS"`
}

type WorkerConfig struct {
	Database PostgresConfig `env:"GAZETTE_POSTGRES"`
	Redis    RedisConfig    `env:"GAZETTE_REDIS"`
}

type OAuthConfig struct {
	ProviderURL  string `env:"PROVIDER_URL"`
	ProviderName string `env:"PROVIDER_NAME"`
	RedirectURL  string `env:"REDIRECT_URL"`
	ClientID     string `env:"CLIENT_ID"`
	ClientSecret string `env:"CLIENT_SECRET"`
}

type PostgresConfig struct {
	Host     string `env:"HOST,required"`
	Port     int    `env:"PORT,default=5432"`
	User     string `env:"USER,required"`
	Password string `env:"PASSWORD,required"`
	DBName   string `env:"DBNAME,required"`
	SSLMode  string `env:"SSLMODE,default=disable"`
}

type RedisConfig struct {
	Addr     string `env:"ADDR,required"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
	DB       int    `env:"DB,default=0"`
}

// LoadServer loads the server configuration from environment variables.
func LoadServer() (ServerConfig, error) {
	var cfg ServerConfig
	err := env.Unmarshal(&cfg)
	return cfg, err
}

// LoadWorker loads the server configuration from environment variables.
func LoadWorker() (WorkerConfig, error) {
	var cfg WorkerConfig
	err := env.Unmarshal(&cfg)
	return cfg, err
}
