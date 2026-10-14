//go:build go1.27

package meridian

import (
	"testing"
	"time"
)

// TestTo verifies the generic To method preserves the instant while changing
// the timezone type, matching FromMoment and supporting chained conversions.
func TestTo(t *testing.T) {
	t.Parallel()

	// 2024-01-15 17:00:00 UTC.
	original := Date[UTC](2024, time.January, 15, 17, 0, 0, 0)

	t.Run("preserves instant across type change", func(t *testing.T) {
		eastern := original.To[EST]()
		pacific := original.To[PST]()

		if !eastern.UTC().Equal(original.UTC()) {
			t.Errorf("To[EST]() UTC = %v, want %v", eastern.UTC(), original.UTC())
		}
		if !pacific.UTC().Equal(original.UTC()) {
			t.Errorf("To[PST]() UTC = %v, want %v", pacific.UTC(), original.UTC())
		}
	})

	t.Run("applies the new timezone", func(t *testing.T) {
		if got, want := original.To[EST]().Hour(), 12; got != want {
			t.Errorf("To[EST]().Hour() = %d, want %d", got, want)
		}
	})

	t.Run("matches FromMoment", func(t *testing.T) {
		if got, want := original.To[PST](), FromMoment[PST](original); got != want {
			t.Errorf("To[PST]() = %v, FromMoment[PST]() = %v; want equal", got, want)
		}
	})

	t.Run("chains and round-trips", func(t *testing.T) {
		if roundTrip := original.To[EST]().To[PST]().To[UTC](); roundTrip != original {
			t.Errorf("round-trip To chain = %v, want %v", roundTrip, original)
		}
	})

	t.Run("converting to same type is identity", func(t *testing.T) {
		if got := original.To[UTC](); got != original {
			t.Errorf("To[UTC]() = %v, want %v", got, original)
		}
	})
}
