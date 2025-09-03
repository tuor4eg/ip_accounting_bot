package service

import (
	"context"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/domain"
	"github.com/tuor4eg/ip_accounting_bot/internal/tax"
)

// IncomeService handles income-related business logic
type IncomeService struct {
	store IncomeStore
}

// PaymentService handles payment-related business logic
type PaymentService struct {
	store PaymentStore
}

// TotalService handles total calculation business logic
type TotalService struct {
	getUserScheme func(ctx context.Context, userID int64) (domain.TaxScheme, error)
	sumIncomes    func(ctx context.Context, userID int64, from, to time.Time) (int64, error)
	sumPayments   func(ctx context.Context, userID int64, from, to time.Time) (int64, int64, error)
	provider      tax.Provider
}
