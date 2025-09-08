package telegram_runner

import (
	"log/slog"

	"github.com/tuor4eg/ip_accounting_bot/internal/bot"
	"github.com/tuor4eg/ip_accounting_bot/internal/telegram"
)

// Runner handles Telegram bot operations
type Runner struct {
	tg      *telegram.Client
	log     *slog.Logger
	botDeps *bot.BotDeps
}
