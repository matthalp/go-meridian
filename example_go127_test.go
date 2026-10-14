//go:build go1.27

package meridian_test

import (
	"fmt"
	"time"

	"github.com/matthalp/go-meridian/v3/timezones/et"
	"github.com/matthalp/go-meridian/v3/timezones/pt"
	"github.com/matthalp/go-meridian/v3/timezones/utc"
)

// ExampleTime_To demonstrates the fluent method form of timezone conversion,
// which needs Go 1.27 or later.
func ExampleTime_To() {
	// Create a time in Eastern timezone
	eastern := et.Date(2024, time.December, 25, 9, 0, 0, 0)

	// Convert with the To method instead of pt.FromMoment(eastern)
	pacific := eastern.To[pt.Timezone]()
	universal := eastern.To[utc.Timezone]()

	// Conversions chain naturally
	roundTrip := eastern.To[utc.Timezone]().To[pt.Timezone]()

	fmt.Println("ET:", eastern.Format("15:04 MST"))
	fmt.Println("PT:", pacific.Format("15:04 MST"))
	fmt.Println("UTC:", universal.Format("15:04 MST"))
	fmt.Println("Same moment:", eastern.Equal(pacific) && eastern.Equal(roundTrip))
	// Output:
	// ET: 09:00 EST
	// PT: 06:00 PST
	// UTC: 14:00 UTC
	// Same moment: true
}
