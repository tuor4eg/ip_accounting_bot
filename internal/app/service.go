package app

import (
	"context"
	"time"
)

type IncomeUsecase interface {
	AddIncome(ctx context.Context, userID int64, at time.Time, amount int64, note string) error
	SumQuarter(ctx context.Context, userID int64, now time.Time) (sum int64, tax int64, qStart time.Time, qEnd time.Time, err error)
}

// SetIncomeUsecase injects domain income usecase into the App and returns the App for chaining.
func (a *App) SetIncomeUsecase(u IncomeUsecase) *App {
	a.income = u
	return a
}
