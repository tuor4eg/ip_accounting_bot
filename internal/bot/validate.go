package bot

import "github.com/tuor4eg/ip_accounting_bot/internal/validate"

func validateEntryInput(amount int64, note string) error {
	const op = "bot.validateEntryInput"

	if err := validate.ValidateAmount(amount); err != nil {
		return validate.Wrap(op, err)
	}

	// TODO: validate note

	return nil
}
