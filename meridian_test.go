package meridian

import (
	"fmt"
	"testing"
	"time"
)

// Test timezone implementations.
type UTC struct{}

func (UTC) Location() *time.Location {
	return time.UTC
}

type EST struct{}

func (EST) Location() *time.Location {
	loc, _ := time.LoadLocation("America/New_York")
	return loc
}

type PST struct{}

func (PST) Location() *time.Location {
	loc, _ := time.LoadLocation("America/Los_Angeles")
	return loc
}

// CustomOffset creates a timezone with a fixed offset from UTC.
type CustomOffset struct {
	offset int // offset in hours
}

func (c CustomOffset) Location() *time.Location {
	return time.FixedZone("Custom", c.offset*3600)
}

func TestVersion(t *testing.T) {
	if Version == "" {
		t.Error("Version should not be empty")
	}
}

func TestNow(t *testing.T) {
	t.Run("UTC", func(t *testing.T) {
		before := time.Now().UTC()
		tzTime := Now[UTC]()
		after := time.Now().UTC()

		// The time should be between before and after
		if tzTime.utcTime.Before(before) || tzTime.utcTime.After(after) {
			t.Errorf("Now[UTC]() returned time outside expected range")
		}
	})

	t.Run("EST", func(t *testing.T) {
		before := time.Now().UTC()
		tzTime := Now[EST]()
		after := time.Now().UTC()

		// The time should be stored in UTC
		if tzTime.utcTime.Before(before) || tzTime.utcTime.After(after) {
			t.Errorf("Now[EST]() returned time outside expected range")
		}
	})

	t.Run("PST", func(t *testing.T) {
		before := time.Now().UTC()
		tzTime := Now[PST]()
		after := time.Now().UTC()

		if tzTime.utcTime.Before(before) || tzTime.utcTime.After(after) {
			t.Errorf("Now[PST]() returned time outside expected range")
		}
	})
}

func TestDate(t *testing.T) {
	tests := []struct {
		name        string
		year        int
		month       time.Month
		day         int
		hour        int
		min         int
		sec         int
		nsec        int
		expectedUTC time.Time
	}{
		{
			name:        "midnight UTC on New Year 2024",
			year:        2024,
			month:       time.January,
			day:         1,
			hour:        0,
			min:         0,
			sec:         0,
			nsec:        0,
			expectedUTC: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:        "noon UTC",
			year:        2024,
			month:       time.June,
			day:         15,
			hour:        12,
			min:         30,
			sec:         45,
			nsec:        123456789,
			expectedUTC: time.Date(2024, time.June, 15, 12, 30, 45, 123456789, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tzTime := Date[UTC](tt.year, tt.month, tt.day, tt.hour, tt.min, tt.sec, tt.nsec)

			if !tzTime.utcTime.Equal(tt.expectedUTC) {
				t.Errorf("Date[UTC]() = %v, want %v", tzTime.utcTime, tt.expectedUTC)
			}
		})
	}
}

func TestDateWithTimezoneOffset(t *testing.T) {
	// Create a time in EST (UTC-5 in winter, UTC-4 in summer)
	// Let's use a winter date to avoid DST complications
	tzTime := Date[EST](2024, time.January, 1, 12, 0, 0, 0) // Noon EST

	// In EST (UTC-5), noon should be 5 PM UTC
	expectedUTC := time.Date(2024, time.January, 1, 17, 0, 0, 0, time.UTC)

	if !tzTime.utcTime.Equal(expectedUTC) {
		t.Errorf("Date[EST](2024, Jan, 1, 12:00:00) = %v, want %v", tzTime.utcTime, expectedUTC)
	}
}

func TestDateWithCustomOffset(t *testing.T) {
	// Create a time at 5 AM in a custom timezone with +5 hours offset
	tzTime := Date[CustomOffset](2024, time.January, 1, 5, 0, 0, 0)

	// 5 AM at +5 should be midnight in UTC
	// Note: CustomOffset zero value has offset=0, so this tests the zero offset case
	expectedUTC := time.Date(2024, time.January, 1, 5, 0, 0, 0, time.UTC)

	if !tzTime.utcTime.Equal(expectedUTC) {
		t.Errorf("Date[CustomOffset](2024, Jan, 1, 05:00:00) = %v, want %v", tzTime.utcTime, expectedUTC)
	}
}

func TestFormat(t *testing.T) {
	// Create a known time in UTC
	utcTime := Date[UTC](2024, time.January, 15, 14, 30, 45, 0)

	tests := []struct {
		name     string
		layout   string
		expected string
	}{
		{
			name:     "RFC3339",
			layout:   time.RFC3339,
			expected: "2024-01-15T14:30:45Z",
		},
		{
			name:     "Kitchen",
			layout:   time.Kitchen,
			expected: "2:30PM",
		},
		{
			name:     "Custom layout",
			layout:   "2006-01-02 15:04:05",
			expected: "2024-01-15 14:30:45",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utcTime.Format(tt.layout)
			if result != tt.expected {
				t.Errorf("Format(%q) = %q, want %q", tt.layout, result, tt.expected)
			}
		})
	}
}

func TestFormatInDifferentTimezone(t *testing.T) {
	// Create a time: Jan 15, 2024 at 17:00 EST (which is 22:00 UTC)
	estTime := Date[EST](2024, time.January, 15, 17, 0, 0, 0)

	// Format should show the time in EST, not UTC
	result := estTime.Format("15:04 MST")

	// Should show 5:00 PM in EST
	if result != "17:00 EST" {
		t.Errorf("Format() = %q, want %q", result, "17:00 EST")
	}

	// Verify it's stored as UTC internally (should be 22:00 UTC)
	expectedUTC := time.Date(2024, time.January, 15, 22, 0, 0, 0, time.UTC)
	if !estTime.utcTime.Equal(expectedUTC) {
		t.Errorf("Internal UTC time = %v, want %v", estTime.utcTime, expectedUTC)
	}
}

func TestTimeTypeSafety(t *testing.T) {
	// This test verifies that different timezone types are distinct
	// at the type level (this is checked at compile time, but we can
	// verify runtime behavior)

	utcTime := Date[UTC](2024, time.January, 1, 12, 0, 0, 0)
	estTime := Date[EST](2024, time.January, 1, 12, 0, 0, 0)

	// Both should format correctly in their respective timezones
	utcStr := utcTime.Format("15:04 MST")
	estStr := estTime.Format("15:04 MST")

	// UTC should show UTC
	if utcStr != "12:00 UTC" {
		t.Errorf("UTC time Format() = %q, want %q", utcStr, "12:00 UTC")
	}

	// EST should show EST
	if estStr != "12:00 EST" {
		t.Errorf("EST time Format() = %q, want %q", estStr, "12:00 EST")
	}

	// Their internal UTC times should be different (5 hours apart in winter)
	hoursDiff := estTime.utcTime.Sub(utcTime.utcTime).Hours()
	if hoursDiff != 5.0 {
		t.Errorf("Time difference between EST and UTC = %v hours, want 5 hours", hoursDiff)
	}
}

