package memstore

import (
	"sync"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/cryptostore"
	"github.com/tuor4eg/ip_accounting_bot/internal/domain"
)

// IncomeRecord represents an income entry in memory storage
type IncomeRecord struct {
	At       time.Time
	Amount   int64
	Note     string
	VoidedAt time.Time
}

// PaymentRecord represents a payment entry in memory storage
type PaymentRecord struct {
	At       time.Time
	Amount   int64
	Note     string
	VoidedAt time.Time
	Type     domain.PaymentType
}

// UserRecord represents a user identity in memory storage
type UserRecord struct {
	UserID int64
	Scheme domain.TaxScheme
}

// Store provides in-memory storage with cryptographic capabilities
type Store struct {
	cryptostore.BaseCryptoStore // Embed crypto capabilities
	mu                          sync.RWMutex
	nextUserID                  int64
	identities                  map[string]UserRecord
	incomes                     map[int64][]IncomeRecord
	payments                    map[int64][]PaymentRecord
}
