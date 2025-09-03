package config

// Config holds all configuration values for the application
type Config struct {
	TelegramToken string `env:"TELEGRAM_TOKEN"`
	LogLevel      string `env:"LOG_LEVEL"`
	LogFormat     string `env:"LOG_FORMAT"`
	DatabaseURL   string `env:"DATABASE_URL"`
	HMACKey       string `env:"HMAC_KEY"`
	AEADKey       string `env:"AEAD_KEY"`
}
