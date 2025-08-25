package memstore

import (
	"context"
	"sync"
	"time"
)

type IncomeRecord struct {
	At       time.Time
	Amount   int64
	Note     string
	VoidedAt time.Time
}

type Store struct {
	mu         sync.RWMutex
	nextUserID int64
	identities map[string]int64
	incomes    map[int64][]IncomeRecord
}

func NewStore() *Store {
	return &Store{
		nextUserID: 1,
		identities: make(map[string]int64),
		incomes:    make(map[int64][]IncomeRecord),
	}
}

func (s *Store) Close(ctx context.Context) error {
	return nil
}
