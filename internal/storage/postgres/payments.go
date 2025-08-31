package postgres

import (
	"context"
	"time"

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
