package bot

import (
	"context"
	"strings"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/validate"
	"github.com/tuor4eg/ip_accounting_bot/pkg/period"
)

func HandleTotal(ctx context.Context, deps *BotDeps, transport, externalID, args string) (string, error) {
	const op = "bot.HandleTotal"

	// TODO: parse args
	_ = strings.TrimSpace(args)

	userID, err := deps.Identities.UpsertIdentity(ctx, transport, externalID, 0)

	if err != nil {
		return "", validate.Wrap(op, err)
	}

	// Clock (UTC)
	now := time.Now
	if deps.Now != nil {
		now = deps.Now
	}

	nowUTC := now().UTC()

	year, quarter := period.QuarterOf(nowUTC)

	QuarterTotals, err := deps.Total.SumQuarter(ctx, userID, nowUTC)

	if err != nil {
		return "", validate.Wrap(op, err)
	}

	YearToDateTotals, err := deps.Total.SumYearToDate(ctx, userID, nowUTC)

	if err != nil {
		return "", validate.Wrap(op, err)
	}

	return TotalText(
		QuarterTotals.IncomeSum,
		QuarterTotals.Tax,
		QuarterTotals.From,
		QuarterTotals.To,
		YearToDateTotals.IncomeSum,
		YearToDateTotals.Tax,
		YearToDateTotals.ContribSum,
		YearToDateTotals.AdvanceSum,
		year, quarter,
	), nil
}
