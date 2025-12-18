package meridian_test

import (
	"fmt"
	"time"

	"github.com/matthalp/go-meridian/v2"
	"github.com/matthalp/go-meridian/v2/timezones/et"
	"github.com/matthalp/go-meridian/v2/timezones/pt"
	"github.com/matthalp/go-meridian/v2/timezones/utc"
)

func ExampleNow() {
	// Get the current time in UTC
	now := utc.Now()

	// Format it
	fmt.Println("Current time format:", now.Format("2006-01-02"))
	// Output will vary, so we can't test exact output
}

func ExampleDate() {
	// Create a specific time in UTC
	t := utc.Date(2024, time.January, 15, 14, 30, 0, 0)

	// Format the time
	fmt.Println(t.Format("2006-01-02 15:04:05"))
	// Output: 2024-01-15 14:30:00
}

func ExampleTime_Format() {
	// Create a specific time in UTC
	t := utc.Date(2024, time.June, 15, 9, 30, 0, 0)

	// Format the time in different layouts
	fmt.Println(t.Format(time.RFC3339))
	fmt.Println(t.Format("Monday, January 2, 2006"))
	// Output:
	// 2024-06-15T09:30:00Z
	// Saturday, June 15, 2024
}

// ExampleFromMoment demonstrates explicit timezone conversions.
func ExampleFromMoment() {
	// Create a time in Eastern timezone
	eastern := et.Date(2024, time.December, 25, 9, 0, 0, 0)

	// Explicitly convert to Pacific time
	pacific := pt.FromMoment(eastern)

	// Explicitly convert to UTC
	universal := utc.FromMoment(eastern)

	// All represent the same moment in time
	fmt.Println("ET:", eastern.Format("15:04 MST"))
	fmt.Println("PT:", pacific.Format("15:04 MST"))
	fmt.Println("UTC:", universal.Format("15:04 MST"))
	fmt.Println("Same moment:", eastern.Equal(pacific) && eastern.Equal(universal))
	// Output:
	// ET: 09:00 EST
	// PT: 06:00 PST
	// UTC: 14:00 UTC
	// Same moment: true
}

// ExampleFromMoment_timeTime demonstrates converting from standard time.Time.
func ExampleFromMoment_timeTime() {
	// Standard library time.Time
	stdTime := time.Date(2024, time.June, 15, 14, 30, 0, 0, time.UTC)

	// Convert to timezone-specific types
	utcTyped := utc.FromMoment(stdTime)
	etTyped := et.FromMoment(stdTime)
	ptTyped := pt.FromMoment(stdTime)

	fmt.Println("UTC:", utcTyped.Format("3:04 PM MST"))
	fmt.Println("ET:", etTyped.Format("3:04 PM MST"))
	fmt.Println("PT:", ptTyped.Format("3:04 PM MST"))
	// Output:
	// UTC: 2:30 PM UTC
	// ET: 10:30 AM EDT
	// PT: 7:30 AM PDT
}

// Example_typeSafety demonstrates compile-time timezone safety.
func Example_typeSafety() {
	// Function that only accepts UTC times
	processUTC := func(t utc.Time) {
		fmt.Println("Processing UTC:", t.Format("15:04 MST"))
	}

	// Function that only accepts ET times
	processET := func(t et.Time) {
		fmt.Println("Processing ET:", t.Format("15:04 MST"))
	}

	utcTime := utc.Date(2024, time.June, 15, 14, 0, 0, 0)
	etTime := et.Date(2024, time.June, 15, 10, 0, 0, 0)

	// These work - types match
	processUTC(utcTime)
	processET(etTime)

	// These would NOT compile (uncomment to see the error):
	// processUTC(etTime)  // Compile error: cannot use et.Time as utc.Time
	// processET(utcTime)  // Compile error: cannot use utc.Time as et.Time

	// To convert, you must be explicit:
	processUTC(utc.FromMoment(etTime))
	processET(et.FromMoment(utcTime))
	// Output:
	// Processing UTC: 14:00 UTC
	// Processing ET: 10:00 EDT
	// Processing UTC: 14:00 UTC
	// Processing ET: 10:00 EDT
}

