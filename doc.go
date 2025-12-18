/*
Package meridian provides type-safe timezone handling for Go using generics.

# The Problem

In Go's standard library, timezone information in time.Time is data, not type.
This means timezone information can be silently lost, leading to bugs:

	func SaveDeadline(t time.Time) {
		// Is this UTC? ET? PT? The compiler can't help you.
		// If someone passes the wrong timezone, it compiles fine but causes bugs.
	}

Meridian makes timezone information immutable by encoding it directly into the type system.

# The Solution

Meridian's Time[TZ] type carries timezone information as a type parameter:

	import (
		"github.com/matthalp/go-meridian/v2/timezones/et"
		"github.com/matthalp/go-meridian/v2/timezones/pt"
		"github.com/matthalp/go-meridian/v2/timezones/utc"
	)

	func SaveDeadline(t utc.Time) {
		// Now the timezone is part of the function signature!
		// The compiler enforces timezone correctness.
	}

Different timezones are different types, so meridian.Time[et.ET] and
meridian.Time[pt.PT] cannot be accidentally mixed:

	func ProcessET(t et.Time) {
		// ... do something ...
	}

	ProcessET(pt.Now())  // Compile error: type mismatch!

# Core Concepts

Type-Safe Timezones: Each timezone is a distinct type. meridian.Time[et.ET]
and meridian.Time[pt.PT] are as different as string and int.

Per-Timezone Packages: Each timezone has its own package (et, pt, utc, etc.)
with convenience functions like et.Now() and pt.Date(...).

Explicit Conversions: Converting between timezones requires explicit function calls:

	eastern := et.Now()
	pacific := pt.FromMoment(eastern)  // Explicit conversion
	utcTime := utc.FromMoment(eastern)  // Convert to UTC for storage

The Moment Interface: Both time.Time and meridian.Time[TZ] implement Moment,
enabling seamless interoperability:

	type Moment interface {
		UTC() time.Time
	}

Internal UTC Storage: All times are stored as UTC internally, making them
safe for database storage and eliminating ambiguity.

# Quick Start

Create times in specific timezones:

	// Current time
	now := utc.Now()

	// Specific date/time
	meeting := et.Date(2024, time.December, 25, 10, 30, 0, 0)

	// Parse from string
	parsed, err := pt.Parse(time.RFC3339, "2024-12-25T10:30:00-08:00")

	// From time.Time
	stdTime := time.Now()
	typed := utc.FromMoment(stdTime)

Convert between timezones explicitly:

	eastern := et.Date(2024, time.December, 25, 9, 0, 0, 0)
	pacific := pt.FromMoment(eastern)  // Same moment, different timezone
	utcTime := utc.FromMoment(eastern)  // Convert to UTC

Work with time.Time seamlessly:

	func processTime(m meridian.Moment) {
		utcTime := m.UTC()  // Works with both time.Time and meridian.Time[TZ]
		// ... do something with utcTime ...
	}

	processTime(time.Now())      // Works
	processTime(et.Now())       // Works
	processTime(pt.Now())       // Works

Write type-safe APIs:

	// Only accepts UTC times
	func SaveToDatabase(t utc.Time) error {
		return db.Save(t.UTC())
	}

	// Generic function that works with any timezone
	func FormatTime[TZ meridian.Timezone](t meridian.Time[TZ]) string {
		return t.Format(time.RFC3339)
	}

# Available Timezones

The package includes these timezone packages:
  - aest: Australian Eastern Time (Australia/Sydney)
  - brt: Brasília Time (America/Sao_Paulo)
  - cet: Central European Time (Europe/Paris)
  - cst: China Standard Time (Asia/Shanghai)
  - ct: Central Time (America/Chicago)
  - et: Eastern Time (America/New_York)
  - gmt: Greenwich Mean Time (Europe/London)
  - hkt: Hong Kong Time (Asia/Hong_Kong)
  - ist: India Standard Time (Asia/Kolkata)
  - jst: Japan Standard Time (Asia/Tokyo)
  - mt: Mountain Time (America/Denver)
  - pt: Pacific Time (America/Los_Angeles)
  - sgt: Singapore Time (Asia/Singapore)
  - utc: Coordinated Universal Time

Additional timezones can be generated using the timezones.yaml configuration.

# Installation

	go get github.com/matthalp/go-meridian/v2

# Design Philosophy

"Make wrong timezone handling impossible to compile."

Meridian prioritizes compile-time safety over convenience and performance.
Explicit timezone conversions make timezone handling visible in code review,
reducing bugs in production.
*/
package meridian
