package bot

import (
	"context"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/domain"
)

type IdentityStore interface {
	UpsertIdentity(ctx context.Context, transport, externalID string) (int64, error)
}

type IncomeAdder interface {
	AddIncome(ctx context.Context, userID int64, at time.Time, amount int64, note string) error
}

type PaymentAdder interface {
	AddPayment(ctx context.Context, userID int64, at time.Time, amount int64, note string, payoutType domain.PaymentType) error
}

type AddDeps struct {
	Identities IdentityStore
	Income     IncomeAdder
	Payment    PaymentAdder
	Now        func() time.Time
}

type QuarterSummer interface {
	// returns (sum, tax, qStart, qEnd)
	SumQuarter(ctx context.Context, userID int64, now time.Time) (int64, int64, time.Time, time.Time, error)
}

type TotalDeps struct {
	Identities IdentityStore
	QuarterSum QuarterSummer
	Now        func() time.Time // optional; если nil — использовать time.Now
}

type undoerPayment interface {
	UndoLastYear(ctx context.Context, userID int64, now time.Time, paymentType domain.PaymentType) (amount int64, at time.Time, note string, pType domain.PaymentType, ok bool, err error)
}

type undoerIncome interface {
	UndoLastQuarter(ctx context.Context, userID int64, now time.Time) (amount int64, at time.Time, note string, ok bool, err error)
}
