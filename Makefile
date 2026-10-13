.PHONY: help test test-coverage lint build clean run-example install-tools generate release-check

# Keep in sync with the golangci-lint version in .github/workflows/ci.yml.
GOLANGCI_LINT_VERSION := v2.14.0

# Default target
help:
	@echo "Available targets:"
	@echo "  make generate       - Generate timezone packages from timezones.yaml"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make lint           - Run linter"
	@echo "  make build          - Build the example binary"
	@echo "  make run-example    - Run the example program"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make install-tools  - Install development tools"
	@echo "  make release-check VERSION=vX.Y.Z - Check the tree is ready to release as VERSION"

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

# Generate timezone packages from timezones.yaml
generate:
	go run ./cmd/generate-timezones

# Keep in sync with GORELEASE_VERSION in .github/workflows/release.yml.
GORELEASE_VERSION := v0.0.0-20261011100332-faecfb24eb49

# Check that the tree is ready to release as $(VERSION), as the Release workflow does.
# GOPROXY=direct keeps gorelease from caching "not found" for VERSION on
# proxy.golang.org; see the Check release step in .github/workflows/release.yml.
release-check:
	@test -n "$(VERSION)" || (echo "usage: make release-check VERSION=vX.Y.Z" && exit 1)
	go install golang.org/x/exp/cmd/gorelease@$(GORELEASE_VERSION)
	GOPROXY=direct python3 .github/scripts/check_release.py $(VERSION) --gorelease "$$(go env GOPATH)/bin/gorelease"

# Install development tools
install-tools:
	@echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION)..."
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	@echo "Tools installed successfully!"
