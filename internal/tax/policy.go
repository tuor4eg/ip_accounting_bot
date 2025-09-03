package tax

import (
	"sort"
	"time"

	"slices"

	"github.com/tuor4eg/ip_accounting_bot/internal/domain"
	"github.com/tuor4eg/ip_accounting_bot/internal/validate"
)

type Policy struct {
	BaseRateBP      int64
	ExcessThreshold int64
	ExcessRateBP    int64
}

type Provider interface {
	// ForDate returns policy for a tax scheme at a given moment.
	// Selection uses inclusive bounds on version intervals.
	ForDate(scheme domain.TaxScheme, date time.Time) (Policy, error)
}

type VersionedPolicy struct {
	ValidFrom time.Time  // inclusive
	ValidTo   *time.Time // inclusive, nil = open-ended
	Policy    Policy
}

type StaticProvider struct {
	versions map[string][]VersionedPolicy // key = scheme code, e.g., "usn_6"
}

func NewStaticProvider(versions map[string][]VersionedPolicy) *StaticProvider {
	out := make(map[string][]VersionedPolicy)

	for k, vs := range versions {
		cp := slices.Clone(vs)

		sort.Slice(cp, func(i, j int) bool { return cp[i].ValidFrom.Before(cp[j].ValidFrom) })
		out[k] = cp
	}
	return &StaticProvider{versions: out}
}

func (p *StaticProvider) ForDate(scheme domain.TaxScheme, date time.Time) (Policy, error) {
	const op = "tax.StaticProvider.ForDate"

	vs, ok := p.versions[string(scheme)]
	if !ok || len(vs) == 0 {
		return Policy{}, validate.Wrap(op, validate.ErrNotFound)
	}

	at := date.UTC()
	for _, v := range vs {
		// Inclusive bounds: [ValidFrom, ValidTo]
		if !at.Before(v.ValidFrom.UTC()) && (v.ValidTo == nil || !at.After(v.ValidTo.UTC())) {
			return v.Policy, nil
		}
	}

	return Policy{}, validate.Wrap(op, validate.ErrNotFound)
}
