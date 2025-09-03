package bot

import (
	"context"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/domain"
)

// IdentityStore abstracts user identity storage.
type IdentityStore interface {
	UpsertIdentity(ctx context.Context, transport, externalID string, chatID int64) (int64, error)
}

// IncomeUsecase abstracts income operations for the bot.
type IncomeUsecase interface {
	AddIncome(ctx context.Context, userID int64, at time.Time, amount int64, note string) error
	UndoLastQuarter(ctx context.Context, userID int64, now time.Time) (amount int64, at time.Time, note string, ok bool, err error)
}

// PaymentUsecase abstracts payment operations for the bot.
type PaymentUsecase interface {
	AddPayment(ctx context.Context, userID int64, at time.Time, amount int64, note string, payoutType domain.PaymentType) error
	UndoLastYear(ctx context.Context, userID int64, now time.Time, paymentType domain.PaymentType) (amount int64, at time.Time, note string, pType domain.PaymentType, ok bool, err error)
}

// TotalUsecase provides totals for quarter and year-to-date.
type TotalUsecase interface {
	SumQuarter(ctx context.Context, userID int64, now time.Time) (domain.Totals, error)
	SumYearToDate(ctx context.Context, userID int64, now time.Time) (domain.Totals, error)
}

// BotDeps contains all dependencies for the bot.
type BotDeps struct {
	Identities IdentityStore
	Income     IncomeUsecase
	Payment    PaymentUsecase
	Total      TotalUsecase
	// Now returns current time; if nil, time.Now is used.
	Now func() time.Time
}
