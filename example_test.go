package meridian_test

import (
	"fmt"
	"time"

	"github.com/matthalp/go-meridian"
)

// UTC timezone for examples.
type UTC struct{}

func (UTC) Location() *time.Location {
	return time.UTC
}

func ExampleNow() {
	// Get the current time in UTC
	now := meridian.Now[UTC]()

	// Format it
	fmt.Println("Current time format:", now.Format("2006-01-02"))
	// Output will vary, so we can't test exact output
}

func ExampleDate() {
	// Create a specific time in UTC
	t := meridian.Date[UTC](2024, time.January, 15, 14, 30, 0, 0)

	// Format the time
	fmt.Println(t.Format("2006-01-02 15:04:05"))
	// Output: 2024-01-15 14:30:00
}

func ExampleTime_Format() {
	// Create a specific time in UTC
	t := meridian.Date[UTC](2024, time.June, 15, 9, 30, 0, 0)

	// Format the time in different layouts
	fmt.Println(t.Format(time.RFC3339))
	fmt.Println(t.Format("Monday, January 2, 2006"))
	// Output:
	// 2024-06-15T09:30:00Z
	// Saturday, June 15, 2024
}
