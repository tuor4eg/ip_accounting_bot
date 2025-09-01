package bot

import (
	"context"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/domain"
	"github.com/tuor4eg/ip_accounting_bot/internal/validate"
)

func HandleAddContrib(ctx context.Context, deps AddDeps, transport, externalID string, args string) (string, error) {
	const op = "bot.HandleAddContrib"
	amount, note, err := ParseAmountAndNote(args)

	if err != nil {
		return "", validate.Wrap(op, err)
	}

	// Resolve or create user identity.
	userID, err := deps.Identities.UpsertIdentity(ctx, transport, externalID)

	if err != nil {
		return "", validate.Wrap(op, err)
	}

	// Use UTC "now"; storage casts to DATE.
	now := time.Now
	if deps.Now != nil {
		now = deps.Now
	}
	at := now().UTC()

	// Persist Contribution.
	if err := deps.Payment.AddPayment(ctx, userID, at, amount, note, domain.PaymentType(domain.PaymentTypeContrib)); err != nil {
		return "", validate.Wrap(op, err)
	}

	return AddContribSuccessText(amount, at, note), nil
}
