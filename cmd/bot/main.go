package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tuor4eg/ip_accounting_bot/config"
	"github.com/tuor4eg/ip_accounting_bot/internal/app"
	telegramrunner "github.com/tuor4eg/ip_accounting_bot/internal/runner/telegram_runner"
	"github.com/tuor4eg/ip_accounting_bot/internal/service"
	"github.com/tuor4eg/ip_accounting_bot/internal/storage/postgres"
	"github.com/tuor4eg/ip_accounting_bot/internal/tax"
	"github.com/tuor4eg/ip_accounting_bot/internal/telegram"
	"github.com/tuor4eg/ip_accounting_bot/pkg/logging"
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

	// Create PostgreSQL storage
	store, err := postgres.Open(ctx, cfg.DatabaseURL)

	if err != nil {
		log.Fatalf("app: postgres store error: %v", err)
	}

	defer func() {
		if err := store.Close(ctx); err != nil {
			log.Printf("failed to close store: %v", err)
		}
	}()

	// Configure crypto keys for PostgreSQL storage
	if err := store.SetCryptoKeys(cfg.HMACKey, 1, cfg.AEADKey, 1); err != nil {
		log.Fatalf("app: set crypto keys error: %v", err)
	}

	// Log crypto keys status
	if store.HasCryptoKeys() {
		log.Printf("PostgreSQL store configured with HMAC key ID: %d, AEAD key ID: %d",
			store.GetHMACKid(), store.GetAEADKid())
	}

	income := service.NewIncomeService(store)
	payment := service.NewPaymentService(store)
	total := service.NewTotalService(
		store.GetUserScheme,
		income.SumIncomes,
		payment.SumPayments,
		tax.NewDefaultProvider())

	a.SetStore(store).SetIncomeUsecase(income).SetPaymentUsecase(payment).SetTotalUsecase(total)

	tg := telegram.New(cfg.TelegramToken, nil)

	botDeps, err := a.BotDeps()

	if err != nil {
		log.Fatalf("app: bot deps error: %v", err)
	}

	a.Register(telegramrunner.NewRunner(tg).SetBotDeps(botDeps))

	if err := a.Run(ctx); err != nil {
		log.Fatalf("app: run error: %v", err)
	}
}
