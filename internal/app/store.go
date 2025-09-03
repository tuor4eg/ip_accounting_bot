package app

import (
	"context"

	"github.com/tuor4eg/ip_accounting_bot/internal/cryptostore"
)

type Store interface {
	Close(ctx context.Context) error
	cryptostore.CryptoStore // Embed crypto capabilities
}

// SetStore injects a storage implementation into the App and returns the App for chaining.
func (a *App) SetStore(s Store) *App {
	a.store = s
	return a
}
