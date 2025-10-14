package meridian_test

import (
	"fmt"
	"time"

	"github.com/matthalp/go-meridian"
	"github.com/matthalp/go-meridian/est"
	"github.com/matthalp/go-meridian/pst"
	"github.com/matthalp/go-meridian/utc"
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
	eastern := est.Date(2024, time.December, 25, 9, 0, 0, 0)

	// Explicitly convert to Pacific time
	pacific := pst.FromMoment(eastern)

	// Explicitly convert to UTC
	universal := utc.FromMoment(eastern)

	// All represent the same moment in time
	fmt.Println("EST:", eastern.Format("15:04 MST"))
	fmt.Println("PST:", pacific.Format("15:04 MST"))
	fmt.Println("UTC:", universal.Format("15:04 MST"))
	fmt.Println("Same moment:", eastern.Equal(pacific) && eastern.Equal(universal))
	// Output:
	// EST: 09:00 EST
	// PST: 06:00 PST
	// UTC: 14:00 UTC
	// Same moment: true
}

// ExampleFromMoment_timeTime demonstrates converting from standard time.Time.
func ExampleFromMoment_timeTime() {
	// Standard library time.Time
	stdTime := time.Date(2024, time.June, 15, 14, 30, 0, 0, time.UTC)

	// Convert to timezone-specific types
	utcTyped := utc.FromMoment(stdTime)
	estTyped := est.FromMoment(stdTime)
	pstTyped := pst.FromMoment(stdTime)

	fmt.Println("UTC:", utcTyped.Format("3:04 PM MST"))
	fmt.Println("EST:", estTyped.Format("3:04 PM MST"))
	fmt.Println("PST:", pstTyped.Format("3:04 PM MST"))
	// Output:
	// UTC: 2:30 PM UTC
	// EST: 10:30 AM EDT
	// PST: 7:30 AM PDT
}

// Example_typeSafety demonstrates compile-time timezone safety.
func Example_typeSafety() {
	// Function that only accepts UTC times
	processUTC := func(t utc.Time) {
		fmt.Println("Processing UTC:", t.Format("15:04 MST"))
	}

	// Function that only accepts EST times
	processEST := func(t est.Time) {
		fmt.Println("Processing EST:", t.Format("15:04 MST"))
	}

	utcTime := utc.Date(2024, time.June, 15, 14, 0, 0, 0)
	estTime := est.Date(2024, time.June, 15, 10, 0, 0, 0)

	// These work - types match
	processUTC(utcTime)
	processEST(estTime)

	// These would NOT compile (uncomment to see the error):
	// processUTC(estTime)  // Compile error: cannot use est.Time as utc.Time
	// processEST(utcTime)  // Compile error: cannot use utc.Time as est.Time

	// To convert, you must be explicit:
	processUTC(utc.FromMoment(estTime))
	processEST(est.FromMoment(utcTime))
	// Output:
	// Processing UTC: 14:00 UTC
	// Processing EST: 10:00 EDT
	// Processing UTC: 14:00 UTC
	// Processing EST: 10:00 EDT
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
	estTime := est.Date(2024, time.June, 15, 10, 30, 0, 0)
	pstTime := pst.Date(2024, time.June, 15, 7, 30, 0, 0)

	fmt.Println("utc.Time:", formatMoment(utcTime))
	fmt.Println("est.Time:", formatMoment(estTime))
	fmt.Println("pst.Time:", formatMoment(pstTime))
	// Output:
	// time.Time: 2024-06-15 14:30:00 UTC
	// utc.Time: 2024-06-15 14:30:00 UTC
	// est.Time: 2024-06-15 14:30:00 UTC
	// pst.Time: 2024-06-15 14:30:00 UTC
}

// Example_multipleTimezones demonstrates working with multiple timezones.
func Example_multipleTimezones() {
	// Schedule a meeting at 2pm EST
	meetingEST := est.Date(2024, time.December, 25, 14, 0, 0, 0)

	// What time is that in other timezones?
	meetingPST := pst.FromMoment(meetingEST)
	meetingUTC := utc.FromMoment(meetingEST)

	fmt.Println("Meeting scheduled:")
	fmt.Println("  Eastern:", meetingEST.Format("3:04 PM MST"))
	fmt.Println("  Pacific:", meetingPST.Format("3:04 PM MST"))
	fmt.Println("  UTC:", meetingUTC.Format("3:04 PM MST"))

	// Add timezone-specific operations
	oneHourLater := meetingEST.Add(time.Hour)
	fmt.Println("  One hour later (EST):", oneHourLater.Format("3:04 PM MST"))
	// Output:
	// Meeting scheduled:
	//   Eastern: 2:00 PM EST
	//   Pacific: 11:00 AM PST
	//   UTC: 7:00 PM UTC
	//   One hour later (EST): 3:00 PM EST
}

// ExampleParse demonstrates parsing time strings in specific timezones.
func ExampleParse() {
	// Parse a time string as EST
	estTime, _ := est.Parse(time.RFC3339, "2024-12-25T09:00:00-05:00")
	fmt.Println("Parsed EST:", estTime.Format("15:04 MST"))

	// Parse a time string as PST
	pstTime, _ := pst.Parse(time.RFC3339, "2024-12-25T06:00:00-08:00")
	fmt.Println("Parsed PST:", pstTime.Format("15:04 MST"))

	// These represent the same moment
	fmt.Println("Same moment:", estTime.Equal(pstTime))
	// Output:
	// Parsed EST: 09:00 EST
	// Parsed PST: 06:00 PST
	// Same moment: true
}

// Example_genericFunction demonstrates writing generic functions that work with any timezone.
func Example_genericFunction() {
	// Generic function that preserves the timezone type parameter
	addBusinessDay := func(t utc.Time) utc.Time {
		return t.Add(24 * time.Hour)
	}

	// Another function for EST times
	addBusinessDayEST := func(t est.Time) est.Time {
		return t.Add(24 * time.Hour)
	}

	// Each function maintains type safety for its timezone
	utcTime := utc.Date(2024, time.June, 14, 12, 0, 0, 0)
	estTime := est.Date(2024, time.June, 14, 8, 0, 0, 0)

	nextUTC := addBusinessDay(utcTime)
	nextEST := addBusinessDayEST(estTime)

	fmt.Println("UTC next day:", nextUTC.Format("Monday, Jan 2"))
	fmt.Println("EST next day:", nextEST.Format("Monday, Jan 2"))
	// Output:
	// UTC next day: Saturday, Jun 15
	// EST next day: Saturday, Jun 15
}

// Example_databaseStorage demonstrates storing times in a database-friendly way.
func Example_databaseStorage() {
	// Create a time in EST
	estTime := est.Date(2024, time.December, 25, 9, 0, 0, 0)

	// Convert to UTC for database storage (standard practice)
	utcForDB := utc.FromMoment(estTime)

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
