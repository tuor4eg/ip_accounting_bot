package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string `env:"TELEGRAM_TOKEN"`
	LogLevel      string `env:"LOG_LEVEL"`
	LogFormat     string `env:"LOG_FORMAT"`
	DatabaseURL   string `env:"DATABASE_URL"`
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	// Get log level from environment, default to "info" if not set
	logLevel := os.Getenv("LOG_LEVEL")
	logFormat := os.Getenv("LOG_FORMAT")

	if logLevel == "" {
		logLevel = "info"
	}

	if logFormat == "" {
		logFormat = "json"
	}

	c := &Config{
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
		LogLevel:      logLevel,
		LogFormat:     logFormat,
		DatabaseURL:   os.Getenv("DATABASE_URL"),
	}

	if c.TelegramToken == "" {
		return nil, errors.New("TELEGRAM_TOKEN is not set")
	}

	return c, nil
}
