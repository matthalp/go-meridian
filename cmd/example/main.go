// Package main provides an example usage of the meridian package.
package main

import (
	"fmt"
	"time"

	"github.com/matthalp/go-meridian"
)

// UTC represents the Coordinated Universal Time timezone.
type UTC struct{}

// Location returns the UTC time.Location.
func (UTC) Location() *time.Location {
	return time.UTC
}

// EST represents the Eastern Standard Time timezone.
type EST struct{}

// Location returns the EST time.Location.
func (EST) Location() *time.Location {
	loc, _ := time.LoadLocation("America/New_York")
	return loc
}

func main() {
	fmt.Println("Meridian Package Example")
	fmt.Println("=========================")
	fmt.Printf("Version: %s\n\n", meridian.Version)

	// Example 1: Get current time in different timezones
	fmt.Println("1. Current Time:")
	utcNow := meridian.Now[UTC]()
	estNow := meridian.Now[EST]()
	fmt.Printf("   UTC: %s\n", utcNow.Format(time.RFC3339))
	fmt.Printf("   EST: %s\n", estNow.Format(time.RFC3339))
	fmt.Println()

	// Example 2: Create a specific date
	fmt.Println("2. Specific Date:")
	meeting := meridian.Date[EST](2024, time.December, 25, 10, 30, 0, 0)
	fmt.Printf("   Meeting time (EST): %s\n", meeting.Format("Monday, January 2, 2006 at 3:04 PM MST"))
	fmt.Println()

	// Example 3: Different time formats
	fmt.Println("3. Different Formats:")
	t := meridian.Date[UTC](2024, time.June, 15, 14, 30, 45, 0)
	fmt.Printf("   RFC3339: %s\n", t.Format(time.RFC3339))
	fmt.Printf("   Kitchen: %s\n", t.Format(time.Kitchen))
	fmt.Printf("   Custom:  %s\n", t.Format("2006-01-02 15:04:05"))
}
