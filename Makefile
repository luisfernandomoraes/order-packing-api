.PHONY: help run build test test-verbose test-coverage lint lint-install clean dev

# Default target
help:
	@echo "Available targets:"
	@echo "  make run           - Run the application"
	@echo "  make build         - Build the application binary"
	@echo "  make test          - Run all tests"
	@echo "  make test-verbose  - Run tests with verbose output"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make bench         - Run benchmark tests"
	@echo "  make lint          - Run linter (golangci-lint)"
	@echo "  make lint-install  - Install golangci-lint"
	@echo "  make fmt           - Format code"
	@echo "  make clean         - Remove build artifacts"

# Run the application
run:
	@echo "Starting server on port 8080..."
	@go run cmd/api/main.go

# Build the application
build:
	@echo "Building application..."
	@go build -o bin/order-packing-api cmd/api/main.go
	@echo "Build complete: bin/order-packing-api"

# Run all tests
test:
	@echo "Running tests..."
	@go test ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests (verbose)..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -cover ./...
	@echo ""
	@echo "Generating HTML coverage report..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmark tests
bench:
	@echo "Running benchmark tests..."
	@go test -bench=. -benchmem ./internal/domain

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Install golangci-lint
lint-install:
	@echo "Installing golangci-lint..."
	@which golangci-lint > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.55.2
	@echo "golangci-lint installed successfully"

# Run linter
lint:
	@if command -v golangci-lint > /dev/null; then \
		echo "Running golangci-lint..."; \
		golangci-lint run ./...; \
	else \
		echo "Error: golangci-lint not found."; \
		echo "Install with: make lint-install"; \
		exit 1; \
	fi

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

# Install dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed"

# Run all checks (fmt, lint, test)
check: fmt lint test
	@echo "All checks passed!"
