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

	s := normalizeInput(input)

	// Reject negatives (incomes only).
	if strings.ContainsRune(s, '-') {
		return 0, fmt.Errorf("%s: %w", op, ErrNegativeNotAllowed)
	}

	// 1) Try explicit rub/kop tokens.
	if rub, kop, ok, err := parseRubKopTokens(s); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	} else if ok {
		return combineRublesAndKopecks(rub, kop)
	}

	// 2) Fallback to generic number parser.
	amt, err := parseGenericNumber(s)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return amt, nil
}

// normalizeInput normalizes the input string by converting NBSP-like spaces to regular spaces
// and replacing '₽' with " руб ".
func normalizeInput(input string) string {
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
	return strings.TrimSpace(s)
}

// combineRublesAndKopecks combines rubles and kopecks into total kopecks with overflow checking.
func combineRublesAndKopecks(rub, kop int64) (int64, error) {
	if rub > math.MaxInt64/100 {
		return 0, ErrOverflow
	}
	rubPart := rub * 100
	if rubPart > math.MaxInt64-kop {
		return 0, ErrOverflow
	}
	return rubPart + kop, nil
}

// parseRubKopTokens parses explicit "rubles/kopecks" tokens from s.
// Supported: "10р 50к", "10 руб 50 коп.", "10 rub 50", "50к", "10р".
// Case-insensitive; Russian abbreviations and Latin RUB/RUR are recognized.
// Returns (rub, kop, ok, err). ok=false means no tokens were found.
// Negative handling is NOT here; caller (ParseAmount) should reject it.
func parseRubKopTokens(s string) (rub, kop int64, ok bool, err error) {
	reRub := regexp.MustCompile(`(?i)(\d+(?:[.,\s]\d{3})*)\s*(?:руб(?:\.|ля|лей)?|р\.?|rub|rur)`)
	reKop := regexp.MustCompile(`(?i)(\d+)\s*(?:коп(?:\.|еек|ейки|ей)?|к\.?)`)

	// Parse rubles tokens
	rub, ok, err = parseRublesTokens(s, reRub)
	if err != nil {
		return 0, 0, false, err
	}

	// Parse kopecks tokens
	kop, kopOk, err := parseKopecksTokens(s, reKop)
	if err != nil {
		return 0, 0, false, err
	}
	ok = ok || kopOk

	// Handle case where there's a rubles token but no kopecks token
	if rub > 0 && kop == 0 {
		kop = findStandaloneKopecks(s, rub)
	}

	return rub, kop, ok, nil
}

// parseRublesTokens parses rubles tokens from the string.
func parseRublesTokens(s string, reRub *regexp.Regexp) (rub int64, ok bool, err error) {
	for _, m := range reRub.FindAllStringSubmatch(s, -1) {
		n, err := parseNumberWithSeparators(m[1])
		if err != nil {
			return 0, false, fmt.Errorf("parse rub: %w", ErrInvalidAmount)
		}

		if rub > math.MaxInt64-n {
			return 0, false, ErrOverflow
		}

		rub += n
		ok = true
	}
	return rub, ok, nil
}

// parseKopecksTokens parses kopecks tokens from the string.
func parseKopecksTokens(s string, reKop *regexp.Regexp) (kop int64, ok bool, err error) {
	for _, m := range reKop.FindAllStringSubmatch(s, -1) {
		n, err := strconv.ParseInt(m[1], 10, 64)
		if err != nil {
			return 0, false, fmt.Errorf("parse kop: %w", ErrInvalidAmount)
		}

		if kop > math.MaxInt64-n {
			return 0, false, ErrOverflow
		}

		kop += n
		ok = true
	}
	return kop, ok, nil
}

// parseNumberWithSeparators parses a number string that may contain thousands separators.
func parseNumberWithSeparators(numStr string) (int64, error) {
	// Remove thousands separators before parsing
	cleanNumStr := strings.ReplaceAll(numStr, ".", "")
	cleanNumStr = strings.ReplaceAll(cleanNumStr, ",", "")
	cleanNumStr = strings.ReplaceAll(cleanNumStr, " ", "")

	return strconv.ParseInt(cleanNumStr, 10, 64)
}

// findStandaloneKopecks looks for standalone numbers that could be kopecks when there's already a rubles token.
func findStandaloneKopecks(s string, rub int64) int64 {
	// Find all numbers in the string (with separators)
	reNum := regexp.MustCompile(`\d+(?:[.,\s]\d{3})*`)
	numbers := reNum.FindAllString(s, -1)

	// If we have exactly 2 numbers and one is already used for rubles,
	// the other should be kopecks
	if len(numbers) == 2 {
		// Find which number is not used for rubles
		for _, numStr := range numbers {
			// Remove separators for comparison
			cleanNumStr := strings.ReplaceAll(numStr, ".", "")
			cleanNumStr = strings.ReplaceAll(cleanNumStr, ",", "")
			cleanNumStr = strings.ReplaceAll(cleanNumStr, " ", "")

			// If this number is not the rubles number, it might be kopecks
			if cleanNumStr != fmt.Sprintf("%d", rub) {
				n, err := strconv.ParseInt(cleanNumStr, 10, 64)
				if err == nil && n < 100 { // Kopecks should be less than 100
					return n
				}
			}
		}
	}
	return 0
}

