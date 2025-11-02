/*
Package est provides Eastern Standard Time timezone support for meridian.

EST represents the America/New_York IANA timezone, which observes Eastern Standard Time (EST) and Eastern Daylight Time (EDT) depending on the time of year.

# Usage

Create EST times:

	now := est.Now()
	meeting := est.Date(2024, time.December, 25, 9, 0, 0, 0)
	parsed, _ := est.Parse(time.RFC3339, "2024-12-25T09:00:00-05:00")

Convert to EST from other timezones:

	pacific := pst.Now()
	eastern := est.FromMoment(pacific)

Convert from standard time.Time:

	stdTime := time.Now()
	typedTime := est.FromMoment(stdTime)

The est.Time type is an alias for meridian.Time[est.Timezone], providing
compile-time timezone safety. Functions that accept est.Time can only receive
times explicitly typed as Eastern Standard Time, preventing timezone confusion.
*/
package est

import (
	"fmt"
	"time"

	"github.com/matthalp/go-meridian"
)

// location is the IANA timezone location, loaded once at package initialization.
var location = mustLoadLocation("America/New_York")

// mustLoadLocation loads a timezone location or panics if it fails.
// This should only fail if the system's timezone database is corrupted or missing.
func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(fmt.Sprintf("failed to load timezone %s: %v", name, err))
	}
	return loc
}

// Timezone represents the Eastern Standard Time timezone.
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

// FromMoment converts any Moment to EST time.
func FromMoment(m meridian.Moment) Time {
	return meridian.FromMoment[Timezone](m)
}

// Parse parses a formatted string and returns the time value it represents in EST.
// The layout defines the format by showing how the reference time would be displayed.
// The time is parsed in the America/New_York location.
func Parse(layout, value string) (Time, error) {
	return meridian.Parse[Timezone](layout, value)
}

// Unix returns the EST time corresponding to the given Unix time,
// sec seconds and nsec nanoseconds since January 1, 1970 UTC.
func Unix(sec, nsec int64) Time {
	return meridian.Unix[Timezone](sec, nsec)
}

// UnixMilli returns the EST time corresponding to the given Unix time,
// msec milliseconds since January 1, 1970 UTC.
func UnixMilli(msec int64) Time {
	return meridian.UnixMilli[Timezone](msec)
}

// UnixMicro returns the EST time corresponding to the given Unix time,
// usec microseconds since January 1, 1970 UTC.
func UnixMicro(usec int64) Time {
	return meridian.UnixMicro[Timezone](usec)
}
