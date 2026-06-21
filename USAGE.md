# Usage Guide

## For Package Users

### Installing the Package

Install the meridian package and timezone subpackages:

```bash
go get github.com/matthalp/go-meridian/v3
```

> **Requires Go 1.27 or later.** v3 uses generic methods (the `To` conversion
> method), introduced in Go 1.27. On older toolchains, use v2.

### Basic Usage - Timezone Packages

The easiest way to use Meridian is through the timezone-specific packages:

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/matthalp/go-meridian/v3/timezones/est"
    "github.com/matthalp/go-meridian/v3/timezones/pst"
    "github.com/matthalp/go-meridian/v3/timezones/utc"
)

func main() {
    // Get current time in different timezones
    utcNow := utc.Now()
    estNow := est.Now()
    pstNow := pst.Now()
    
    fmt.Printf("UTC: %s\n", utcNow.Format(time.RFC3339))
    fmt.Printf("EST: %s\n", estNow.Format(time.RFC3339))
    fmt.Printf("PST: %s\n", pstNow.Format(time.RFC3339))
    
    // Create a specific date/time
    meeting := est.Date(2024, time.December, 25, 10, 30, 0, 0)
    fmt.Printf("Meeting: %s\n", meeting.Format(time.Kitchen))
}
```

### Type-Safe Function Signatures

Use timezone-specific types in your function signatures for compiler-enforced correctness:

```go
package main

import (
    "database/sql"
    "github.com/matthalp/go-meridian/v3/timezones/utc"
    "github.com/matthalp/go-meridian/v3/timezones/est"
)

// Function only accepts UTC times
func storeInDatabase(db *sql.DB, t utc.Time) error {
    // You know for certain this is UTC
    return db.Exec("INSERT INTO events (timestamp) VALUES (?)", 
                   t.Format(time.RFC3339))
}

// Function only accepts EST times for display
func displayToUser(t est.Time) string {
    return t.Format("3:04 PM MST")
}

func main() {
    // ✅ This compiles
    storeInDatabase(db, utc.Now())
    
    // ❌ This won't compile - type safety!
    // storeInDatabase(db, est.Now())
    
    // ✅ This compiles
    displayToUser(est.Now())
    
    // ❌ This won't compile
    // displayToUser(utc.Now())
}
```

### Converting Between Timezones

For `Time` → `Time` conversions, use the `To` method (a Go 1.27 generic method).
It is the fluent equivalent of a package's `FromMoment` function and chains:

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/matthalp/go-meridian/v3/timezones/est"
    "github.com/matthalp/go-meridian/v3/timezones/pst"
    "github.com/matthalp/go-meridian/v3/timezones/utc"
)

func main() {
    // Create a time in EST
    meeting := est.Date(2024, time.December, 25, 10, 30, 0, 0)
    
    // Convert to other timezones with the To method
    utcMeeting := meeting.To[utc.Timezone]()
    pstMeeting := meeting.To[pst.Timezone]()
    
    // Equivalent package-function form:
    //   utcMeeting := utc.FromMoment(meeting)
    //   pstMeeting := pst.FromMoment(meeting)
    
    fmt.Printf("EST: %s\n", meeting.Format(time.Kitchen))    // 10:30AM
    fmt.Printf("UTC: %s\n", utcMeeting.Format(time.Kitchen)) // 3:30PM
    fmt.Printf("PST: %s\n", pstMeeting.Format(time.Kitchen)) // 7:30AM
    
    // All represent the same moment in time
    fmt.Println(meeting.UTC().Equal(utcMeeting.UTC()))  // true
    fmt.Println(meeting.UTC().Equal(pstMeeting.UTC()))  // true
}
```

### Converting from time.Time

The `Moment` interface allows seamless conversion from standard `time.Time`.
Use `FromMoment` for this — a plain `time.Time` has no `To` method:

```go
package main

import (
    "time"
    "github.com/matthalp/go-meridian/v3/timezones/utc"
    "github.com/matthalp/go-meridian/v3/timezones/est"
)

func processStandardTime(stdTime time.Time) {
    // Convert to type-safe timezone types
    utcTime := utc.FromMoment(stdTime)
    estTime := est.FromMoment(stdTime)
    
    // Now you have type-safe times for your functions
    storeInDatabase(utcTime)     // Function requires utc.Time
    displayToUser(estTime)       // Function requires est.Time
}

func storeInDatabase(t utc.Time) { /* ... */ }
func displayToUser(t est.Time) { /* ... */ }
```

### Creating Times from Various Sources

Each timezone package provides multiple factory methods for creating times:

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/matthalp/go-meridian/v3/timezones/est"
    "github.com/matthalp/go-meridian/v3/timezones/pst"
    "github.com/matthalp/go-meridian/v3/timezones/utc"
)

func main() {
    // From current time
    now := utc.Now()
    
    // From specific date/time
    meeting := est.Date(2024, time.December, 25, 10, 30, 0, 0)
    
    // From formatted string (parsed in the timezone's location)
    parsed, err := utc.Parse(time.RFC3339, "2024-01-15T12:00:00Z")
    if err != nil {
        panic(err)
    }
    fmt.Println(parsed.Format(time.Kitchen))
    
    // From Unix timestamp (seconds + nanoseconds)
    t1 := utc.Unix(1705320000, 0)
    
    // From Unix milliseconds
    t2 := pst.UnixMilli(1705320000000)
    
    // From Unix microseconds
    t3 := est.UnixMicro(1705320000000000)
    
    // All timestamps represent the same moment
    fmt.Println(t1.UTC().Equal(t2.UTC())) // true
    fmt.Println(t2.UTC().Equal(t3.UTC())) // true
}
```

**Important Note**: `ParseInLocation` from the standard `time` package is not needed in Meridian timezone packages because the location is already determined by the package (e.g., `est.Parse` always parses in EST, `utc.Parse` in UTC).

### Timezone-Specific Parsing

The `Parse` function in each timezone package interprets the input string in that timezone's location:

```go
// Parse the same clock time in different timezones
estTime, _ := est.Parse("2006-01-02 15:04:05", "2024-01-15 12:00:00")
pstTime, _ := pst.Parse("2006-01-02 15:04:05", "2024-01-15 12:00:00")
utcTime, _ := utc.Parse("2006-01-02 15:04:05", "2024-01-15 12:00:00")

