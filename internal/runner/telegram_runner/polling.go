package telegram_runner

import (
	"context"
	"errors"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/bot"
)

const (
	tgGetUpdatesTimeoutSec = 30
	tgPollReqTimeout       = 35 * time.Second
	tgSendTimeout          = 5 * time.Second
	tgPingTimeout          = 8 * time.Second
)

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
	botDeps *bot.BotDeps,
) error {
	self = NormalizeSelf(self)

	offset := int64(0)

	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		ctxPoll, cancel := context.WithTimeout(ctx, tgPollReqTimeout)

		updates, err := client.GetUpdates(ctxPoll, offset, tgGetUpdatesTimeoutSec)

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

			if err := HandleTelegramUpdate(ctx, self, update, sender, botDeps); err != nil {

				continue
			}
		}
	}
}
