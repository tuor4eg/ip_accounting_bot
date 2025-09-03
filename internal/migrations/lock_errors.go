package migrations

import "errors"

var (
	ErrInvalidPool = errors.New("invalid pool")
	ErrEmptyName   = errors.New("empty name")
)
