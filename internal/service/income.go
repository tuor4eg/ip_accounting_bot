package service

import (
	"context"
	"strings"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/period"
	"github.com/tuor4eg/ip_accounting_bot/internal/validate"
)

func NewIncomeService(store IncomeStore) *IncomeService {
	return &IncomeService{store: store}
}

// AddIncome validates input and persists a single income record.
// - userID must be > 0
// - amount is in minor units (e.g., kopecks) and must be >= 0
// - at must be a non-zero time; the date part is persisted (storage casts to DATE)
// - note is trimmed; empty string is stored as NULL (handled by storage)
func (s *IncomeService) AddIncome(ctx context.Context, userID int64, at time.Time, amount int64, note string) error {
	const op = "service.IncomeService.AddIncome"

	if err := validate.ValidateUserID(userID); err != nil {
		return validate.Wrap(op, err)
	}

	if err := validate.ValidateAmount(amount); err != nil {
		return validate.Wrap(op, err)
	}

	if err := validate.ValidateDate(at); err != nil {
		return validate.Wrap(op, err)
	}
	note = strings.TrimSpace(note)

	// Delegate to storage; it applies DATE cast and NULLIF on note.
	if err := s.store.InsertIncome(ctx, userID, at, amount, note); err != nil {
		return validate.Wrap(op, err)
	}
	return nil
}

// UndoLastQuarter deletes the last quarter's income records.
// It's a no-op if there are no records to delete.
func (s *IncomeService) UndoLastQuarter(ctx context.Context, userID int64, now time.Time) (int64, time.Time, string, bool, error) {
	const op = "service.IncomeService.UndoLastQuarter"

	nowUTC := now.UTC()
	qStart, qEnd := period.QuarterBounds(now.UTC())

	amount, at, note, ok, err := s.store.VoidLastIncomeInRange(ctx, userID, qStart, qEnd, nowUTC)

	if err != nil {
		return 0, time.Time{}, "", false, validate.Wrap(op, err)
	}

	return amount, at, note, ok, nil
}

func (s *IncomeService) SumIncomes(ctx context.Context, userID int64, from, to time.Time) (int64, error) {
	const op = "service.IncomeService.SumIncomes"

	if err := validate.ValidateUserID(userID); err != nil {
		return 0, validate.Wrap(op, err)
	}

	if err := validate.ValidateDateRangeUTC(from, to); err != nil {
		return 0, validate.Wrap(op, err)
	}

	sum, err := s.store.SumIncomes(ctx, userID, from, to)

	if err != nil {
		return 0, validate.Wrap(op, err)
	}

	return sum, nil
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
		return 0, 0, qStart, qEnd, validate.Wrap(op, err)
	}

	// Deterministic 6% tax: floor division, no floats.
	tax = (sum * 6) / 100

	return sum, tax, qStart, qEnd, nil
}
