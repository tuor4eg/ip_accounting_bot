package validate

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/domain"
)

var (
	ErrInvalidUserID      = errors.New("userID must be positive")
	ErrInvalidAmount      = errors.New("amount must be positive")
	ErrNoAllowedValues    = errors.New("no allowed values")
	ErrNotAllowed         = errors.New("value not allowed")
	ErrInvalidDateRange   = errors.New("invalid date range")
	ErrInvalidDateUTC     = errors.New("date is not UTC")
	ErrInvalidPaymentType = errors.New("invalid payment type")
	ErrInvalidTransport   = errors.New("invalid transport")
	ErrInvalidExternalID  = errors.New("invalid external ID")
	ErrInvalidDate        = errors.New("invalid date")
	ErrEmptyString        = errors.New("empty string")
	ErrNotFound           = errors.New("not found")
)

func IsUTC(t time.Time) bool { _, off := t.Zone(); return off == 0 }

func ValidateUserID(userID int64) error {
	if userID <= 0 {
		return ErrInvalidUserID
	}
	return nil
}

func ValidateAmount(amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	return nil
}

// OneOf reports an error if v is not one of allowed.
// Generic and independent from concrete domain types.
func OneOf[T comparable](v T, allowed ...T) error {
	if len(allowed) == 0 {
		return ErrNoAllowedValues
	}
	if slices.Contains(allowed, v) {
		return nil
	}
	// Wrap with context for better diagnostics.
	return fmt.Errorf("%w: got=%v allowed=%v", ErrNotAllowed, v, allowed)
}

func ValidatePaymentType(paymentType domain.PaymentType) error {
	if err := OneOf(paymentType, domain.PaymentTypeContrib, domain.PaymentTypeAdvance); err != nil {
		return ErrInvalidPaymentType
	}
	return nil
}

func ValidateDateRangeUTC(from, to time.Time) error {
	if from.IsZero() || to.IsZero() {
		return ErrInvalidDateRange
	}

	if !IsUTC(from) || !IsUTC(to) {
		return ErrInvalidDateUTC
	}

	if to.Before(from) {
		return ErrInvalidDateRange
	}
	return nil
}

func ValidateTransport(transport string) error {
	if transport == "" {
		return ErrInvalidTransport
	}
	return nil
}

func ValidateExternalID(externalID string) error {
	if externalID == "" {
		return ErrInvalidExternalID
	}
	return nil
}

func ValidateDate(date time.Time) error {
	if date.IsZero() {
		return ErrInvalidDate
	}
	return nil
}

func ValidateEmptyString(s string) error {
	if s == "" {
		return ErrEmptyString
	}
	return nil
}
