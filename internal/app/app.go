package app

import (
	"context"
	"log/slog"

	"github.com/tuor4eg/ip_accounting_bot/internal/config"
	"github.com/tuor4eg/ip_accounting_bot/internal/logging"
)

type App struct {
	cfg     *config.Config
	runners []Runner
	log     *slog.Logger
}

func New(cfg *config.Config) *App {
	return &App{
		cfg: cfg,
		log: logging.WithPackage(),
	}
}

func (a *App) Register(runner Runner) {
	a.runners = append(a.runners, runner)
}

func (a *App) Run(ctx context.Context) error {
	a.log.Info("IP Accounting Bot: starting...")
	defer a.log.Info("IP Accounting Bot: stopped.")

	return runAll(ctx, a.runners)
}
