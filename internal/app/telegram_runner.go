package app

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/logging"
	"github.com/tuor4eg/ip_accounting_bot/internal/telegram"
)

const (
	codeTGStarted          = "tg_started"
	codeTGGetMeFailed      = "tg_getme_failed"
	codeTGGetUpdatesFailed = "tg_getupdates_failed"
	codeTGSendFailed       = "tg_send_failed"
)

type TelegramRunner struct {
	tg  *telegram.Client
	log *slog.Logger
}

func NewTelegramRunner(tg *telegram.Client) *TelegramRunner {
	return &TelegramRunner{
		tg:  tg,
		log: logging.WithPackage(),
	}
}

func (r *TelegramRunner) Name() string {
	return "telegram"
}

func (r *TelegramRunner) Run(ctx context.Context) error {
	pingCtx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	me, err := r.tg.GetMe(pingCtx)

	if err != nil {
		r.log.Error("failed to get me", "code", codeTGGetMeFailed, "error", err)
		return err
	}

	r.log.Info("bot started", "username", me.Username, "id", me.ID)

	var offset int64

	for {
		if ctx.Err() != nil {
			return nil
		}

		callCtx, cancel := context.WithTimeout(ctx, 35*time.Second)

		updates, err := r.tg.GetUpdates(callCtx, telegram.GetUpdatesParams{
			Offset:         offset,
			Timeout:        30,
			AllowedUpdates: []string{"message"},
		})
		cancel()

		if err != nil {
			if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
				r.log.Error("getUpdates error", "code", codeTGGetUpdatesFailed, "error", err)

				return err
			}

			select {
			case <-time.After(75 * time.Millisecond):
				continue
			case <-ctx.Done():
				return nil
			}
		}

		if len(updates) == 0 {
			continue
		}

		for _, u := range updates {
			if u.Message != nil && u.Message.Text != "" {
				msg := u.Message

				sentCtx, cancel := context.WithTimeout(ctx, 5*time.Second)

				_, sendErr := r.tg.SendMessage(sentCtx, telegram.SendMessageParams{
					ChatID:           msg.Chat.ID,
					Text:             msg.Text,
					ReplyToMessageID: msg.MessageID,
				})

				cancel()

				if sendErr != nil && ctx.Err() == nil {
					r.log.Error("failed to send message", "code", codeTGSendFailed, "error", sendErr)
				}
			}

			if u.UpdateID >= offset {
				offset = u.UpdateID + 1
			}
		}
	}

}
