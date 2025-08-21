package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/period"
)

type IncomeStore interface {
	InsertIncome(ctx context.Context, userID int64, at time.Time, amount int64, note string) error
	SumIncomes(ctx context.Context, userID int64, from, to time.Time) (int64, error)
}

type IncomeService struct {
	store IncomeStore
}

func NewIncomeService(store IncomeStore) *IncomeService {
	return &IncomeService{store: store}
}

// AddIncome validates input and persists a single income record.
// - userID must be > 0
// - amount is in minor units (e.g., kopecks) and must be >= 0
// - at must be a non-zero time; the date part is persisted (storage casts to DATE)
// - note is trimmed; empty string is stored as NULL (handled by storage)
func (s *IncomeService) AddIncome(ctx context.Context, userID int64, at time.Time, amount int64, note string) error {
	if userID <= 0 {
		return fmt.Errorf("add income: invalid userID")
	}
	if at.IsZero() {
		return fmt.Errorf("add income: at is zero")
	}
	if amount < 0 {
		return fmt.Errorf("add income: negative amount")
	}
	note = strings.TrimSpace(note)

	// Delegate to storage; it applies DATE cast and NULLIF on note.
	if err := s.store.InsertIncome(ctx, userID, at, amount, note); err != nil {
		return fmt.Errorf("add income: %w", err)
	}
	return nil
}

// SumQuarter returns total income and 6% tax for the current quarter.
// All amounts are int64 minor units (kopecks). No floats.
// Tax is computed with floor division for determinism.
func (s *IncomeService) SumQuarter(
	ctx context.Context,
	userID int64,
	now time.Time,
) (sum int64, tax int64, qStart time.Time, qEnd time.Time, err error) {
	const op = "service.IncomeService.SumQuarter"

	// Determine quarter bounds (inclusive) in UTC.
	qStart, qEnd = period.QuarterBounds(now.UTC())

	// Aggregate sum of incomes for the user within the quarter.
	sum, err = s.store.SumIncomes(ctx, userID, qStart, qEnd)
	if err != nil {
		return 0, 0, qStart, qEnd, fmt.Errorf("%s: SumIncomes: %w", op, err)
	}

	// Deterministic 6% tax: floor division, no floats.
	tax = (sum * 6) / 100

	return sum, tax, qStart, qEnd, nil
}
