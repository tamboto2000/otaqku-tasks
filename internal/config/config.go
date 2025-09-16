package config

import (
	"fmt"
	"log/slog"

	"github.com/tamboto2000/otaqku-tasks/pkg/config"
)

type Database struct {
	Host         string `env:"DB_HOST"`
	Port         string `env:"DB_PORT"`
	Username     string `env:"DB_USERNAME"`
	Password     string `env:"DB_PASSWORD"`
	Database     string `env:"DB_DATABASE"`
	SSLMode      string `env:"DB_SSL_MODE"`
	MigrationDir string `env:"DB_MIGRATION_DIR"`
}

type HTTPServer struct {
	Port string `env:"HTTP_SERVER_PORT"`
}

type Logging struct {
	Level      slog.Level `env:"LOGGING_LEVEL"`
	WithSource bool       `env:"LOGGING_WITH_SOURCE"`
}

type JWT struct {
	AccessTokenDuration  config.MinuteDuration   `env:"JWT_ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration config.MinuteDuration   `env:"JWT_REFRESH_TOKEN_DURATION"`
	SigningKey           config.RawBase64Encoded `env:"JWT_SIGNING_KEY"`
}

type Config struct {
	Database   Database
	HTTPServer HTTPServer
	Logging    Logging
	JWT        JWT
}

func LoadConfig() (Config, error) {
	var cfg Config
	if err := config.LoadEnv(&cfg); err != nil {
		return Config{}, fmt.Errorf("loading config error: %v", err)
	}

	return cfg, nil
}
