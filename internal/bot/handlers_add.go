package bot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/money"
)

type IdentityStore interface {
	UpsertIdentity(ctx context.Context, transport, externalID string) (int64, error)
}

type IncomeAdder interface {
	AddIncome(ctx context.Context, userID int64, at time.Time, amount int64, note string) error
}

type AddDeps struct {
	Identities IdentityStore
	Income     IncomeAdder
	Now        func() time.Time
}

func HandleAdd(ctx context.Context, deps AddDeps, transport, externalID string, args string) (string, error) {
	const op = "bot.HandleAdd"

	args = strings.TrimSpace(args)

	if args == "" {
		return "", fmt.Errorf("%s: no arguments provided: %w", op, ErrBadInput)
	}

	toks := strings.Fields(args)

	var (
		amount int64
		cut    = -1
	)

	for i := 1; i <= len(toks); i++ {
		prefix := strings.Join(toks[:i], " ")
		v, err := money.ParseAmount(prefix)
		if err == nil {
			amount = v
			cut = i
		}
	}
	if cut == -1 {
		return "", fmt.Errorf("%s: parse amount: %w", op, ErrBadInput)
	}
	note := strings.TrimSpace(strings.Join(toks[cut:], " "))

	// Resolve or create user identity.
	userID, err := deps.Identities.UpsertIdentity(ctx, transport, externalID)
	if err != nil {
		return "", fmt.Errorf("%s: upsert identity: %w", op, err)
	}

	// Use UTC "now"; storage casts to DATE.
	now := time.Now
	if deps.Now != nil {
		now = deps.Now
	}
	at := now().UTC()

	// Persist income.
	if err := deps.Income.AddIncome(ctx, userID, at, amount, note); err != nil {
		return "", fmt.Errorf("%s: add income: %w", op, err)
	}

	return AddSuccessText(amount, at, note), nil
}
