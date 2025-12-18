package pst

import (
	"testing"
	"time"

	"github.com/matthalp/go-meridian/v2/timezones/pt"
	"github.com/matthalp/go-meridian/v2/timezones/utc"
)

func TestPSTLocation(t *testing.T) {
	var tz Timezone
	loc := tz.Location()
	if loc.String() != "America/Los_Angeles" {
		t.Errorf("Timezone.Location() = %v, want America/Los_Angeles", loc.String())
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
	// Create a time: Jan 15, 2024 at noon PST
	tzTime := Date(2024, time.January, 15, 12, 0, 0, 0)

	// Format should show the time in PST
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
	// Create a time in PST (UTC offset varies by timezone and DST)
	// Noon PST should have corresponding UTC offset
	tzTime := Date(2024, time.January, 1, 12, 0, 0, 0)

	// Parse the formatted time and convert to UTC to verify
	parsed, err := time.Parse(time.RFC3339, tzTime.Format(time.RFC3339))
	if err != nil {
		t.Fatalf("Failed to parse result: %v", err)
	}
	utcTime := parsed.UTC()

	// Verify that the hour in PST location is 12
	locationTime := utcTime.In(location)
	if locationTime.Hour() != 12 {
		t.Errorf("Date() hour in PST = %v, want 12", locationTime.Hour())
	}
}

func TestFromMoment(t *testing.T) {
	t.Run("from time.Time", func(t *testing.T) {
		// Test converting from standard time.Time in UTC
		stdTime := time.Date(2024, time.January, 15, 17, 0, 0, 0, time.UTC)
		pstTime := FromMoment(stdTime)

		// Verify the conversion - should represent same moment
		if !pstTime.UTC().Equal(stdTime) {
			t.Errorf("FromMoment(time.Time) UTC = %v, want %v", pstTime.UTC(), stdTime)
		}
	})

	t.Run("from UTC", func(t *testing.T) {
		// Create 17:00 UTC
		utcTime := utc.Date(2024, time.January, 15, 17, 0, 0, 0)

		// Convert to PST
		pstTime := FromMoment(utcTime)

		// Verify same moment in time
		if !pstTime.UTC().Equal(utcTime.UTC()) {
			t.Error("Converted time doesn't represent same moment")
		}
	})

	t.Run("from PT", func(t *testing.T) {
		// Create 9:00 PT
		ptTime := pt.Date(2024, time.January, 15, 9, 0, 0, 0)

		// Convert to PST
		pstTime := FromMoment(ptTime)

		// Verify same moment in time
		if !pstTime.UTC().Equal(ptTime.UTC()) {
			t.Error("Converted time doesn't represent same moment")
		}
	})

	t.Run("round trip conversion", func(t *testing.T) {
		// Create time in PST
		original := Date(2024, time.January, 15, 14, 30, 0, 0)

		// Convert to UTC and back
		viaUTC := FromMoment(utc.FromMoment(original))

		// Should represent the same moment
		if !viaUTC.UTC().Equal(original.UTC()) {
			t.Error("Round trip conversion changed the moment in time")
		}

		// Should format the same
		if viaUTC.Format(time.RFC3339) != original.Format(time.RFC3339) {
			t.Errorf("Round trip format = %q, want %q",
				viaUTC.Format(time.RFC3339), original.Format(time.RFC3339))
		}
	})
}

func TestParse(t *testing.T) {
	t.Run("RFC3339 format", func(t *testing.T) {
		// Parse a time string without timezone, should be interpreted as PST
		parsed, err := Parse("2006-01-02 15:04:05", "2024-01-15 12:00:00")
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}

		// Should be interpreted as 12:00 PST
		expected := Date(2024, time.January, 15, 12, 0, 0, 0)
		if parsed.Format(time.RFC3339) != expected.Format(time.RFC3339) {
			t.Errorf("Parse() = %v, want %v", parsed.Format(time.RFC3339), expected.Format(time.RFC3339))
		}
	})

	t.Run("timezone specific interpretation", func(t *testing.T) {
		// Parse same clock time in PST during summer (July) to ensure DST offset
		pstParsed, err := Parse("2006-01-02 15:04:05", "2024-07-15 12:00:00")
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}

		// Same clock time parsed in UTC
		utcParsed, err := utc.Parse("2006-01-02 15:04:05", "2024-07-15 12:00:00")
		if err != nil {
			t.Fatalf("utc.Parse() error = %v", err)
		}

		// During summer, most timezones have DST offset from UTC, so they should represent different moments
		// For timezones without DST (like some Asian/African zones), this may still pass if offset != 0
		if pstParsed.UTC().Equal(utcParsed.UTC()) {
			t.Error("PST and UTC parse of same clock time should be different moments")
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
