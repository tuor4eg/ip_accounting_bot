package telegram_runner

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/tuor4eg/ip_accounting_bot/internal/bot"
	"github.com/tuor4eg/ip_accounting_bot/internal/telegram"
	"github.com/tuor4eg/ip_accounting_bot/internal/validate"
)

// Normalize self: drop leading '@' if provided
func NormalizeSelf(self string) string {
	return strings.TrimPrefix(self, "@")
}

func HandleTelegramUpdate(
	ctx context.Context,
	self string,
	upd telegram.Update,
	sender TelegramSender,
	botDeps *bot.BotDeps,
) error {
	op := "telegram.HandleTelegramUpdate"

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

	reply, handled, err := bot.DispatchCommand(ctx, text, self, "telegram", externalID, botDeps)

	if !handled {
		if sendErr := sender.SendMessage(ctx, chatID, bot.UnknownCommandText()); sendErr != nil {
			return validate.Wrap(op, sendErr)
		}

		return nil
	}

	if err != nil {
		if errors.Is(err, bot.ErrBadInput) {
			if sendErr := sender.SendMessage(ctx, chatID, bot.BadAmountHintText()); sendErr != nil {
				return validate.Wrap(op, sendErr)
			}

			return nil
		}

		if errors.Is(err, bot.ErrAmountIsZero) {
			if sendErr := sender.SendMessage(ctx, chatID, bot.AmountIsZeroText()); sendErr != nil {
				return validate.Wrap(op, sendErr)
			}

			return nil
		}

		if strings.Contains(err.Error(), "unknown command") {
			if sendErr := sender.SendMessage(ctx, chatID, bot.UnknownCommandText()); sendErr != nil {
				return validate.Wrap(op, sendErr)
			}

			return nil
		}
		if sendErr := sender.SendMessage(ctx, chatID, bot.ErrorText()); sendErr != nil {
			return validate.Wrap(op, sendErr)
		}

		return nil
	}

	if err := sender.SendMessage(ctx, chatID, reply); err != nil {
		return validate.Wrap(op, err)
	}

	return nil
}
