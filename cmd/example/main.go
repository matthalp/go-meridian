// Package main provides an example usage of the meridian package.
package main

import (
	"fmt"
	"time"

	"github.com/matthalp/go-meridian"
	"github.com/matthalp/go-meridian/est"
	"github.com/matthalp/go-meridian/pst"
	"github.com/matthalp/go-meridian/utc"
)

func main() {
	fmt.Println("Meridian Package Example")
	fmt.Println("=========================")
	fmt.Printf("Version: %s\n\n", meridian.Version)

	// Example 1: Get current time in different timezones
	fmt.Println("1. Current Time:")
	utcNow := utc.Now()
	estNow := est.Now()
	pstNow := pst.Now()
	fmt.Printf("   UTC: %s\n", utcNow.Format(time.RFC3339))
	fmt.Printf("   EST: %s\n", estNow.Format(time.RFC3339))
	fmt.Printf("   PST: %s\n", pstNow.Format(time.RFC3339))
	fmt.Println()

	// Example 2: Create a specific date
	fmt.Println("2. Specific Date:")
	meeting := est.Date(2024, time.December, 25, 10, 30, 0, 0)
	fmt.Printf("   Meeting time (EST): %s\n", meeting.Format("Monday, January 2, 2006 at 3:04 PM MST"))
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
	printESTTime(est.Now())
	// Note: printUTCTime(est.Now()) would not compile due to type safety
}

// printUTCTime accepts only UTC times.
func printUTCTime(t utc.Time) {
	fmt.Printf("   UTC Time: %s\n", t.Format(time.RFC3339))
}

// printESTTime accepts only EST times.
func printESTTime(t est.Time) {
	fmt.Printf("   EST Time: %s\n", t.Format(time.RFC3339))
}
