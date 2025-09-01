package bot

import (
	"strconv"
	"strings"

	"github.com/tuor4eg/ip_accounting_bot/internal/money"
	"github.com/tuor4eg/ip_accounting_bot/internal/validate"
)

func ParseSlashCommand(text string, self string) (cmd string, args string, ok bool) {
	text = strings.TrimSpace(text)

	if !strings.HasPrefix(text, "/") {
		return "", "", false
	}

	after := text[1:] // remove leading '/'

	// find first whitespace to split token and args
	ws := strings.IndexFunc(after, func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n' || r == '\r'
	})

	token := after

	if ws >= 0 {
		token = after[:ws]
		args = strings.TrimLeft(after[ws:], " \t\r\n")
	} else {
		args = ""
	}

	if token == "" {
		return "", "", false
	}

	// optional @bot suffix
	cmdPart := token

	if at := strings.IndexByte(token, '@'); at >= 0 {
		cmdPart = token[:at]
		username := token[at+1:]
		// normalize self (allow passing with '@' by mistake)
		if self != "" {
			selfName := strings.TrimPrefix(self, "@")
			if username == "" || !strings.EqualFold(username, selfName) {
				return "", "", false
			}
		}
	}

	if cmdPart == "" {
		return "", "", false
	}

	// validate [A-Za-z0-9_]+ and lowercase
	for i := range len(cmdPart) {
		c := cmdPart[i]
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_') {
			return "", "", false
		}
	}
	cmd = strings.ToLower(cmdPart)

	return cmd, args, true
}

func ParseAmountAndNote(args string) (amount int64, note string, err error) {
	const op = "bot.parseAmountAndNote"

	args = strings.TrimSpace(args)

	if args == "" {
		return 0, "", validate.Wrap(op, ErrBadInput)
	}

	toks := strings.Fields(args)

	var (
		amountValue int64
		cut         = -1
	)

	// Try to find the best split point for amount and comment
	for i := 1; i <= len(toks); i++ {
		prefix := strings.Join(toks[:i], " ")
		v, err := money.ParseAmount(prefix)
		if err == nil {
			amountValue = v
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
		return 0, "", validate.Wrap(op, ErrBadInput)
	}

	note = strings.TrimSpace(strings.Join(toks[cut:], " "))

	return amountValue, note, nil
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
