package tax

import (
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/domain"
)

// Policy defines tax rates and thresholds
type Policy struct {
	BaseRateBP      int64
	ExcessThreshold int64
	ExcessRateBP    int64
}

// Provider interface for getting tax policies
type Provider interface {
	// ForDate returns policy for a tax scheme at a given moment.
	// Selection uses inclusive bounds on version intervals.
	ForDate(scheme domain.TaxScheme, date time.Time) (Policy, error)
}

// VersionedPolicy represents a policy valid for a specific time period
type VersionedPolicy struct {
	ValidFrom time.Time  // inclusive
	ValidTo   *time.Time // inclusive, nil = open-ended
	Policy    Policy
}

// StaticProvider provides static tax policies
type StaticProvider struct {
	versions map[string][]VersionedPolicy // key = scheme code, e.g., "usn_6"
}
