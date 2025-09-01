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
}

var (
	ErrStoreNotSet                        = errors.New("store is not set")
	ErrIncomeUsecaseNotSet                = errors.New("income usecase is not set")
	ErrStoreDoesNotImplementIdentityStore = errors.New("store does not implement IdentityStore")
	ErrPaymentUsecaseNotSet               = errors.New("payment usecase is not set")
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

func (a *App) BotDeps() (add bot.AddDeps, total bot.TotalDeps, err error) {
	op := "app.BotDeps"

	if a.store == nil {
		return add, total, validate.Wrap(op, ErrStoreNotSet)
	}

	if a.income == nil {
		return add, total, validate.Wrap(op, ErrIncomeUsecaseNotSet)
	}

	if a.payment == nil {
		return add, total, validate.Wrap(op, ErrPaymentUsecaseNotSet)
	}

	ids, ok := a.store.(bot.IdentityStore)

	if !ok {
		return add, total, validate.Wrap(op, ErrStoreDoesNotImplementIdentityStore)
	}

	add = bot.AddDeps{
		Identities: ids,
		Income:     a.income,
		Payment:    a.payment,
		Now:        time.Now,
	}

	total = bot.TotalDeps{
		Identities: ids,
		QuarterSum: a.income,
		Now:        time.Now,
	}

	return add, total, nil
}

func (a *App) Run(ctx context.Context) error {
	a.log.Info("IP Accounting Bot: starting...")
	defer a.log.Info("IP Accounting Bot: stopped.")

	return runAll(ctx, a.runners)
}
