// Package utc provides Coordinated Universal Time timezone support for meridian.
package utc

import (
	"time"

	"github.com/matthalp/go-meridian"
)

// location is the IANA timezone location, loaded once at package initialization.
var location = time.UTC

// Timezone represents the Coordinated Universal Time timezone.
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

// Convert converts any Moment to UTC time.
func Convert(m meridian.Moment) Time {
	return meridian.FromMoment[Timezone](m)
}
