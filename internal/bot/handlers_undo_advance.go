package bot

import (
	"context"
	"strings"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/domain"
	"github.com/tuor4eg/ip_accounting_bot/internal/validate"
)

func HandleUndoAdvance(ctx context.Context, deps *BotDeps, transport, externalID string, args string) (string, error) {
	const op = "bot.HandleUndoAdvance"

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

	amount, at, note, _, ok, err := deps.Payment.UndoLastYear(ctx, userID, nowUTC, domain.PaymentTypeAdvance)

	if err != nil {
		return "", validate.Wrap(op, err)
	}

	if !ok {
		return UndoNoAdvanceText(), nil
	}

	return UndoAdvanceSuccessText(amount, at, note), nil
}
