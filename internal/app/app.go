package app

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/bot"
	"github.com/tuor4eg/ip_accounting_bot/internal/config"
	"github.com/tuor4eg/ip_accounting_bot/internal/logging"
	"github.com/tuor4eg/ip_accounting_bot/internal/validate"
)

type App struct {
	cfg     *config.Config
	runners []Runner
	log     *slog.Logger
	store   Store
	income  IncomeUsecase
	payment PaymentUsecase
	total   TotalUsecase
}

var (
	ErrStoreNotSet                        = errors.New("store is not set")
	ErrIncomeUsecaseNotSet                = errors.New("income usecase is not set")
	ErrStoreDoesNotImplementIdentityStore = errors.New("store does not implement IdentityStore")
	ErrPaymentUsecaseNotSet               = errors.New("payment usecase is not set")
	ErrTotalUsecaseNotSet                 = errors.New("total usecase is not set")
)

func New(cfg *config.Config) *App {
	return &App{
		cfg: cfg,
		log: logging.WithPackage(),
	}
}

func (a *App) Register(runner Runner) {
	a.runners = append(a.runners, runner)
}

func (a *App) BotDeps() (*bot.BotDeps, error) {
	op := "app.BotDeps"

	if a.store == nil {
		return nil, validate.Wrap(op, ErrStoreNotSet)
	}

	if a.income == nil {
		return nil, validate.Wrap(op, ErrIncomeUsecaseNotSet)
	}

	if a.payment == nil {
		return nil, validate.Wrap(op, ErrPaymentUsecaseNotSet)
	}

	if a.total == nil {
		return nil, validate.Wrap(op, ErrTotalUsecaseNotSet)
	}

	// Check that store implements the required interface
	ids, ok := a.store.(interface {
		UpsertIdentity(ctx context.Context, transport, externalID string, chatID int64) (int64, error)
	})

	if !ok {
		return nil, validate.Wrap(op, ErrStoreDoesNotImplementIdentityStore)
	}

	return bot.NewBotDeps(ids, a.income, a.payment, a.total, time.Now), nil
}

func (a *App) Run(ctx context.Context) error {
	a.log.Info("IP Accounting Bot: starting...")
	defer a.log.Info("IP Accounting Bot: stopped.")

	return runAll(ctx, a.runners)
}
