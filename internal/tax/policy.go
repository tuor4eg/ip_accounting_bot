package tax

import (
	"sort"
	"time"

	"slices"

	"github.com/tuor4eg/ip_accounting_bot/internal/domain"
	"github.com/tuor4eg/ip_accounting_bot/internal/validate"
)

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
