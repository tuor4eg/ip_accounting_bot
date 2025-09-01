package memstore

import (
	"context"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/domain"
	"github.com/tuor4eg/ip_accounting_bot/internal/validate"
)

func (s *Store) InsertPayment(ctx context.Context, userID int64, at time.Time, amount int64, note string, paymentType domain.PaymentType) error {
	const op = "memstore.InsertPayment"

	if err := validate.ValidateUserID(userID); err != nil {
		return validate.Wrap(op, err)
	}

	if err := validate.ValidateAmount(amount); err != nil {
		return validate.Wrap(op, err)
	}

	if err := validate.ValidatePaymentType(domain.PaymentType(paymentType)); err != nil {
		return validate.Wrap(op, err)
	}

	if err := validate.ValidateDateRangeUTC(at, at); err != nil {
		return validate.Wrap(op, err)
	}

	at = at.UTC()
	day := time.Date(at.Year(), at.Month(), at.Day(), 0, 0, 0, 0, time.UTC)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.payments[userID] = append(s.payments[userID], PaymentRecord{
		At:     day,
		Amount: amount,
		Note:   note,
		Type:   domain.PaymentType(paymentType),
	})

	return nil
}

func (s *Store) VoidLastPaymentInRange(ctx context.Context, userID int64, from, to, now time.Time, paymentType domain.PaymentType) (
	amount int64, at time.Time, note string, pType domain.PaymentType, ok bool, err error,
) {
	const op = "memstore.VoidLastPaymentInRange"

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := validate.ValidateDateRangeUTC(from, to); err != nil {
		return 0, time.Time{}, "", "", false, validate.Wrap(op, err)
	}

	payments := s.payments[userID]

	if len(payments) == 0 {
		return 0, time.Time{}, "", "", false, nil
	}

	bestIdx := -1
	bestAt := time.Time{}

	for i, payment := range payments {
		if !payment.At.Before(from) && !payment.At.After(to) && payment.VoidedAt.IsZero() && payment.Type == domain.PaymentType(paymentType) {
			if bestIdx == -1 || payment.At.After(bestAt) || (payment.At.Equal(bestAt) && i > bestIdx) {
				bestIdx = i
				bestAt = payment.At
			}
		}
	}

	if bestIdx == -1 {
		return 0, time.Time{}, "", "", false, nil
	}

	payments[bestIdx].VoidedAt = now
	return payments[bestIdx].Amount, payments[bestIdx].At, payments[bestIdx].Note, payments[bestIdx].Type, true, nil
}

func (s *Store) SumPayments(ctx context.Context, userID int64, from, to time.Time) (int64, int64, error) {
	const op = "memstore.SumPayments"

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := validate.ValidateDateRangeUTC(from, to); err != nil {
		return 0, 0, validate.Wrap(op, err)
	}

	payments := s.payments[userID]

	sumContrib := int64(0)
	sumAdvance := int64(0)

	for _, payment := range payments {
		if !payment.At.Before(from) && !payment.At.After(to) && payment.VoidedAt.IsZero() {
			if payment.Type == domain.PaymentTypeContrib {
				sumContrib += payment.Amount
			} else {
				sumAdvance += payment.Amount
			}
		}
	}

	return sumContrib, sumAdvance, nil
}
