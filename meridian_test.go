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
