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

func TestFromMoment(t *testing.T) {
	t.Run("from time.Time", func(t *testing.T) {
		// Test converting from standard time.Time
		stdTime := time.Date(2024, time.January, 15, 12, 0, 0, 0, time.UTC)
		utcTime := FromMoment(stdTime)

		// Verify the conversion
		if !utcTime.UTC().Equal(stdTime) {
			t.Errorf("FromMoment(time.Time) = %v, want %v", utcTime.UTC(), stdTime)
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
		utcTime := FromMoment(estTime)

		// 12:00 EST = 17:00 UTC in winter
		expected := time.Date(2024, time.January, 15, 17, 0, 0, 0, time.UTC)
		if !utcTime.UTC().Equal(expected) {
			t.Errorf("FromMoment(EST) = %v, want %v", utcTime.UTC(), expected)
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
		utcTime := FromMoment(pstTime)

		// 12:00 PST = 20:00 UTC in winter
		expected := time.Date(2024, time.January, 15, 20, 0, 0, 0, time.UTC)
		if !utcTime.UTC().Equal(expected) {
			t.Errorf("FromMoment(PST) = %v, want %v", utcTime.UTC(), expected)
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
		utcFromEST := FromMoment(estTime)
		utcFromPST := FromMoment(pstTime)

		// All should be equal
		if !utcFromEST.UTC().Equal(utcTime.UTC()) {
			t.Error("EST conversion doesn't match expected UTC time")
		}
		if !utcFromPST.UTC().Equal(utcTime.UTC()) {
			t.Error("PST conversion doesn't match expected UTC time")
		}
	})
}

func TestParse(t *testing.T) {
	t.Run("RFC3339 format", func(t *testing.T) {
		parsed, err := Parse(time.RFC3339, "2024-01-15T12:00:00Z")
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}

		expected := Date(2024, time.January, 15, 12, 0, 0, 0)
		if parsed.Format(time.RFC3339) != expected.Format(time.RFC3339) {
			t.Errorf("Parse() = %v, want %v", parsed.Format(time.RFC3339), expected.Format(time.RFC3339))
		}
	})

	t.Run("custom format", func(t *testing.T) {
		parsed, err := Parse("2006-01-02 15:04:05", "2024-01-15 12:00:00")
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}

		result := parsed.Format("2006-01-02 15:04:05")
		if result != "2024-01-15 12:00:00" {
			t.Errorf("Parse() formatted = %v, want %v", result, "2024-01-15 12:00:00")
		}
	})

	t.Run("invalid format", func(t *testing.T) {
		_, err := Parse(time.RFC3339, "invalid-time-string")
		if err == nil {
			t.Error("Parse() expected error for invalid input, got nil")
		}
	})
}

func TestUnix(t *testing.T) {
	t.Run("epoch", func(t *testing.T) {
		epoch := Unix(0, 0)
		expected := "1970-01-01T00:00:00Z"
		if epoch.Format(time.RFC3339) != expected {
			t.Errorf("Unix(0, 0) = %v, want %v", epoch.Format(time.RFC3339), expected)
		}
	})

	t.Run("known timestamp", func(t *testing.T) {
		// 2024-01-15 12:00:00 UTC
		result := Unix(1705320000, 0)
		expected := "2024-01-15T12:00:00Z"
		if result.Format(time.RFC3339) != expected {
			t.Errorf("Unix(1705320000, 0) = %v, want %v", result.Format(time.RFC3339), expected)
		}
	})

	t.Run("with nanoseconds", func(t *testing.T) {
		result := Unix(1705320000, 500000000)
		// Should be 12:00:00.5
		if !result.UTC().Equal(time.Unix(1705320000, 500000000)) {
			t.Errorf("Unix with nanoseconds didn't match expected time")
		}
	})
}

func TestUnixMilli(t *testing.T) {
	t.Run("known millisecond timestamp", func(t *testing.T) {
		// 2024-01-15 12:00:00.000 UTC
		msec := int64(1705320000000)
		result := UnixMilli(msec)
		expected := "2024-01-15T12:00:00Z"
		if result.Format(time.RFC3339) != expected {
			t.Errorf("UnixMilli(%d) = %v, want %v", msec, result.Format(time.RFC3339), expected)
		}
	})

	t.Run("with milliseconds precision", func(t *testing.T) {
		// 2024-01-15 12:00:00.123 UTC
		msec := int64(1705320000123)
		result := UnixMilli(msec)
		stdTime := time.UnixMilli(msec)
		if !result.UTC().Equal(stdTime) {
			t.Errorf("UnixMilli precision mismatch")
		}
	})
}

func TestUnixMicro(t *testing.T) {
	t.Run("known microsecond timestamp", func(t *testing.T) {
		// 2024-01-15 12:00:00.000000 UTC
		usec := int64(1705320000000000)
		result := UnixMicro(usec)
		expected := "2024-01-15T12:00:00Z"
		if result.Format(time.RFC3339) != expected {
			t.Errorf("UnixMicro(%d) = %v, want %v", usec, result.Format(time.RFC3339), expected)
		}
	})

	t.Run("with microseconds precision", func(t *testing.T) {
		// 2024-01-15 12:00:00.123456 UTC
		usec := int64(1705320000123456)
		result := UnixMicro(usec)
		stdTime := time.UnixMicro(usec)
		if !result.UTC().Equal(stdTime) {
			t.Errorf("UnixMicro precision mismatch")
		}
	})
}
