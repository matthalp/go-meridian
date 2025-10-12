.PHONY: help test test-coverage lint build clean run-example install-tools

# Default target
help:
	@echo "Available targets:"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make lint           - Run linter"
	@echo "  make build          - Build the example binary"
	@echo "  make run-example    - Run the example program"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make install-tools  - Install development tools"

# Run tests
test:
	go test -v -race ./...

# Run tests with coverage
test-coverage:
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linter
lint:
	golangci-lint run

# Build the example binary
build:
	go build -o bin/example ./cmd/example

# Run the example program
run-example:
	go run ./cmd/example

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean

# Install development tools
install-tools:
	@echo "Installing golangci-lint..."
	@which golangci-lint > /dev/null || \
		(echo "Installing golangci-lint..." && \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2)
	@echo "Tools installed successfully!"

