// Package pst provides Pacific Standard Time timezone support for meridian.
package pst

import (
	"fmt"
	"time"

	"github.com/matthalp/go-meridian"
)

// location is the IANA timezone location, loaded once at package initialization.
var location = mustLoadLocation("America/Los_Angeles")

// mustLoadLocation loads a timezone location or panics if it fails.
// This should only fail if the system's timezone database is corrupted or missing.
func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(fmt.Sprintf("failed to load timezone %s: %v", name, err))
	}
	return loc
}

// Timezone represents the Pacific Standard Time timezone.
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

// Convert converts any Moment to PST time.
func Convert(m meridian.Moment) Time {
	return meridian.FromMoment[Timezone](m)
}

// Parse parses a formatted string and returns the time value it represents in PST.
// The layout defines the format by showing how the reference time would be displayed.
// Note: ParseInLocation is not needed as the location is already PST.
func Parse(layout, value string) (Time, error) {
	t, err := time.ParseInLocation(layout, value, location)
	if err != nil {
		return Time{}, err
	}
	return meridian.FromMoment[Timezone](t), nil
}

// Unix returns the PST time corresponding to the given Unix time,
// sec seconds and nsec nanoseconds since January 1, 1970 UTC.
func Unix(sec, nsec int64) Time {
	return meridian.FromMoment[Timezone](time.Unix(sec, nsec))
}

// UnixMilli returns the PST time corresponding to the given Unix time,
// msec milliseconds since January 1, 1970 UTC.
func UnixMilli(msec int64) Time {
	return meridian.FromMoment[Timezone](time.UnixMilli(msec))
}

// UnixMicro returns the PST time corresponding to the given Unix time,
// usec microseconds since January 1, 1970 UTC.
func UnixMicro(usec int64) Time {
	return meridian.FromMoment[Timezone](time.UnixMicro(usec))
}
