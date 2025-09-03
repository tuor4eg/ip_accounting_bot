package period_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/tuor4eg/ip_accounting_bot/internal/period"
)

func TestQuarterBounds(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input         time.Time
		expectedStart time.Time
		expectedEnd   time.Time
		desc          string
	}{
		// Q1 (January-March)
		{
			time.Date(2024, 1, 1, 12, 30, 45, 123456789, time.UTC),
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC),
			"Q1 start of year",
		},
		{
			time.Date(2024, 2, 15, 8, 15, 30, 0, time.UTC),
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC),
			"Q1 middle of February",
		},
		{
			time.Date(2024, 3, 31, 23, 59, 59, 999999999, time.UTC),
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC),
			"Q1 end of March",
		},

		// Q2 (April-June)
		{
			time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC),
			"Q2 start of April",
		},
		{
			time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC),
			time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC),
			"Q2 middle of May",
		},
		{
			time.Date(2024, 6, 30, 23, 59, 59, 0, time.UTC),
			time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC),
			"Q2 end of June",
		},

		// Q3 (July-September)
		{
			time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 9, 30, 0, 0, 0, 0, time.UTC),
			"Q3 start of July",
		},
		{
			time.Date(2024, 8, 15, 15, 30, 45, 0, time.UTC),
			time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 9, 30, 0, 0, 0, 0, time.UTC),
			"Q3 middle of August",
		},
		{
			time.Date(2024, 9, 30, 23, 59, 59, 999999999, time.UTC),
			time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 9, 30, 0, 0, 0, 0, time.UTC),
			"Q3 end of September",
		},

		// Q4 (October-December)
		{
			time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
			"Q4 start of October",
		},
		{
			time.Date(2024, 11, 15, 10, 20, 30, 0, time.UTC),
			time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
			"Q4 middle of November",
		},
		{
			time.Date(2024, 12, 31, 23, 59, 59, 999999999, time.UTC),
			time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
			"Q4 end of year",
		},

		// Leap year tests
		{
			time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC),
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC),
			"leap year February 29",
		},

		// Non-leap year tests
		{
			time.Date(2023, 2, 28, 12, 0, 0, 0, time.UTC),
			time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 3, 31, 0, 0, 0, 0, time.UTC),
			"non-leap year February 28",
		},

		// Different timezones (should be normalized to UTC)
		{
			time.Date(2024, 6, 15, 12, 0, 0, 0, time.FixedZone("EST", -5*3600)),
			time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC),
			"different timezone (EST)",
		},
		{
			time.Date(2024, 6, 15, 12, 0, 0, 0, time.FixedZone("PST", -8*3600)),
			time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC),
			"different timezone (PST)",
		},

		// Edge cases around quarter boundaries
		{
			time.Date(2024, 3, 31, 23, 59, 59, 999999999, time.UTC),
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC),
			"Q1 boundary - last moment of March",
		},
		{
			time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC),
			"Q2 boundary - first moment of April",
		},
		{
			time.Date(2024, 6, 30, 23, 59, 59, 999999999, time.UTC),
			time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC),
			"Q2 boundary - last moment of June",
		},
		{
			time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 9, 30, 0, 0, 0, 0, time.UTC),
			"Q3 boundary - first moment of July",
		},
		{
			time.Date(2024, 9, 30, 23, 59, 59, 999999999, time.UTC),
			time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 9, 30, 0, 0, 0, 0, time.UTC),
			"Q3 boundary - last moment of September",
		},
		{
			time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
			"Q4 boundary - first moment of October",
		},
		{
			time.Date(2024, 12, 31, 23, 59, 59, 999999999, time.UTC),
			time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
			"Q4 boundary - last moment of December",
		},

		// Year boundary tests
		{
			time.Date(2023, 12, 31, 23, 59, 59, 999999999, time.UTC),
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
			"year boundary - last moment of 2023",
		},
		{
			time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2025, 3, 31, 0, 0, 0, 0, time.UTC),
			"year boundary - first moment of 2025",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			start, end := period.QuarterBounds(tc.input)

			if !start.Equal(tc.expectedStart) {
				t.Fatalf("QuarterBounds(%v) start = %v, want %v", tc.input, start, tc.expectedStart)
			}

			if !end.Equal(tc.expectedEnd) {
				t.Fatalf("QuarterBounds(%v) end = %v, want %v", tc.input, end, tc.expectedEnd)
			}

			// Additional validation: start should be before or equal to end
			if start.After(end) {
				t.Fatalf("QuarterBounds(%v) start (%v) is after end (%v)", tc.input, start, end)
			}

			// Additional validation: input should be within the returned bounds
			// Note: end time is at 00:00:00, so we need to check if input is before the next day
			nextDayAfterEnd := end.AddDate(0, 0, 1)
			if tc.input.Before(start) || !tc.input.Before(nextDayAfterEnd) {
				t.Fatalf("QuarterBounds(%v) input is outside returned bounds [%v, %v)", tc.input, start, nextDayAfterEnd)
			}

			// Additional validation: start and end should be at 00:00:00 UTC
			if start.Hour() != 0 || start.Minute() != 0 || start.Second() != 0 || start.Nanosecond() != 0 {
				t.Fatalf("QuarterBounds(%v) start time is not 00:00:00: %v", tc.input, start)
			}

			if end.Hour() != 0 || end.Minute() != 0 || end.Second() != 0 || end.Nanosecond() != 0 {
				t.Fatalf("QuarterBounds(%v) end time is not 00:00:00: %v", tc.input, end)
			}

			// Additional validation: start and end should be in UTC
			if start.Location() != time.UTC {
				t.Fatalf("QuarterBounds(%v) start is not in UTC: %v", tc.input, start.Location())
			}

			if end.Location() != time.UTC {
				t.Fatalf("QuarterBounds(%v) end is not in UTC: %v", tc.input, end.Location())
			}
		})
	}
}

