package money

import "errors"

var (
	// ErrInvalidFormat is returned when the input cannot be parsed as a money amount.
	ErrInvalidFormat = errors.New("invalid amount format")
	// ErrOverflow is returned when the parsed number would overflow int64.
	ErrOverflow = errors.New("amount overflow")
)
