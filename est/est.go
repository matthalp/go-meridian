// Package est provides Eastern Standard Time timezone support for meridian.
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

// Convert converts any Moment to EST time.
func Convert(m meridian.Moment) Time {
	return meridian.FromMoment[Timezone](m)
}
