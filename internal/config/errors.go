package config

import "errors"

var (
	ErrTelegramTokenNotSet = errors.New("TELEGRAM_TOKEN is not set")
	ErrHMACKeyNotSet       = errors.New("HMAC_KEY is not set")
	ErrAEADKeyNotSet       = errors.New("AEAD_KEY is not set")
)
