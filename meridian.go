// Package meridian provides first-class, type-safe timezones for Go.
// Because timezone information shouldn't be optional.
package meridian

import "time"

// Version is the current version of the meridian package.
const Version = "0.0.0"

// Timezone interface that all timezone types must implement.
type Timezone interface {
	Location() *time.Location
}

// Moment represents a moment in time that can be converted to UTC.
type Moment interface {
	UTC() time.Time
}

// Now returns the current time in the specified timezone.
func Now[TZ Timezone]() Time[TZ] {
	return Time[TZ]{utcTime: time.Now().UTC()}
}

// Date returns the Time corresponding to the specified date and time
// in the specified timezone.
func Date[TZ Timezone](year int, month time.Month, day, hour, minute, sec, nsec int) Time[TZ] {
	loc := getLocation[TZ]()
	t := time.Date(year, month, day, hour, minute, sec, nsec, loc)
	return Time[TZ]{utcTime: t.UTC()}
}

// FromMoment creates a Time[TZ] from any Moment (e.g., time.Time or another Time[TZ]).
// The Moment is converted to UTC and wrapped in the specified timezone type.
func FromMoment[TZ Timezone](m Moment) Time[TZ] {
	return Time[TZ]{utcTime: m.UTC()}
}

// getLocation extracts the *time.Location from a timezone type.
func getLocation[TZ Timezone]() *time.Location {
	var tz TZ
	return tz.Location()
}

// Time is a time.Time wrapper that carries timezone information in its type.
type Time[TZ Timezone] struct {
	// utcTime is the internal representation of time, stored in UTC.
	// We use UTC internally because the zero value of time.Time in Go is UTC,
	// which ensures our zero values have well-defined behavior.
	utcTime time.Time
}

// Format is a wrapper around time.Time.Format that returns the time in the timezone's location.
func (t Time[TZ]) Format(layout string) string {
	return t.nativeTimeInLocation().Format(layout)
}

// UTC returns the time as a standard time.Time in UTC.
func (t Time[TZ]) UTC() time.Time {
	return t.utcTime
}

// nativeTimeInLocation returns the native time in the location of the timezone.
func (t Time[TZ]) nativeTimeInLocation() time.Time {
	// This is a bit of a hack to get the timezone's location.
	// We're using a type assertion to get the timezone type and then calling the Location method.
	loc := getLocation[TZ]()
	return t.utcTime.In(loc)
}
