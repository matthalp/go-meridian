# Usage Guide

## For Package Users

### Installing the Package

Install the meridian package and timezone subpackages:

```bash
go get github.com/matthalp/go-meridian
```

### Basic Usage - Timezone Packages

The easiest way to use Meridian is through the timezone-specific packages:

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/matthalp/go-meridian/est"
    "github.com/matthalp/go-meridian/pst"
    "github.com/matthalp/go-meridian/utc"
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
    "github.com/matthalp/go-meridian/utc"
    "github.com/matthalp/go-meridian/est"
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

Meridian provides `Convert()` functions in each timezone package to convert between timezones:

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/matthalp/go-meridian/est"
    "github.com/matthalp/go-meridian/pst"
    "github.com/matthalp/go-meridian/utc"
)

func main() {
    // Create a time in EST
    meeting := est.Date(2024, time.December, 25, 10, 30, 0, 0)
    
    // Convert to other timezones
    utcMeeting := utc.Convert(meeting)
    pstMeeting := pst.Convert(meeting)
    
    fmt.Printf("EST: %s\n", meeting.Format(time.Kitchen))    // 10:30AM
    fmt.Printf("UTC: %s\n", utcMeeting.Format(time.Kitchen)) // 3:30PM
    fmt.Printf("PST: %s\n", pstMeeting.Format(time.Kitchen)) // 7:30AM
    
    // All represent the same moment in time
    fmt.Println(meeting.UTC().Equal(utcMeeting.UTC()))  // true
    fmt.Println(meeting.UTC().Equal(pstMeeting.UTC()))  // true
}
```

### Converting from time.Time

The `Moment` interface allows seamless conversion from standard `time.Time`:

```go
package main

import (
    "time"
    "github.com/matthalp/go-meridian/utc"
    "github.com/matthalp/go-meridian/est"
)

func processStandardTime(stdTime time.Time) {
    // Convert to type-safe timezone types
    utcTime := utc.Convert(stdTime)
    estTime := est.Convert(stdTime)
    
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
    
    "github.com/matthalp/go-meridian/est"
    "github.com/matthalp/go-meridian/pst"
    "github.com/matthalp/go-meridian/utc"
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
    "github.com/matthalp/go-meridian"
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
├── meridian.go          # Core generic types and functions
├── meridian_test.go     # Core package tests
├── example_test.go      # Testable examples (appear in docs)
├── doc.go               # Package-level documentation
├── cmd/example/         # Example program using the package
├── utc/                 # UTC timezone package
│   ├── utc.go
│   └── utc_test.go
├── est/                 # EST timezone package
│   ├── est.go
│   └── est_test.go
├── pst/                 # PST timezone package
│   ├── pst.go
│   └── pst_test.go
└── ...                  # CI/CD and config files
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

To add a new timezone package (e.g., `jst` for Japan Standard Time):

1. Create a new directory `jst/` with `jst.go`:
   ```go
   package jst
   
   import (
       "fmt"
       "time"
       "github.com/matthalp/go-meridian"
   )
   
   var location = mustLoadLocation("Asia/Tokyo")
   
   func mustLoadLocation(name string) *time.Location {
       loc, err := time.LoadLocation(name)
       if err != nil {
           panic(fmt.Sprintf("failed to load timezone %s: %v", name, err))
       }
       return loc
   }
   
   type Timezone struct{}
   
   func (Timezone) Location() *time.Location {
       return location
   }
   
   type Time = meridian.Time[Timezone]
   
   func Now() Time {
       return meridian.Now[Timezone]()
   }
   
   func Date(year int, month time.Month, day, hour, minute, sec, nsec int) Time {
       return meridian.Date[Timezone](year, month, day, hour, minute, sec, nsec)
   }
   
   func Convert(m meridian.Moment) Time {
       return meridian.FromMoment[Timezone](m)
   }
   
   func Parse(layout, value string) (Time, error) {
       t, err := time.ParseInLocation(layout, value, location)
       if err != nil {
           return Time{}, err
       }
       return meridian.FromMoment[Timezone](t), nil
   }
   
   func Unix(sec, nsec int64) Time {
       return meridian.FromMoment[Timezone](time.Unix(sec, nsec))
   }
   
   func UnixMilli(msec int64) Time {
       return meridian.FromMoment[Timezone](time.UnixMilli(msec))
   }
   
   func UnixMicro(usec int64) Time {
       return meridian.FromMoment[Timezone](time.UnixMicro(usec))
   }
   ```

2. Add tests in `jst/jst_test.go` following the pattern in `utc/utc_test.go`

3. Update documentation to include the new timezone package

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

