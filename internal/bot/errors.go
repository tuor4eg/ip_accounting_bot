package bot

import "errors"

// ErrBadInput marks user-facing input errors (e.g., cannot parse amount).
var ErrBadInput = errors.New("bad input")
