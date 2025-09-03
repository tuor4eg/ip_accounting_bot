// internal/bot/deps.go
package bot

import (
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/domain"
)

// NewBotDeps wires dependencies for the bot.
// If now is nil, time.Now will be used.
func NewBotDeps(identities domain.IdentityStore, income domain.IncomeUsecase, payment domain.PaymentUsecase, total domain.TotalUsecase, now func() time.Time) *BotDeps {
	if now == nil {
		now = time.Now
	}
	return &BotDeps{
		Identities: identities,
		Income:     income,
		Payment:    payment,
		Total:      total,
		Now:        now,
	}
}
