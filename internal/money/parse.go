package money

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	// ErrInvalidAmount is returned when the input cannot be parsed as a money amount.
	ErrInvalidAmount = errors.New("invalid amount")
	// ErrNegativeNotAllowed is returned when a negative amount is provided.
	ErrNegativeNotAllowed = errors.New("negative amount is not allowed")
	// ErrOverflow is returned when the parsed number would overflow int64.
	ErrOverflow = errors.New("amount overflow")
)

// ParseAmount converts human input like "1 234,56", "1.234,56", "1234.5", "1234", "1 234 ₽"
// into minor units (kopecks) as int64 without using floats.
// Rules:
//   - Ignores RUB tokens (₽, RUB, РУБ, р., руб.) at both ends.
//   - Treats spaces/NBSP/underscores/apostrophes as thousands separators.
//   - The LAST '.' or ',' is the decimal separator; others are grouping.
//   - Up to 2 decimal digits; fewer are padded (e.g., "12," -> 12.00).
//   - Negative values are rejected.
func ParseAmount(input string) (int64, error) {
	const op = "money.ParseAmount"
	if input == "" {
		return 0, fmt.Errorf("%w: empty input", ErrInvalidAmount)
	}

	// Normalize NBSP-like spaces to regular spaces and replace '₽' with " руб ".
	var sb strings.Builder
	sb.Grow(len(input))
	for _, r := range input {
		switch r {
		case '\u00A0', '\u202F', '\u2007':
			sb.WriteRune(' ')
		default:
			sb.WriteRune(r)
		}
	}
	s := strings.ReplaceAll(sb.String(), "₽", " руб ")
	s = strings.TrimSpace(s)

	// Reject negatives (incomes only).
	if strings.ContainsRune(s, '-') {
		return 0, fmt.Errorf("%s: %w", op, ErrNegativeNotAllowed)
	}

	// 1) Try explicit rub/kop tokens.
	if rub, kop, ok, err := parseRubKopTokens(s); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	} else if ok {
		if rub > math.MaxInt64/100 {
			return 0, fmt.Errorf("%s: %w", op, ErrOverflow)
		}
		rubPart := rub * 100
		if rubPart > math.MaxInt64-kop {
			return 0, fmt.Errorf("%s: %w", op, ErrOverflow)
		}
		return rubPart + kop, nil
	}

	// 2) Fallback to generic number parser.
	amt, err := parseGenericNumber(s)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return amt, nil
}

// parseRubKopTokens parses explicit "rubles/kopecks" tokens from s.
// Supported: "10р 50к", "10 руб 50 коп.", "10 rub 50", "50к", "10р".
// Case-insensitive; Russian abbreviations and Latin RUB/RUR are recognized.
// Returns (rub, kop, ok, err). ok=false means no tokens were found.
// Negative handling is NOT here; caller (ParseAmount) should reject it.
func parseRubKopTokens(s string) (rub, kop int64, ok bool, err error) {
	reRub := regexp.MustCompile(`(?i)(^|[\s])(\d+)\s*(?:руб(?:\.|ля|лей)?|р\.?|rub|rur)\b`)
	reKop := regexp.MustCompile(`(?i)(^|[\s])(\d+)\s*(?:коп(?:\.|еек|ейки|ей)?|к\.?)\b`)

	// Sum rubles tokens
	for _, m := range reRub.FindAllStringSubmatch(s, -1) {
		n, perr := strconv.ParseInt(m[2], 10, 64)

		if perr != nil {
			return 0, 0, false, fmt.Errorf("parse rub: %w", ErrInvalidAmount)
		}

		if rub > math.MaxInt64-n {
			return 0, 0, false, ErrOverflow
		}

		rub += n
		ok = true
	}

	// Sum kopecks tokens
	for _, m := range reKop.FindAllStringSubmatch(s, -1) {
		n, perr := strconv.ParseInt(m[2], 10, 64)

		if perr != nil {
			return 0, 0, false, fmt.Errorf("parse kop: %w", ErrInvalidAmount)
		}

		if kop > math.MaxInt64-n {
			return 0, 0, false, ErrOverflow
		}

		kop += n
		ok = true
	}

	return rub, kop, ok, nil
}

// parseGenericNumber parses amounts like "1 234,56", "1.234,56", "1234.5", "1234 руб".
// It treats the LAST '.' or ',' as the decimal separator (others as grouping).
// Grouping separators allowed: any Unicode space, underscore '_', apostrophe '\”.
// Currency tokens at the ends (руб, р., RUB, RUR) are stripped here.
// Negative handling is NOT here; caller should reject it.
// Returns amount in kopecks.
func parseGenericNumber(s string) (int64, error) {
	// Strip common RUB tokens from both ends (case-insensitive).
	tokens := []string{"руб.", "руб", "rur", "rub", "р.", "р"}
	for {
		ls := strings.ToLower(strings.TrimSpace(s))
		changed := false
		for _, t := range tokens {
			if strings.HasPrefix(ls, t) {
				s = strings.TrimSpace(s[len(t):])
				changed = true
				break
			}
			if strings.HasSuffix(ls, t) {
				s = strings.TrimSpace(s[:len(s)-len(t)])
				changed = true
				break
			}
		}
		if !changed {
			break
		}
	}

	s = strings.TrimSpace(s)
	if s == "" {
		return 0, ErrInvalidAmount
	}

	// Find last decimal separator among '.' and ','.
	runes := []rune(s)
	lastSep := -1
	for i, ch := range runes {
		if ch == '.' || ch == ',' {
			lastSep = i
		}
	}

	var total int64
	decimals := 0
	decimalMode := false

	for i, ch := range runes {
		switch {
		case ch >= '0' && ch <= '9':
			d := int64(ch - '0')
			if total > (math.MaxInt64-d)/10 {
				return 0, ErrOverflow
			}
			total = total*10 + d
			if decimalMode {
				decimals++
			}
		case ch == '.' || ch == ',':
			if i == lastSep {
				if decimalMode {
					return 0, ErrInvalidAmount
				}
				decimalMode = true
			}
			// else: treat as thousands separator — ignore
		default:
			// Ignore grouping separators: any space, underscore, apostrophe.
			if ch == '_' || ch == '\'' || unicode.IsSpace(ch) {
				continue
			}
			return 0, ErrInvalidAmount
		}
	}

	// Scale to kopecks deterministically.
	switch decimals {
	case 0:
		if total > math.MaxInt64/100 {
			return 0, ErrOverflow
		}
		return total * 100, nil
	case 1:
		if total > math.MaxInt64/10 {
			return 0, ErrOverflow
		}
		return total * 10, nil
	case 2:
		return total, nil
	default:
		return 0, ErrInvalidAmount
	}
}
