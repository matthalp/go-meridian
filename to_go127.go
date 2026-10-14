//go:build go1.27

package meridian

// To converts t to the same instant expressed in the timezone TZ2, preserving
// the moment in time (UTC equality) while changing the timezone type.
//
// To is the method form of the package-level FromMoment function and reads
// fluently at the call site:
//
//	eastern := et.Now()
//	pacific := eastern.To[pt.Timezone]()  // same instant, Pacific type
//
// To is a generic method, a language feature added in Go 1.27, so it is only
// available when building with Go 1.27 or later. The rest of the package
// supports older Go releases; FromMoment works everywhere. Because To operates
// on an existing Time[TZ], it cannot start from a plain time.Time; use a
// timezone package's FromMoment (e.g. utc.FromMoment(stdTime)) for that.
func (t Time[TZ]) To[TZ2 Timezone]() Time[TZ2] {
	return Time[TZ2](t)
}
