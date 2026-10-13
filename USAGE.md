# Usage Guide

## For Package Users

### Installing the Package

Install the meridian package and timezone subpackages:

```bash
go get github.com/matthalp/go-meridian/v2
```

### Basic Usage - Timezone Packages

The easiest way to use Meridian is through the timezone-specific packages:

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/matthalp/go-meridian/v2/timezones/et"
    "github.com/matthalp/go-meridian/v2/timezones/pt"
    "github.com/matthalp/go-meridian/v2/timezones/utc"
)

func main() {
    // Get current time in different timezones
    utcNow := utc.Now()
    etNow := et.Now()
    ptNow := pt.Now()
    
    fmt.Printf("UTC: %s\n", utcNow.Format(time.RFC3339))
    fmt.Printf("ET:  %s\n", etNow.Format(time.RFC3339))
    fmt.Printf("PT:  %s\n", ptNow.Format(time.RFC3339))
    
    // Create a specific date/time
    meeting := et.Date(2024, time.December, 25, 10, 30, 0, 0)
    fmt.Printf("Meeting: %s\n", meeting.Format(time.Kitchen))
}
```

### Type-Safe Function Signatures

Use timezone-specific types in your function signatures for compiler-enforced correctness:

```go
package events

import (
    "database/sql"

    "github.com/matthalp/go-meridian/v2/timezones/et"
    "github.com/matthalp/go-meridian/v2/timezones/utc"
)

// Function only accepts UTC times
func storeInDatabase(db *sql.DB, t utc.Time) error {
    // You know for certain this is UTC
    _, err := db.Exec("INSERT INTO events (timestamp) VALUES (?)", t)
    return err
}

// Function only accepts Eastern Time values for display
func displayToUser(t et.Time) string {
    return t.Format("3:04 PM MST")
}

func run(db *sql.DB) error {
    // ✅ This compiles
    if err := storeInDatabase(db, utc.Now()); err != nil {
        return err
    }

    // ❌ This won't compile - type safety!
    // storeInDatabase(db, et.Now())
    
    // ✅ This compiles
    displayToUser(et.Now())
    
    // ❌ This won't compile
    // displayToUser(utc.Now())

    return nil
}
```

### Converting Between Timezones

Each timezone package provides a `FromMoment()` function to convert into that timezone:

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/matthalp/go-meridian/v2/timezones/et"
    "github.com/matthalp/go-meridian/v2/timezones/pt"
    "github.com/matthalp/go-meridian/v2/timezones/utc"
)

func main() {
    // Create a time in Eastern Time
    meeting := et.Date(2024, time.December, 25, 10, 30, 0, 0)
    
    // Convert to other timezones
    utcMeeting := utc.FromMoment(meeting)
    ptMeeting := pt.FromMoment(meeting)
    
    fmt.Printf("ET:  %s\n", meeting.Format(time.Kitchen))    // 10:30AM
    fmt.Printf("UTC: %s\n", utcMeeting.Format(time.Kitchen)) // 3:30PM
    fmt.Printf("PT:  %s\n", ptMeeting.Format(time.Kitchen)) // 7:30AM
    
    // All represent the same moment in time
    fmt.Println(meeting.UTC().Equal(utcMeeting.UTC()))  // true
    fmt.Println(meeting.UTC().Equal(ptMeeting.UTC()))  // true
}
```

### Converting from time.Time

The `Moment` interface allows seamless conversion from standard `time.Time`:

```go
package main

import (
    "time"
    "github.com/matthalp/go-meridian/v2/timezones/utc"
    "github.com/matthalp/go-meridian/v2/timezones/et"
)

func processStandardTime(stdTime time.Time) {
    // Convert to type-safe timezone types
    utcTime := utc.FromMoment(stdTime)
    etTime := et.FromMoment(stdTime)
    
    // Now you have type-safe times for your functions
    storeInDatabase(utcTime)     // Function requires utc.Time
    displayToUser(etTime)       // Function requires et.Time
}

func storeInDatabase(t utc.Time) { /* ... */ }
func displayToUser(t et.Time) { /* ... */ }
```

### Creating Times from Various Sources

Each timezone package provides multiple factory methods for creating times:

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/matthalp/go-meridian/v2/timezones/et"
    "github.com/matthalp/go-meridian/v2/timezones/pt"
    "github.com/matthalp/go-meridian/v2/timezones/utc"
)

