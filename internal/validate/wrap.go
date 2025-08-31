package validate

import "fmt"

// Wrap adds operation context to err using %w.
// If err is nil, Wrap returns nil.
func Wrap(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", op, err)
}
