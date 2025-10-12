package est

import (
	"testing"
	"time"

	"github.com/matthalp/go-meridian/pst"
	"github.com/matthalp/go-meridian/utc"
)

func TestESTLocation(t *testing.T) {
	var tz Timezone
	loc := tz.Location()
	if loc.String() != "America/New_York" {
		t.Errorf("Timezone.Location() = %v, want America/New_York", loc.String())
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
	// Create a time: Jan 15, 2024 at noon EST
	tzTime := Date(2024, time.January, 15, 12, 0, 0, 0)

	// Format should show the time in EST
	result := tzTime.Format("15:04 MST")

	// Should show noon in EST
	if result != "12:00 EST" {
		t.Errorf("Format() = %q, want %q", result, "12:00 EST")
	}
}

func TestDateWithOffset(t *testing.T) {
	// Create a time in EST (UTC-5 in winter)
	// Noon EST should be 5 PM UTC
	tzTime := Date(2024, time.January, 1, 12, 0, 0, 0)

	// Parse the formatted time and convert to UTC to verify
	parsed, err := time.Parse(time.RFC3339, tzTime.Format(time.RFC3339))
	if err != nil {
		t.Fatalf("Failed to parse result: %v", err)
	}
	utcTime := parsed.UTC()

	// 12:00 EST + 5 hours = 17:00 UTC
	expected := time.Date(2024, time.January, 1, 17, 0, 0, 0, time.UTC)

	if !utcTime.Equal(expected) {
		t.Errorf("Date() in UTC = %v, want %v", utcTime, expected)
	}
}

func TestConvert(t *testing.T) {
	t.Run("from time.Time", func(t *testing.T) {
		// Test converting from standard time.Time in UTC
		stdTime := time.Date(2024, time.January, 15, 17, 0, 0, 0, time.UTC)
		estTime := Convert(stdTime)

		// Verify the conversion - should represent same moment
		if !estTime.UTC().Equal(stdTime) {
			t.Errorf("Convert(time.Time) UTC = %v, want %v", estTime.UTC(), stdTime)
		}

		// Verify formatting shows EST (17:00 UTC = 12:00 EST)
		result := estTime.Format("15:04 MST")
		if result != "12:00 EST" {
			t.Errorf("Formatted time = %q, want %q", result, "12:00 EST")
		}
	})

	t.Run("from UTC", func(t *testing.T) {
		// Create 17:00 UTC
		utcTime := utc.Date(2024, time.January, 15, 17, 0, 0, 0)

		// Convert to EST
		estTime := Convert(utcTime)

		// 17:00 UTC = 12:00 EST in winter
		result := estTime.Format("15:04 MST")
		if result != "12:00 EST" {
			t.Errorf("Formatted EST time = %q, want %q", result, "12:00 EST")
		}

		// Verify same moment in time
		if !estTime.UTC().Equal(utcTime.UTC()) {
			t.Error("Converted time doesn't represent same moment")
		}
	})

	t.Run("from PST", func(t *testing.T) {
		// Create 9:00 PST
		pstTime := pst.Date(2024, time.January, 15, 9, 0, 0, 0)

		// Convert to EST
		estTime := Convert(pstTime)

		// 9:00 PST = 12:00 EST (3 hour difference)
		result := estTime.Format("15:04 MST")
		if result != "12:00 EST" {
			t.Errorf("Formatted EST time = %q, want %q", result, "12:00 EST")
		}

		// Verify same moment in time
		if !estTime.UTC().Equal(pstTime.UTC()) {
			t.Error("Converted time doesn't represent same moment")
		}
	})

	t.Run("round trip conversion", func(t *testing.T) {
		// Create time in EST
		original := Date(2024, time.January, 15, 14, 30, 0, 0)

		// Convert to UTC and back
		viaUTC := Convert(utc.Convert(original))

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
