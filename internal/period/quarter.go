package period

import "time"

// QuarterBounds returns the inclusive [start, end] date range for the calendar quarter of t.
// Dates are returned at 00:00:00 in UTC; storage casts dates to DATE, so time-of-day is irrelevant.
func QuarterBounds(t time.Time) (start time.Time, end time.Time) {
	// Normalize to UTC to avoid DST artifacts when doing date math.
	tt := t.In(time.UTC)
	y, m, _ := tt.Date()

	// Determine the first month of the quarter: Jan(1), Apr(4), Jul(7), Oct(10).
	qStartMonth := time.Month(((int(m)-1)/3)*3 + 1)

	// Start = first day of the quarter.
	start = time.Date(y, qStartMonth, 1, 0, 0, 0, 0, time.UTC)

	// End = last day of the quarter: next quarter's first day minus 1 day.
	nextQStart := start.AddDate(0, 3, 0)
	end = nextQStart.AddDate(0, 0, -1)

	return start, end
}
