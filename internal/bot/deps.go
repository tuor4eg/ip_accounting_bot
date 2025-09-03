// internal/bot/deps.go
package bot

import (
	"time"
)

// NewBotDeps wires dependencies for the bot.
// If now is nil, time.Now will be used.
func NewBotDeps(identities IdentityStore, income IncomeUsecase, payment PaymentUsecase, total TotalUsecase, now func() time.Time) *BotDeps {
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
