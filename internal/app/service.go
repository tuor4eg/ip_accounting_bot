package app

import (
	"context"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/domain"
)

type IncomeUsecase interface {
	AddIncome(ctx context.Context, userID int64, at time.Time, amount int64, note string) error
	UndoLastQuarter(ctx context.Context, userID int64, now time.Time) (int64, time.Time, string, bool, error)
}

type PaymentUsecase interface {
	AddPayment(ctx context.Context, userID int64, at time.Time, amount int64, note string, payoutType domain.PaymentType) error
	UndoLastYear(ctx context.Context, userID int64, now time.Time, paymentType domain.PaymentType) (int64, time.Time, string, domain.PaymentType, bool, error)
}

type TotalUsecase interface {
	SumQuarter(ctx context.Context, userID int64, now time.Time) (domain.Totals, error)
	SumYearToDate(ctx context.Context, userID int64, now time.Time) (domain.Totals, error)
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

// SetTotalUsecase injects domain total usecase into the App and returns the App for chaining.
func (a *App) SetTotalUsecase(u TotalUsecase) *App {
	a.total = u
	return a
}
