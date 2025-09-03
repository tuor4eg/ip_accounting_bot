package bot

import (
	"context"
	"strings"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/validate"
)

func HandleUndo(ctx context.Context, deps *BotDeps, transport, externalID string, args string) (string, error) {
	const op = "bot.HandleUndo"

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

	amount, at, note, ok, err := deps.Income.UndoLastQuarter(ctx, userID, nowUTC)

	if err != nil {
		return "", validate.Wrap(op, err)
	}

	if !ok {
		return UndoNoIncomeText(), nil
	}

	return UndoSuccessText(amount, at, note), nil
}
