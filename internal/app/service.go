package app

import (
	"context"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/domain"
)

type IncomeUsecase interface {
	AddIncome(ctx context.Context, userID int64, at time.Time, amount int64, note string) error
	SumQuarter(ctx context.Context, userID int64, now time.Time) (sum int64, tax int64, qStart time.Time, qEnd time.Time, err error)
}

type PaymentUsecase interface {
	AddPayment(ctx context.Context, userID int64, at time.Time, amount int64, note string, payoutType domain.PaymentType) error
}

// SetIncomeUsecase injects domain income usecase into the App and returns the App for chaining.
func (a *App) SetIncomeUsecase(u IncomeUsecase) *App {
	a.income = u
	return a
}

// SetPaymentUsecase injects domain payment usecase into the App and returns the App for chaining.
func (a *App) SetPaymentUsecase(u PaymentUsecase) *App {
	a.payment = u
	return a
}
