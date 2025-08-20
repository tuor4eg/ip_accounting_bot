package app

import (
	"context"
	"fmt"
	"sync"
)

type Runner interface {
	Name() string
	Run(ctx context.Context) error
}

func runAll(ctx context.Context, runners []Runner) error {
	if len(runners) == 0 {
		<-ctx.Done()
		return nil
	}

	var wg sync.WaitGroup

	wg.Add(len(runners))

	errOnce := make(chan error, 1)

	for _, r := range runners {
		r := r

		go func() {
			defer wg.Done()

			if err := r.Run(ctx); err != nil {
				select {
				case errOnce <- fmt.Errorf("runner %s failed: %w", r.Name(), err):
				default:
				}
			}
		}()
	}

	done := make(chan struct{})

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errOnce:
		return err
	case <-done:
		return nil
	}
}
