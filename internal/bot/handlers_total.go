package bot

import (
	"context"
	"strings"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/validate"
)

type QuarterSummer interface {
	// returns (sum, tax, qStart, qEnd)
	SumQuarter(ctx context.Context, userID int64, now time.Time) (int64, int64, time.Time, time.Time, error)
}

type TotalDeps struct {
	Identities IdentityStore
	QuarterSum QuarterSummer
	Now        func() time.Time // optional; если nil — использовать time.Now
}

func HandleTotal(ctx context.Context, deps TotalDeps, transport, externalID, args string) (string, error) {
	const op = "bot.HandleTotal"

	// TODO: parse args
	_ = strings.TrimSpace(args)

	userID, err := deps.Identities.UpsertIdentity(ctx, transport, externalID)

	if err != nil {
		return "", validate.Wrap(op, err)
	}

	// Clock (UTC)
	now := time.Now
	if deps.Now != nil {
		now = deps.Now
	}

	nowUTC := now().UTC()

	sum, tax, qStart, qEnd, err := deps.QuarterSum.SumQuarter(ctx, userID, nowUTC)

	if err != nil {
		return "", validate.Wrap(op, err)
	}

	return TotalText(sum, tax, qStart, qEnd), nil
}
