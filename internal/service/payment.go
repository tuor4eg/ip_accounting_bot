package service

import (
	"context"
	"strings"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/domain"
	"github.com/tuor4eg/ip_accounting_bot/internal/period"
	"github.com/tuor4eg/ip_accounting_bot/internal/validate"
)

type PaymentStore interface {
	InsertPayment(ctx context.Context, userID int64, at time.Time, amount int64, note string, payoutType domain.PaymentType) error
	VoidLastPaymentInRange(ctx context.Context, userID int64, from, to, now time.Time, paymentType domain.PaymentType) (
		amount int64, at time.Time, note string, pType domain.PaymentType, ok bool, err error,
	)
	SumPayments(ctx context.Context, userID int64, from, to time.Time) (int64, int64, error)
}

type PaymentService struct {
	store PaymentStore
}

func NewPaymentService(store PaymentStore) *PaymentService {
	return &PaymentService{store: store}
}

func (s *PaymentService) AddPayment(ctx context.Context, userID int64, at time.Time, amount int64, note string, payoutType domain.PaymentType) error {
	const op = "service.PaymentService.AddPayment"

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
	if err := s.store.InsertPayment(ctx, userID, at, amount, note, payoutType); err != nil {
		return validate.Wrap(op, err)
	}
	return nil
}

func (s *PaymentService) UndoLastYear(ctx context.Context, userID int64, now time.Time, paymentType domain.PaymentType) (int64, time.Time, string, domain.PaymentType, bool, error) {
	const op = "service.PaymentService.UndoLastYear"

	nowUTC := now.UTC()
	yStart, yEnd := period.YearBounds(now.UTC())

	amount, at, note, pType, ok, err := s.store.VoidLastPaymentInRange(ctx, userID, yStart, yEnd, nowUTC, paymentType)

	if err != nil {
		return 0, time.Time{}, "", "", false, validate.Wrap(op, err)
	}

	return amount, at, note, pType, ok, nil
}

func (s *PaymentService) SumPayments(ctx context.Context, userID int64, from, to time.Time) (int64, int64, error) {
	const op = "service.PaymentService.SumPayments"

	if err := validate.ValidateUserID(userID); err != nil {
		return 0, 0, validate.Wrap(op, err)
	}

	if err := validate.ValidateDateRangeUTC(from, to); err != nil {
		return 0, 0, validate.Wrap(op, err)
	}

	sumContrib, sumAdvance, err := s.store.SumPayments(ctx, userID, from, to)

	if err != nil {
		return 0, 0, validate.Wrap(op, err)
	}

	return sumContrib, sumAdvance, nil
}
