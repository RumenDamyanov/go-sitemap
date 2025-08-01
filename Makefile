.PHONY: test test-verbose test-coverage clean build lint fmt examples

# Default target
all: test build

# Run tests
test:
	go test ./...

# Run tests with verbose output
test-verbose:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -cover ./...

# Generate detailed coverage report
coverage-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean generated files
clean:
	rm -f coverage.out coverage.html
	go clean ./...

# Build examples
build-examples:
	@echo "Building examples..."
	cd examples/nethttp && go build .
	@echo "Built examples successfully"

# Format code
fmt:
	go fmt ./...

# Run linter (if golangci-lint is installed)
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, skipping..."; \
	fi

# Run go vet
vet:
	go vet ./...

# Download dependencies
deps:
	go mod download
	go mod tidy

# Run all checks
check: fmt vet lint test

# Help
help:
	@echo "Available targets:"
	@echo "  test            - Run tests"
	@echo "  test-verbose    - Run tests with verbose output"
	@echo "  test-coverage   - Run tests with coverage"
	@echo "  coverage-html   - Generate HTML coverage report"
	@echo "  build-examples  - Build example projects"
	@echo "  fmt             - Format code"
	@echo "  lint            - Run linter"
	@echo "  vet             - Run go vet"
	@echo "  deps            - Download and tidy dependencies"
	@echo "  check           - Run all checks (fmt, vet, lint, test)"
	@echo "  clean           - Clean generated files"
	@echo "  help            - Show this help"
