package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/tuor4eg/ip_accounting_bot/internal/domain"
	"github.com/tuor4eg/ip_accounting_bot/internal/validate"
)

func (s *Store) InsertPayment(ctx context.Context, userID int64, at time.Time, amount int64, note string, paymentType domain.PaymentType) error {
	const op = "postgres.InsertPayment"

	if err := validate.ValidateUserID(userID); err != nil {
		return validate.Wrap(op, err)
	}

	if err := validate.ValidateAmount(amount); err != nil {
		return validate.Wrap(op, err)
	}

	if err := validate.OneOf(paymentType, domain.PaymentTypeContrib, domain.PaymentTypeAdvance); err != nil {
		return validate.Wrap(op, err)
	}

	if err := validate.ValidateDateRangeUTC(at, at); err != nil {
		return validate.Wrap(op, err)
	}

	_, err := s.Pool.Exec(ctx, `
		INSERT INTO payments (user_id, at, amount, note, type)
		VALUES ($1, $2::date, $3, NULLIF($4, ''), $5)
	`, userID, at, amount, note, paymentType)
	if err != nil {
		return validate.Wrap(op, err)
	}

	return nil
}

func (s *Store) VoidLastPaymentInRange(ctx context.Context, userID int64, from, to, now time.Time, paymentType domain.PaymentType) (
	amount int64, at time.Time, note string, pType domain.PaymentType, ok bool, err error,
) {
	const op = "postgres.VoidLastPaymentInRange"

	if err := validate.ValidateUserID(userID); err != nil {
		return 0, time.Time{}, "", domain.PaymentType(""), false, validate.Wrap(op, err)
	}

	if err := validate.OneOf(paymentType, domain.PaymentTypeContrib, domain.PaymentTypeAdvance); err != nil {
		return 0, time.Time{}, "", domain.PaymentType(""), false, validate.Wrap(op, err)
	}

	if err := validate.ValidateDateRangeUTC(from, to); err != nil {
		return 0, time.Time{}, "", domain.PaymentType(""), false, validate.Wrap(op, err)
	}

	const q = `
	WITH cand AS (
		SELECT id
		FROM payments
		WHERE user_id = $1
		  AND at BETWEEN $2 AND $3
		  AND type = $4
		  AND voided_at IS NULL
		ORDER BY at DESC, created_at DESC, id DESC
		LIMIT 1
	)
	UPDATE payments AS p
	SET voided_at = $5
	FROM cand
	WHERE p.id = cand.id
	RETURNING p.amount, p.at, p.note, p.type;
	`

	row := s.Pool.QueryRow(ctx, q, userID, from, to, paymentType, now)

	var noteNull sql.NullString

	if err := row.Scan(&amount, &at, &noteNull, &pType); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, time.Time{}, "", domain.PaymentType(""), false, nil
		}
		return 0, time.Time{}, "", domain.PaymentType(""), false, validate.Wrap(op, err)
	}

	if noteNull.Valid {
		note = noteNull.String
	}
	return amount, at, note, pType, true, nil
}

func (s *Store) SumPayments(ctx context.Context, userID int64, from, to time.Time) (int64, int64, error) {
	const op = "postgres.SumPayments"

	if err := validate.ValidateUserID(userID); err != nil {
		return 0, 0, validate.Wrap(op, err)
	}

	if err := validate.ValidateDateRangeUTC(from, to); err != nil {
		return 0, 0, validate.Wrap(op, err)
	}

	sumContrib := int64(0)
	sumAdvance := int64(0)

	if err := s.Pool.QueryRow(ctx, `
	SELECT
		COALESCE(SUM(amount) FILTER (WHERE type = $1), 0)::bigint AS sum_contrib,
		COALESCE(SUM(amount) FILTER (WHERE type = $2), 0)::bigint AS sum_advance
	FROM payments
	WHERE user_id = $3
		AND at BETWEEN $4::date AND $5::date
		AND voided_at IS NULL

	`, domain.PaymentTypeContrib, domain.PaymentTypeAdvance, userID, from, to).Scan(&sumContrib, &sumAdvance); err != nil {
		return 0, 0, validate.Wrap(op, err)
	}
	return sumContrib, sumAdvance, nil
}
