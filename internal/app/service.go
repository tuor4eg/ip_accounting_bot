package app

import (
	"github.com/tuor4eg/ip_accounting_bot/internal/domain"
)

// SetIncomeUsecase injects domain income usecase into the App and returns the App for chaining.
func (a *App) SetIncomeUsecase(u domain.IncomeUsecase) *App {
	a.income = u
	return a
}

// SetPaymentUsecase injects domain payment usecase into the App and returns the App for chaining.
func (a *App) SetPaymentUsecase(u domain.PaymentUsecase) *App {
	a.payment = u
	return a
}

// SetTotalUsecase injects domain total usecase into the App and returns the App for chaining.
func (a *App) SetTotalUsecase(u domain.TotalUsecase) *App {
	a.total = u
	return a
}
