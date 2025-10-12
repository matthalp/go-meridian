package utc

import (
	"testing"
	"time"

	"github.com/matthalp/go-meridian/est"
	"github.com/matthalp/go-meridian/pst"
)

func TestUTCLocation(t *testing.T) {
	var tz Timezone
	loc := tz.Location()
	if loc != time.UTC {
		t.Errorf("Timezone.Location() = %v, want %v", loc, time.UTC)
	}
}

func TestNow(t *testing.T) {
	before := time.Now().UTC()
	tzTime := Now()
	after := time.Now().UTC()

	// Format to get the underlying time for comparison
	// Parse back to verify it's within range
	parsed, err := time.Parse(time.RFC3339, tzTime.Format(time.RFC3339))
	if err != nil {
		t.Fatalf("Failed to parse result: %v", err)
	}

	if parsed.Before(before.Add(-time.Second)) || parsed.After(after.Add(time.Second)) {
		t.Errorf("Now() returned time outside expected range: got %v, expected between %v and %v", parsed, before, after)
	}
}

func TestDate(t *testing.T) {
	tests := []struct {
		name     string
		year     int
		month    time.Month
		day      int
		hour     int
		min      int
		sec      int
		nsec     int
		expected string
	}{
		{
			name:     "midnight on New Year 2024",
			year:     2024,
			month:    time.January,
			day:      1,
			hour:     0,
			min:      0,
			sec:      0,
			nsec:     0,
			expected: "2024-01-01T00:00:00Z",
		},
		{
			name:     "noon",
			year:     2024,
			month:    time.June,
			day:      15,
			hour:     12,
			min:      30,
			sec:      45,
			nsec:     0,
			expected: "2024-06-15T12:30:45Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tzTime := Date(tt.year, tt.month, tt.day, tt.hour, tt.min, tt.sec, tt.nsec)
			result := tzTime.Format(time.RFC3339)
			if result != tt.expected {
				t.Errorf("Date() formatted as %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestConvert(t *testing.T) {
	t.Run("from time.Time", func(t *testing.T) {
		// Test converting from standard time.Time
		stdTime := time.Date(2024, time.January, 15, 12, 0, 0, 0, time.UTC)
		utcTime := Convert(stdTime)

		// Verify the conversion
		if !utcTime.UTC().Equal(stdTime) {
			t.Errorf("Convert(time.Time) = %v, want %v", utcTime.UTC(), stdTime)
		}

		// Verify formatting shows UTC
		result := utcTime.Format("15:04 MST")
		if result != "12:00 UTC" {
			t.Errorf("Formatted time = %q, want %q", result, "12:00 UTC")
		}
	})

	t.Run("from EST", func(t *testing.T) {
		// Create noon EST
		estTime := est.Date(2024, time.January, 15, 12, 0, 0, 0)

		// Convert to UTC
		utcTime := Convert(estTime)

		// 12:00 EST = 17:00 UTC in winter
		expected := time.Date(2024, time.January, 15, 17, 0, 0, 0, time.UTC)
		if !utcTime.UTC().Equal(expected) {
			t.Errorf("Convert(EST) = %v, want %v", utcTime.UTC(), expected)
		}

		// Verify it displays as UTC
		result := utcTime.Format("15:04 MST")
		if result != "17:00 UTC" {
			t.Errorf("Formatted UTC time = %q, want %q", result, "17:00 UTC")
		}
	})

	t.Run("from PST", func(t *testing.T) {
		// Create noon PST
		pstTime := pst.Date(2024, time.January, 15, 12, 0, 0, 0)

		// Convert to UTC
		utcTime := Convert(pstTime)

		// 12:00 PST = 20:00 UTC in winter
		expected := time.Date(2024, time.January, 15, 20, 0, 0, 0, time.UTC)
		if !utcTime.UTC().Equal(expected) {
			t.Errorf("Convert(PST) = %v, want %v", utcTime.UTC(), expected)
		}

		// Verify it displays as UTC
		result := utcTime.Format("15:04 MST")
		if result != "20:00 UTC" {
			t.Errorf("Formatted UTC time = %q, want %q", result, "20:00 UTC")
		}
	})

	t.Run("preserves moment in time", func(t *testing.T) {
		// All these should represent the same moment
		estTime := est.Date(2024, time.January, 15, 12, 0, 0, 0)
		pstTime := pst.Date(2024, time.January, 15, 9, 0, 0, 0) // 3 hours earlier
		utcTime := Date(2024, time.January, 15, 17, 0, 0, 0)    // EST + 5 hours

		// Convert all to UTC
		utcFromEST := Convert(estTime)
		utcFromPST := Convert(pstTime)

		// All should be equal
		if !utcFromEST.UTC().Equal(utcTime.UTC()) {
			t.Error("EST conversion doesn't match expected UTC time")
		}
		if !utcFromPST.UTC().Equal(utcTime.UTC()) {
			t.Error("PST conversion doesn't match expected UTC time")
		}
	})
}
