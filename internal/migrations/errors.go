package migrations

import "errors"

var (
	ErrInvalidFS                 = errors.New("fs is nil")
	ErrDuplicateMigrationVersion = errors.New("duplicate migration version")
)
