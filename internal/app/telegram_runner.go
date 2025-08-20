package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/telegram"
)

type TelegramRunner struct {
	tg *telegram.Client
}

func NewTelegramRunner(tg *telegram.Client) *TelegramRunner {
	return &TelegramRunner{tg: tg}
}

func (r *TelegramRunner) Name() string {
	return "telegram"
}

func (r *TelegramRunner) Run(ctx context.Context) error {
	pingCtx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	me, err := r.tg.GetMe(pingCtx)

	if err != nil {
		return fmt.Errorf("failed to get me: %w", err)
	}

	fmt.Printf("telegram: bot @%s (id=%d) is running\n", me.Username, me.ID)

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
				return fmt.Errorf("telegram: getUpdates error: %v\n", err)
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
					fmt.Printf("telegram: failed to send message: %v\n", sendErr)
				}
			}

			if u.UpdateID >= offset {
				offset = u.UpdateID + 1
			}
		}
	}

}
