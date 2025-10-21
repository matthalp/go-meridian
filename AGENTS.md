# AGENTS.md

## Project Overview

**go-meridian** is a Go library that provides type-safe timezone handling using Go generics. It solves a fundamental problem: timezone information in `time.Time` is data, not type, and can be lost without the compiler noticing. Meridian makes timezone information immutable by encoding it directly into the type system.

**Core Philosophy**: "Make wrong timezone handling impossible to compile."

This is a **distributable library** intended for consumption by other Go projects, not an application. All changes should maintain backwards compatibility and API stability.

## Key Concepts

### 1. Timezones as Types
- `meridian.Time[TZ]` carries timezone information in its type parameter
- `meridian.Time[est.EST]` and `meridian.Time[pst.PST]` are **different types**
- The compiler prevents accidental timezone mixing or loss

### 2. Per-Timezone Packages
- Each timezone lives in its own package: `aest`, `brt`, `cet`, `cst`, `ct`, `et`, `gmt`, `hkt`, `ist`, `jst`, `mt`, `pt`, `sgt`, `utc`, etc.
- Timezone packages provide helper functions: `et.Now()`, `pt.Date(...)`, etc.
- Type aliases enable clean signatures: `utc.Time`, `et.Time`, `pt.Time`
- Package name conveys timezone, type is always `Timezone`

### 3. Explicit Conversions
- Timezone conversions must be explicit: `est.FromMoment(pacificTime)`
- Conversions use the `Moment` interface, supporting both `time.Time` and `meridian.Time[TZ]`
- This makes timezone handling visible in code review
- No silent timezone changes or data loss
- All conversions preserve the moment in time (UTC equality)

### 4. Internal UTC Storage
- All times are stored as UTC internally (`utcTime time.Time`)
- Timezone is applied during operations (display, hour extraction, etc.)
- Database-friendly and eliminates ambiguity

## Project Structure

```
.
├── meridian.go              # Core package: Time[TZ] type and core functions
├── meridian_test.go         # Core package tests
├── example_test.go          # Testable examples (appear in godoc)
├── doc.go                   # Package-level documentation
├── cmd/example/main.go      # Example usage program
├── est/                     # Eastern Time timezone package
│   ├── est.go
│   └── est_test.go
├── pst/                     # Pacific Time timezone package
│   ├── pst.go
│   └── pst_test.go
├── utc/                     # UTC timezone package
│   ├── utc.go
│   └── utc_test.go
├── .golangci.yml            # Linter configuration
└── Makefile                 # Development commands
```

## Development Commands

### Quick Reference
```bash
make test           # Run tests with race detection
make test-coverage  # Generate coverage report (coverage.html)
make lint           # Run golangci-lint
make run-example    # Run the example program
make clean          # Remove build artifacts
make install-tools  # Install development dependencies
```

### Full Commands
```bash
# Testing
go test -v -race ./...                                    # Run all tests
go test -v -race -coverprofile=coverage.out ./...        # With coverage
go tool cover -html=coverage.out                         # View coverage

# Linting
golangci-lint run                                         # Run all linters

# Running examples
go run cmd/example/main.go                                # Run example
```

## Testing Requirements

### Must-Have for All Changes
1. **Tests must pass**: `make test` must succeed (includes race detection)
2. **Coverage should not decrease**: Run `make test-coverage` to verify
3. **All linters must pass**: `make lint` must show no errors
4. **Examples must compile**: Testable examples in `example_test.go` must be valid

### Testing Philosophy
- Write table-driven tests for multiple scenarios
- Test edge cases: DST transitions, leap seconds, year boundaries
- Use `t.Parallel()` for independent tests
- Test both the generic API (`meridian.Time[TZ]`) and timezone-specific helpers

### Race Detection
Always run tests with `-race` flag. This is a library dealing with time, which often involves concurrency. The CI pipeline enforces this.

## Code Style Guidelines

### Enforced by Linters
- **gofmt**: Standard Go formatting
- **goimports**: Proper import ordering and grouping
- **revive**: Naming conventions, exported types must have comments
- **godot**: Comments must end with punctuation
- **errcheck**: All errors must be handled
- **gosec**: Security best practices
- **gocyclo**: Max complexity 15

### Type-Specific Conventions

#### Generic Functions
```go
// Good: Generic function with timezone type parameter
func Now[TZ Timezone]() Time[TZ] {
    return Time[TZ]{utcTime: time.Now().UTC()}
}

// Good: Preserving timezone types through operations
func Add[TZ Timezone](t Time[TZ], d time.Duration) Time[TZ] {
    return Time[TZ]{utcTime: t.utcTime.Add(d)}
}
```

#### Timezone Package Functions
```go
// Good: Timezone-specific helper (in est/est.go)
func Now() meridian.Time[EST] {
    return meridian.Now[EST]()
}

// Good: Timezone-specific constructor
func Date(year int, month time.Month, day, hour, min, sec, nsec int) meridian.Time[EST] {
    return meridian.Date[EST](year, month, day, hour, min, sec, nsec)
}
```

