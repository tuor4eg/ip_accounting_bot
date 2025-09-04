package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/tuor4eg/ip_accounting_bot/internal/validate"
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
