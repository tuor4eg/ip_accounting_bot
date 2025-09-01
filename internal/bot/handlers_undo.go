package bot

import (
	"context"
	"strings"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/validate"
)

func HandleUndo(ctx context.Context, deps AddDeps, transport, externalID string, args string) (string, error) {
	const op = "bot.HandleUndo"

	_ = strings.TrimSpace(args) // args

	userID, err := deps.Identities.UpsertIdentity(ctx, transport, externalID)

	if err != nil {
		return "", validate.Wrap(op, err)
	}

	now := time.Now

	if deps.Now != nil {
		now = deps.Now
	}

	nowUTC := now().UTC()

	u, ok := deps.Income.(undoerIncome)

	if !ok {
		return "", validate.Wrap(op, ErrServiceDoesNotSupportUndo)
	}

	amount, at, note, ok, err := u.UndoLastQuarter(ctx, userID, nowUTC)

	if err != nil {
		return "", validate.Wrap(op, err)
	}

	if !ok {
		return UndoNoIncomeText(), nil
	}

	return UndoSuccessText(amount, at, note), nil
}
