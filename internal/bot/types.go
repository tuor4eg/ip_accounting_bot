package bot

import (
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/domain"
)

// BotDeps contains all dependencies for the bot.
type BotDeps struct {
	Identities domain.IdentityStore
	Income     domain.IncomeUsecase
	Payment    domain.PaymentUsecase
	Total      domain.TotalUsecase
	// Now returns current time; if nil, time.Now is used.
	Now func() time.Time
}
