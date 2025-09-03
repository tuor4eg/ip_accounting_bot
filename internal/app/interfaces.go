package app

import (
	"context"

	"github.com/tuor4eg/ip_accounting_bot/internal/cryptostore"
	"github.com/tuor4eg/ip_accounting_bot/internal/telegram"
)

type Store interface {
	Close(ctx context.Context) error
	cryptostore.CryptoStore
}

type Runner interface {
	Name() string
	Run(ctx context.Context) error
}

type TelegramUpdateGetter interface {
	GetUpdates(ctx context.Context, offset int64, timeoutSec int) ([]telegram.Update, error)
}

type TelegramSender interface {
	SendMessage(ctx context.Context, chatID int64, text string) error
}
