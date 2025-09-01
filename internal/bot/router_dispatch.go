package bot

import (
	"context"
	"fmt"

	"github.com/tuor4eg/ip_accounting_bot/internal/validate"
)

func DispatchCommand(
	ctx context.Context,
	text string,
	self string,
	transport string,
	externalID string,
	addDeps AddDeps,
	totalDeps TotalDeps,
) (reply string, handled bool, err error) {
	const op = "bot.DispatchCommand"

	cmd, args, ok := ParseSlashCommand(text, self)
	if !ok {
		// Not a slash-command (or addressed to another bot)
		return "", false, nil
	}

	switch cmd {
	case "start":
		return HandleStart(ctx), true, nil
	case "help":
		return HandleHelp(ctx), true, nil
	case "add":
		reply, err := HandleAdd(ctx, addDeps, transport, externalID, args)
		if err != nil {
			return "", true, validate.Wrap(op, err)
		}
		return reply, true, nil
	case "undo":
		reply, err := HandleUndo(ctx, addDeps, transport, externalID, args)
		if err != nil {
			return "", true, validate.Wrap(op, err)
		}
		return reply, true, nil
	case "add_contrib":
		reply, err := HandleAddContrib(ctx, addDeps, transport, externalID, args)
		if err != nil {
			return "", true, validate.Wrap(op, err)
		}
		return reply, true, nil
	case "undo_contrib":
		reply, err := HandleUndoContrib(ctx, addDeps, transport, externalID, args)
		if err != nil {
			fmt.Println("error in undo_contrib", err)
			return "", true, validate.Wrap(op, err)
		}
		return reply, true, nil
	case "add_advance":
		reply, err := HandleAddAdvance(ctx, addDeps, transport, externalID, args)
		if err != nil {
			return "", true, validate.Wrap(op, err)
		}
		return reply, true, nil
	case "undo_advance":
		reply, err := HandleUndoAdvance(ctx, addDeps, transport, externalID, args)
		if err != nil {
			return "", true, validate.Wrap(op, err)
		}
		return reply, true, nil
	case "total":
		reply, err := HandleTotal(ctx, totalDeps, transport, externalID, args)
		if err != nil {
			return "", true, validate.Wrap(op, err)
		}
		return reply, true, nil
	default:
		// Unknown command: handled=true
		return "", true, validate.Wrap(op, ErrUnknownCommand)
	}
}
