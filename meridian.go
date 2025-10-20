/*
Package meridian provides first-class, type-safe timezones for Go.

# The Problem: Timezones Are Data, Not Types

In Go's standard time package, timezone information in time.Time is optional data
that can be silently lost or misinterpreted. The compiler provides no protection:

	func ProcessDeadline(deadline time.Time) {
		// Is deadline in UTC? EST? The user's local time?
		// If the caller passes the wrong timezone, it compiles but causes bugs.
	}

This leads to subtle bugs in production where times are interpreted in the wrong
timezone, causing incorrect behavior.

# The Solution: Timezone as Type

Meridian encodes timezone information directly in the type system using Go generics.
Time[TZ] is a time.Time wrapper where TZ is a timezone type parameter:

	import "github.com/matthalp/go-meridian/utc"

	func ProcessDeadline(deadline utc.Time) {
		// Now the timezone is guaranteed by the compiler!
		// Only UTC times can be passed, preventing timezone bugs.
	}

Different timezones are incompatible types. meridian.Time[et.ET] and
meridian.Time[pt.PT] cannot be mixed without explicit conversion:

	func ProcessET(t et.Time) {
		// ... do something ...
	}

	ProcessET(pt.Now())  // Compile error: cannot use pt.Time as et.Time

# Core Design

Type-Safe Times: Time[TZ] carries timezone information as a type parameter,
making timezone part of the type system rather than runtime data.

Explicit Conversions: Converting between timezones requires explicit function calls
using the FromMoment function, making timezone handling visible in code review:

	eastern := et.Now()
	pacific := pt.FromMoment(eastern)  // Explicit and reviewable

Moment Interface: Both time.Time and Time[TZ] implement the Moment interface,
enabling seamless interoperability with existing code:

	type Moment interface {
		UTC() time.Time
	}

Internal UTC Storage: All times are stored as UTC internally, eliminating DST
ambiguity and making database storage straightforward. The timezone is applied
during display and component extraction operations.

Per-Timezone Packages: Each timezone has its own package (et, pt, utc) with
convenience functions like et.Now() and pt.Date(...), plus a Time alias
for meridian.Time[Timezone].

# Philosophy

"Make wrong timezone handling impossible to compile."

Meridian prioritizes compile-time safety over convenience and performance.
Explicit conversions and type-safe APIs make timezone handling visible and
prevent an entire class of bugs from reaching production.
*/
package meridian

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"time"
)

// Version is the current version of the meridian package.
const Version = "1.5.0"

// Timezone interface that all timezone types must implement.
// Each timezone package defines its own Timezone type that satisfies this interface,
// enabling Time[TZ] to be parameterized with type-safe timezone information.
type Timezone interface {
	Location() *time.Location
}

// Moment represents a moment in time that can be converted to UTC.
// Both time.Time and Time[TZ] implement this interface, enabling functions
// to accept either type while maintaining interoperability with the standard library.
type Moment interface {
	UTC() time.Time
}

// Now returns the current time in the specified timezone.
// The timezone type parameter TZ is typically inferred from context or explicitly
// specified. For most use cases, prefer timezone-specific helpers like est.Now()
// or utc.Now() for better readability.
func Now[TZ Timezone]() Time[TZ] {
	return Time[TZ]{utcTime: time.Now().UTC()}
}

// Date returns the Time corresponding to the specified date and time
// in the specified timezone. The date components are interpreted in the timezone's
// location, then stored internally as UTC. The timezone type is preserved in the
// return type, ensuring type-safe handling. For most use cases, prefer timezone-specific
// helpers like est.Date() or utc.Date() for better readability.
func Date[TZ Timezone](year int, month time.Month, day, hour, minute, sec, nsec int) Time[TZ] {
	loc := getLocation[TZ]()
	t := time.Date(year, month, day, hour, minute, sec, nsec, loc)
	return Time[TZ]{utcTime: t.UTC()}
}

// FromMoment creates a Time[TZ] from any Moment (e.g., time.Time or another Time[TZ]).
// This is the primary way to convert between timezones explicitly. The conversion
// preserves the moment in time (UTC equality) but changes the timezone type, making
// the conversion visible in code review. For most use cases, prefer timezone-specific
// helpers like est.FromMoment() or pst.FromMoment() for better readability.
func FromMoment[TZ Timezone](m Moment) Time[TZ] {
	return Time[TZ]{utcTime: m.UTC()}
}

