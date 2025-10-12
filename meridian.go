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

// Time Arithmetic & Manipulation

// Add returns the time t+d, preserving the timezone type.
func (t Time[TZ]) Add(d time.Duration) Time[TZ] {
	return Time[TZ]{utcTime: t.utcTime.Add(d)}
}

// AddDate returns the time corresponding to adding the given number of years,
// months, and days to t, preserving the timezone type.
func (t Time[TZ]) AddDate(years, months, days int) Time[TZ] {
	return Time[TZ]{utcTime: t.utcTime.AddDate(years, months, days)}
}

// Sub returns the duration t-u. If the result exceeds the maximum (or minimum)
// value that can be stored in a Duration, the maximum (or minimum) duration
// will be returned. The parameter u can be any Moment (time.Time or Time[TZ]).
func (t Time[TZ]) Sub(u Moment) time.Duration {
	return t.utcTime.Sub(u.UTC())
}

// Round returns the result of rounding t to the nearest multiple of d (since the zero time),
// preserving the timezone type.
func (t Time[TZ]) Round(d time.Duration) Time[TZ] {
	return Time[TZ]{utcTime: t.utcTime.Round(d)}
}

// Truncate returns the result of rounding t down to a multiple of d (since the zero time),
// preserving the timezone type.
func (t Time[TZ]) Truncate(d time.Duration) Time[TZ] {
	return Time[TZ]{utcTime: t.utcTime.Truncate(d)}
}

// Comparisons & Validation

// After reports whether the time instant t is after u.
// The parameter u can be any Moment (time.Time or Time[TZ]).
func (t Time[TZ]) After(u Moment) bool {
	return t.utcTime.After(u.UTC())
}

// Before reports whether the time instant t is before u.
// The parameter u can be any Moment (time.Time or Time[TZ]).
func (t Time[TZ]) Before(u Moment) bool {
	return t.utcTime.Before(u.UTC())
}

// Equal reports whether t and u represent the same time instant.
// The parameter u can be any Moment (time.Time or Time[TZ]).
func (t Time[TZ]) Equal(u Moment) bool {
	return t.utcTime.Equal(u.UTC())
}

// Compare compares the time instant t with u. If t is before u, it returns -1;
// if t is after u, it returns +1; if they're the same, it returns 0.
// The parameter u can be any Moment (time.Time or Time[TZ]).
func (t Time[TZ]) Compare(u Moment) int {
	return t.utcTime.Compare(u.UTC())
}

// IsZero reports whether t represents the zero time instant,
// January 1, year 1, 00:00:00 UTC.
func (t Time[TZ]) IsZero() bool {
	return t.utcTime.IsZero()
}

// nativeTimeInLocation returns the native time in the location of the timezone.
func (t Time[TZ]) nativeTimeInLocation() time.Time {
	// This is a bit of a hack to get the timezone's location.
	// We're using a type assertion to get the timezone type and then calling the Location method.
	loc := getLocation[TZ]()
	return t.utcTime.In(loc)
}