#### Documentation
- All exported types, functions, and methods **must** have godoc comments
- Start comments with the name of the item: `// Time represents a moment in time...`
- Explain **why** timezone type safety matters in the package docs
- Include examples in `example_test.go` for common use cases

### Import Grouping
```go
import (
    // Standard library
    "fmt"
    "time"
    
    // External dependencies (none currently)
    
    // Internal packages
    "github.com/matthalp/go-meridian"
    "github.com/matthalp/go-meridian/est"
)
```

## Common Operations & Patterns

### Creating Times
```go
// From current time
now := est.Now()

// From specific date/time
meeting := pst.Date(2024, 3, 15, 9, 0, 0, 0)

// From formatted string (parsed in the timezone's location)
parsed, _ := est.Parse(time.RFC3339, "2024-03-15T09:00:00-04:00")

// From Unix timestamps
unixTime := utc.Unix(1705320000, 0)           // seconds + nanoseconds
milliTime := pst.UnixMilli(1705320000000)     // milliseconds
microTime := est.UnixMicro(1705320000000000)  // microseconds

// From standard time.Time
stdTime := time.Now()
typed := utc.FromMoment(stdTime)  // Convert using Moment interface
```

**Note**: `ParseInLocation` from the `time` package is not needed in timezone packages because the location is already determined by the package (e.g., `utc.Parse` always parses in UTC, `est.Parse` in EST, `pst.Parse` in PST).

### Converting Between Timezones
```go
// Using timezone package FromMoment functions
eastern := est.Now()
pacific := pst.FromMoment(eastern)  // Explicit conversion
utcTime := utc.FromMoment(eastern)  // To UTC for storage

// Convert from time.Time
stdTime := time.Now()
utcTyped := utc.FromMoment(stdTime)
estTyped := est.FromMoment(stdTime)

// All conversions preserve the moment in time
fmt.Println(eastern.UTC().Equal(pacific.UTC()))  // true
```

### The Moment Interface
```go
// Both time.Time and meridian.Time[TZ] implement Moment
type Moment interface {
    UTC() time.Time
}

// Functions can accept any Moment
func processTime(m Moment) {
    utcTime := m.UTC()  // Get underlying UTC time
    // ... process ...
}

// Works with both types
processTime(time.Now())
processTime(est.Now())
```

### Preserving Types Through Operations
```go
func GetBusinessDayEnd[TZ Timezone](t meridian.Time[TZ]) meridian.Time[TZ] {
    // Return type preserves the input timezone type
    return t.Add(8 * time.Hour)  // Assuming 8-hour business day
}
```

## API Stability & Versioning

