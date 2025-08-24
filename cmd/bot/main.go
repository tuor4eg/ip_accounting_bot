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
	"github.com/tuor4eg/ip_accounting_bot/internal/service"
	"github.com/tuor4eg/ip_accounting_bot/internal/storage/postgres"
	"github.com/tuor4eg/ip_accounting_bot/internal/telegram"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()

	if err != nil {
		log.Fatalf("app: config error: %v", err)
	}

	logging.InitFromEnv(cfg.LogLevel, cfg.LogFormat)

	a := app.New(cfg)

	store, err := postgres.Open(ctx, cfg.DatabaseURL)

	if err != nil {
		log.Fatalf("app: postgres store error: %v", err)
	}

	defer func() {
		if err := store.Close(ctx); err != nil {
			log.Printf("failed to close store: %v", err)
		}
	}()

	income := service.NewIncomeService(store)
	a.SetStore(store).SetIncomeUsecase(income)

	tg := telegram.New(cfg.TelegramToken, nil)

	addDeps, totalDeps, err := a.BotDeps()

	if err != nil {
		log.Fatalf("app: bot deps error: %v", err)
	}

	a.Register(app.NewTelegramRunner(tg).SetBotDeps(addDeps, totalDeps))

	if err := a.Run(ctx); err != nil {
		log.Fatalf("app: run error: %v", err)
	}
}
