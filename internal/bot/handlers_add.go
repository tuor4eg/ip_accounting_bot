package bot

import (
	"context"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/validate"
)

func HandleAdd(ctx context.Context, deps *BotDeps, transport, externalID string, args string) (string, error) {
	const op = "bot.HandleAdd"

	amount, note, err := ParseAmountAndNote(args)

	if err != nil {
		return "", validate.Wrap(op, err)
	}

	if err := validateEntryInput(amount, note); err != nil {
		return "", validate.Wrap(op, err)
	}

	// Resolve or create user identity.
	userID, err := deps.Identities.UpsertIdentity(ctx, transport, externalID, 0)

	if err != nil {
		return "", validate.Wrap(op, err)
	}

	// Use UTC "now"; storage casts to DATE.
	now := time.Now
	if deps.Now != nil {
		now = deps.Now
	}
	at := now().UTC()

	// Persist income.
	if err := deps.Income.AddIncome(ctx, userID, at, amount, note); err != nil {
		return "", validate.Wrap(op, err)
	}

	return AddSuccessText(amount, at, note), nil
}
