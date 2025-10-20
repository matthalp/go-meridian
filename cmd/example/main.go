// Package main provides an example usage of the meridian package.
package main

import (
	"fmt"
	"time"

	"github.com/matthalp/go-meridian"
	"github.com/matthalp/go-meridian/et"
	"github.com/matthalp/go-meridian/pt"
	"github.com/matthalp/go-meridian/utc"
)

func main() {
	fmt.Println("Meridian Package Example")
	fmt.Println("=========================")
	fmt.Printf("Version: %s\n\n", meridian.Version)

	// Example 1: Get current time in different timezones
	fmt.Println("1. Current Time:")
	utcNow := utc.Now()
	etNow := et.Now()
	ptNow := pt.Now()
	fmt.Printf("   UTC: %s\n", utcNow.Format(time.RFC3339))
	fmt.Printf("   ET: %s\n", etNow.Format(time.RFC3339))
	fmt.Printf("   PT: %s\n", ptNow.Format(time.RFC3339))
	fmt.Println()

	// Example 2: Create a specific date
	fmt.Println("2. Specific Date:")
	meeting := et.Date(2024, time.December, 25, 10, 30, 0, 0)
	fmt.Printf("   Meeting time (ET): %s\n", meeting.Format("Monday, January 2, 2006 at 3:04 PM MST"))
	fmt.Println()

	// Example 3: Different time formats
	fmt.Println("3. Different Formats:")
	t := utc.Date(2024, time.June, 15, 14, 30, 45, 0)
	fmt.Printf("   RFC3339: %s\n", t.Format(time.RFC3339))
	fmt.Printf("   Kitchen: %s\n", t.Format(time.Kitchen))
	fmt.Printf("   Custom:  %s\n", t.Format("2006-01-02 15:04:05"))
	fmt.Println()

	// Example 4: Type-safe function signatures using timezone-specific types
	fmt.Println("4. Type-Safe Function Signatures:")
	printUTCTime(utc.Now())
	printETTime(et.Now())
	// Note: printUTCTime(et.Now()) would not compile due to type safety
	fmt.Println()

	// Example 5: Converting between timezones
	fmt.Println("5. Timezone Conversion:")
	etMeeting := et.Date(2024, time.December, 25, 10, 30, 0, 0)
	utcMeeting := utc.FromMoment(etMeeting)
	ptMeeting := pt.FromMoment(etMeeting)

	fmt.Printf("   Meeting ET: %s\n", etMeeting.Format(time.Kitchen))
	fmt.Printf("   Meeting UTC: %s\n", utcMeeting.Format(time.Kitchen))
	fmt.Printf("   Meeting PT: %s\n", ptMeeting.Format(time.Kitchen))
	fmt.Println()

	// Example 6: Converting from standard time.Time
	fmt.Println("6. Converting from time.Time:")
	stdTime := time.Date(2024, time.June, 15, 18, 0, 0, 0, time.UTC)
	fmt.Printf("   Standard time.Time: %s\n", stdTime.Format(time.RFC3339))

	// Convert to timezone-specific types
	utcFromStd := utc.FromMoment(stdTime)
	etFromStd := et.FromMoment(stdTime)
	ptFromStd := pt.FromMoment(stdTime)

	fmt.Printf("   As UTC: %s\n", utcFromStd.Format("3:04 PM MST"))
	fmt.Printf("   As ET: %s\n", etFromStd.Format("3:04 PM MST"))
	fmt.Printf("   As PT: %s\n", ptFromStd.Format("3:04 PM MST"))
}

// printUTCTime accepts only UTC times.
func printUTCTime(t utc.Time) {
	fmt.Printf("   UTC Time: %s\n", t.Format(time.RFC3339))
}

// printETTime accepts only ET times.
func printETTime(t et.Time) {
	fmt.Printf("   ET Time: %s\n", t.Format(time.RFC3339))
}