func main() {
    // From current time
    now := utc.Now()
    fmt.Println(now.Format(time.RFC3339))
    
    // From specific date/time
    meeting := et.Date(2024, time.December, 25, 10, 30, 0, 0)
    fmt.Println(meeting.Format(time.Kitchen))

    // From formatted string (parsed in the timezone's location)
    parsed, err := utc.Parse(time.RFC3339, "2024-01-15T12:00:00Z")
    if err != nil {
        panic(err)
    }
    fmt.Println(parsed.Format(time.Kitchen))
    
    // From Unix timestamp (seconds + nanoseconds)
    t1 := utc.Unix(1705320000, 0)
    
    // From Unix milliseconds
    t2 := pt.UnixMilli(1705320000000)
    
    // From Unix microseconds
    t3 := et.UnixMicro(1705320000000000)
    
    // All timestamps represent the same moment
    fmt.Println(t1.UTC().Equal(t2.UTC())) // true
    fmt.Println(t2.UTC().Equal(t3.UTC())) // true
}
```

**Important Note**: `ParseInLocation` from the standard `time` package is not needed in Meridian timezone packages because the location is already determined by the package (e.g., `et.Parse` always parses in Eastern Time, `utc.Parse` in UTC).

### Timezone-Specific Parsing

The `Parse` function in each timezone package interprets the input string in that timezone's location:

```go
// Parse the same clock time in different timezones
etTime, _ := et.Parse("2006-01-02 15:04:05", "2024-01-15 12:00:00")
ptTime, _ := pt.Parse("2006-01-02 15:04:05", "2024-01-15 12:00:00")
utcTime, _ := utc.Parse("2006-01-02 15:04:05", "2024-01-15 12:00:00")

// These represent different moments in time!
// Eastern noon (EST in January) happens 5 hours after UTC noon
// Pacific noon (PST in January) happens 8 hours after UTC noon
fmt.Println(etTime.UTC()) // 2024-01-15 17:00:00 +0000 UTC
fmt.Println(ptTime.UTC()) // 2024-01-15 20:00:00 +0000 UTC
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
logEvent(et.Now(), "completed")
```

### Advanced Usage - Generic API

For custom timezones or advanced usage, use the generic API:

```go
package main

import (
    "fmt"
    "time"

    "github.com/matthalp/go-meridian/v2"
)

// Define a custom timezone. Load the location once, not on every call.
var seoul = mustLoadLocation("Asia/Seoul")

func mustLoadLocation(name string) *time.Location {
    loc, err := time.LoadLocation(name)
    if err != nil {
        panic(err)
    }
    return loc
}

type KST struct{}

func (KST) Location() *time.Location { return seoul }

func main() {
    // Use the generic API with your custom timezone
    now := meridian.Now[KST]()
    meeting := meridian.Date[KST](2024, time.June, 15, 14, 30, 0, 0)
    fmt.Println(now, meeting)
}
```

## For Package Developers

### Project Structure

```
go-meridian/
├── meridian.go               # Core generic types and functions
├── meridian_test.go          # Core package tests
├── example_test.go           # Testable examples (appear in docs)
├── doc.go                    # Package-level documentation
├── cmd/example/              # Example program using the package
├── cmd/generate-timezones/   # Timezone package generator
├── timezones.yaml            # Timezone definitions (generator input)
├── timezones/                # Generated timezone packages
│   ├── et/
│   │   ├── et.go
│   │   └── et_test.go
│   ├── pt/
│   ├── utc/
│   └── ...
└── ...                       # CI/CD and config files
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
5. **Add a CHANGELOG entry** under `## [Unreleased]`
6. **Open a pull request.** Releases are cut separately (see "Publishing Updates")

### Adding New Timezone Packages

Timezone packages are generated; don't write them by hand. To add one (e.g., `kst` for Korea Standard Time):

1. Add an entry to `timezones.yaml`:
   ```yaml
   - name: kst
     location: Asia/Seoul
     description: Korea Standard Time
   ```

2. Run `make generate`. This writes `timezones/kst/kst.go` and `timezones/kst/kst_test.go`.

3. Run `make test` and `make lint`, then commit the YAML change together with the generated files. CI fails if generated code is out of sync with `timezones.yaml`.

4. Add the package to the timezone lists in `README.md` and `doc.go`.

### Publishing Updates

See "Publishing a New Version" in the README. Releases are cut by the Release
workflow, not by pushing tags by hand.

### Best Practices

- ✅ Keep the API surface small and focused
- ✅ Write clear documentation comments
- ✅ Maintain 100% test coverage for public APIs
- ✅ Follow Go naming conventions
- ✅ Use semantic versioning
- ✅ Keep backward compatibility within major versions
- ✅ Document breaking changes clearly

