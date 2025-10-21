/*
Package ct provides Central Time timezone support for meridian.

CT represents the America/Chicago IANA timezone, which observes Central Time depending on the time of year.

# Usage

Create CT times:

	now := ct.Now()
	specific := ct.Date(2024, time.December, 25, 10, 30, 0, 0)
	parsed, _ := ct.Parse(time.RFC3339, "2024-12-25T10:30:00Z")

Convert to CT from other timezones:

	eastern := est.Now()
	pacific := ct.FromMoment(eastern)

Convert from standard time.Time:

	stdTime := time.Now()
	typedTime := ct.FromMoment(stdTime)

The ct.Time type is an alias for meridian.Time[ct.Timezone], providing
compile-time timezone safety. Functions that accept ct.Time can only receive
times explicitly typed as Central Time, preventing timezone confusion.
*/
package ct

import (
	"fmt"
	"time"

	"github.com/matthalp/go-meridian"
)

// location is the IANA timezone location, loaded once at package initialization.
var location = mustLoadLocation("America/Chicago")

// mustLoadLocation loads a timezone location or panics if it fails.
// This should only fail if the system's timezone database is corrupted or missing.
func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(fmt.Sprintf("failed to load timezone %s: %v", name, err))
	}
	return loc
}

// Timezone represents the Central Time timezone.
type Timezone struct{}

// Location returns the IANA timezone location.
func (Timezone) Location() *time.Location {
	return location
}

// Time is a convenience alias for meridian.Time[Timezone].
type Time = meridian.Time[Timezone]

// Now returns the current time in this timezone.
func Now() Time {
	return meridian.Now[Timezone]()
}

// Date creates a new time in this timezone with the specified date and time components.
func Date(year int, month time.Month, day, hour, minute, sec, nsec int) Time {
	return meridian.Date[Timezone](year, month, day, hour, minute, sec, nsec)
}

// FromMoment converts any Moment to CT time.
func FromMoment(m meridian.Moment) Time {
	return meridian.FromMoment[Timezone](m)
}

// Parse parses a formatted string and returns the time value it represents in CT.
// The layout defines the format by showing how the reference time would be displayed.
// The time is parsed in the America/Chicago location.
func Parse(layout, value string) (Time, error) {
	return meridian.Parse[Timezone](layout, value)
}

// Unix returns the CT time corresponding to the given Unix time,
// sec seconds and nsec nanoseconds since January 1, 1970 UTC.
func Unix(sec, nsec int64) Time {
	return meridian.Unix[Timezone](sec, nsec)
}

// UnixMilli returns the CT time corresponding to the given Unix time,
// msec milliseconds since January 1, 1970 UTC.
func UnixMilli(msec int64) Time {
	return meridian.UnixMilli[Timezone](msec)
}

// UnixMicro returns the CT time corresponding to the given Unix time,
// usec microseconds since January 1, 1970 UTC.
func UnixMicro(usec int64) Time {
	return meridian.UnixMicro[Timezone](usec)
}
