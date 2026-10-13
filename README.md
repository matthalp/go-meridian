# go-meridian

[![CI](https://github.com/matthalp/go-meridian/actions/workflows/ci.yml/badge.svg)](https://github.com/matthalp/go-meridian/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/matthalp/go-meridian/branch/main/graph/badge.svg)](https://codecov.io/gh/matthalp/go-meridian)
[![Go Report Card](https://goreportcard.com/badge/github.com/matthalp/go-meridian)](https://goreportcard.com/report/github.com/matthalp/go-meridian)

**Type-safe timezone handling for Go using generics.**

Meridian solves a fundamental problem: timezone information in `time.Time` is data, not type, and can be lost without the compiler noticing. With Meridian, timezone information is encoded directly into the type system, making wrong timezone handling impossible to compile.

## Features

- ✅ **Type-safe timezones**: `utc.Time` and `et.Time` are different types
- ✅ **Compiler-enforced correctness**: Prevents accidental timezone mixing
- ✅ **Clean, ergonomic API**: `utc.Now()`, `et.Date(...)`, `pt.Time`
- ✅ **16 built-in timezone packages**: Covers major global business centers
- ✅ **Extensible**: Easy to add custom timezone packages
- ✅ Full GitHub Actions CI/CD pipeline
- ✅ Automated testing with coverage reports
- ✅ Race condition detection

## Installation

Install the package in your Go project:

```bash
go get github.com/matthalp/go-meridian/v2
```

Timezone packages load their location from the host's timezone database when
first imported, and panic if it is missing. On hosts without one (for example
`scratch` or distroless container images), embed Go's copy by adding this to
your `main` package:

```go
import _ "time/tzdata"
```

## Quick Start

Import and use timezone-specific packages:

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/matthalp/go-meridian/v2/timezones/et"
    "github.com/matthalp/go-meridian/v2/timezones/utc"
)

func main() {
    // Get current time in different timezones
    now := utc.Now()
    fmt.Println(now.Format(time.RFC3339))
    
    // Create a specific date/time
    meeting := et.Date(2024, time.December, 25, 10, 30, 0, 0)
    fmt.Println(meeting.Format(time.Kitchen))
    
    // Type-safe function signatures
    storeInDatabase(utc.Now())  // ✅ Compiles
    // storeInDatabase(et.Now()) // ❌ Won't compile!
}

// Functions can require specific timezones
func storeInDatabase(t utc.Time) {
    // Always receives UTC time, guaranteed by the compiler
}
```

## Why Meridian?

**Problem**: Standard Go `time.Time` loses timezone information easily:
```go
t := time.Now().UTC()
// Later in code...
formatted := t.Format(time.Kitchen) // What timezone is this? 🤷
```

**Solution**: Meridian encodes timezone in the type:
```go
t := utc.Now()
// Later in code...
formatted := t.Format(time.Kitchen) // Definitely UTC! ✅
```

## Available Timezone Packages

All timezone packages live under `timezones/`:

- `github.com/matthalp/go-meridian/v2/timezones/aest` - Australian Eastern Time (Australia/Sydney)
- `github.com/matthalp/go-meridian/v2/timezones/brt` - Brasília Time (America/Sao_Paulo)
- `github.com/matthalp/go-meridian/v2/timezones/cet` - Central European Time (Europe/Paris)
- `github.com/matthalp/go-meridian/v2/timezones/cst` - China Standard Time (Asia/Shanghai)
- `github.com/matthalp/go-meridian/v2/timezones/ct` - Central Time (America/Chicago)
- `github.com/matthalp/go-meridian/v2/timezones/est` - **Deprecated**, use `et` (same zone; follows daylight saving time despite the name). Removed in v3.
- `github.com/matthalp/go-meridian/v2/timezones/et` - Eastern Time (America/New_York)
- `github.com/matthalp/go-meridian/v2/timezones/gmt` - UK time (Europe/London): GMT in winter, BST (UTC+1) in summer. For a fixed UTC+0 offset, use `utc`
- `github.com/matthalp/go-meridian/v2/timezones/hkt` - Hong Kong Time (Asia/Hong_Kong)
- `github.com/matthalp/go-meridian/v2/timezones/ist` - India Standard Time (Asia/Kolkata)
- `github.com/matthalp/go-meridian/v2/timezones/jst` - Japan Standard Time (Asia/Tokyo)
- `github.com/matthalp/go-meridian/v2/timezones/mt` - Mountain Time (America/Denver)
- `github.com/matthalp/go-meridian/v2/timezones/pst` - **Deprecated**, use `pt` (same zone; follows daylight saving time despite the name). Removed in v3.
- `github.com/matthalp/go-meridian/v2/timezones/pt` - Pacific Time (America/Los_Angeles)
- `github.com/matthalp/go-meridian/v2/timezones/sgt` - Singapore Time (Asia/Singapore)
- `github.com/matthalp/go-meridian/v2/timezones/utc` - Coordinated Universal Time

Packages are generated from `timezones.yaml`; see [Adding Custom Timezones](#adding-custom-timezones) below.

### Package API

Each timezone package provides:
- `Now()` - Get current time in that timezone
- `Date()` - Create a specific date/time
- `Parse()` - Parse a formatted string in that timezone
- `Unix()`, `UnixMilli()`, `UnixMicro()` - Create from Unix timestamps
- `FromMoment()` - Convert any time to that timezone
- `Time` - Type alias for clean function signatures

Note: `ParseInLocation` is not needed as timezone packages already have their location built-in.

## Converting Between Timezones

Meridian provides seamless timezone conversion while preserving type safety:

```go
// Convert between timezone types
etTime := et.Date(2024, time.December, 25, 10, 30, 0, 0)
utcTime := utc.FromMoment(etTime)  // Same moment, displayed as UTC
ptTime := pt.FromMoment(etTime)  // Same moment, displayed as PT

// Convert from standard time.Time
stdTime := time.Now()
typedTime := utc.FromMoment(stdTime)  // Now type-safe!

// All conversions preserve the moment in time
fmt.Println(etTime.UTC().Equal(utcTime.UTC()))  // true
```

The `Moment` interface allows both `time.Time` and `meridian.Time[TZ]` to be used interchangeably for conversions, providing flexibility while maintaining type safety where it matters.

## Adding Custom Timezones

Timezone packages are generated from the `timezones.yaml` configuration file. To add a new timezone:

1. **Edit `timezones.yaml`** and add your timezone definition:
   ```yaml
   timezones:
     - name: kst
       location: Asia/Seoul
       description: Korea Standard Time
   ```

2. **Generate the package**:
   ```bash
   make generate
   ```

3. **Import and use**:
   ```go
   import "github.com/matthalp/go-meridian/v2/timezones/kst"

   now := kst.Now()
   ```

The generator creates both the package implementation and comprehensive tests automatically. For more details, see `AGENTS.md`.

### Running the Example

This repository includes an example program demonstrating usage:

```bash
go run cmd/example/main.go
```

## Development

### Running Tests

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out
```

### Running Linter

```bash
# Install golangci-lint (if not already installed)
# https://golangci-lint.run/usage/install/

# Run linter
golangci-lint run
```

## CI/CD Pipeline

The project includes a comprehensive GitHub Actions workflow that:

1. **Test Job**: Runs unit tests with race detection and generates coverage reports
2. **Lint Job**: Runs golangci-lint to ensure code quality
3. **Build Job**: Verifies the project builds successfully and go.mod is tidy

Coverage reports are automatically uploaded to Codecov for tracking test coverage over time.

## Project Structure

```
.
├── .github/
│   └── workflows/
│       └── ci.yml          # GitHub Actions workflow
├── cmd/
│   ├── example/
│   │   └── main.go         # Example usage program
│   └── generate-timezones/
│       └── main.go         # Timezone package generator
├── timezones/              # Generated timezone packages (one per timezone)
│   ├── et/
│   │   ├── et.go
│   │   └── et_test.go
│   ├── pt/
│   ├── utc/
│   └── ...
├── .golangci.yml           # Linter configuration
├── doc.go                  # Package documentation
├── example_test.go         # Testable examples
├── go.mod                  # Go module file
├── meridian.go             # Core generic types and functions
├── meridian_test.go        # Core package tests
├── timezones.yaml          # Timezone definitions for generator
├── Makefile                # Development tasks
└── README.md               # This file
```

## API Documentation

For detailed API documentation, see [pkg.go.dev](https://pkg.go.dev/github.com/matthalp/go-meridian/v2).

## Publishing a New Version

To make your package available for others to use:

1. **Push to GitHub**: 
   ```bash
   git add .
   git commit -m "Initial commit"
   git remote add origin https://github.com/matthalp/go-meridian.git
   git push -u origin main
   ```

2. **Create a version tag**:
   ```bash
   git tag v2.0.0
   git push origin v2.0.0
   ```

3. **The package will be automatically available** via `go get`:
   - Others can install with: `go get github.com/matthalp/go-meridian/v2@v2.0.0`
   - Documentation will appear on [pkg.go.dev](https://pkg.go.dev/github.com/matthalp/go-meridian/v2) within minutes

4. **Update the version** in `meridian.go` when releasing new versions

### Versioning

This project follows [Semantic Versioning](https://semver.org/):
- **MAJOR** version (v1.0.0) for incompatible API changes
- **MINOR** version (v0.1.0) for new functionality in a backwards compatible manner
- **PATCH** version (v0.0.1) for backwards compatible bug fixes

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

All PRs must pass CI checks before merging.

## License

This project is licensed under the MIT License.

