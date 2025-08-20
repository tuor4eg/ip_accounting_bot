package app

import (
	"context"
	"fmt"

	"github.com/tuor4eg/ip_accounting_bot/internal/config"
)

type App struct {
	cfg     *config.Config
	runners []Runner
}

func New(cfg *config.Config) *App {
	return &App{cfg: cfg}
}

func (a *App) Register(runner Runner) {
	a.runners = append(a.runners, runner)
}

func (a *App) Run(ctx context.Context) error {
	fmt.Println("IP Accounting Bot: starting...")
	defer fmt.Println("IP Accounting Bot: stopped.")

	return runAll(ctx, a.runners)
}
