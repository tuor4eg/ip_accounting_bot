package memstore

import (
	"context"
)

func NewStore() *Store {
	return &Store{
		nextUserID: 1,
		identities: make(map[string]UserRecord),
		incomes:    make(map[int64][]IncomeRecord),
		payments:   make(map[int64][]PaymentRecord),
	}
}

func (s *Store) Close(ctx context.Context) error {
	return nil
}
