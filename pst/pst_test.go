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

func TestConvert(t *testing.T) {
	t.Run("from time.Time", func(t *testing.T) {
		// Test converting from standard time.Time in UTC
		stdTime := time.Date(2024, time.January, 15, 20, 0, 0, 0, time.UTC)
		pstTime := Convert(stdTime)

		// Verify the conversion - should represent same moment
		if !pstTime.UTC().Equal(stdTime) {
			t.Errorf("Convert(time.Time) UTC = %v, want %v", pstTime.UTC(), stdTime)
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
		pstTime := Convert(utcTime)

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
		pstTime := Convert(estTime)

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
		viaEST := Convert(est.Convert(original))

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