// Parse parses a formatted string and returns the time value it represents in the specified timezone.
// The layout defines the format by showing how the reference time would be displayed.
func Parse[TZ Timezone](layout, value string) (Time[TZ], error) {
	loc := getLocation[TZ]()
	t, err := time.ParseInLocation(layout, value, loc)
	if err != nil {
		return Time[TZ]{}, err
	}
	return Time[TZ]{utcTime: t.UTC()}, nil
}

// Unix returns the Time corresponding to the given Unix time,
// sec seconds and nsec nanoseconds since January 1, 1970 UTC,
// in the specified timezone.
func Unix[TZ Timezone](sec, nsec int64) Time[TZ] {
	return Time[TZ]{utcTime: time.Unix(sec, nsec).UTC()}
}

// UnixMilli returns the Time corresponding to the given Unix time,
// msec milliseconds since January 1, 1970 UTC, in the specified timezone.
func UnixMilli[TZ Timezone](msec int64) Time[TZ] {
	return Time[TZ]{utcTime: time.UnixMilli(msec).UTC()}
}

// UnixMicro returns the Time corresponding to the given Unix time,
// usec microseconds since January 1, 1970 UTC, in the specified timezone.
func UnixMicro[TZ Timezone](usec int64) Time[TZ] {
	return Time[TZ]{utcTime: time.UnixMicro(usec).UTC()}
}

// getLocation extracts the *time.Location from a timezone type.
func getLocation[TZ Timezone]() *time.Location {
	var tz TZ
	return tz.Location()
}

// Time is a time.Time wrapper that carries timezone information in its type parameter.
// Unlike time.Time where timezone is optional data, Time[TZ] makes timezone part of
// the type system, providing compile-time safety. Different timezone types are
// incompatible, preventing accidental timezone mixing.
type Time[TZ Timezone] struct {
	// utcTime is the internal representation of time, stored in UTC.
	// We use UTC internally because the zero value of time.Time in Go is UTC,
	// which ensures our zero values have well-defined behavior. The timezone
	// type parameter TZ is applied during display and component extraction.
	utcTime time.Time
}

// Compile-time interface assertions.
var (
	_ fmt.Stringer               = Time[Timezone]{}
	_ fmt.GoStringer             = Time[Timezone]{}
	_ json.Marshaler             = Time[Timezone]{}
	_ json.Unmarshaler           = (*Time[Timezone])(nil)
	_ encoding.TextMarshaler     = Time[Timezone]{}
	_ encoding.TextUnmarshaler   = (*Time[Timezone])(nil)
	_ encoding.BinaryMarshaler   = Time[Timezone]{}
	_ encoding.BinaryUnmarshaler = (*Time[Timezone])(nil)
	_ driver.Valuer              = Time[Timezone]{}
	_ sql.Scanner                = (*Time[Timezone])(nil)
)

// Formatting & String Output

// Format is a wrapper around time.Time.Format that returns the time in the timezone's location.
func (t Time[TZ]) Format(layout string) string {
	return t.nativeTimeInLocation().Format(layout)
}

// AppendFormat is like Format but appends the textual representation to b and returns
// the extended buffer.
func (t Time[TZ]) AppendFormat(b []byte, layout string) []byte {
	return t.nativeTimeInLocation().AppendFormat(b, layout)
}

// String returns the time formatted using the RFC3339 layout with the timezone's location.
// It implements the fmt.Stringer interface.
func (t Time[TZ]) String() string {
	return t.nativeTimeInLocation().String()
}

// GoString returns a string representation of the Time value in Go syntax.
// It implements the fmt.GoStringer interface for use in debugging.
func (t Time[TZ]) GoString() string {
	return fmt.Sprintf("meridian.Time[%s]{%s}", t.Location().String(), t.Format(time.RFC3339Nano))
}

// UTC returns the time as a standard time.Time in UTC.
// This method implements the Moment interface, enabling interoperability with
// both time.Time and other Time[TZ] types. The returned time.Time is always in UTC.
func (t Time[TZ]) UTC() time.Time {
	return t.utcTime
}

