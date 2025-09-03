package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tuor4eg/ip_accounting_bot/internal/cryptostore"
)

// Store provides PostgreSQL storage with cryptographic capabilities
type Store struct {
	cryptostore.BaseCryptoStore // Embed crypto capabilities
	Pool                        *pgxpool.Pool
}
