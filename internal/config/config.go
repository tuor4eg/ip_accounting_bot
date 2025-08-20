package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string `env:"TELEGRAM_TOKEN"`
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	c := &Config{
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
	}

	if c.TelegramToken == "" {
		return nil, errors.New("TELEGRAM_TOKEN is not set")
	}

	return c, nil
}