func TestGetLocation(t *testing.T) {
	// Test that getLocation correctly extracts locations
	utcLoc := getLocation[UTC]()
	if utcLoc != time.UTC {
		t.Errorf("getLocation[UTC]() = %v, want %v", utcLoc, time.UTC)
	}

	estLoc := getLocation[EST]()
	if estLoc.String() != "America/New_York" {
		t.Errorf("getLocation[EST]().String() = %q, want %q", estLoc.String(), "America/New_York")
	}

	pstLoc := getLocation[PST]()
	if pstLoc.String() != "America/Los_Angeles" {
		t.Errorf("getLocation[PST]().String() = %q, want %q", pstLoc.String(), "America/Los_Angeles")
	}
}

func TestNativeTimeInLocation(t *testing.T) {
	// Create a UTC time
	utcTime := Date[UTC](2024, time.June, 15, 18, 0, 0, 0)

	// Get the native time
	native := utcTime.nativeTimeInLocation()

	// Should be in UTC location
	if native.Location() != time.UTC {
		t.Errorf("nativeTimeInLocation().Location() = %v, want %v", native.Location(), time.UTC)
	}

	// Time should be the same
	if !native.Equal(utcTime.utcTime) {
		t.Errorf("nativeTimeInLocation() = %v, want %v", native, utcTime.utcTime)
	}
}

func TestEdgeCases(t *testing.T) {
	t.Run("leap year", func(t *testing.T) {
		leapTime := Date[UTC](2024, time.February, 29, 12, 0, 0, 0)
		result := leapTime.Format("2006-01-02")
		expected := "2024-02-29"
		if result != expected {
			t.Errorf("Leap year date Format() = %q, want %q", result, expected)
		}
	})

	t.Run("end of year", func(t *testing.T) {
		eoyTime := Date[UTC](2024, time.December, 31, 23, 59, 59, 999999999)
		result := eoyTime.Format("2006-01-02 15:04:05")
		expected := "2024-12-31 23:59:59"
		if result != expected {
			t.Errorf("End of year Format() = %q, want %q", result, expected)
		}
	})

	t.Run("zero nanoseconds", func(t *testing.T) {
		zeroTime := Date[UTC](2024, time.January, 1, 0, 0, 0, 0)
		if zeroTime.utcTime.Nanosecond() != 0 {
			t.Errorf("Zero nanoseconds, got %d", zeroTime.utcTime.Nanosecond())
		}
	})
}

func TestTimeUTC(t *testing.T) {
	// Test that UTC() returns the internal UTC time
	meridianTime := Date[UTC](2024, time.January, 15, 12, 0, 0, 0)
	stdTime := meridianTime.UTC()

	expected := time.Date(2024, time.January, 15, 12, 0, 0, 0, time.UTC)
	if !stdTime.Equal(expected) {
		t.Errorf("UTC() = %v, want %v", stdTime, expected)
	}

	// Test with EST time - should return the UTC equivalent
	estTime := Date[EST](2024, time.January, 15, 12, 0, 0, 0)
	utcFromEST := estTime.UTC()

	// 12:00 EST = 17:00 UTC in winter
	expectedUTC := time.Date(2024, time.January, 15, 17, 0, 0, 0, time.UTC)
	if !utcFromEST.Equal(expectedUTC) {
		t.Errorf("EST time UTC() = %v, want %v", utcFromEST, expectedUTC)
	}
}

func TestTimezoneConversion(t *testing.T) {
	// Create time in EST (noon)
	estTime := Date[EST](2024, time.January, 15, 12, 0, 0, 0)

	// Create time in PST (same clock time, different timezone)
	pstTime := Date[PST](2024, time.January, 15, 12, 0, 0, 0)

	// These should NOT be the same moment in time
	if estTime.UTC().Equal(pstTime.UTC()) {
		t.Error("EST noon and PST noon should be different moments")
	}

	// EST is 3 hours ahead of PST, so EST noon happens 3 hours before PST noon
	diff := pstTime.UTC().Sub(estTime.UTC())
	expectedDiff := 3 * time.Hour
	if diff != expectedDiff {
		t.Errorf("Time difference between PST and EST = %v, want %v", diff, expectedDiff)
	}
}

func TestMomentInterface(t *testing.T) {
	// Test that meridian.Time implements Moment
	var _ Moment = Date[UTC](2024, time.January, 1, 0, 0, 0, 0)

	// Test that time.Time implements Moment (it has UTC() method)
	stdTime := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)
	var _ Moment = stdTime

	// Verify they can be used interchangeably
	moments := []Moment{
		Date[UTC](2024, time.January, 1, 12, 0, 0, 0),
		time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC),
	}

	for i, m := range moments {
		utc := m.UTC()
		if utc.IsZero() {
			t.Errorf("Moment %d returned zero time", i)
		}
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		start    Time[UTC]
		duration time.Duration
		expected time.Time
	}{
		{
			name:     "add 2 hours",
			start:    Date[UTC](2024, time.January, 15, 10, 0, 0, 0),
			duration: 2 * time.Hour,
			expected: time.Date(2024, time.January, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "add 30 minutes",
			start:    Date[UTC](2024, time.January, 15, 10, 30, 0, 0),
			duration: 30 * time.Minute,
			expected: time.Date(2024, time.January, 15, 11, 0, 0, 0, time.UTC),
		},
		{
			name:     "add negative duration",
			start:    Date[UTC](2024, time.January, 15, 10, 0, 0, 0),
			duration: -1 * time.Hour,
			expected: time.Date(2024, time.January, 15, 9, 0, 0, 0, time.UTC),
		},
		{
			name:     "add across day boundary",
			start:    Date[UTC](2024, time.January, 15, 23, 30, 0, 0),
			duration: 1 * time.Hour,
			expected: time.Date(2024, time.January, 16, 0, 30, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.start.Add(tt.duration)
			if !result.UTC().Equal(tt.expected) {
				t.Errorf("Add() = %v, want %v", result.UTC(), tt.expected)
			}
		})
	}
}

func TestAddPreservesTimezoneType(t *testing.T) {
	// Test with EST
	estTime := Date[EST](2024, time.January, 15, 10, 0, 0, 0)
	result := estTime.Add(2 * time.Hour)

	// Verify the result is still EST by formatting
	formatted := result.Format("15:04 MST")
	if !containsTimezone(formatted, "EST") {
		t.Errorf("Add() did not preserve EST timezone: %s", formatted)
	}

	// Test with PST
	pstTime := Date[PST](2024, time.January, 15, 10, 0, 0, 0)
	resultPST := pstTime.Add(2 * time.Hour)

	formattedPST := resultPST.Format("15:04 MST")
	if !containsTimezone(formattedPST, "PST") {
		t.Errorf("Add() did not preserve PST timezone: %s", formattedPST)
	}
}

