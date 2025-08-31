package memstore

import (
	"context"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/validate"
)

func (s *Store) InsertIncome(ctx context.Context, userID int64, at time.Time, amount int64, note string) error {
	const op = "memstore.InsertIncome"

	if err := validate.ValidateAmount(amount); err != nil {
		return validate.Wrap(op, err)
	}

	utc := at.UTC()
	day := time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.incomes[userID] = append(s.incomes[userID], IncomeRecord{
		At:     day,
		Amount: amount,
		Note:   note,
	})

	return nil
}

func (s *Store) VoidLastIncomeInRange(ctx context.Context, userID int64, from, to, now time.Time) (
	amount int64, at time.Time, note string, ok bool, err error,
) {
	const op = "memstore.VoidLastIncomeInRange"

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := validate.ValidateDateRangeUTC(from, to); err != nil {
		return 0, time.Time{}, "", false, validate.Wrap(op, err)
	}

	incomes := s.incomes[userID]

	if len(incomes) == 0 {
		return 0, time.Time{}, "", false, nil
	}

	bestIdx := -1
	bestAt := time.Time{}

	for i, income := range incomes {
		if !income.At.Before(from) && !income.At.After(to) && income.VoidedAt.IsZero() {

			if bestIdx == -1 || income.At.After(bestAt) || (income.At.Equal(bestAt) && i > bestIdx) {
				bestIdx = i
				bestAt = income.At
			}
		}
	}

	if bestIdx == -1 {
		return 0, time.Time{}, "", false, nil
	}

	incomes[bestIdx].VoidedAt = now

	return incomes[bestIdx].Amount, incomes[bestIdx].At, incomes[bestIdx].Note, true, nil
}

func (s *Store) SumIncomes(ctx context.Context, userID int64, from, to time.Time) (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sum := int64(0)
	for _, income := range s.incomes[userID] {
		if !income.At.After(from) && !income.At.Before(to) && income.VoidedAt.IsZero() {
			sum += income.Amount
		}
	}

	return sum, nil
}
