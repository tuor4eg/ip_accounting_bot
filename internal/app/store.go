package app

import (
	"context"
)

type Store interface {
	Close(ctx context.Context) error
}

// SetStore injects a storage implementation into the App and returns the App for chaining.
func (a *App) SetStore(s Store) *App {
	a.store = s
	return a
}
