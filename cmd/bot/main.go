package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tuor4eg/ip_accounting_bot/internal/app"
	"github.com/tuor4eg/ip_accounting_bot/internal/config"
	"github.com/tuor4eg/ip_accounting_bot/internal/logging"
	"github.com/tuor4eg/ip_accounting_bot/internal/telegram"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()

	if err != nil {
		panic(err)
	}

	logging.InitFromEnv(cfg.LogLevel, cfg.LogFormat)

	a := app.New(cfg)

	tg := telegram.New(cfg.TelegramToken, nil)

	a.Register(app.NewTelegramRunner(tg))

	if err := a.Run(ctx); err != nil {
		log.Fatalf("app: run error: %v", err)
	}
}
