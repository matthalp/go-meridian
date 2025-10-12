# Usage Guide

## For Package Users

### Installing the Package

Once published to GitHub, users can install your package with:

```bash
go get github.com/matthalp/go-meridian
```

### Using in Their Code

```go
package main

import (
    "fmt"
    "github.com/matthalp/go-meridian"
)

func main() {
    // Use any exported function from the package
    greeting := meridian.Greet("Developer")
    fmt.Println(greeting)
    
    // Access exported constants
    fmt.Printf("Using meridian version: %s\n", meridian.Version)
}
```

## For Package Developers

### Project Structure

```
go-meridian/
├── meridian.go          # Main package code with exported functions
├── meridian_test.go     # Unit tests
├── example_test.go      # Testable examples (appear in docs)
├── doc.go               # Package-level documentation
├── cmd/example/         # Example program using the package
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

### Adding New Functionality

To add a new function that others can use:

1. Add the function to `meridian.go`:
   ```go
   // NewFunction does something useful.
   // Provide a clear description of what it does.
   func NewFunction(param string) string {
       // Implementation
       return result
   }
   ```

2. Add tests in `meridian_test.go`:
   ```go
   func TestNewFunction(t *testing.T) {
       result := NewFunction("test")
       if result != expected {
           t.Errorf("got %v, want %v", result, expected)
       }
   }
   ```

3. Add example in `example_test.go`:
   ```go
   func ExampleNewFunction() {
       result := meridian.NewFunction("input")
       fmt.Println(result)
       // Output: expected output
   }
   ```

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

