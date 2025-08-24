package bot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/money"
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
		return "", fmt.Errorf("%s: upsert identity: %w", op, err)
	}

	// Clock (UTC)
	now := time.Now
	if deps.Now != nil {
		now = deps.Now
	}

	nowUTC := now().UTC()

	sum, tax, qStart, qEnd, err := deps.QuarterSum.SumQuarter(ctx, userID, nowUTC)

	if err != nil {
		return "", fmt.Errorf("%s: sum quarter: %w", op, err)
	}

	var b strings.Builder
	b.WriteString("Сумма за квартал: ")
	b.WriteString(qStart.Format("2006-01-02"))
	b.WriteString(" - ")
	b.WriteString(qEnd.Format("2006-01-02"))
	b.WriteString("\nСумма: ")
	b.WriteString(money.FormatAmountShort(sum))
	b.WriteString("\nНалог: ")
	b.WriteString(money.FormatAmountShort(tax))

	return b.String(), nil
}
