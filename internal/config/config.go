package config

import (
	"github.com/syntaqx/env"
)

type ServerConfig struct {
	Port      int            `env:"GAZETTE_PORT,default=8080"`
	SecretKey string         `env:"GAZETTE_SECRET_KEY,required"`
	Database  PostgresConfig `env:"GAZETTE_POSTGRES"`
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

// Load loads the server configuration from environment variables.
func Load() (ServerConfig, error) {
	var cfg ServerConfig
	err := env.Unmarshal(&cfg)
	return cfg, err
}
