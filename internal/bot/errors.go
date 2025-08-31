package bot

import "errors"

// ErrBadInput marks user-facing input errors (e.g., cannot parse amount).
var (
	ErrBadInput                  = errors.New("bad input")
	ErrAmountIsZero              = errors.New("amount is zero")
	ErrUnknownCommand            = errors.New("unknown command")
	ErrServiceDoesNotSupportUndo = errors.New("service does not support undo")
)
