package period_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/pkg/period"
)

// Test year bounds functionality
func TestYearBounds(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input         time.Time
		expectedStart time.Time
		expectedEnd   time.Time
		desc          string
	}{
		// Regular year tests
		{
			time.Date(2024, 1, 1, 12, 30, 45, 123456789, time.UTC),
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
			"start of year 2024",
		},
		{
			time.Date(2024, 6, 15, 8, 15, 30, 0, time.UTC),
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
			"middle of year 2024",
		},
		{
			time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
			"end of year 2024",
		},

		// Leap year tests
		{
			time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC),
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
			"leap year February 29",
		},

		// Non-leap year tests
		{
			time.Date(2023, 2, 28, 12, 0, 0, 0, time.UTC),
			time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC),
			"non-leap year February 28",
		},

		// Different timezones (should be normalized to UTC)
		{
			time.Date(2024, 6, 15, 12, 0, 0, 0, time.FixedZone("EST", -5*3600)),
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
			"different timezone (EST)",
		},
		{
			time.Date(2024, 6, 15, 12, 0, 0, 0, time.FixedZone("PST", -8*3600)),
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
			"different timezone (PST)",
		},

		// Edge cases around year boundaries
		{
			time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC),
			time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC),
			"year boundary - last moment of 2023",
		},
		{
			time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC),
			"year boundary - first moment of 2025",
		},

		// Century and millennium boundaries
		{
			time.Date(2000, 6, 15, 12, 0, 0, 0, time.UTC),
			time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2000, 12, 31, 23, 59, 59, 0, time.UTC),
			"century boundary year 2000",
		},
		{
			time.Date(1999, 12, 31, 23, 59, 59, 0, time.UTC),
			time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(1999, 12, 31, 23, 59, 59, 0, time.UTC),
			"century boundary - last moment of 1999",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			start, end := period.YearBounds(tc.input)

			if !start.Equal(tc.expectedStart) {
				t.Fatalf("YearBounds(%v) start = %v, want %v", tc.input, start, tc.expectedStart)
			}

			if !end.Equal(tc.expectedEnd) {
				t.Fatalf("YearBounds(%v) end = %v, want %v", tc.input, end, tc.expectedEnd)
			}

			// Additional validation: start should be before end
			if start.After(end) {
				t.Fatalf("YearBounds(%v) start (%v) is after end (%v)", tc.input, start, end)
			}

			// Additional validation: input should be within the returned bounds
			// Note: end time is at 23:59:59, so we need to check if input is before or equal to end
			if tc.input.Before(start) || tc.input.After(end) {
				t.Fatalf("YearBounds(%v) input is outside returned bounds [%v, %v]", tc.input, start, end)
			}

			// Additional validation: start should be at 00:00:00 UTC
			if start.Hour() != 0 || start.Minute() != 0 || start.Second() != 0 || start.Nanosecond() != 0 {
				t.Fatalf("YearBounds(%v) start time is not 00:00:00: %v", tc.input, start)
			}

			// Additional validation: end should be at 23:59:59 UTC
			if end.Hour() != 23 || end.Minute() != 59 || end.Second() != 59 || end.Nanosecond() != 0 {
				t.Fatalf("YearBounds(%v) end time is not 23:59:59: %v", tc.input, end)
			}

			// Additional validation: start and end should be in UTC
			if start.Location() != time.UTC {
				t.Fatalf("YearBounds(%v) start is not in UTC: %v", tc.input, start.Location())
			}

			if end.Location() != time.UTC {
				t.Fatalf("YearBounds(%v) end is not in UTC: %v", tc.input, end.Location())
			}
		})
	}
}

// Test year boundary consistency
func TestYearBoundaryConsistency(t *testing.T) {
	t.Parallel()

	// Test that consecutive years don't overlap and have no gaps
	testCases := []time.Time{
		time.Date(2023, 6, 15, 12, 0, 0, 0, time.UTC),
		time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC),
		time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC),
	}

	for i := 0; i < len(testCases)-1; i++ {
		t.Run(fmt.Sprintf("consecutive_years_%d_to_%d", 2023+i, 2024+i), func(t *testing.T) {
			_, end1 := period.YearBounds(testCases[i])
			start2, _ := period.YearBounds(testCases[i+1])

			// End of current year should be one second before start of next year
			expectedNextStart := end1.Add(time.Second)
			if !start2.Equal(expectedNextStart) {
				t.Fatalf("Year %d ends at %v, but year %d starts at %v (expected %v)",
					2023+i, end1, 2024+i, start2, expectedNextStart)
			}
		})
	}
}

// Test year duration consistency
func TestYearDurationConsistency(t *testing.T) {
	t.Parallel()

	testCases := []time.Time{
		time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC), // leap year
		time.Date(2023, 6, 15, 12, 0, 0, 0, time.UTC), // non-leap year
		time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC), // non-leap year
	}

	expectedDurations := []int{
		366, // 2024: leap year (366 days)
		365, // 2023: non-leap year (365 days)
		365, // 2025: non-leap year (365 days)
	}

	for i, input := range testCases {
		t.Run(fmt.Sprintf("year_duration_%d", input.Year()), func(t *testing.T) {
			start, end := period.YearBounds(input)
			duration := end.Sub(start)
			actualDays := int(duration/(24*time.Hour)) + 1 // +1 because bounds are inclusive

			if actualDays != expectedDurations[i] {
				t.Fatalf("Year %d duration = %d days, want %d days", input.Year(), actualDays, expectedDurations[i])
			}
		})
	}
}

// Test timezone normalization for years
func TestYearTimezoneNormalization(t *testing.T) {
	t.Parallel()

	// Test that different timezones produce the same year bounds
	baseTime := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	timezones := []*time.Location{
		time.UTC,
		time.FixedZone("EST", -5*3600),
		time.FixedZone("PST", -8*3600),
		time.FixedZone("JST", 9*3600),
		time.FixedZone("IST", 5*3600+30*60),
	}

	expectedStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedEnd := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	for _, tz := range timezones {
		t.Run(fmt.Sprintf("timezone_%s", tz.String()), func(t *testing.T) {
			input := baseTime.In(tz)
			start, end := period.YearBounds(input)

			if !start.Equal(expectedStart) {
				t.Fatalf("YearBounds(%v in %s) start = %v, want %v", baseTime, tz, start, expectedStart)
			}

			if !end.Equal(expectedEnd) {
				t.Fatalf("YearBounds(%v in %s) end = %v, want %v", baseTime, tz, end, expectedEnd)
			}
		})
	}
}