// Example_momentInterface demonstrates using the Moment interface for flexibility.
func Example_momentInterface() {
	// Function that accepts any Moment (time.Time or meridian.Time[TZ])
	formatMoment := func(m meridian.Moment) string {
		// Convert to UTC for consistent formatting
		return m.UTC().Format("2006-01-02 15:04:05 MST")
	}

	// Works with time.Time
	stdTime := time.Date(2024, time.June, 15, 14, 30, 0, 0, time.UTC)
	fmt.Println("time.Time:", formatMoment(stdTime))

	// Works with any meridian.Time[TZ]
	utcTime := utc.Date(2024, time.June, 15, 14, 30, 0, 0)
	etTime := et.Date(2024, time.June, 15, 10, 30, 0, 0)
	ptTime := pt.Date(2024, time.June, 15, 7, 30, 0, 0)

	fmt.Println("utc.Time:", formatMoment(utcTime))
	fmt.Println("et.Time:", formatMoment(etTime))
	fmt.Println("pt.Time:", formatMoment(ptTime))
	// Output:
	// time.Time: 2024-06-15 14:30:00 UTC
	// utc.Time: 2024-06-15 14:30:00 UTC
	// et.Time: 2024-06-15 14:30:00 UTC
	// pt.Time: 2024-06-15 14:30:00 UTC
}

// Example_multipleTimezones demonstrates working with multiple timezones.
func Example_multipleTimezones() {
	// Schedule a meeting at 2pm ET
	meetingET := et.Date(2024, time.December, 25, 14, 0, 0, 0)

	// What time is that in other timezones?
	meetingPT := pt.FromMoment(meetingET)
	meetingUTC := utc.FromMoment(meetingET)

	fmt.Println("Meeting scheduled:")
	fmt.Println("  Eastern:", meetingET.Format("3:04 PM MST"))
	fmt.Println("  Pacific:", meetingPT.Format("3:04 PM MST"))
	fmt.Println("  UTC:", meetingUTC.Format("3:04 PM MST"))

	// Add timezone-specific operations
	oneHourLater := meetingET.Add(time.Hour)
	fmt.Println("  One hour later (ET):", oneHourLater.Format("3:04 PM MST"))
	// Output:
	// Meeting scheduled:
	//   Eastern: 2:00 PM EST
	//   Pacific: 11:00 AM PST
	//   UTC: 7:00 PM UTC
	//   One hour later (ET): 3:00 PM EST
}

// ExampleParse demonstrates parsing time strings in specific timezones.
func ExampleParse() {
	// Parse a time string as ET
	etTime, _ := et.Parse(time.RFC3339, "2024-12-25T09:00:00-05:00")
	fmt.Println("Parsed ET:", etTime.Format("15:04 MST"))

	// Parse a time string as PT
	ptTime, _ := pt.Parse(time.RFC3339, "2024-12-25T06:00:00-08:00")
	fmt.Println("Parsed PT:", ptTime.Format("15:04 MST"))

	// These represent the same moment
	fmt.Println("Same moment:", etTime.Equal(ptTime))
	// Output:
	// Parsed ET: 09:00 EST
	// Parsed PT: 06:00 PST
	// Same moment: true
}

// Example_genericFunction demonstrates writing generic functions that work with any timezone.
func Example_genericFunction() {
	// Generic function that preserves the timezone type parameter
	addBusinessDay := func(t utc.Time) utc.Time {
		return t.Add(24 * time.Hour)
	}

	// Another function for ET times
	addBusinessDayET := func(t et.Time) et.Time {
		return t.Add(24 * time.Hour)
	}

	// Each function maintains type safety for its timezone
	utcTime := utc.Date(2024, time.June, 14, 12, 0, 0, 0)
	etTime := et.Date(2024, time.June, 14, 8, 0, 0, 0)

	nextUTC := addBusinessDay(utcTime)
	nextET := addBusinessDayET(etTime)

	fmt.Println("UTC next day:", nextUTC.Format("Monday, Jan 2"))
	fmt.Println("ET next day:", nextET.Format("Monday, Jan 2"))
	// Output:
	// UTC next day: Saturday, Jun 15
	// ET next day: Saturday, Jun 15
}

// Example_databaseStorage demonstrates storing times in a database-friendly way.
func Example_databaseStorage() {
	// Create a time in ET
	etTime := et.Date(2024, time.December, 25, 9, 0, 0, 0)

	// Convert to UTC for database storage (standard practice)
	utcForDB := utc.FromMoment(etTime)

	// Get the underlying time.Time for database operations
	dbValue := utcForDB.UTC()

	fmt.Println("Store in DB:", dbValue.Format(time.RFC3339))

	// When reading from DB, convert back to typed timezone
	retrieved := utc.FromMoment(dbValue)
	fmt.Println("Retrieved:", retrieved.Format(time.RFC3339))
	// Output:
	// Store in DB: 2024-12-25T14:00:00Z
	// Retrieved: 2024-12-25T14:00:00Z
}
