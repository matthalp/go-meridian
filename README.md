# go-meridian

[![CI](https://github.com/matthalp/go-meridian/actions/workflows/ci.yml/badge.svg)](https://github.com/matthalp/go-meridian/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/matthalp/go-meridian/branch/main/graph/badge.svg)](https://codecov.io/gh/matthalp/go-meridian)
[![Go Report Card](https://goreportcard.com/badge/github.com/matthalp/go-meridian)](https://goreportcard.com/report/github.com/matthalp/go-meridian)

A Go 1.20 package/library with comprehensive CI/CD pipeline.

## Features

- ✅ Full GitHub Actions CI/CD pipeline
- ✅ Automated testing with coverage reports
- ✅ Code linting with golangci-lint
- ✅ Race condition detection
- ✅ Coverage reports uploaded to Codecov
- ✅ Importable as a Go module

## Installation

Install the package in your Go project:

```bash
go get github.com/matthalp/go-meridian
```

## Usage

Import and use the package in your code:

```go
package main

import (
    "fmt"
    "github.com/matthalp/go-meridian"
)

func main() {
    // Use the Greet function
    greeting := meridian.Greet("World")
    fmt.Println(greeting) // Output: Hello, World!

    // Check the version
    fmt.Printf("Version: %s\n", meridian.Version)
}
```

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
│   └── example/
│       └── main.go         # Example usage program
├── .golangci.yml           # Linter configuration
├── doc.go                  # Package documentation
├── example_test.go         # Testable examples
├── go.mod                  # Go module file
├── meridian.go             # Main package code
├── meridian_test.go        # Package tests
├── Makefile                # Development tasks
└── README.md               # This file
```

## API Documentation

For detailed API documentation, see [pkg.go.dev](https://pkg.go.dev/github.com/matthalp/go-meridian) once the package is published.

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
   git tag v0.1.0
   git push origin v0.1.0
   ```

3. **The package will be automatically available** via `go get`:
   - Others can install with: `go get github.com/matthalp/go-meridian@v0.1.0`
   - Documentation will appear on [pkg.go.dev](https://pkg.go.dev/github.com/matthalp/go-meridian) within minutes

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

