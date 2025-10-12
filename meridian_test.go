package meridian

import (
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