// Test quarter boundary consistency
func TestQuarterBoundaryConsistency(t *testing.T) {
	t.Parallel()

	// Test that consecutive quarters don't overlap and have no gaps
	testCases := []time.Time{
		time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),  // Q1
		time.Date(2024, 4, 15, 12, 0, 0, 0, time.UTC),  // Q2
		time.Date(2024, 7, 15, 12, 0, 0, 0, time.UTC),  // Q3
		time.Date(2024, 10, 15, 12, 0, 0, 0, time.UTC), // Q4
	}

	for i := 0; i < len(testCases)-1; i++ {
		t.Run(fmt.Sprintf("consecutive_quarters_%d_to_%d", i+1, i+2), func(t *testing.T) {
			_, end1 := period.QuarterBounds(testCases[i])
			start2, _ := period.QuarterBounds(testCases[i+1])

			// End of current quarter should be one day before start of next quarter
			expectedNextStart := end1.AddDate(0, 0, 1)
			if !start2.Equal(expectedNextStart) {
				t.Fatalf("Quarter %d ends at %v, but quarter %d starts at %v (expected %v)",
					i+1, end1, i+2, start2, expectedNextStart)
			}
		})
	}
}

// Test quarter duration consistency
func TestQuarterDurationConsistency(t *testing.T) {
	t.Parallel()

	testCases := []time.Time{
		time.Date(2024, 2, 15, 12, 0, 0, 0, time.UTC),  // Q1 (leap year)
		time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC),  // Q2
		time.Date(2024, 8, 15, 12, 0, 0, 0, time.UTC),  // Q3
		time.Date(2024, 11, 15, 12, 0, 0, 0, time.UTC), // Q4
	}

	expectedDurations := []int{
		91, // Q1: Jan(31) + Feb(29) + Mar(31) = 91 days (leap year)
		91, // Q2: Apr(30) + May(31) + Jun(30) = 91 days
		92, // Q3: Jul(31) + Aug(31) + Sep(30) = 92 days
		92, // Q4: Oct(31) + Nov(30) + Dec(31) = 92 days
	}

	for i, input := range testCases {
		t.Run(fmt.Sprintf("quarter_duration_%d", i+1), func(t *testing.T) {
			start, end := period.QuarterBounds(input)
			duration := end.Sub(start) / (24 * time.Hour)
			actualDays := int(duration) + 1 // +1 because bounds are inclusive

			if actualDays != expectedDurations[i] {
				t.Fatalf("Quarter %d duration = %d days, want %d days", i+1, actualDays, expectedDurations[i])
			}
		})
	}
}

