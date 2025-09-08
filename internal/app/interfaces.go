package app

import (
	"context"

	"github.com/tuor4eg/ip_accounting_bot/internal/cryptostore"
)

type Store interface {
	Close(ctx context.Context) error
	cryptostore.CryptoStore
}

type Runner interface {
	Name() string
	Run(ctx context.Context) error
}
