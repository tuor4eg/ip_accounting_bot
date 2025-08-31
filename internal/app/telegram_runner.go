package app

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/bot"
	"github.com/tuor4eg/ip_accounting_bot/internal/logging"
	"github.com/tuor4eg/ip_accounting_bot/internal/telegram"
)

const (
	codeTGStarted            = "tg_started"
	codeTGGetMeFailed        = "tg_getme_failed"
	codeTGGetUpdatesFailed   = "tg_getupdates_failed"
	codeTGSendFailed         = "tg_send_failed"
	codeTGHandleUpdateFailed = "tg_handle_update_failed"
)

type TelegramRunner struct {
	tg        *telegram.Client
	log       *slog.Logger
	addDeps   bot.AddDeps
	totalDeps bot.TotalDeps
}

func NewTelegramRunner(tg *telegram.Client) *TelegramRunner {
	tgRunner := &TelegramRunner{
		tg:  tg,
		log: logging.WithPackage(),
	}

	return tgRunner.SetBotDeps(bot.AddDeps{}, bot.TotalDeps{})
}

func (r *TelegramRunner) Name() string {
	return "telegram"
}

// SetBotDeps injects bot dependencies into the TelegramRunner and returns the runner for chaining.
func (r *TelegramRunner) SetBotDeps(add bot.AddDeps, total bot.TotalDeps) *TelegramRunner {
	r.addDeps = add
	r.totalDeps = total

	return r
}

func (r *TelegramRunner) SendMessage(ctx context.Context, chatID int64, text string) error {
	sentCtx, cancel := context.WithTimeout(ctx, tgSendTimeout)
	defer cancel()

	_, err := r.tg.SendMessage(sentCtx, telegram.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})

	return err
}

func (r *TelegramRunner) Run(ctx context.Context) error {
	pingCtx, cancel := context.WithTimeout(ctx, tgPingTimeout)
	defer cancel()

	me, err := r.tg.GetMe(pingCtx)

	if err != nil {
		r.log.Error("failed to get me", "code", codeTGGetMeFailed, "error", err)
		return err
	}

	self := NormalizeSelf(me.Username)

	r.log.Info("bot started", "username", self, "id", me.ID)

	var offset int64

	for {
		if ctx.Err() != nil {
			return nil
		}

		callCtx, cancel := context.WithTimeout(ctx, tgPollReqTimeout)

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
			if err := HandleTelegramUpdate(ctx, self, u, r, r.addDeps, r.totalDeps); err != nil {
				r.log.Error("handle telegram update", "code", codeTGHandleUpdateFailed, "error", err)
			}

			if u.UpdateID >= offset {
				offset = u.UpdateID + 1
			}
		}
	}

}
