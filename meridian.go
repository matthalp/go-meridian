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

// Component Extraction

// Clock returns the hour, minute, and second within the day specified by t,
// in the timezone's location.
func (t Time[TZ]) Clock() (hour, minute, sec int) {
	return t.nativeTimeInLocation().Clock()
}

// Date returns the year, month, and day in which t occurs, in the timezone's location.
func (t Time[TZ]) Date() (year int, month time.Month, day int) {
	return t.nativeTimeInLocation().Date()
}

// Year returns the year in which t occurs, in the timezone's location.
func (t Time[TZ]) Year() int {
	return t.nativeTimeInLocation().Year()
}

// Month returns the month of the year specified by t, in the timezone's location.
func (t Time[TZ]) Month() time.Month {
	return t.nativeTimeInLocation().Month()
}

// Day returns the day of the month specified by t, in the timezone's location.
func (t Time[TZ]) Day() int {
	return t.nativeTimeInLocation().Day()
}

// Hour returns the hour within the day specified by t, in the range [0, 23],
// in the timezone's location.
func (t Time[TZ]) Hour() int {
	return t.nativeTimeInLocation().Hour()
}

// Minute returns the minute offset within the hour specified by t, in the range [0, 59],
// in the timezone's location.
func (t Time[TZ]) Minute() int {
	return t.nativeTimeInLocation().Minute()
}

// Second returns the second offset within the minute specified by t, in the range [0, 59],
// in the timezone's location.
func (t Time[TZ]) Second() int {
	return t.nativeTimeInLocation().Second()
}

// Nanosecond returns the nanosecond offset within the second specified by t,
// in the range [0, 999999999], in the timezone's location.
func (t Time[TZ]) Nanosecond() int {
	return t.nativeTimeInLocation().Nanosecond()
}

// Weekday returns the day of the week specified by t, in the timezone's location.
func (t Time[TZ]) Weekday() time.Weekday {
	return t.nativeTimeInLocation().Weekday()
}

// YearDay returns the day of the year specified by t, in the range [1, 365] for non-leap years,
// and [1, 366] in leap years, in the timezone's location.
func (t Time[TZ]) YearDay() int {
	return t.nativeTimeInLocation().YearDay()
}

// ISOWeek returns the ISO 8601 year and week number in which t occurs.
// Week ranges from 1 to 53. Jan 01 to Jan 03 of year n might belong to
// week 52 or 53 of year n-1, and Dec 29 to Dec 31 might belong to week 1
// of year n+1. Computed in the timezone's location.
func (t Time[TZ]) ISOWeek() (year, week int) {
	return t.nativeTimeInLocation().ISOWeek()
}

// Timezone & Location

// In returns a standard time.Time representing the same time instant as t,
// but with the specified location. This is useful for converting to arbitrary
// timezones without type safety.
func (t Time[TZ]) In(loc *time.Location) time.Time {
	return t.utcTime.In(loc)
}

// Local returns a standard time.Time representing the same time instant as t,
// but with the system's local timezone.
func (t Time[TZ]) Local() time.Time {
	return t.utcTime.Local()
}

// Time returns a standard time.Time representing the time instant in the
// timezone's location. This is useful for interoperating with code that
// expects time.Time.
func (t Time[TZ]) Time() time.Time {
	return t.nativeTimeInLocation()
}

// Location returns the time zone location associated with the timezone type.
func (t Time[TZ]) Location() *time.Location {
	return getLocation[TZ]()
}

// Zone computes the time zone name and its offset in seconds east of UTC
// at the time t in the timezone's location.
func (t Time[TZ]) Zone() (name string, offset int) {
	return t.nativeTimeInLocation().Zone()
}

// ZoneBounds returns the bounds of the time zone in effect at time t.
// The zone begins at start and the next zone begins at end.
// If the zone begins at the beginning of time, start will be returned as zero.
// If the zone goes on forever, end will be returned as zero.
func (t Time[TZ]) ZoneBounds() (start, end time.Time) {
	return t.nativeTimeInLocation().ZoneBounds()
}

// IsDST reports whether the time in the timezone's location is in
// Daylight Saving Time.
func (t Time[TZ]) IsDST() bool {
	return t.nativeTimeInLocation().IsDST()
}

// nativeTimeInLocation returns the native time in the location of the timezone.
func (t Time[TZ]) nativeTimeInLocation() time.Time {
	// This is a bit of a hack to get the timezone's location.
	// We're using a type assertion to get the timezone type and then calling the Location method.
	loc := getLocation[TZ]()
	return t.utcTime.In(loc)
}
