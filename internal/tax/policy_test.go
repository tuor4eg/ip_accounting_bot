package tax_test

import (
	"testing"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/domain"
	"github.com/tuor4eg/ip_accounting_bot/internal/tax"
)

func TestStaticProvider_ForDate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scheme domain.TaxScheme
		date   time.Time
		want   tax.Policy
	}{
		{domain.TaxSchemeUSN6, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), tax.Policy{BaseRateBP: 100, ExcessThreshold: 100000, ExcessRateBP: 100}},
		// Regular case with different date
		{domain.TaxSchemeUSN6, time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC), tax.Policy{BaseRateBP: 200, ExcessThreshold: 200000, ExcessRateBP: 150}},
		// Edge case - last day of year
		{domain.TaxSchemeUSN6, time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC), tax.Policy{BaseRateBP: 300, ExcessThreshold: 300000, ExcessRateBP: 200}},
		// Edge case - first day of year
		{domain.TaxSchemeUSN6, time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), tax.Policy{BaseRateBP: 400, ExcessThreshold: 400000, ExcessRateBP: 250}},
		// Error case - non-existent scheme
		{"invalid_scheme", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), tax.Policy{}},
		// Error case - date before any policy
		{domain.TaxSchemeUSN6, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), tax.Policy{}},
		// Error case - date after policy expiration
		{domain.TaxSchemeUSN6, time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC), tax.Policy{}},
	}

	for _, test := range tests {
		t.Run(string(test.scheme), func(t *testing.T) {
			got, err := tax.NewStaticProvider(map[string][]tax.VersionedPolicy{
				string(test.scheme): {{ValidFrom: test.date, ValidTo: &test.date, Policy: test.want}},
			}).ForDate(test.scheme, test.date)
			if err != nil {
				t.Errorf("ForDate() error = %v", err)
			}

			if got != test.want {
				t.Errorf("ForDate() got = %v, want %v", got, test.want)
			}
		})
	}
}
