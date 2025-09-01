package tax_test

import (
	"fmt"
	"testing"

	"github.com/tuor4eg/ip_accounting_bot/internal/tax"
)

func TestExtraOverThreshold(t *testing.T) {
	tests := []struct {
		yearIncome int64
		threshold  int64
		rateBP     int64
		want       int64
	}{
		{100000, 100000, 100, 0},
		{100000, 100000, 100, 0},
		{100000, 100000, 100, 0},
		{100000, 100000, 100, 0},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("yearIncome=%d, threshold=%d, rateBP=%d", test.yearIncome, test.threshold, test.rateBP), func(t *testing.T) {
			got := tax.ExtraOverThreshold(test.yearIncome, test.threshold, test.rateBP)
			if got != test.want {
				t.Errorf("ExtraOverThreshold(%d, %d, %d) = %d, want %d", test.yearIncome, test.threshold, test.rateBP, got, test.want)
			}
		})
	}
}