// Test non-leap year quarter durations
func TestNonLeapYearQuarterDurations(t *testing.T) {
	t.Parallel()

	testCases := []time.Time{
		time.Date(2023, 2, 15, 12, 0, 0, 0, time.UTC),  // Q1 (non-leap year)
		time.Date(2023, 5, 15, 12, 0, 0, 0, time.UTC),  // Q2
		time.Date(2023, 8, 15, 12, 0, 0, 0, time.UTC),  // Q3
		time.Date(2023, 11, 15, 12, 0, 0, 0, time.UTC), // Q4
	}

	expectedDurations := []int{
		90, // Q1: Jan(31) + Feb(28) + Mar(31) = 90 days (non-leap year)
		91, // Q2: Apr(30) + May(31) + Jun(30) = 91 days
		92, // Q3: Jul(31) + Aug(31) + Sep(30) = 92 days
		92, // Q4: Oct(31) + Nov(30) + Dec(31) = 92 days
	}

	for i, input := range testCases {
		t.Run(fmt.Sprintf("non_leap_quarter_duration_%d", i+1), func(t *testing.T) {
			start, end := period.QuarterBounds(input)
			duration := end.Sub(start) / (24 * time.Hour)
			actualDays := int(duration) + 1 // +1 because bounds are inclusive

			if actualDays != expectedDurations[i] {
				t.Fatalf("Quarter %d duration = %d days, want %d days", i+1, actualDays, expectedDurations[i])
			}
		})
	}
}

// Test timezone normalization
func TestTimezoneNormalization(t *testing.T) {
	t.Parallel()

	// Test that different timezones produce the same quarter bounds
	baseTime := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	timezones := []*time.Location{
		time.UTC,
		time.FixedZone("EST", -5*3600),
		time.FixedZone("PST", -8*3600),
		time.FixedZone("JST", 9*3600),
		time.FixedZone("IST", 5*3600+30*60),
	}

	expectedStart := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)
	expectedEnd := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)

	for _, tz := range timezones {
		t.Run(fmt.Sprintf("timezone_%s", tz.String()), func(t *testing.T) {
			input := baseTime.In(tz)
			start, end := period.QuarterBounds(input)

			if !start.Equal(expectedStart) {
				t.Fatalf("QuarterBounds(%v in %s) start = %v, want %v", baseTime, tz, start, expectedStart)
			}

			if !end.Equal(expectedEnd) {
				t.Fatalf("QuarterBounds(%v in %s) end = %v, want %v", baseTime, tz, end, expectedEnd)
			}
		})
	}
}

// Test edge cases around DST transitions
func TestDSTEdgeCases(t *testing.T) {
	t.Parallel()

	// Test dates around DST transitions in different timezones
	testCases := []struct {
		input         time.Time
		expectedStart time.Time
		expectedEnd   time.Time
		desc          string
	}{
		{
			time.Date(2024, 3, 10, 2, 30, 0, 0, time.FixedZone("EST", -5*3600)), // DST start
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC),
			"DST start in EST",
		},
		{
			time.Date(2024, 11, 3, 2, 30, 0, 0, time.FixedZone("EST", -5*3600)), // DST end
			time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
			"DST end in EST",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			start, end := period.QuarterBounds(tc.input)

			if !start.Equal(tc.expectedStart) {
				t.Fatalf("QuarterBounds(%v) start = %v, want %v", tc.input, start, tc.expectedStart)
			}

			if !end.Equal(tc.expectedEnd) {
				t.Fatalf("QuarterBounds(%v) end = %v, want %v", tc.input, end, tc.expectedEnd)
			}
		})
	}
}

