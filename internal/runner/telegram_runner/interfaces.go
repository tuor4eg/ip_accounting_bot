package telegram_runner

import (
	"context"

	"github.com/tuor4eg/ip_accounting_bot/internal/telegram"
)

// TelegramUpdateGetter defines the interface for getting Telegram updates
type TelegramUpdateGetter interface {
	GetUpdates(ctx context.Context, offset int64, timeoutSec int) ([]telegram.Update, error)
}

// TelegramSender defines the interface for sending Telegram messages
type TelegramSender interface {
	SendMessage(ctx context.Context, chatID int64, text string) error
}
