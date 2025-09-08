package runner

import "context"

// Runner is the interface that all transport runners must implement
type Runner interface {
	Name() string
	Run(ctx context.Context) error
}