// Test QuarterOf function
func TestQuarterOf(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input           time.Time
		expectedYear    int
		expectedQuarter int
		desc            string
	}{
		// Q1 (January-March)
		{
			time.Date(2024, 1, 1, 12, 30, 45, 123456789, time.UTC),
			2024, 1,
			"Q1 start of year",
		},
		{
			time.Date(2024, 2, 15, 8, 15, 30, 0, time.UTC),
			2024, 1,
			"Q1 middle of February",
		},
		{
			time.Date(2024, 3, 31, 23, 59, 59, 999999999, time.UTC),
			2024, 1,
			"Q1 end of March",
		},

		// Q2 (April-June)
		{
			time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
			2024, 2,
			"Q2 start of April",
		},
		{
			time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC),
			2024, 2,
			"Q2 middle of May",
		},
		{
			time.Date(2024, 6, 30, 23, 59, 59, 0, time.UTC),
			2024, 2,
			"Q2 end of June",
		},

		// Q3 (July-September)
		{
			time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
			2024, 3,
			"Q3 start of July",
		},
		{
			time.Date(2024, 8, 15, 15, 30, 45, 0, time.UTC),
			2024, 3,
			"Q3 middle of August",
		},
		{
			time.Date(2024, 9, 30, 23, 59, 59, 999999999, time.UTC),
			2024, 3,
			"Q3 end of September",
		},

		// Q4 (October-December)
		{
			time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			2024, 4,
			"Q4 start of October",
		},
		{
			time.Date(2024, 11, 15, 10, 20, 30, 0, time.UTC),
			2024, 4,
			"Q4 middle of November",
		},
		{
			time.Date(2024, 12, 31, 23, 59, 59, 999999999, time.UTC),
			2024, 4,
			"Q4 end of year",
		},

		// Leap year tests
		{
			time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC),
			2024, 1,
			"leap year February 29",
		},

		// Non-leap year tests
		{
			time.Date(2023, 2, 28, 12, 0, 0, 0, time.UTC),
			2023, 1,
			"non-leap year February 28",
		},

		// Different timezones (should be normalized to UTC)
		{
			time.Date(2024, 6, 15, 12, 0, 0, 0, time.FixedZone("EST", -5*3600)),
			2024, 2,
			"different timezone (EST)",
		},
		{
			time.Date(2024, 6, 15, 12, 0, 0, 0, time.FixedZone("PST", -8*3600)),
			2024, 2,
			"different timezone (PST)",
		},

		// Edge cases around quarter boundaries
		{
			time.Date(2024, 3, 31, 23, 59, 59, 999999999, time.UTC),
			2024, 1,
			"Q1 boundary - last moment of March",
		},
		{
			time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
			2024, 2,
			"Q2 boundary - first moment of April",
		},
		{
			time.Date(2024, 6, 30, 23, 59, 59, 999999999, time.UTC),
			2024, 2,
			"Q2 boundary - last moment of June",
		},
		{
			time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
			2024, 3,
			"Q3 boundary - first moment of July",
		},
		{
			time.Date(2024, 9, 30, 23, 59, 59, 999999999, time.UTC),
			2024, 3,
			"Q3 boundary - last moment of September",
		},
		{
			time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			2024, 4,
			"Q4 boundary - first moment of October",
		},
		{
			time.Date(2024, 12, 31, 23, 59, 59, 999999999, time.UTC),
			2024, 4,
			"Q4 boundary - last moment of December",
		},

		// Year boundary tests
		{
			time.Date(2023, 12, 31, 23, 59, 59, 999999999, time.UTC),
			2023, 4,
			"year boundary - last moment of 2023",
		},
		{
			time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			2025, 1,
			"year boundary - first moment of 2025",
		},

		// Different years
		{
			time.Date(2020, 6, 15, 12, 0, 0, 0, time.UTC),
			2020, 2,
			"different year 2020",
		},
		{
			time.Date(2030, 9, 15, 12, 0, 0, 0, time.UTC),
			2030, 3,
			"different year 2030",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			year, quarter := period.QuarterOf(tc.input)

			if year != tc.expectedYear {
				t.Fatalf("QuarterOf(%v) year = %d, want %d", tc.input, year, tc.expectedYear)
			}

			if quarter != tc.expectedQuarter {
				t.Fatalf("QuarterOf(%v) quarter = %d, want %d", tc.input, quarter, tc.expectedQuarter)
			}

			// Additional validation: quarter should be in range [1, 4]
			if quarter < 1 || quarter > 4 {
				t.Fatalf("QuarterOf(%v) quarter = %d, should be in range [1, 4]", tc.input, quarter)
			}

			// Additional validation: year should be reasonable
			if year < 1900 || year > 2100 {
				t.Fatalf("QuarterOf(%v) year = %d, seems unreasonable", tc.input, year)
			}
		})
	}
}

