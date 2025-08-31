package memstore

import (
	"context"
	"sync"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/domain"
)

type IncomeRecord struct {
	At       time.Time
	Amount   int64
	Note     string
	VoidedAt time.Time
}

type PaymentRecord struct {
	At       time.Time
	Amount   int64
	Note     string
	VoidedAt time.Time
	Type     domain.PaymentType
}

type Store struct {
	mu         sync.RWMutex
	nextUserID int64
	identities map[string]int64
	incomes    map[int64][]IncomeRecord
	payments   map[int64][]PaymentRecord
}

func NewStore() *Store {
	return &Store{
		nextUserID: 1,
		identities: make(map[string]int64),
		incomes:    make(map[int64][]IncomeRecord),
		payments:   make(map[int64][]PaymentRecord),
	}
}

func (s *Store) Close(ctx context.Context) error {
	return nil
}
