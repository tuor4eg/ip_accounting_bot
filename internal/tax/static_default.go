// internal/tax/static_default.go
package tax

import "time"

// NewDefaultProvider returns a static Provider with a single open-ended policy
// for scheme "usn_6": base rate 6% (600 bp), extra threshold 300_000 ₽,
// extra rate 1% (100 bp). Boundaries are inclusive.
func NewDefaultProvider() Provider {
	v := VersionedPolicy{
		ValidFrom: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), // inclusive
		ValidTo:   nil,                                         // open-ended
		Policy: Policy{
			BaseRateBP:      600,        // 6% in basis points
			ExcessThreshold: 300_000_00, // 300,000 RUB in kopecks
			ExcessRateBP:    100,        // 1% in basis points
		},
	}
	return NewStaticProvider(map[string][]VersionedPolicy{
		// If you have a domain constant (e.g., domain.TaxSchemeUSN6),
		// you can replace the string literal with it.
		"usn_6": {v},
	})
}
