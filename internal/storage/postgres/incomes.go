package postgres

import (
	"context"
	"fmt"
	"time"
)

// InsertIncome inserts a single income record.
// 'amount' is in minor currency units (e.g., kopecks), must be >= 0.
// 'at' is the income date; only the date part is stored (cast to DATE in SQL).
func (s *Store) InsertIncome(ctx context.Context, userID int64, at time.Time, amount int64, note string) error {
	if s == nil || s.Pool == nil {
		return fmt.Errorf("insert income: store not initialized")
	}
	if userID <= 0 {
		return fmt.Errorf("insert income: invalid userID")
	}
	if amount < 0 {
		return fmt.Errorf("insert income: negative amount")
	}

	// Persist only the calendar day for 'at'; NULLIF trims empty notes to NULL.
	_, err := s.Pool.Exec(ctx, `
		INSERT INTO incomes (user_id, at, amount, note)
		VALUES ($1, $2::date, $3, NULLIF($4, ''))
	`, userID, at, amount, note)
	if err != nil {
		return fmt.Errorf("insert income: %w", err)
	}
	return nil
}

// SumIncomes returns the total amount (in minor units) for a user in [from..to] inclusive.
// 'from' and 'to' are interpreted by their calendar dates (cast to DATE in SQL).
func (s *Store) SumIncomes(ctx context.Context, userID int64, from, to time.Time) (int64, error) {
	if userID <= 0 {
		return 0, fmt.Errorf("sum incomes: invalid userID")
	}
	if from.IsZero() || to.IsZero() {
		return 0, fmt.Errorf("sum incomes: from/to must be set")
	}
	if from.After(to) {
		return 0, fmt.Errorf("sum incomes: from > to")
	}

	var sum int64
	if err := s.Pool.
		QueryRow(ctx, `
			SELECT COALESCE(SUM(amount), 0)
			FROM incomes
			WHERE user_id = $1
			  AND at >= $2::date
			  AND at <= $3::date
		`, userID, from, to).
		Scan(&sum); err != nil {
		return 0, fmt.Errorf("sum incomes: %w", err)
	}
	return sum, nil
}