// These represent different moments in time!
// EST noon happens 5 hours after UTC noon
// PST noon happens 8 hours after UTC noon
fmt.Println(estTime.UTC()) // 2024-01-15 17:00:00 +0000 UTC
fmt.Println(pstTime.UTC()) // 2024-01-15 20:00:00 +0000 UTC
fmt.Println(utcTime.UTC()) // 2024-01-15 12:00:00 +0000 UTC
```

### The Moment Interface

Both `time.Time` and `meridian.Time[TZ]` implement the `Moment` interface:

```go
type Moment interface {
    UTC() time.Time
}
```

This allows functions to accept times from any source:

```go
func logEvent(m meridian.Moment, event string) {
    // Works with time.Time or any meridian.Time[TZ]
    utcTime := m.UTC()
    fmt.Printf("[%s] %s\n", utcTime.Format(time.RFC3339), event)
}

// Can be called with either
logEvent(time.Now(), "started")
logEvent(utc.Now(), "processing")
logEvent(est.Now(), "completed")
```

### Advanced Usage - Generic API

For custom timezones or advanced usage, use the generic API:

```go
package main

import (
    "time"
    "github.com/matthalp/go-meridian/v3"
)

// Define a custom timezone
type JST struct{}

func (JST) Location() *time.Location {
    loc, _ := time.LoadLocation("Asia/Tokyo")
    return loc
}

func main() {
    // Use the generic API with your custom timezone
    now := meridian.Now[JST]()
    meeting := meridian.Date[JST](2024, time.June, 15, 14, 30, 0, 0)
}
```

## For Package Developers

### Project Structure

```
go-meridian/
├── meridian.go              # Core generic types and functions
├── meridian_test.go         # Core package tests
├── example_test.go          # Testable examples (appear in docs)
├── doc.go                   # Package-level documentation
├── timezones.yaml           # Timezone definitions (source of generated packages)
├── cmd/example/             # Example program using the package
├── cmd/generate-timezones/  # Code generator for timezone packages
├── timezones/               # Generated timezone packages
│   ├── utc/                 # UTC timezone package (utc.go, utc_test.go)
│   ├── est/                 # EST timezone package
│   ├── pst/                 # PST timezone package
│   └── ...                  # et, pt, ct, mt, jst, ist, and more
└── ...                      # CI/CD and config files
```

### Key Concepts

1. **Exported vs Unexported**
   - Exported (public): Start with uppercase letter (e.g., `Greet`, `Version`)
   - Unexported (private): Start with lowercase letter (e.g., `helper`)

2. **Package Documentation**
   - Add comments above functions, types, constants
   - Comments starting with the name will appear in docs
   - Use `doc.go` for package-level documentation

3. **Testing**
   - Unit tests in `*_test.go` files
   - Example tests must start with `Example` and include `// Output:` comments
   - Run tests: `make test` or `go test ./...`

### Development Workflow

1. **Write code** in `meridian.go`
2. **Write tests** in `meridian_test.go`
3. **Write examples** in `example_test.go`
4. **Test locally**:
   ```bash
   make test           # Run all tests
   make test-coverage  # Generate coverage report
   make lint           # Run linter
   make run-example    # Test the example program
   ```
5. **Update version** in `meridian.go` and `CHANGELOG.md`
6. **Commit and tag**:
   ```bash
   git add .
   git commit -m "Add new feature"
   git tag v0.2.0
   git push origin main --tags
   ```

### Adding New Timezone Packages

Timezone packages are **generated** from `timezones.yaml` — you do not write them
by hand. To add a new timezone (e.g., `jst` for Japan Standard Time):

1. Add an entry to `timezones.yaml`:
   ```yaml
   timezones:
     - name: jst
       location: Asia/Tokyo
       description: Japan Standard Time
   ```

2. Regenerate the packages:
   ```bash
   make generate
   ```
   This creates `timezones/jst/jst.go` and `timezones/jst/jst_test.go`, each
   wired to the core `meridian` API (`Now`, `Date`, `FromMoment`, `Parse`,
   `Unix`, `UnixMilli`, `UnixMicro`).

3. Verify:
   ```bash
   make test   # generated tests pass
   make lint   # code quality
   ```

CI fails if the generated packages are out of sync with `timezones.yaml`, so
always commit the regenerated files. To convert between any two generated
timezones, use the `To` method (e.g. `jst.Now().To[utc.Timezone]()`).

### Publishing Updates

1. Update version number in `meridian.go`
2. Update `CHANGELOG.md` with changes
3. Create a git tag matching the version:
   ```bash
   git tag v0.2.0
   git push origin v0.2.0
   ```
4. GitHub Actions will run CI automatically
5. The package will be available on pkg.go.dev within minutes

### Best Practices

- ✅ Keep the API surface small and focused
- ✅ Write clear documentation comments
- ✅ Maintain 100% test coverage for public APIs
- ✅ Follow Go naming conventions
- ✅ Use semantic versioning
- ✅ Keep backward compatibility within major versions
- ✅ Document breaking changes clearly

