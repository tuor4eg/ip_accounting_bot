package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
	"github.com/tuor4eg/ip_accounting_bot/internal/validate"
)

type Config struct {
	TelegramToken string `env:"TELEGRAM_TOKEN"`
	LogLevel      string `env:"LOG_LEVEL"`
	LogFormat     string `env:"LOG_FORMAT"`
	DatabaseURL   string `env:"DATABASE_URL"`
	HMACKey       string `env:"HMAC_KEY"`
	AEADKey       string `env:"AEAD_KEY"`
}

var (
	ErrTelegramTokenNotSet = errors.New("TELEGRAM_TOKEN is not set")
	ErrHMACKeyNotSet       = errors.New("HMAC_KEY is not set")
	ErrAEADKeyNotSet       = errors.New("AEAD_KEY is not set")
)

func Load() (*Config, error) {
	const op = "config.Load"

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
		HMACKey:       os.Getenv("HMAC_KEY"),
		AEADKey:       os.Getenv("AEAD_KEY"),
	}

	if c.TelegramToken == "" {
		return nil, validate.Wrap(op, ErrTelegramTokenNotSet)
	}
	if c.HMACKey == "" {
		return nil, validate.Wrap(op, ErrHMACKeyNotSet)
	}
	if c.AEADKey == "" {
		return nil, validate.Wrap(op, ErrAEADKeyNotSet)
	}

	return c, nil
}
