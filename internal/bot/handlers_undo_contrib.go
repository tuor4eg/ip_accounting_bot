package bot

import (
	"context"
	"strings"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/domain"
	"github.com/tuor4eg/ip_accounting_bot/internal/validate"
)

func HandleUndoContrib(ctx context.Context, deps *BotDeps, transport, externalID string, args string) (string, error) {
	const op = "bot.HandleUndoContrib"

	_ = strings.TrimSpace(args) // args

	userID, err := deps.Identities.UpsertIdentity(ctx, transport, externalID, 0)

	if err != nil {
		return "", validate.Wrap(op, err)
	}

	now := time.Now

	if deps.Now != nil {
		now = deps.Now
	}

	nowUTC := now().UTC()

	amount, at, note, _, ok, err := deps.Payment.UndoLastYear(ctx, userID, nowUTC, domain.PaymentTypeContrib)

	if err != nil {
		return "", validate.Wrap(op, err)
	}

	if !ok {
		return UndoNoContribText(), nil
	}

	return UndoContribSuccessText(amount, at, note), nil
}