// Test QuarterOf consistency with QuarterBounds
func TestQuarterOfConsistencyWithQuarterBounds(t *testing.T) {
	t.Parallel()

	testCases := []time.Time{
		time.Date(2024, 2, 15, 12, 0, 0, 0, time.UTC),  // Q1
		time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC),  // Q2
		time.Date(2024, 8, 15, 12, 0, 0, 0, time.UTC),  // Q3
		time.Date(2024, 11, 15, 12, 0, 0, 0, time.UTC), // Q4
	}

	for _, input := range testCases {
		t.Run(fmt.Sprintf("consistency_%v", input.Format("2006-01-02")), func(t *testing.T) {
			year, quarter := period.QuarterOf(input)
			start, end := period.QuarterBounds(input)

			// Verify that the returned year matches the year from QuarterBounds
			if start.Year() != year {
				t.Fatalf("QuarterOf(%v) year = %d, but QuarterBounds start year = %d",
					input, year, start.Year())
			}

			// Verify that the input date falls within the quarter bounds
			if input.Before(start) || input.After(end) {
				t.Fatalf("QuarterOf(%v) returned quarter %d, but input is outside QuarterBounds [%v, %v]",
					input, quarter, start, end)
			}

			// Verify quarter number consistency
			expectedQuarter := (int(input.Month())-1)/3 + 1
			if quarter != expectedQuarter {
				t.Fatalf("QuarterOf(%v) quarter = %d, but calculated quarter = %d",
					input, quarter, expectedQuarter)
			}
		})
	}
}

// Test QuarterOf timezone normalization
func TestQuarterOfTimezoneNormalization(t *testing.T) {
	t.Parallel()

	// Test that different timezones produce the same quarter result
	baseTime := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	timezones := []*time.Location{
		time.UTC,
		time.FixedZone("EST", -5*3600),
		time.FixedZone("PST", -8*3600),
		time.FixedZone("JST", 9*3600),
		time.FixedZone("IST", 5*3600+30*60),
	}

	expectedYear := 2024
	expectedQuarter := 2

	for _, tz := range timezones {
		t.Run(fmt.Sprintf("timezone_%s", tz.String()), func(t *testing.T) {
			input := baseTime.In(tz)
			year, quarter := period.QuarterOf(input)

			if year != expectedYear {
				t.Fatalf("QuarterOf(%v in %s) year = %d, want %d", baseTime, tz, year, expectedYear)
			}

			if quarter != expectedQuarter {
				t.Fatalf("QuarterOf(%v in %s) quarter = %d, want %d", baseTime, tz, quarter, expectedQuarter)
			}
		})
	}
}

// Test QuarterOf edge cases around DST transitions
func TestQuarterOfDSTEdgeCases(t *testing.T) {
	t.Parallel()

	// Test dates around DST transitions in different timezones
	testCases := []struct {
		input           time.Time
		expectedYear    int
		expectedQuarter int
		desc            string
	}{
		{
			time.Date(2024, 3, 10, 2, 30, 0, 0, time.FixedZone("EST", -5*3600)), // DST start
			2024, 1,
			"DST start in EST",
		},
		{
			time.Date(2024, 11, 3, 2, 30, 0, 0, time.FixedZone("EST", -5*3600)), // DST end
			2024, 4,
			"DST end in EST",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			year, quarter := period.QuarterOf(tc.input)

			if year != tc.expectedYear {
				t.Fatalf("QuarterOf(%v) year = %d, want %d", tc.input, year, tc.expectedYear)
			}

			if quarter != tc.expectedQuarter {
				t.Fatalf("QuarterOf(%v) quarter = %d, want %d", tc.input, quarter, tc.expectedQuarter)
			}
		})
	}
}
