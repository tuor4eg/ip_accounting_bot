package bot

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/money"
	"github.com/tuor4eg/ip_accounting_bot/internal/validate"
)

type IdentityStore interface {
	UpsertIdentity(ctx context.Context, transport, externalID string) (int64, error)
}

type IncomeAdder interface {
	AddIncome(ctx context.Context, userID int64, at time.Time, amount int64, note string) error
}

type AddDeps struct {
	Identities IdentityStore
	Income     IncomeAdder
	Now        func() time.Time
}

func HandleAdd(ctx context.Context, deps AddDeps, transport, externalID string, args string) (string, error) {
	const op = "bot.HandleAdd"

	args = strings.TrimSpace(args)

	if args == "" {
		return "", validate.Wrap(op, ErrBadInput)
	}

	toks := strings.Fields(args)

	var (
		amount int64
		cut    = -1
	)

	// Try to find the best split point for amount and comment
	for i := 1; i <= len(toks); i++ {
		prefix := strings.Join(toks[:i], " ")
		v, err := money.ParseAmount(prefix)
		if err == nil {
			amount = v
			cut = i

			// Check if the next token looks like a comment (not a number or currency token)
			if i < len(toks) {
				nextToken := toks[i]
				// If next token is not a number and doesn't look like currency token,
				// this is likely the end of amount
				if !isCurrencyToken(nextToken) && !isNumber(nextToken) {
					break
				}
			}
		}
	}
	if cut == -1 {
		return "", validate.Wrap(op, ErrBadInput)
	}

	if amount == 0 {
		return "", validate.Wrap(op, ErrAmountIsZero)
	}

	note := strings.TrimSpace(strings.Join(toks[cut:], " "))

	// Resolve or create user identity.
	userID, err := deps.Identities.UpsertIdentity(ctx, transport, externalID)
	if err != nil {
		return "", validate.Wrap(op, err)
	}

	// Use UTC "now"; storage casts to DATE.
	now := time.Now
	if deps.Now != nil {
		now = deps.Now
	}
	at := now().UTC()

	// Persist income.
	if err := deps.Income.AddIncome(ctx, userID, at, amount, note); err != nil {
		return "", validate.Wrap(op, err)
	}

	return AddSuccessText(amount, at, note), nil
}

// isCurrencyToken checks if a token looks like a currency token
func isCurrencyToken(token string) bool {
	token = strings.ToLower(strings.TrimSpace(token))
	currencyTokens := []string{"р", "руб", "руб.", "rub", "rur", "к", "коп", "коп.", "копеек", "копейки"}
	for _, ct := range currencyTokens {
		if token == ct {
			return true
		}
	}
	// Also check if token ends with currency token (for cases like "50к", "100р")
	for _, ct := range currencyTokens {
		if strings.HasSuffix(token, ct) {
			return true
		}
	}
	return false
}

// isNumber checks if a token looks like a number
func isNumber(token string) bool {
	token = strings.TrimSpace(token)
	// Check if it's a pure number
	if _, err := strconv.ParseInt(token, 10, 64); err == nil {
		return true
	}
	// Check if it's a number with separators (like 1,234 or 1.234)
	// Remove common separators and check if the result is a number
	clean := strings.ReplaceAll(token, ",", "")
	clean = strings.ReplaceAll(clean, ".", "")
	clean = strings.ReplaceAll(clean, " ", "")
	if _, err := strconv.ParseInt(clean, 10, 64); err == nil {
		return true
	}
	return false
}