// Time Arithmetic & Manipulation

// Add returns the time t+d, preserving the timezone type.
// The timezone type is maintained in the return value, ensuring that operations
// on typed times continue to provide type-safe timezone guarantees.
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

// Unix Timestamp Conversion

// Unix returns t as a Unix time, the number of seconds elapsed since
// January 1, 1970 UTC.
func (t Time[TZ]) Unix() int64 {
	return t.utcTime.Unix()
}

// UnixMilli returns t as a Unix time, the number of milliseconds elapsed since
// January 1, 1970 UTC.
func (t Time[TZ]) UnixMilli() int64 {
	return t.utcTime.UnixMilli()
}

// UnixMicro returns t as a Unix time, the number of microseconds elapsed since
// January 1, 1970 UTC.
func (t Time[TZ]) UnixMicro() int64 {
	return t.utcTime.UnixMicro()
}

// UnixNano returns t as a Unix time, the number of nanoseconds elapsed since
// January 1, 1970 UTC. The result is undefined if the Unix time in nanoseconds
// cannot be represented by an int64 (a date before the year 1678 or after 2262).
func (t Time[TZ]) UnixNano() int64 {
	return t.utcTime.UnixNano()
}

// Serialization Interfaces

// MarshalJSON implements the json.Marshaler interface.
// The time is formatted as an RFC 3339 string in the timezone's location.
func (t Time[TZ]) MarshalJSON() ([]byte, error) {
	return t.nativeTimeInLocation().MarshalJSON()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is parsed and stored as UTC internally.
func (t *Time[TZ]) UnmarshalJSON(data []byte) error {
	var stdTime time.Time
	if err := stdTime.UnmarshalJSON(data); err != nil {
		return err
	}
	t.utcTime = stdTime.UTC()
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
// The time is formatted as an RFC 3339 string in the timezone's location.
func (t Time[TZ]) MarshalText() ([]byte, error) {
	return t.nativeTimeInLocation().MarshalText()
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The time is parsed and stored as UTC internally.
func (t *Time[TZ]) UnmarshalText(data []byte) error {
	var stdTime time.Time
	if err := stdTime.UnmarshalText(data); err != nil {
		return err
	}
	t.utcTime = stdTime.UTC()
	return nil
}

// AppendText appends the textual representation of t to b and returns the extended buffer.
// The time is formatted as an RFC 3339 string in the timezone's location.
func (t Time[TZ]) AppendText(b []byte) ([]byte, error) {
	return t.nativeTimeInLocation().AppendFormat(b, time.RFC3339Nano), nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (t Time[TZ]) MarshalBinary() ([]byte, error) {
	return t.utcTime.MarshalBinary()
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (t *Time[TZ]) UnmarshalBinary(data []byte) error {
	return t.utcTime.UnmarshalBinary(data)
}

// AppendBinary appends the binary representation of t to b and returns the extended buffer.
func (t Time[TZ]) AppendBinary(b []byte) ([]byte, error) {
	enc, err := t.utcTime.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return append(b, enc...), nil
}

// GobEncode implements the gob.GobEncoder interface.
func (t Time[TZ]) GobEncode() ([]byte, error) {
	return t.utcTime.GobEncode()
}

// GobDecode implements the gob.GobDecoder interface.
func (t *Time[TZ]) GobDecode(data []byte) error {
	return t.utcTime.GobDecode(data)
}

// Database/SQL Support

// Value implements the driver.Valuer interface for database/sql.
// The time is stored as UTC in the database.
func (t Time[TZ]) Value() (driver.Value, error) {
	return t.utcTime, nil
}

// Scan implements the sql.Scanner interface for database/sql.
// It accepts time.Time values and stores them as UTC internally.
func (t *Time[TZ]) Scan(value interface{}) error {
	if value == nil {
		t.utcTime = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		t.utcTime = v.UTC()
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into meridian.Time", value)
	}
}

// nativeTimeInLocation returns the native time in the location of the timezone.
func (t Time[TZ]) nativeTimeInLocation() time.Time {
	// This is a bit of a hack to get the timezone's location.
	// We're using a type assertion to get the timezone type and then calling the Location method.
	loc := getLocation[TZ]()
	return t.utcTime.In(loc)
}
