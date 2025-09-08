package app

import (
	"log/slog"

	"github.com/tuor4eg/ip_accounting_bot/config"
	"github.com/tuor4eg/ip_accounting_bot/internal/domain"
)

// App is the main application that manages all components
type App struct {
	cfg     *config.Config
	runners []Runner
	log     *slog.Logger
	store   Store
	income  domain.IncomeUsecase
	payment domain.PaymentUsecase
	total   domain.TotalUsecase
}