### This is a Library
- Breaking changes require major version bumps (v1.0.0 → v2.0.0)
- New functionality should be backwards compatible (minor version bump)
- Bug fixes only affect patch version
- Follow [Semantic Versioning](https://semver.org/) strictly

### Before Adding New Features
1. Is it consistent with the type-safety philosophy?
2. Does it maintain explicit conversion requirements?
3. Does it work with the generic `Time[TZ]` type?
4. Is the API intuitive for library consumers?

### Publishing Checklist
1. Update `Version` constant in `meridian.go`
2. Update `CHANGELOG.md` with changes
3. Ensure all tests pass: `make test`
4. Ensure linting passes: `make lint`
5. Run coverage check: `make test-coverage`
6. Create git tag: `git tag v0.x.x`
7. Push tag: `git push origin v0.x.x`
8. Documentation automatically appears on pkg.go.dev

## Common Pitfalls to Avoid

### ❌ Don't: Lose timezone information
```go
// Bad: Stripping timezone info
func GetUTC[TZ Timezone](t Time[TZ]) time.Time {
    return t.utcTime  // Exposes internal UTC, loses type
}
```

### ✅ Do: Maintain type-safety
```go
// Good: Explicit conversion to UTC type
func GetUTC[TZ Timezone](t Time[TZ]) utc.Time {
    return utc.FromMoment(t)
}

// Or use the Moment interface for flexibility
func GetUTC(m Moment) utc.Time {
    return utc.FromMoment(m)
}
```

### ❌ Don't: Allow implicit timezone mixing
```go
// Bad: Accepting any time without timezone awareness
func ScheduleJob(t time.Time) { ... }
```

### ✅ Do: Require explicit timezone types
```go
// Good: Timezone is part of the contract
func ScheduleJob(t Time[utc.UTC]) { ... }

// Or generic if it works with any timezone
func ScheduleJob[TZ Timezone](t Time[TZ]) { ... }
```

### ❌ Don't: Panic on timezone errors
```go
// Bad: Library code should not panic
loc := time.LoadLocation("America/New_York")  // Can panic
```

### ✅ Do: Handle errors gracefully
```go
// Good: Return errors or use validated timezones
loc, err := time.LoadLocation("America/New_York")
if err != nil {
    return Time[TZ]{}, err
}
```

## Security Considerations

- No user input should directly create `time.Location` without validation
- Be cautious with timezone strings from external sources
- Always validate timezone names against IANA database
- DST transitions can cause time ambiguity—test these cases

## Performance Notes

- Timezone resolution happens at compile time (zero runtime overhead for types)
- Times are stored as UTC internally (single representation)
- Timezone conversion is just a `time.In()` call (standard library performance)
- Generic type parameters have no runtime cost in Go 1.18+

## CI/CD Pipeline

The GitHub Actions workflow enforces:
1. **Tests pass** on Go 1.20+ with race detection
2. **Linting passes** with golangci-lint
3. **go.mod is tidy**: `go mod tidy` produces no changes
4. **Generated code is in sync**: Timezone packages match `timezones.yaml`
5. **Coverage tracking**: Reports uploaded to Codecov

All PRs must pass these checks. If CI fails, the PR cannot merge.

**Note on generated code**: If you modify `timezones.yaml`, you must run `make generate` and commit the generated files. The CI will fail if generated packages are out of sync with the configuration.

## IDE Setup

### Recommended Tools
- **gopls**: Go language server (built into most IDEs)
- **golangci-lint**: Install with `make install-tools`
- **go test**: Built-in to Go toolchain

### VS Code Settings
```json
{
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "workspace",
  "go.testFlags": ["-race"],
  "go.coverOnSave": true
}
```

## When Adding New Timezones

Timezone packages are **automatically generated** from the `timezones.yaml` configuration file. To add a new timezone:

### Step 1: Add to timezones.yaml

Edit `timezones.yaml` in the project root and add your timezone:

```yaml
timezones:
  - name: jst
    location: Asia/Tokyo
    description: Japan Standard Time
```

**Configuration fields:**
- `name`: Package name (lowercase, will be the directory name: `jst/`)
- `location`: IANA timezone name (e.g., `Asia/Tokyo`, `America/Chicago`)
- `description`: Human-readable timezone name for documentation

### Step 2: Generate the package

```bash
make generate
```

This will create:
- `jst/jst.go` - Package implementation with all timezone methods
- `jst/jst_test.go` - Complete test suite

### Step 3: Verify

```bash
make test       # Ensure all tests pass
make lint       # Ensure code quality standards are met
```

**Important**: The CI pipeline will verify that generated code is in sync with `timezones.yaml`. If you forget to run `make generate` after modifying the YAML file, the CI build will fail with a helpful error message showing which files are out of sync.

### Generated Code Structure

Each generated timezone package includes:

**Package file (`{name}/{name}.go`):**
- `type Timezone struct{}` - Timezone type
- `type Time = meridian.Time[Timezone]` - Convenience alias
- `func Location() *time.Location` - Returns IANA location
- `func Now() Time` - Current time in this timezone
- `func Date(...) Time` - Create time from components
- `func FromMoment(Moment) Time` - Convert from any timezone
- `func Parse(layout, value) (Time, error)` - Parse time string
- `func Unix(sec, nsec) Time` - From Unix timestamp
- `func UnixMilli(msec) Time` - From Unix milliseconds
- `func UnixMicro(usec) Time` - From Unix microseconds

**Test file (`{name}/{name}_test.go`):**
- Location verification test
- Now() functionality test
- Date() with timezone offset tests
- Conversion tests (from time.Time, from other timezones, round-trip)
- Parse tests (format handling, timezone interpretation)
- Unix timestamp tests (all precisions)

### Special Cases

**UTC timezone** is handled automatically:
- Uses `time.UTC` directly instead of `time.LoadLocation`
- No `mustLoadLocation` helper needed
- Simplified Parse implementation using `time.Parse`

### Generator Implementation

The code generator lives at `cmd/generate-timezones/main.go`:
- Reads `timezones.yaml` using `gopkg.in/yaml.v3`
- Uses Go templates for package and test generation
- Formats output with `goimports` for proper import ordering
- Handles conditional logic for UTC and import cycle prevention

## Questions to Ask Before Committing

1. ✅ Do all tests pass with race detection?
2. ✅ Does golangci-lint pass without warnings?
3. ✅ Are all exported items documented?
4. ✅ Does this maintain type-safety guarantees?
5. ✅ Is this a breaking change? (If yes, needs major version bump)
6. ✅ Are there tests covering the new functionality?
7. ✅ Does the example code still compile?
8. ✅ Is `go.mod` still tidy?

## Additional Resources

- **Go Time Package**: https://pkg.go.dev/time
- **Go Generics**: https://go.dev/doc/tutorial/generics
- **IANA Time Zone Database**: https://www.iana.org/time-zones
- **Semantic Versioning**: https://semver.org/

