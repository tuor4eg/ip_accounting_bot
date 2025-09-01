package period

import "time"

func YearBounds(t time.Time) (start time.Time, end time.Time) {
	year := t.Year()
	start = time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	end = time.Date(year, time.December, 31, 23, 59, 59, 0, time.UTC)
	return start, end
}
