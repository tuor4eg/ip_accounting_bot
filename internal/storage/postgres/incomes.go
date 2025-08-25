package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/tuor4eg/ip_accounting_bot/internal/helper"
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

// VoidLastIncomeInRange marks the newest "active" income in [from,to] as voided (soft-delete).
// "Newest" is determined by (at DESC, created_at DESC, id DESC).
// Returns the voided record's (amount, at, note). ok=false if nothing to void.
func (s *Store) VoidLastIncomeInRange(ctx context.Context, userID int64, from, to, now time.Time) (
	amount int64, at time.Time, note string, ok bool, err error,
) {
	const op = "postgres.VoidLastIncomeInRange"

	if !helper.IsUTC(from) || !helper.IsUTC(to) || !helper.IsUTC(now) {
		return 0, time.Time{}, "", false, fmt.Errorf("%s: non-UTC time", op)
	}

	const q = `
	WITH cand AS (
		SELECT id
		FROM incomes
		WHERE user_id = $1
		  AND at BETWEEN $2 AND $3
		  AND voided_at IS NULL
		ORDER BY at DESC, created_at DESC, id DESC
		LIMIT 1
	)
	UPDATE incomes AS i
	SET voided_at = $4
	FROM cand
	WHERE i.id = cand.id
	RETURNING i.amount, i.at, i.note;
	`

	row := s.Pool.QueryRow(ctx, q, userID, from, to, now)

	var noteNull sql.NullString

	if err := row.Scan(&amount, &at, &noteNull); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, time.Time{}, "", false, nil
		}
		return 0, time.Time{}, "", false, fmt.Errorf("%s: %w", op, err)
	}

	if noteNull.Valid {
		note = noteNull.String
	}

	return amount, at, note, true, nil
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
			  AND voided_at IS NULL
		`, userID, from, to).
		Scan(&sum); err != nil {
		return 0, fmt.Errorf("sum incomes: %w", err)
	}
	return sum, nil
}