// parseGenericNumber parses amounts like "1 234,56", "1.234,56", "1234.5", "1234 руб".
// It treats the LAST '.' or ',' as the decimal separator (others as grouping).
// Grouping separators allowed: any Unicode space, underscore '_', apostrophe '\".
// Currency tokens at the ends (руб, р., RUB, RUR) are stripped here.
// Negative handling is NOT here; caller should reject it.
// Returns amount in kopecks.
func parseGenericNumber(s string) (int64, error) {
	s = stripCurrencyTokens(s)
	if s == "" {
		return 0, ErrInvalidAmount
	}

	runes := []rune(s)
	separators := findSeparators(runes)

	// Handle different separator scenarios
	switch len(separators) {
	case 0:
		return parseIntegerOnly(runes)
	case 1:
		return parseSingleSeparator(runes, separators[0])
	default:
		return parseMultipleSeparators(runes, separators)
	}
}

// stripCurrencyTokens removes currency tokens from the beginning and end of the string.
func stripCurrencyTokens(s string) string {
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
	return strings.TrimSpace(s)
}

// findSeparators finds all '.' and ',' separators in the rune slice.
func findSeparators(runes []rune) []int {
	var separators []int
	for i, ch := range runes {
		if ch == '.' || ch == ',' {
			separators = append(separators, i)
		}
	}
	return separators
}

// parseIntegerOnly parses a string with no separators as an integer.
func parseIntegerOnly(runes []rune) (int64, error) {
	var total int64
	for _, ch := range runes {
		if ch >= '0' && ch <= '9' {
			d := int64(ch - '0')
			if total > (math.MaxInt64-d)/10 {
				return 0, ErrOverflow
			}
			total = total*10 + d
		} else if ch == '_' || ch == '\'' || unicode.IsSpace(ch) {
			continue
		} else {
			return 0, ErrInvalidAmount
		}
	}
	if total > math.MaxInt64/100 {
		return 0, ErrOverflow
	}
	return total * 100, nil
}

// parseSingleSeparator handles the case with exactly one separator.
func parseSingleSeparator(runes []rune, sepPos int) (int64, error) {
	digitCount := countDigitsAfterSeparator(runes, sepPos)

	// If 1-2 digits after, treat as decimal separator
	if digitCount <= 2 {
		return parseDecimalNumber(runes, sepPos)
	} else {
		// Treat as thousands separator - remove it
		normalized := removeSeparatorAtPosition(runes, sepPos)
		return parseGenericNumber(normalized)
	}
}

// parseMultipleSeparators handles the case with multiple separators.
func parseMultipleSeparators(runes []rune, separators []int) (int64, error) {
	lastSep := separators[len(separators)-1]
	digitCount := countDigitsAfterSeparator(runes, lastSep)

	if digitCount <= 2 {
		// Last separator is decimal, others are thousands separators
		normalized := keepOnlyLastSeparator(runes, lastSep)
		return parseGenericNumber(normalized)
	} else {
		// All separators are thousands separators - remove them all
		normalized := removeAllSeparators(runes)
		return parseGenericNumber(normalized)
	}
}

// countDigitsAfterSeparator counts digits after a separator position.
func countDigitsAfterSeparator(runes []rune, sepPos int) int {
	digitCount := 0
	for i := sepPos + 1; i < len(runes); i++ {
		if runes[i] >= '0' && runes[i] <= '9' {
			digitCount++
		} else if unicode.IsSpace(runes[i]) || runes[i] == '_' || runes[i] == '\'' {
			continue
		} else {
			break
		}
	}
	return digitCount
}

// parseDecimalNumber parses a number with a decimal separator.
func parseDecimalNumber(runes []rune, sepPos int) (int64, error) {
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
			if i == sepPos {
				if decimalMode {
					return 0, ErrInvalidAmount
				}
				decimalMode = true
			}
		default:
			if ch == '_' || ch == '\'' || unicode.IsSpace(ch) {
				continue
			}
			return 0, ErrInvalidAmount
		}
	}

	return scaleToKopecks(total, decimals)
}

// scaleToKopecks scales the total to kopecks based on the number of decimal places.
func scaleToKopecks(total int64, decimals int) (int64, error) {
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

// removeSeparatorAtPosition removes a separator at a specific position.
func removeSeparatorAtPosition(runes []rune, sepPos int) string {
	var normalized strings.Builder
	for i, ch := range runes {
		if i != sepPos {
			normalized.WriteRune(ch)
		}
	}
	return normalized.String()
}

// keepOnlyLastSeparator keeps only the last separator and removes all others.
func keepOnlyLastSeparator(runes []rune, lastSep int) string {
	var normalized strings.Builder
	for i, ch := range runes {
		if ch == '.' || ch == ',' {
			if i == lastSep {
				normalized.WriteRune(ch) // Keep decimal separator
			}
			// Remove thousands separators
		} else {
			normalized.WriteRune(ch)
		}
	}
	return normalized.String()
}

// removeAllSeparators removes all separators from the string.
func removeAllSeparators(runes []rune) string {
	var normalized strings.Builder
	for _, ch := range runes {
		if ch != '.' && ch != ',' {
			normalized.WriteRune(ch)
		}
	}
	return normalized.String()
}
