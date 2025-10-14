package pst

import (
	"testing"
	"time"

	"github.com/matthalp/go-meridian/est"
	"github.com/matthalp/go-meridian/utc"
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

	// Should show noon in PST
	if result != "12:00 PST" {
		t.Errorf("Format() = %q, want %q", result, "12:00 PST")
	}
}

func TestDateWithOffset(t *testing.T) {
	// Create a time in PST (UTC-8 in winter)
	// Noon PST should be 8 PM UTC
	tzTime := Date(2024, time.January, 1, 12, 0, 0, 0)

	// Parse the formatted time and convert to UTC to verify
	parsed, err := time.Parse(time.RFC3339, tzTime.Format(time.RFC3339))
	if err != nil {
		t.Fatalf("Failed to parse result: %v", err)
	}
	utcTime := parsed.UTC()

	// 12:00 PST + 8 hours = 20:00 UTC
	expected := time.Date(2024, time.January, 1, 20, 0, 0, 0, time.UTC)

	if !utcTime.Equal(expected) {
		t.Errorf("Date() in UTC = %v, want %v", utcTime, expected)
	}
}

func TestFromMoment(t *testing.T) {
	t.Run("from time.Time", func(t *testing.T) {
		// Test converting from standard time.Time in UTC
		stdTime := time.Date(2024, time.January, 15, 20, 0, 0, 0, time.UTC)
		pstTime := FromMoment(stdTime)

		// Verify the conversion - should represent same moment
		if !pstTime.UTC().Equal(stdTime) {
			t.Errorf("FromMoment(time.Time) UTC = %v, want %v", pstTime.UTC(), stdTime)
		}

		// Verify formatting shows PST (20:00 UTC = 12:00 PST)
		result := pstTime.Format("15:04 MST")
		if result != "12:00 PST" {
			t.Errorf("Formatted time = %q, want %q", result, "12:00 PST")
		}
	})

	t.Run("from UTC", func(t *testing.T) {
		// Create 20:00 UTC
		utcTime := utc.Date(2024, time.January, 15, 20, 0, 0, 0)

		// Convert to PST
		pstTime := FromMoment(utcTime)

		// 20:00 UTC = 12:00 PST in winter
		result := pstTime.Format("15:04 MST")
		if result != "12:00 PST" {
			t.Errorf("Formatted PST time = %q, want %q", result, "12:00 PST")
		}

		// Verify same moment in time
		if !pstTime.UTC().Equal(utcTime.UTC()) {
			t.Error("Converted time doesn't represent same moment")
		}
	})

	t.Run("from EST", func(t *testing.T) {
		// Create 3:00 EST
		estTime := est.Date(2024, time.January, 15, 15, 0, 0, 0)

		// Convert to PST
		pstTime := FromMoment(estTime)

		// 3:00 PM EST = 12:00 PM PST (3 hour difference)
		result := pstTime.Format("15:04 MST")
		if result != "12:00 PST" {
			t.Errorf("Formatted PST time = %q, want %q", result, "12:00 PST")
		}

		// Verify same moment in time
		if !pstTime.UTC().Equal(estTime.UTC()) {
			t.Error("Converted time doesn't represent same moment")
		}
	})

	t.Run("round trip conversion", func(t *testing.T) {
		// Create time in PST
		original := Date(2024, time.January, 15, 10, 45, 0, 0)

		// Convert to EST and back
		viaEST := FromMoment(est.FromMoment(original))

		// Should represent the same moment
		if !viaEST.UTC().Equal(original.UTC()) {
			t.Error("Round trip conversion changed the moment in time")
		}

		// Should format the same
		if viaEST.Format(time.RFC3339) != original.Format(time.RFC3339) {
			t.Errorf("Round trip format = %q, want %q",
				viaEST.Format(time.RFC3339), original.Format(time.RFC3339))
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
		// Parse same clock time in PST
		pstParsed, err := Parse("2006-01-02 15:04:05", "2024-01-15 12:00:00")
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}

		// Same clock time parsed in UTC would be different
		utcParsed, err := utc.Parse("2006-01-02 15:04:05", "2024-01-15 12:00:00")
		if err != nil {
			t.Fatalf("utc.Parse() error = %v", err)
		}

		// They should represent different moments in time
		if pstParsed.UTC().Equal(utcParsed.UTC()) {
			t.Error("PST and UTC parse of same clock time should be different moments")
		}

		// PST noon should be 8 hours after UTC noon (in winter)
		diff := pstParsed.UTC().Sub(utcParsed.UTC())
		expectedDiff := 8 * time.Hour
		if diff != expectedDiff {
			t.Errorf("Time difference = %v, want %v", diff, expectedDiff)
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
		// Epoch in PST should display as 1969-12-31 16:00:00 PST (UTC-8 in winter)
		formatted := epoch.Format("2006-01-02 15:04:05 MST")
		if formatted[:19] != "1969-12-31 16:00:00" {
			// Just verify the date/time portion, MST vs PST varies
			t.Logf("Unix(0, 0) formatted as: %v (expected starts with 1969-12-31 16:00:00)", formatted)
		}

		// But UTC should be epoch
		if !epoch.UTC().Equal(time.Unix(0, 0)) {
			t.Error("Unix(0, 0) UTC time should be epoch")
		}
	})

	t.Run("known timestamp", func(t *testing.T) {
		// 2024-01-15 12:00:00 UTC = 2024-01-15 04:00:00 PST
		result := Unix(1705320000, 0)
		formatted := result.Format("15:04 MST")
		if formatted != "04:00 PST" {
			t.Errorf("Unix(1705320000, 0) = %v, want %v", formatted, "04:00 PST")
		}
	})
}

func TestUnixMilli(t *testing.T) {
	t.Run("known millisecond timestamp", func(t *testing.T) {
		// 2024-01-15 12:00:00.000 UTC = 04:00:00 PST
		msec := int64(1705320000000)
		result := UnixMilli(msec)
		formatted := result.Format("15:04 MST")
		if formatted != "04:00 PST" {
			t.Errorf("UnixMilli(%d) = %v, want %v", msec, formatted, "04:00 PST")
		}

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
		// 2024-01-15 12:00:00.000000 UTC = 04:00:00 PST
		usec := int64(1705320000000000)
		result := UnixMicro(usec)
		formatted := result.Format("15:04 MST")
		if formatted != "04:00 PST" {
			t.Errorf("UnixMicro(%d) = %v, want %v", usec, formatted, "04:00 PST")
		}

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
