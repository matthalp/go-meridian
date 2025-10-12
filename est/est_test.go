package est

import (
	"testing"
	"time"
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
