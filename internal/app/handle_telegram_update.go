package app

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/tuor4eg/ip_accounting_bot/internal/bot"
	"github.com/tuor4eg/ip_accounting_bot/internal/telegram"
)

type TelegramSender interface {
	SendMessage(ctx context.Context, chatID int64, text string) error
}

// Normalize self: drop leading '@' if provided
func NormalizeSelf(self string) string {
	return strings.TrimPrefix(self, "@")
}

func HandleTelegramUpdate(
	ctx context.Context,
	self string,
	upd telegram.Update,
	sender TelegramSender,
	addDeps bot.AddDeps,
	totalDeps bot.TotalDeps,
) error {
	op := "app.HandleTelegramUpdate"

	if upd.Message == nil {
		return nil
	}

	text := strings.TrimSpace(upd.Message.Text)

	if text == "" {
		return nil
	}

	chatID := upd.Message.Chat.ID
	externalID := strconv.FormatInt(upd.Message.From.ID, 10)

	self = NormalizeSelf(self)

	reply, handled, err := bot.DispatchCommand(ctx, text, self, "telegram", externalID, addDeps, totalDeps)

	if !handled {
		return nil
	}

	if err != nil {
		if errors.Is(err, bot.ErrBadInput) {
			if sendErr := sender.SendMessage(ctx, chatID, bot.BadAmountHintText()); sendErr != nil {
				return fmt.Errorf("%s: send message: %w", op, sendErr)
			}

			return nil
		}

		if strings.Contains(err.Error(), "unknown command") {
			if sendErr := sender.SendMessage(ctx, chatID, bot.UnknownCommandText()); sendErr != nil {
				return fmt.Errorf("%s: send message: %w", op, sendErr)
			}

			return nil
		}
		if sendErr := sender.SendMessage(ctx, chatID, bot.ErrorText()); sendErr != nil {
			return fmt.Errorf("%s: send message: %w", op, sendErr)
		}

		return nil
	}

	if err := sender.SendMessage(ctx, chatID, reply); err != nil {
		return fmt.Errorf("%s: send message: %w", op, err)
	}

	return nil
}
