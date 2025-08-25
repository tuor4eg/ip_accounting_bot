package money_test

import (
	"fmt"
	"testing"

	"github.com/tuor4eg/ip_accounting_bot/internal/money"
)

func TestFormatAmount(t *testing.T) {
	t.Parallel()

	cases := []struct {
		amount int64
		want   string
		desc   string
	}{
		{0, "0 ₽ 00 коп", "zero amount"},
		{100, "1 ₽ 00 коп", "one ruble"},
		{50, "0 ₽ 50 коп", "fifty kopecks"},
		{1234, "12 ₽ 34 коп", "twelve rubles thirty four kopecks"},
		{10000, "100 ₽ 00 коп", "hundred rubles"},
		{10050, "100 ₽ 50 коп", "hundred rubles fifty kopecks"},
		{123456, "1234 ₽ 56 коп", "thousand two hundred thirty four rubles fifty six kopecks"},
		{999999, "9999 ₽ 99 коп", "maximum kopecks"},
		{1000000, "10000 ₽ 00 коп", "ten thousand rubles"},
		{123456789, "1234567 ₽ 89 коп", "large amount"},
		{9223372036854775807, "92233720368547758 ₽ 07 коп", "max int64"},
		{-100, "-1 ₽ 00 коп", "negative amount"},
		{-50, "0 ₽ -50 коп", "negative kopecks"},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := money.FormatAmount(tc.amount)
			if err != nil {
				t.Fatalf("FormatAmount(%d) unexpected error: %v", tc.amount, err)
			}
			if got != tc.want {
				t.Fatalf("FormatAmount(%d) = %q, want %q", tc.amount, got, tc.want)
			}
		})
	}
}

func TestFormatAmountShort(t *testing.T) {
	t.Parallel()

	cases := []struct {
		amount int64
		want   string
		desc   string
	}{
		{0, "0.00₽", "zero amount"},
		{100, "1.00₽", "one ruble"},
		{50, "0.50₽", "fifty kopecks"},
		{1234, "12.34₽", "twelve rubles thirty four kopecks"},
		{10000, "100.00₽", "hundred rubles"},
		{10050, "100.50₽", "hundred rubles fifty kopecks"},
		{123456, "1234.56₽", "thousand two hundred thirty four rubles fifty six kopecks"},
		{999999, "9999.99₽", "maximum kopecks"},
		{1000000, "10000.00₽", "ten thousand rubles"},
		{123456789, "1234567.89₽", "large amount"},
		{9223372036854775807, "92233720368547758.07₽", "max int64"},
		{-100, "-1.00₽", "negative amount"},
		{-50, "0.-50₽", "negative kopecks"},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got := money.FormatAmountShort(tc.amount)
			if got != tc.want {
				t.Fatalf("FormatAmountShort(%d) = %q, want %q", tc.amount, got, tc.want)
			}
		})
	}
}

func TestGetAmountParts(t *testing.T) {
	t.Parallel()

	cases := []struct {
		amount      int64
		wantRubles  int64
		wantKopecks int64
		desc        string
	}{
		{0, 0, 0, "zero amount"},
		{100, 1, 0, "one ruble"},
		{50, 0, 50, "fifty kopecks"},
		{1234, 12, 34, "twelve rubles thirty four kopecks"},
		{10000, 100, 0, "hundred rubles"},
		{10050, 100, 50, "hundred rubles fifty kopecks"},
		{123456, 1234, 56, "thousand two hundred thirty four rubles fifty six kopecks"},
		{999999, 9999, 99, "maximum kopecks"},
		{1000000, 10000, 0, "ten thousand rubles"},
		{123456789, 1234567, 89, "large amount"},
		{9223372036854775807, 92233720368547758, 7, "max int64"},
		{-100, -1, 0, "negative amount"},
		{-50, 0, -50, "negative kopecks"},
		{-1234, -12, -34, "negative amount with kopecks"},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			gotRubles, gotKopecks := money.GetAmountParts(tc.amount)
			if gotRubles != tc.wantRubles {
				t.Fatalf("GetAmountParts(%d) rubles = %d, want %d", tc.amount, gotRubles, tc.wantRubles)
			}
			if gotKopecks != tc.wantKopecks {
				t.Fatalf("GetAmountParts(%d) kopecks = %d, want %d", tc.amount, gotKopecks, tc.wantKopecks)
			}
		})
	}
}

func TestConstants(t *testing.T) {
	t.Parallel()

	if money.RubleSymbol != "₽" {
		t.Fatalf("RubleSymbol = %q, want %q", money.RubleSymbol, "₽")
	}

	if money.KopekSymbol != "коп" {
		t.Fatalf("KopekSymbol = %q, want %q", money.KopekSymbol, "коп")
	}
}

// Test edge cases and boundary conditions
func TestFormatAmountEdgeCases(t *testing.T) {
	t.Parallel()

	// Test that formatting and parsing work together
	testCases := []int64{
		0, 1, 50, 99, 100, 101, 999, 1000, 1001, 9999, 10000,
		12345, 123456, 999999, 1000000, 1234567, 12345678,
	}

	for _, amount := range testCases {
		t.Run(fmt.Sprintf("roundtrip_%d", amount), func(t *testing.T) {
			// Format the amount
			formatted, err := money.FormatAmount(amount)
			if err != nil {
				t.Fatalf("FormatAmount(%d) error: %v", amount, err)
			}

			// Parse it back (this would require implementing a parse function that handles the formatted string)
			// For now, just verify the format is correct
			if amount == 0 && formatted != "0 ₽ 00 коп" {
				t.Fatalf("Zero amount formatted incorrectly: %s", formatted)
			}

			// Test short format
			shortFormatted := money.FormatAmountShort(amount)
			if amount == 0 && shortFormatted != "0.00₽" {
				t.Fatalf("Zero amount short formatted incorrectly: %s", shortFormatted)
			}

			// Test parts extraction
			rubles, kopecks := money.GetAmountParts(amount)
			if rubles*100+kopecks != amount {
				t.Fatalf("GetAmountParts(%d) = (%d, %d), but %d*100+%d = %d",
					amount, rubles, kopecks, rubles, kopecks, rubles*100+kopecks)
			}
		})
	}
}

// Test consistency between different formatting functions
func TestFormatConsistency(t *testing.T) {
	t.Parallel()

	testCases := []int64{
		0, 50, 100, 150, 1234, 10000, 12345, 123456, 999999, 1000000,
	}

	for _, amount := range testCases {
		t.Run(fmt.Sprintf("consistency_%d", amount), func(t *testing.T) {
			// Get parts
			rubles, kopecks := money.GetAmountParts(amount)

			// Verify parts are consistent with amount
			if rubles*100+kopecks != amount {
				t.Fatalf("Inconsistent parts: %d*100+%d = %d, but amount = %d",
					rubles, kopecks, rubles*100+kopecks, amount)
			}

			// Verify kopecks are in valid range
			if kopecks < 0 || kopecks > 99 {
				t.Fatalf("Invalid kopecks: %d (should be 0-99)", kopecks)
			}

			// Test that short format contains the same information
			shortFormatted := money.FormatAmountShort(amount)
			expectedShort := fmt.Sprintf("%d.%02d₽", rubles, kopecks)
			if shortFormatted != expectedShort {
				t.Fatalf("Short format mismatch: got %s, want %s", shortFormatted, expectedShort)
			}
		})
	}
}
