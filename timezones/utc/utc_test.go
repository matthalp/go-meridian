package utc

import (
	"testing"
	"time"

	"github.com/matthalp/go-meridian/pt"
)

func TestUTCLocation(t *testing.T) {
	var tz Timezone
	loc := tz.Location()
	if loc.String() != "UTC" {
		t.Errorf("Timezone.Location() = %v, want UTC", loc.String())
	}
}

func TestNow(t *testing.T) {
	before := time.Now().UTC()
	tzTime := Now()
	after := time.Now().UTC()

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
	// Create a time: Jan 15, 2024 at noon UTC
	tzTime := Date(2024, time.January, 15, 12, 0, 0, 0)

	// Format should show the time in UTC
	result := tzTime.Format("15:04 MST")

	// January 15 is during winter, so should show standard time abbreviation
	// The IANA database provides timezone-specific abbreviations (EST, PST, etc.)
	// We just verify it contains the expected hour
	if !contains(result, "12:00") {
		t.Errorf("Format() = %q, expected to contain 12:00", result)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr))
}

func TestDateWithOffset(t *testing.T) {
	// Create a time in UTC (UTC offset varies by timezone and DST)
	// Noon UTC should have corresponding UTC offset
	tzTime := Date(2024, time.January, 1, 12, 0, 0, 0)

	// Parse the formatted time and convert to UTC to verify
	parsed, err := time.Parse(time.RFC3339, tzTime.Format(time.RFC3339))
	if err != nil {
		t.Fatalf("Failed to parse result: %v", err)
	}
	utcTime := parsed.UTC()

	// Verify that the hour in UTC location is 12
	locationTime := utcTime.In(location)
	if locationTime.Hour() != 12 {
		t.Errorf("Date() hour in UTC = %v, want 12", locationTime.Hour())
	}
}

func TestFromMoment(t *testing.T) {
	t.Run("from time.Time", func(t *testing.T) {
		// Test converting from standard time.Time in UTC
		stdTime := time.Date(2024, time.January, 15, 17, 0, 0, 0, time.UTC)
		utcTime := FromMoment(stdTime)

		// Verify the conversion - should represent same moment
		if !utcTime.UTC().Equal(stdTime) {
			t.Errorf("FromMoment(time.Time) UTC = %v, want %v", utcTime.UTC(), stdTime)
		}
	})

	t.Run("from PT", func(t *testing.T) {
		// Create 9:00 PT
		ptTime := pt.Date(2024, time.January, 15, 9, 0, 0, 0)

		// Convert to UTC
		utcTime := FromMoment(ptTime)

		// Verify same moment in time
		if !utcTime.UTC().Equal(ptTime.UTC()) {
			t.Error("Converted time doesn't represent same moment")
		}
	})
}

func TestParse(t *testing.T) {
	t.Run("RFC3339 format", func(t *testing.T) {
		// Parse a time string without timezone, should be interpreted as UTC
		parsed, err := Parse("2006-01-02 15:04:05", "2024-01-15 12:00:00")
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}

		// Should be interpreted as 12:00 UTC
		expected := Date(2024, time.January, 15, 12, 0, 0, 0)
		if parsed.Format(time.RFC3339) != expected.Format(time.RFC3339) {
			t.Errorf("Parse() = %v, want %v", parsed.Format(time.RFC3339), expected.Format(time.RFC3339))
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

		// But UTC should be epoch
		if !epoch.UTC().Equal(time.Unix(0, 0)) {
			t.Error("Unix(0, 0) UTC time should be epoch")
		}
	})

	t.Run("known timestamp", func(t *testing.T) {
		// 2024-01-15 12:00:00 UTC
		result := Unix(1705320000, 0)

		// Verify UTC equivalence
		if !result.UTC().Equal(time.Unix(1705320000, 0)) {
			t.Error("Unix timestamp doesn't match")
		}
	})
}

func TestUnixMilli(t *testing.T) {
	t.Run("known millisecond timestamp", func(t *testing.T) {
		// 2024-01-15 12:00:00.000 UTC
		msec := int64(1705320000000)
		result := UnixMilli(msec)

		// Verify UTC equivalence
		if !result.UTC().Equal(time.UnixMilli(msec)) {
			t.Error("UnixMilli UTC time doesn't match")
		}
	})

	t.Run("with milliseconds precision", func(t *testing.T) {
		msec := int64(1705320000123)
		result := UnixMilli(msec)
		if !result.UTC().Equal(time.UnixMilli(msec)) {
			t.Errorf("UnixMilli precision mismatch")
		}
	})
}

func TestUnixMicro(t *testing.T) {
	t.Run("known microsecond timestamp", func(t *testing.T) {
		// 2024-01-15 12:00:00.000000 UTC
		usec := int64(1705320000000000)
		result := UnixMicro(usec)

		// Verify UTC equivalence
		if !result.UTC().Equal(time.UnixMicro(usec)) {
			t.Error("UnixMicro UTC time doesn't match")
		}
	})

	t.Run("with microseconds precision", func(t *testing.T) {
		usec := int64(1705320000123456)
		result := UnixMicro(usec)
		if !result.UTC().Equal(time.UnixMicro(usec)) {
			t.Errorf("UnixMicro precision mismatch")
		}
	})
}