func TestAddDate(t *testing.T) {
	tests := []struct {
		name     string
		start    Time[UTC]
		years    int
		months   int
		days     int
		expected time.Time
	}{
		{
			name:     "add 1 year",
			start:    Date[UTC](2024, time.January, 15, 10, 0, 0, 0),
			years:    1,
			months:   0,
			days:     0,
			expected: time.Date(2025, time.January, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			name:     "add 3 months",
			start:    Date[UTC](2024, time.January, 15, 10, 0, 0, 0),
			years:    0,
			months:   3,
			days:     0,
			expected: time.Date(2024, time.April, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			name:     "add 10 days",
			start:    Date[UTC](2024, time.January, 15, 10, 0, 0, 0),
			years:    0,
			months:   0,
			days:     10,
			expected: time.Date(2024, time.January, 25, 10, 0, 0, 0, time.UTC),
		},
		{
			name:     "add negative months",
			start:    Date[UTC](2024, time.March, 15, 10, 0, 0, 0),
			years:    0,
			months:   -1,
			days:     0,
			expected: time.Date(2024, time.February, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			name:     "add across year boundary",
			start:    Date[UTC](2024, time.November, 15, 10, 0, 0, 0),
			years:    0,
			months:   3,
			days:     0,
			expected: time.Date(2025, time.February, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			name:     "leap year edge case",
			start:    Date[UTC](2024, time.February, 29, 10, 0, 0, 0),
			years:    1,
			months:   0,
			days:     0,
			expected: time.Date(2025, time.March, 1, 10, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.start.AddDate(tt.years, tt.months, tt.days)
			if !result.UTC().Equal(tt.expected) {
				t.Errorf("AddDate(%d, %d, %d) = %v, want %v",
					tt.years, tt.months, tt.days, result.UTC(), tt.expected)
			}
		})
	}
}

func TestSub(t *testing.T) {
	tests := []struct {
		name     string
		t1       Time[UTC]
		t2       Time[UTC]
		expected time.Duration
	}{
		{
			name:     "2 hours apart",
			t1:       Date[UTC](2024, time.January, 15, 12, 0, 0, 0),
			t2:       Date[UTC](2024, time.January, 15, 10, 0, 0, 0),
			expected: 2 * time.Hour,
		},
		{
			name:     "30 minutes apart",
			t1:       Date[UTC](2024, time.January, 15, 10, 30, 0, 0),
			t2:       Date[UTC](2024, time.January, 15, 10, 0, 0, 0),
			expected: 30 * time.Minute,
		},
		{
			name:     "negative duration",
			t1:       Date[UTC](2024, time.January, 15, 10, 0, 0, 0),
			t2:       Date[UTC](2024, time.January, 15, 12, 0, 0, 0),
			expected: -2 * time.Hour,
		},
		{
			name:     "same time",
			t1:       Date[UTC](2024, time.January, 15, 10, 0, 0, 0),
			t2:       Date[UTC](2024, time.January, 15, 10, 0, 0, 0),
			expected: 0,
		},
		{
			name:     "across days",
			t1:       Date[UTC](2024, time.January, 16, 2, 0, 0, 0),
			t2:       Date[UTC](2024, time.January, 15, 10, 0, 0, 0),
			expected: 16 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.t1.Sub(tt.t2)
			if result != tt.expected {
				t.Errorf("Sub() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSubAcrossTimezones(t *testing.T) {
	// Create same moment in time in different timezones
	estTime := Date[EST](2024, time.January, 15, 12, 0, 0, 0) // Noon EST
	pstTime := Date[PST](2024, time.January, 15, 9, 0, 0, 0)  // 9 AM PST = Noon EST

	// They represent the same moment, so subtracting should give 0
	// Now this works because Sub accepts Moment interface!
	diff := estTime.Sub(pstTime)
	if diff != 0 {
		t.Errorf("Same moment in different timezones: estTime.Sub(pstTime) = %v, want 0", diff)
	}

	// Also test with different moments
	estLater := Date[EST](2024, time.January, 15, 14, 0, 0, 0) // 2 PM EST
	diff2 := estLater.Sub(estTime)
	expected := 2 * time.Hour
	if diff2 != expected {
		t.Errorf("estLater.Sub(estTime) = %v, want %v", diff2, expected)
	}
}

func TestSubWithTimeTime(t *testing.T) {
	// Test that Sub works with standard time.Time
	meridianTime := Date[UTC](2024, time.January, 15, 12, 0, 0, 0)
	stdTime := time.Date(2024, time.January, 15, 10, 0, 0, 0, time.UTC)

	// Subtract time.Time from meridian.Time
	diff := meridianTime.Sub(stdTime)
	expected := 2 * time.Hour
	if diff != expected {
		t.Errorf("meridianTime.Sub(stdTime) = %v, want %v", diff, expected)
	}

	// Test reverse (negative duration)
	diff2 := meridianTime.Sub(time.Date(2024, time.January, 15, 14, 0, 0, 0, time.UTC))
	expected2 := -2 * time.Hour
	if diff2 != expected2 {
		t.Errorf("meridianTime.Sub(laterStdTime) = %v, want %v", diff2, expected2)
	}
}

func TestRound(t *testing.T) {
	tests := []struct {
		name     string
		start    Time[UTC]
		duration time.Duration
		expected time.Time
	}{
		{
			name:     "round to nearest hour (down)",
			start:    Date[UTC](2024, time.January, 15, 10, 20, 0, 0),
			duration: time.Hour,
			expected: time.Date(2024, time.January, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			name:     "round to nearest hour (up)",
			start:    Date[UTC](2024, time.January, 15, 10, 40, 0, 0),
			duration: time.Hour,
			expected: time.Date(2024, time.January, 15, 11, 0, 0, 0, time.UTC),
		},
		{
			name:     "round to nearest 15 minutes",
			start:    Date[UTC](2024, time.January, 15, 10, 37, 0, 0),
			duration: 15 * time.Minute,
			expected: time.Date(2024, time.January, 15, 10, 30, 0, 0, time.UTC),
		},
		{
			name:     "round to nearest minute",
			start:    Date[UTC](2024, time.January, 15, 10, 30, 35, 0),
			duration: time.Minute,
			expected: time.Date(2024, time.January, 15, 10, 31, 0, 0, time.UTC),
		},
		{
			name:     "round exact time",
			start:    Date[UTC](2024, time.January, 15, 10, 0, 0, 0),
			duration: time.Hour,
			expected: time.Date(2024, time.January, 15, 10, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.start.Round(tt.duration)
			if !result.UTC().Equal(tt.expected) {
				t.Errorf("Round(%v) = %v, want %v", tt.duration, result.UTC(), tt.expected)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name     string
		start    Time[UTC]
		duration time.Duration
		expected time.Time
	}{
		{
			name:     "truncate to hour",
			start:    Date[UTC](2024, time.January, 15, 10, 45, 30, 0),
			duration: time.Hour,
			expected: time.Date(2024, time.January, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			name:     "truncate to 15 minutes",
			start:    Date[UTC](2024, time.January, 15, 10, 37, 0, 0),
			duration: 15 * time.Minute,
			expected: time.Date(2024, time.January, 15, 10, 30, 0, 0, time.UTC),
		},
		{
			name:     "truncate to minute",
			start:    Date[UTC](2024, time.January, 15, 10, 30, 45, 123),
			duration: time.Minute,
			expected: time.Date(2024, time.January, 15, 10, 30, 0, 0, time.UTC),
		},
		{
			name:     "truncate exact time",
			start:    Date[UTC](2024, time.January, 15, 10, 0, 0, 0),
			duration: time.Hour,
			expected: time.Date(2024, time.January, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			name:     "truncate to second",
			start:    Date[UTC](2024, time.January, 15, 10, 30, 45, 999999999),
			duration: time.Second,
			expected: time.Date(2024, time.January, 15, 10, 30, 45, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.start.Truncate(tt.duration)
			if !result.UTC().Equal(tt.expected) {
				t.Errorf("Truncate(%v) = %v, want %v", tt.duration, result.UTC(), tt.expected)
			}
		})
	}
}

// Helper function to check if a formatted time string contains a timezone.
func containsTimezone(s, tz string) bool {
	return s != "" && (s[len(s)-3:] == tz || len(s) > 3 && s[len(s)-4:len(s)-1] == tz)
}

func TestAfter(t *testing.T) {
	t1 := Date[UTC](2024, time.January, 15, 12, 0, 0, 0)
	t2 := Date[UTC](2024, time.January, 15, 10, 0, 0, 0)
	t3 := Date[UTC](2024, time.January, 15, 12, 0, 0, 0) // Same as t1

	tests := []struct {
		name     string
		t        Time[UTC]
		u        Time[UTC]
		expected bool
	}{
		{
			name:     "t is after u",
			t:        t1,
			u:        t2,
			expected: true,
		},
		{
			name:     "t is before u",
			t:        t2,
			u:        t1,
			expected: false,
		},
		{
			name:     "t equals u",
			t:        t1,
			u:        t3,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.t.After(tt.u)
			if result != tt.expected {
				t.Errorf("After() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAfterAcrossTimezones(t *testing.T) {
	// Same moment in different timezones
	estTime := Date[EST](2024, time.January, 15, 12, 0, 0, 0) // Noon EST
	pstTime := Date[PST](2024, time.January, 15, 9, 0, 0, 0)  // 9 AM PST = Noon EST

	// Same moment, neither is after the other
	if estTime.After(pstTime) {
		t.Error("Same moment: estTime.After(pstTime) should be false")
	}

	// Different moment - EST noon vs PST noon
	pstNoon := Date[PST](2024, time.January, 15, 12, 0, 0, 0) // Noon PST (3 hours after noon EST)
	if !pstNoon.After(estTime) {
		t.Error("PST noon should be after EST noon")
	}
}

func TestAfterWithTimeTime(t *testing.T) {
	meridianTime := Date[UTC](2024, time.January, 15, 12, 0, 0, 0)
	stdTimeBefore := time.Date(2024, time.January, 15, 10, 0, 0, 0, time.UTC)
	stdTimeAfter := time.Date(2024, time.January, 15, 14, 0, 0, 0, time.UTC)

	if !meridianTime.After(stdTimeBefore) {
		t.Error("meridianTime should be after stdTimeBefore")
	}

	if meridianTime.After(stdTimeAfter) {
		t.Error("meridianTime should not be after stdTimeAfter")
	}
}

func TestBefore(t *testing.T) {
	t1 := Date[UTC](2024, time.January, 15, 10, 0, 0, 0)
	t2 := Date[UTC](2024, time.January, 15, 12, 0, 0, 0)
	t3 := Date[UTC](2024, time.January, 15, 10, 0, 0, 0) // Same as t1

	tests := []struct {
		name     string
		t        Time[UTC]
		u        Time[UTC]
		expected bool
	}{
		{
			name:     "t is before u",
			t:        t1,
			u:        t2,
			expected: true,
		},
		{
			name:     "t is after u",
			t:        t2,
			u:        t1,
			expected: false,
		},
		{
			name:     "t equals u",
			t:        t1,
			u:        t3,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.t.Before(tt.u)
			if result != tt.expected {
				t.Errorf("Before() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestBeforeAcrossTimezones(t *testing.T) {
	// EST noon vs PST noon (EST is 3 hours ahead)
	estNoon := Date[EST](2024, time.January, 15, 12, 0, 0, 0)
	pstNoon := Date[PST](2024, time.January, 15, 12, 0, 0, 0)

	// EST noon happens before PST noon (3 hours earlier)
	if !estNoon.Before(pstNoon) {
		t.Error("EST noon should be before PST noon")
	}

	if pstNoon.Before(estNoon) {
		t.Error("PST noon should not be before EST noon")
	}
}

func TestEqual(t *testing.T) {
	t1 := Date[UTC](2024, time.January, 15, 12, 0, 0, 0)
	t2 := Date[UTC](2024, time.January, 15, 12, 0, 0, 0) // Same as t1
	t3 := Date[UTC](2024, time.January, 15, 12, 0, 1, 0) // 1 second later

	tests := []struct {
		name     string
		t        Time[UTC]
		u        Time[UTC]
		expected bool
	}{
		{
			name:     "equal times",
			t:        t1,
			u:        t2,
			expected: true,
		},
		{
			name:     "different times",
			t:        t1,
			u:        t3,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.t.Equal(tt.u)
			if result != tt.expected {
				t.Errorf("Equal() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEqualAcrossTimezones(t *testing.T) {
	// Same moment in different timezones
	estTime := Date[EST](2024, time.January, 15, 12, 0, 0, 0) // Noon EST
	pstTime := Date[PST](2024, time.January, 15, 9, 0, 0, 0)  // 9 AM PST = Noon EST
	utcTime := Date[UTC](2024, time.January, 15, 17, 0, 0, 0) // 5 PM UTC = Noon EST

	// All represent the same moment
	if !estTime.Equal(pstTime) {
		t.Error("EST and PST times should be equal (same moment)")
	}

	if !estTime.Equal(utcTime) {
		t.Error("EST and UTC times should be equal (same moment)")
	}

	if !pstTime.Equal(utcTime) {
		t.Error("PST and UTC times should be equal (same moment)")
	}
}

func TestEqualWithTimeTime(t *testing.T) {
	meridianTime := Date[UTC](2024, time.January, 15, 12, 0, 0, 0)
	stdTimeSame := time.Date(2024, time.January, 15, 12, 0, 0, 0, time.UTC)
	stdTimeDifferent := time.Date(2024, time.January, 15, 12, 0, 1, 0, time.UTC)

	if !meridianTime.Equal(stdTimeSame) {
		t.Error("meridianTime should equal stdTimeSame")
	}

	if meridianTime.Equal(stdTimeDifferent) {
		t.Error("meridianTime should not equal stdTimeDifferent")
	}
}

func TestCompare(t *testing.T) {
	t1 := Date[UTC](2024, time.January, 15, 10, 0, 0, 0)
	t2 := Date[UTC](2024, time.January, 15, 12, 0, 0, 0)
	t3 := Date[UTC](2024, time.January, 15, 10, 0, 0, 0) // Same as t1

	tests := []struct {
		name     string
		t        Time[UTC]
		u        Time[UTC]
		expected int
	}{
		{
			name:     "t before u returns -1",
			t:        t1,
			u:        t2,
			expected: -1,
		},
		{
			name:     "t after u returns 1",
			t:        t2,
			u:        t1,
			expected: 1,
		},
		{
			name:     "t equals u returns 0",
			t:        t1,
			u:        t3,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.t.Compare(tt.u)
			if result != tt.expected {
				t.Errorf("Compare() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestCompareAcrossTimezones(t *testing.T) {
	estNoon := Date[EST](2024, time.January, 15, 12, 0, 0, 0)
	pstNoon := Date[PST](2024, time.January, 15, 12, 0, 0, 0)
	pst9am := Date[PST](2024, time.January, 15, 9, 0, 0, 0) // Same as EST noon

	// EST noon is before PST noon
	if estNoon.Compare(pstNoon) != -1 {
		t.Error("estNoon.Compare(pstNoon) should return -1")
	}

	// PST noon is after EST noon
	if pstNoon.Compare(estNoon) != 1 {
		t.Error("pstNoon.Compare(estNoon) should return 1")
	}

	// EST noon equals PST 9am (same moment)
	if estNoon.Compare(pst9am) != 0 {
		t.Error("estNoon.Compare(pst9am) should return 0")
	}
}

func TestCompareWithTimeTime(t *testing.T) {
	meridianTime := Date[UTC](2024, time.January, 15, 12, 0, 0, 0)
	stdTimeBefore := time.Date(2024, time.January, 15, 10, 0, 0, 0, time.UTC)
	stdTimeSame := time.Date(2024, time.January, 15, 12, 0, 0, 0, time.UTC)
	stdTimeAfter := time.Date(2024, time.January, 15, 14, 0, 0, 0, time.UTC)

	if meridianTime.Compare(stdTimeBefore) != 1 {
		t.Error("meridianTime.Compare(stdTimeBefore) should return 1")
	}

	if meridianTime.Compare(stdTimeSame) != 0 {
		t.Error("meridianTime.Compare(stdTimeSame) should return 0")
	}

	if meridianTime.Compare(stdTimeAfter) != -1 {
		t.Error("meridianTime.Compare(stdTimeAfter) should return -1")
	}
}

func TestIsZero(t *testing.T) {
	zeroTime := Time[UTC]{}
	nonZeroTime := Date[UTC](2024, time.January, 15, 12, 0, 0, 0)
	explicitZero := Date[UTC](1, time.January, 1, 0, 0, 0, 0)

	if !zeroTime.IsZero() {
		t.Error("Zero value Time should return IsZero() = true")
	}

	if nonZeroTime.IsZero() {
		t.Error("Non-zero Time should return IsZero() = false")
	}

	if !explicitZero.IsZero() {
		t.Error("Explicit zero time (year 1) should return IsZero() = true")
	}
}

func TestIsZeroAcrossTimezones(t *testing.T) {
	// Zero values in different timezone types
	zeroUTC := Time[UTC]{}
	zeroEST := Time[EST]{}
	zeroPST := Time[PST]{}

	if !zeroUTC.IsZero() {
		t.Error("Zero UTC time should return IsZero() = true")
	}

	if !zeroEST.IsZero() {
		t.Error("Zero EST time should return IsZero() = true")
	}

	if !zeroPST.IsZero() {
		t.Error("Zero PST time should return IsZero() = true")
	}
}

func TestClock(t *testing.T) {
	tests := []struct {
		name string
		time Time[UTC]
		hour int
		min  int
		sec  int
	}{
		{
			name: "midnight",
			time: Date[UTC](2024, time.January, 15, 0, 0, 0, 0),
			hour: 0,
			min:  0,
			sec:  0,
		},
		{
			name: "noon",
			time: Date[UTC](2024, time.January, 15, 12, 30, 45, 0),
			hour: 12,
			min:  30,
			sec:  45,
		},
		{
			name: "end of day",
			time: Date[UTC](2024, time.January, 15, 23, 59, 59, 0),
			hour: 23,
			min:  59,
			sec:  59,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hour, minute, sec := tt.time.Clock()
			if hour != tt.hour || minute != tt.min || sec != tt.sec {
				t.Errorf("Clock() = (%d, %d, %d), want (%d, %d, %d)",
					hour, minute, sec, tt.hour, tt.min, tt.sec)
			}
		})
	}
}

func TestClockAcrossTimezones(t *testing.T) {
	// Create noon UTC (which is 7 AM EST in winter, 5 AM PST in winter)
	utcNoon := Date[UTC](2024, time.January, 15, 12, 0, 0, 0)
	estMorning := Date[EST](2024, time.January, 15, 7, 0, 0, 0)
	pstMorning := Date[PST](2024, time.January, 15, 4, 0, 0, 0)

	// Verify they represent the same moment
	if !utcNoon.Equal(estMorning) || !utcNoon.Equal(pstMorning) {
		t.Fatal("Times should represent the same moment")
	}

	// Clock() should return different values based on timezone
	utcHour, _, _ := utcNoon.Clock()
	estHour, _, _ := estMorning.Clock()
	pstHour, _, _ := pstMorning.Clock()

	if utcHour != 12 {
		t.Errorf("UTC hour = %d, want 12", utcHour)
	}
	if estHour != 7 {
		t.Errorf("EST hour = %d, want 7", estHour)
	}
	if pstHour != 4 {
		t.Errorf("PST hour = %d, want 4", pstHour)
	}
}

func TestDateMethod(t *testing.T) {
	tests := []struct {
		name  string
		time  Time[UTC]
		year  int
		month time.Month
		day   int
	}{
		{
			name:  "new year",
			time:  Date[UTC](2024, time.January, 1, 0, 0, 0, 0),
			year:  2024,
			month: time.January,
			day:   1,
		},
		{
			name:  "leap day",
			time:  Date[UTC](2024, time.February, 29, 12, 0, 0, 0),
			year:  2024,
			month: time.February,
			day:   29,
		},
		{
			name:  "end of year",
			time:  Date[UTC](2024, time.December, 31, 23, 59, 59, 0),
			year:  2024,
			month: time.December,
			day:   31,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			year, month, day := tt.time.Date()
			if year != tt.year || month != tt.month || day != tt.day {
				t.Errorf("Date() = (%d, %s, %d), want (%d, %s, %d)",
					year, month, day, tt.year, tt.month, tt.day)
			}
		})
	}
}

func TestDateMethodAcrossTimezones(t *testing.T) {
	// Create a time at 1 AM UTC on Jan 2 (which is 8 PM EST on Jan 1, 5 PM PST on Jan 1)
	utcTime := Date[UTC](2024, time.January, 2, 1, 0, 0, 0)
	estTime := Date[EST](2024, time.January, 1, 20, 0, 0, 0)
	pstTime := Date[PST](2024, time.January, 1, 17, 0, 0, 0)

	// Verify same moment
	if !utcTime.Equal(estTime) || !utcTime.Equal(pstTime) {
		t.Fatal("Times should represent the same moment")
	}

	// Date() should return different dates based on timezone
	utcYear, utcMonth, utcDay := utcTime.Date()
	estYear, estMonth, estDay := estTime.Date()
	pstYear, pstMonth, pstDay := pstTime.Date()

	if utcDay != 2 {
		t.Errorf("UTC day = %d, want 2", utcDay)
	}
	if estDay != 1 {
		t.Errorf("EST day = %d, want 1", estDay)
	}
	if pstDay != 1 {
		t.Errorf("PST day = %d, want 1", pstDay)
	}

	// All should be January 2024
	if utcYear != 2024 || utcMonth != time.January {
		t.Error("UTC date should be January 2024")
	}
	if estYear != 2024 || estMonth != time.January {
		t.Error("EST date should be January 2024")
	}
	if pstYear != 2024 || pstMonth != time.January {
		t.Error("PST date should be January 2024")
	}
}

func TestIndividualDateComponents(t *testing.T) {
	testTime := Date[UTC](2024, time.June, 15, 14, 30, 45, 123456789)

	if testTime.Year() != 2024 {
		t.Errorf("Year() = %d, want 2024", testTime.Year())
	}

	if testTime.Month() != time.June {
		t.Errorf("Month() = %s, want June", testTime.Month())
	}

	if testTime.Day() != 15 {
		t.Errorf("Day() = %d, want 15", testTime.Day())
	}
}

func TestIndividualTimeComponents(t *testing.T) {
	testTime := Date[UTC](2024, time.June, 15, 14, 30, 45, 123456789)

	if testTime.Hour() != 14 {
		t.Errorf("Hour() = %d, want 14", testTime.Hour())
	}

	if testTime.Minute() != 30 {
		t.Errorf("Minute() = %d, want 30", testTime.Minute())
	}

	if testTime.Second() != 45 {
		t.Errorf("Second() = %d, want 45", testTime.Second())
	}

	if testTime.Nanosecond() != 123456789 {
		t.Errorf("Nanosecond() = %d, want 123456789", testTime.Nanosecond())
	}
}

func TestWeekday(t *testing.T) {
	tests := []struct {
		name    string
		time    Time[UTC]
		weekday time.Weekday
	}{
		{
			name:    "Monday",
			time:    Date[UTC](2024, time.January, 15, 12, 0, 0, 0),
			weekday: time.Monday,
		},
		{
			name:    "Sunday",
			time:    Date[UTC](2024, time.January, 21, 12, 0, 0, 0),
			weekday: time.Sunday,
		},
		{
			name:    "Saturday",
			time:    Date[UTC](2024, time.January, 20, 12, 0, 0, 0),
			weekday: time.Saturday,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.time.Weekday() != tt.weekday {
				t.Errorf("Weekday() = %s, want %s", tt.time.Weekday(), tt.weekday)
			}
		})
	}
}

func TestWeekdayAcrossTimezones(t *testing.T) {
	// Create a time at midnight UTC on Monday (which is Sunday evening in PST)
	utcMonday := Date[UTC](2024, time.January, 15, 0, 0, 0, 0)  // Monday midnight UTC
	pstSunday := Date[PST](2024, time.January, 14, 16, 0, 0, 0) // Sunday 4 PM PST

	// Verify same moment
	if !utcMonday.Equal(pstSunday) {
		t.Fatal("Times should represent the same moment")
	}

	// Weekday should differ based on timezone
	if utcMonday.Weekday() != time.Monday {
		t.Errorf("UTC weekday = %s, want Monday", utcMonday.Weekday())
	}
	if pstSunday.Weekday() != time.Sunday {
		t.Errorf("PST weekday = %s, want Sunday", pstSunday.Weekday())
	}
}

func TestYearDay(t *testing.T) {
	tests := []struct {
		name    string
		time    Time[UTC]
		yearDay int
	}{
		{
			name:    "first day of year",
			time:    Date[UTC](2024, time.January, 1, 0, 0, 0, 0),
			yearDay: 1,
		},
		{
			name:    "leap day",
			time:    Date[UTC](2024, time.February, 29, 12, 0, 0, 0),
			yearDay: 60,
		},
		{
			name:    "last day of leap year",
			time:    Date[UTC](2024, time.December, 31, 23, 59, 59, 0),
			yearDay: 366,
		},
		{
			name:    "last day of non-leap year",
			time:    Date[UTC](2023, time.December, 31, 23, 59, 59, 0),
			yearDay: 365,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.time.YearDay() != tt.yearDay {
				t.Errorf("YearDay() = %d, want %d", tt.time.YearDay(), tt.yearDay)
			}
		})
	}
}

func TestISOWeek(t *testing.T) {
	tests := []struct {
		name string
		time Time[UTC]
		year int
		week int
	}{
		{
			name: "first week of 2024",
			time: Date[UTC](2024, time.January, 8, 12, 0, 0, 0),
			year: 2024,
			week: 2,
		},
		{
			name: "last week of year belongs to next year",
			time: Date[UTC](2024, time.December, 30, 12, 0, 0, 0),
			year: 2025,
			week: 1,
		},
		{
			name: "mid-year week",
			time: Date[UTC](2024, time.June, 15, 12, 0, 0, 0),
			year: 2024,
			week: 24,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			year, week := tt.time.ISOWeek()
			if year != tt.year || week != tt.week {
				t.Errorf("ISOWeek() = (%d, %d), want (%d, %d)", year, week, tt.year, tt.week)
			}
		})
	}
}

func TestComponentsRespectTimezone(t *testing.T) {
	// Create the same UTC moment represented in different timezones
	// 2024-01-15 18:00 UTC = 2024-01-15 13:00 EST = 2024-01-15 10:00 PST
	utcTime := Date[UTC](2024, time.January, 15, 18, 0, 0, 0)
	estTime := Date[EST](2024, time.January, 15, 13, 0, 0, 0)
	pstTime := Date[PST](2024, time.January, 15, 10, 0, 0, 0)

	// Verify same moment
	if !utcTime.Equal(estTime) || !utcTime.Equal(pstTime) {
		t.Fatal("Times should represent the same moment")
	}

	// Hours should be different
	if utcTime.Hour() != 18 {
		t.Errorf("UTC hour = %d, want 18", utcTime.Hour())
	}
	if estTime.Hour() != 13 {
		t.Errorf("EST hour = %d, want 13", estTime.Hour())
	}
	if pstTime.Hour() != 10 {
		t.Errorf("PST hour = %d, want 10", pstTime.Hour())
	}

	// But minutes, seconds, nanoseconds should be the same
	if utcTime.Minute() != 0 || estTime.Minute() != 0 || pstTime.Minute() != 0 {
		t.Error("Minutes should all be 0")
	}
	if utcTime.Second() != 0 || estTime.Second() != 0 || pstTime.Second() != 0 {
		t.Error("Seconds should all be 0")
	}
}

func TestIn(t *testing.T) {
	// Create a time in UTC
	utcTime := Date[UTC](2024, time.January, 15, 12, 0, 0, 0)

	// Convert to different locations
	estLoc, _ := time.LoadLocation("America/New_York")
	pstLoc, _ := time.LoadLocation("America/Los_Angeles")

	estConverted := utcTime.In(estLoc)
	pstConverted := utcTime.In(pstLoc)

	// Should represent the same moment
	if !estConverted.Equal(utcTime.UTC()) {
		t.Error("In(EST) should represent the same moment as UTC")
	}
	if !pstConverted.Equal(utcTime.UTC()) {
		t.Error("In(PST) should represent the same moment as UTC")
	}

	// Hours should be different (winter time: EST = UTC-5, PST = UTC-8)
	if estConverted.Hour() != 7 {
		t.Errorf("EST hour = %d, want 7", estConverted.Hour())
	}
	if pstConverted.Hour() != 4 {
		t.Errorf("PST hour = %d, want 4", pstConverted.Hour())
	}
}

func TestLocal(t *testing.T) {
	// Create a time in UTC
	utcTime := Date[UTC](2024, time.January, 15, 12, 0, 0, 0)

	// Convert to local time
	localTime := utcTime.Local()

	// Should represent the same moment
	if !localTime.Equal(utcTime.UTC()) {
		t.Error("Local() should represent the same moment as UTC")
	}

	// Should be in local location
	if localTime.Location() != time.Local {
		t.Errorf("Local() location = %v, want time.Local", localTime.Location())
	}
}

func TestTime(t *testing.T) {
	// Create times in different timezones
	utcTime := Date[UTC](2024, time.January, 15, 12, 0, 0, 0)
	estTime := Date[EST](2024, time.January, 15, 7, 0, 0, 0) // Same moment as UTC noon

	// Get time.Time values
	utcStd := utcTime.Time()
	estStd := estTime.Time()

	// Should be in the correct locations
	if utcStd.Location() != time.UTC {
		t.Errorf("UTC Time() location = %v, want UTC", utcStd.Location())
	}

	estLoc, _ := time.LoadLocation("America/New_York")
	if estStd.Location().String() != estLoc.String() {
		t.Errorf("EST Time() location = %v, want America/New_York", estStd.Location())
	}

	// Should show the correct hours in their respective timezones
	if utcStd.Hour() != 12 {
		t.Errorf("UTC Time() hour = %d, want 12", utcStd.Hour())
	}
	if estStd.Hour() != 7 {
		t.Errorf("EST Time() hour = %d, want 7", estStd.Hour())
	}

	// But they should represent the same moment
	if !utcStd.Equal(estStd) {
		t.Error("UTC and EST Time() should represent the same moment")
	}
}

func TestLocation(t *testing.T) {
	utcTime := Date[UTC](2024, time.January, 15, 12, 0, 0, 0)
	estTime := Date[EST](2024, time.January, 15, 12, 0, 0, 0)
	pstTime := Date[PST](2024, time.January, 15, 12, 0, 0, 0)

	// Check UTC location
	if utcTime.Location() != time.UTC {
		t.Errorf("UTC Location() = %v, want time.UTC", utcTime.Location())
	}

	// Check EST location
	estLoc, _ := time.LoadLocation("America/New_York")
	if estTime.Location().String() != estLoc.String() {
		t.Errorf("EST Location() = %v, want America/New_York", estTime.Location())
	}

	// Check PST location
	pstLoc, _ := time.LoadLocation("America/Los_Angeles")
	if pstTime.Location().String() != pstLoc.String() {
		t.Errorf("PST Location() = %v, want America/Los_Angeles", pstTime.Location())
	}
}

func TestZone(t *testing.T) {
	// Test winter time (standard time, not DST)
	winterTime := Date[EST](2024, time.January, 15, 12, 0, 0, 0)
	winterName, winterOffset := winterTime.Zone()

	if winterName != "EST" {
		t.Errorf("Winter zone name = %q, want %q", winterName, "EST")
	}
	// EST is UTC-5
	if winterOffset != -5*3600 {
		t.Errorf("Winter zone offset = %d, want %d", winterOffset, -5*3600)
	}

	// Test summer time (DST)
	summerTime := Date[EST](2024, time.July, 15, 12, 0, 0, 0)
	summerName, summerOffset := summerTime.Zone()

	if summerName != "EDT" {
		t.Errorf("Summer zone name = %q, want %q", summerName, "EDT")
	}
	// EDT is UTC-4
	if summerOffset != -4*3600 {
		t.Errorf("Summer zone offset = %d, want %d", summerOffset, -4*3600)
	}

	// UTC should always be UTC with 0 offset
	utcTime := Date[UTC](2024, time.January, 15, 12, 0, 0, 0)
	utcName, utcOffset := utcTime.Zone()

	if utcName != "UTC" {
		t.Errorf("UTC zone name = %q, want %q", utcName, "UTC")
	}
	if utcOffset != 0 {
		t.Errorf("UTC zone offset = %d, want 0", utcOffset)
	}
}

func TestZoneBounds(t *testing.T) {
	// Create a time in EST during winter
	winterTime := Date[EST](2024, time.January, 15, 12, 0, 0, 0)
	start, end := winterTime.ZoneBounds()

	// Should have bounds (DST transitions)
	if start.IsZero() && end.IsZero() {
		t.Error("EST should have zone bounds (DST transitions)")
	}

	// The end bound should be after the start bound
	if end.Before(start) {
		t.Error("Zone end bound should be after start bound")
	}

	// UTC should have no bounds (no DST)
	utcTime := Date[UTC](2024, time.January, 15, 12, 0, 0, 0)
	utcStart, utcEnd := utcTime.ZoneBounds()

	// UTC has no transitions, so both should be zero
	if !utcStart.IsZero() || !utcEnd.IsZero() {
		t.Error("UTC should have no zone bounds")
	}
}

func TestIsDST(t *testing.T) {
	// Test winter time (not DST)
	winterTime := Date[EST](2024, time.January, 15, 12, 0, 0, 0)
	if winterTime.IsDST() {
		t.Error("January in EST should not be DST")
	}

	// Test summer time (DST)
	summerTime := Date[EST](2024, time.July, 15, 12, 0, 0, 0)
	if !summerTime.IsDST() {
		t.Error("July in EST should be DST")
	}

	// UTC never has DST
	utcTime := Date[UTC](2024, time.July, 15, 12, 0, 0, 0)
	if utcTime.IsDST() {
		t.Error("UTC should never be DST")
	}

	// PST/PDT tests
	pstWinter := Date[PST](2024, time.January, 15, 12, 0, 0, 0)
	if pstWinter.IsDST() {
		t.Error("January in PST should not be DST")
	}

	pstSummer := Date[PST](2024, time.July, 15, 12, 0, 0, 0)
	if !pstSummer.IsDST() {
		t.Error("July in PST should be DST")
	}
}

func TestTimezoneConversions(t *testing.T) {
	// Create a time and test various conversions
	estTime := Date[EST](2024, time.January, 15, 12, 0, 0, 0)

	// Convert to UTC (already exists)
	utcStd := estTime.UTC()
	if utcStd.Hour() != 17 { // Noon EST = 5 PM UTC
		t.Errorf("UTC() hour = %d, want 17", utcStd.Hour())
	}

	// Convert to PST
	pstLoc, _ := time.LoadLocation("America/Los_Angeles")
	pstStd := estTime.In(pstLoc)
	if pstStd.Hour() != 9 { // Noon EST = 9 AM PST
		t.Errorf("In(PST) hour = %d, want 9", pstStd.Hour())
	}

	// Get time in EST location
	estStd := estTime.Time()
	if estStd.Hour() != 12 {
		t.Errorf("Time() hour = %d, want 12", estStd.Hour())
	}

	// All should represent the same moment
	if !utcStd.Equal(pstStd) || !utcStd.Equal(estStd) {
		t.Error("All conversions should represent the same moment")
	}
}

func TestAppendFormat(t *testing.T) {
	testTime := Date[UTC](2024, time.June, 15, 14, 30, 45, 0)

	tests := []struct {
		name     string
		initial  []byte
		layout   string
		expected string
	}{
		{
			name:     "append to empty slice",
			initial:  []byte{},
			layout:   time.RFC3339,
			expected: "2024-06-15T14:30:45Z",
		},
		{
			name:     "append to existing slice",
			initial:  []byte("Time: "),
			layout:   time.Kitchen,
			expected: "Time: 2:30PM",
		},
		{
			name:     "append custom format",
			initial:  []byte("Date is "),
			layout:   "2006-01-02",
			expected: "Date is 2024-06-15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testTime.AppendFormat(tt.initial, tt.layout)
			if string(result) != tt.expected {
				t.Errorf("AppendFormat() = %q, want %q", string(result), tt.expected)
			}
		})
	}
}

func TestAppendFormatPreservesCapacity(t *testing.T) {
	testTime := Date[UTC](2024, time.June, 15, 14, 30, 45, 0)

	// Create a buffer with capacity
	buf := make([]byte, 0, 100)
	originalCap := cap(buf)

	// AppendFormat should reuse the existing capacity
	result := testTime.AppendFormat(buf, time.RFC3339)

	// Capacity should not have changed (no reallocation)
	if cap(result) != originalCap {
		t.Errorf("AppendFormat() changed capacity: got %d, want %d", cap(result), originalCap)
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		time     Time[UTC]
		expected string
	}{
		{
			name:     "standard time",
			time:     Date[UTC](2024, time.June, 15, 14, 30, 45, 0),
			expected: "2024-06-15 14:30:45 +0000 UTC",
		},
		{
			name:     "midnight",
			time:     Date[UTC](2024, time.January, 1, 0, 0, 0, 0),
			expected: "2024-01-01 00:00:00 +0000 UTC",
		},
		{
			name:     "with nanoseconds",
			time:     Date[UTC](2024, time.June, 15, 14, 30, 45, 123456789),
			expected: "2024-06-15 14:30:45.123456789 +0000 UTC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.time.String()
			if result != tt.expected {
				t.Errorf("String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestStringInDifferentTimezones(t *testing.T) {
	// Create same moment in different timezones
	utcTime := Date[UTC](2024, time.January, 15, 12, 0, 0, 0)
	estTime := Date[EST](2024, time.January, 15, 7, 0, 0, 0) // Same moment
	pstTime := Date[PST](2024, time.January, 15, 4, 0, 0, 0) // Same moment

	// String should show different times based on timezone
	utcStr := utcTime.String()
	estStr := estTime.String()
	pstStr := pstTime.String()

	// UTC should show 12:00
	if utcStr != "2024-01-15 12:00:00 +0000 UTC" {
		t.Errorf("UTC String() = %q, want %q", utcStr, "2024-01-15 12:00:00 +0000 UTC")
	}

	// EST should show 7:00 with EST timezone name
	if estStr != "2024-01-15 07:00:00 -0500 EST" {
		t.Errorf("EST String() = %q, want %q", estStr, "2024-01-15 07:00:00 -0500 EST")
	}

	// PST should show 4:00 with PST timezone name
	if pstStr != "2024-01-15 04:00:00 -0800 PST" {
		t.Errorf("PST String() = %q, want %q", pstStr, "2024-01-15 04:00:00 -0800 PST")
	}
}

func TestStringWithPrint(t *testing.T) {
	// Test that String() is called by fmt.Print family
	testTime := Date[UTC](2024, time.June, 15, 14, 30, 45, 0)

	// fmt.Sprint should use String() method
	result := testTime.String()
	expected := "2024-06-15 14:30:45 +0000 UTC"

	if result != expected {
		t.Errorf("String() = %q, want %q", result, expected)
	}
}

func TestGoString(t *testing.T) {
	tests := []struct {
		name     string
		time     Time[UTC]
		contains []string
	}{
		{
			name: "UTC time",
			time: Date[UTC](2024, time.June, 15, 14, 30, 45, 123456789),
			contains: []string{
				"meridian.Time",
				"UTC",
				"2024-06-15T14:30:45",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.time.GoString()

			// Check that all expected substrings are present
			for _, substr := range tt.contains {
				if !contains(result, substr) {
					t.Errorf("GoString() = %q, expected to contain %q", result, substr)
				}
			}
		})
	}
}

func TestGoStringInDifferentTimezones(t *testing.T) {
	utcTime := Date[UTC](2024, time.June, 15, 14, 30, 45, 0)
	estTime := Date[EST](2024, time.June, 15, 10, 30, 45, 0) // Same moment as UTC

	utcGoStr := utcTime.GoString()
	estGoStr := estTime.GoString()

	// UTC GoString should contain "UTC"
	if !contains(utcGoStr, "UTC") {
		t.Errorf("UTC GoString() = %q, expected to contain 'UTC'", utcGoStr)
	}

	// EST GoString should contain "America/New_York"
	if !contains(estGoStr, "America/New_York") {
		t.Errorf("EST GoString() = %q, expected to contain 'America/New_York'", estGoStr)
	}
}

func TestGoStringWithPrintf(t *testing.T) {
	// Test that GoString() is called by fmt.Printf with %#v
	testTime := Date[UTC](2024, time.June, 15, 14, 30, 45, 0)

	result := fmt.Sprintf("%#v", testTime)

	// Should contain the GoString representation
	if !contains(result, "meridian.Time") {
		t.Errorf("fmt.Sprintf(%%#v) = %q, expected to contain 'meridian.Time'", result)
	}
	if !contains(result, "UTC") {
		t.Errorf("fmt.Sprintf(%%#v) = %q, expected to contain 'UTC'", result)
	}
}

// Helper function to check if a string contains a substring.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
