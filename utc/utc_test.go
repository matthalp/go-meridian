package utc

import (
	"testing"
	"time"
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
