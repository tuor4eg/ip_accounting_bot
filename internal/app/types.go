package app

import (
	"log/slog"

	"github.com/tuor4eg/ip_accounting_bot/internal/bot"
	"github.com/tuor4eg/ip_accounting_bot/internal/config"
	"github.com/tuor4eg/ip_accounting_bot/internal/telegram"
)

// App is the main application that manages all components
type App struct {
	cfg     *config.Config
	runners []Runner
	log     *slog.Logger
	store   Store
	income  IncomeUsecase
	payment PaymentUsecase
	total   TotalUsecase
}

// TelegramRunner handles Telegram bot operations
type TelegramRunner struct {
	tg      *telegram.Client
	log     *slog.Logger
	botDeps *bot.BotDeps
}
