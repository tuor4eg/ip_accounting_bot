package service

import (
	"context"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/domain"
)

type IncomeStore interface {
	InsertIncome(ctx context.Context, userID int64, at time.Time, amount int64, note string) error
	VoidLastIncomeInRange(ctx context.Context, userID int64, from, to, now time.Time) (
		amount int64, at time.Time, note string, ok bool, err error,
	)
	SumIncomes(ctx context.Context, userID int64, from, to time.Time) (int64, error)
}

type PaymentStore interface {
	InsertPayment(ctx context.Context, userID int64, at time.Time, amount int64, note string, payoutType domain.PaymentType) error
	VoidLastPaymentInRange(ctx context.Context, userID int64, from, to, now time.Time, payoutType domain.PaymentType) (
		amount int64, at time.Time, note string, pType domain.PaymentType, ok bool, err error,
	)
	SumPayments(ctx context.Context, userID int64, from, to time.Time) (int64, int64, error)
}
