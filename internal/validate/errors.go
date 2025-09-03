package validate

import "errors"

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
