// internal/domain/totals.go
package domain

import "time"

// Totals is a DTO with period results for /total.
type Totals struct {
	From           time.Time // inclusive; UTC DATE is enforced in storage
	To             time.Time // inclusive
	IncomeSum      int64     // kopecks
	Tax            int64     // BaseRateBP% of IncomeSum, integer math
	ContribSum     int64     // payments type=contrib in [From,To]
	AdvanceSum     int64     // payments type=advance in [From,To]
	ContribApplied int64     // min(Tax, ContribSum)
	Due            int64     // max(0, Tax - ContribApplied - AdvanceSum)
}
