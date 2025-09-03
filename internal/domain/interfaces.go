package domain

import (
	"context"
	"time"
)

type IncomeUsecase interface {
	AddIncome(ctx context.Context, userID int64, at time.Time, amount int64, note string) error
	UndoLastQuarter(ctx context.Context, userID int64, now time.Time) (int64, time.Time, string, bool, error)
}

type PaymentUsecase interface {
	AddPayment(ctx context.Context, userID int64, at time.Time, amount int64, note string, payoutType PaymentType) error
	UndoLastYear(ctx context.Context, userID int64, now time.Time, paymentType PaymentType) (int64, time.Time, string, PaymentType, bool, error)
}

type TotalUsecase interface {
	SumQuarter(ctx context.Context, userID int64, now time.Time) (Totals, error)
	SumYearToDate(ctx context.Context, userID int64, now time.Time) (Totals, error)
}

type IdentityStore interface {
	UpsertIdentity(ctx context.Context, transport, externalID string, chatID int64) (int64, error)
}
