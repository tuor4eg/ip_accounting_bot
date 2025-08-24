package app

import (
	"context"
	"errors"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/bot"
	"github.com/tuor4eg/ip_accounting_bot/internal/telegram"
)

type TelegramUpdateGetter interface {
	GetUpdates(ctx context.Context, offset int64, timeoutSec int) ([]telegram.Update, error)
}

func nextOffset(cur int64, updID int64) int64 {
	if v := updID + 1; v > cur {
		return v
	}
	return cur
}

func RunTelegramPolling(
	ctx context.Context,
	self string,
	client TelegramUpdateGetter,
	sender TelegramSender,
	addDeps bot.AddDeps,
	totalDeps bot.TotalDeps,
) error {
	self = NormalizeSelf(self)

	offset := int64(0)

	timeoutSec := 30

	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		ctxPoll, cancel := context.WithTimeout(ctx, 35*time.Second)

		updates, err := client.GetUpdates(ctxPoll, offset, timeoutSec)

		cancel()

		if err != nil {
			if errors.Is(ctx.Err(), context.Canceled) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return ctx.Err()
			}

			time.Sleep(1 * time.Second)

			continue
		}

		for _, update := range updates {
			offset = nextOffset(offset, update.UpdateID)

			if err := HandleTelegramUpdate(ctx, self, update, sender, addDeps, totalDeps); err != nil {

				continue
			}
		}
	}
}
