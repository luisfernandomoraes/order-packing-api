.PHONY: help run build test test-verbose test-coverage lint lint-install clean dev swagger build-container run-container tools deps

# Development tools versions
GOLANGCI_LINT_VERSION := v2.5.0
SWAG_VERSION := v1.16.3

IMAGE_NAME ?= order-packing-api
IMAGE_TAG ?= latest
PORT ?= 8080

# Default target
help:
	@echo "Available targets:"
	@echo "  make run           - Run the application"
	@echo "  make build         - Build the application binary"
	@echo "  make test          - Run all tests"
	@echo "  make test-verbose  - Run tests with verbose output"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make bench         - Run benchmark tests"
	@echo "  make deps          - Download Go dependencies"
	@echo "  make tools         - Install development tools (golangci-lint, swag)"
	@echo "  make fmt           - Format code (requires tools)"
	@echo "  make lint          - Run linter (requires tools)"
	@echo "  make swagger       - Generate Swagger documentation (requires tools)"
	@echo "  make check         - Run fmt, lint, and test"
	@echo "  make clean         - Remove build artifacts"
	@echo "  make build-container - Build Docker image $(IMAGE_NAME):$(IMAGE_TAG)"
	@echo "  make run-container   - Run Docker container mapping port $(PORT)->8080"

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
fmt: tools
	@echo "Formatting code with gofmt and gci..."
	@gofmt -s -w .
	@$$(go env GOPATH)/bin/golangci-lint fmt

# Install development tools
tools:
	@echo "🔍 Checking development tools..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "⬇️  Installing golangci-lint $(GOLANGCI_LINT_VERSION)..."; \
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | \
        sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION); \
	else \
		echo "✅ golangci-lint already installed."; \
	fi
	@if ! command -v swag >/dev/null 2>&1; then \
		echo "⬇️  Installing swag $(SWAG_VERSION)..."; \
		go install github.com/swaggo/swag/cmd/swag@$(SWAG_VERSION); \
	else \
		echo "✅ swag already installed."; \
	fi
	@echo "🚀 All development tools are ready!"

# Install swag CLI tool
swagger-install:
	@echo "Installing swag CLI tool..."
	@which swag > /dev/null || go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger documentation
swagger: tools
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/api/main.go -o docs
	@echo "Swagger documentation generated at docs/"

# Install golangci-lint
lint-install:
	@echo "Installing golangci-lint v2.5.0..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | \
		sh -s -- -b $$(go env GOPATH)/bin v2.5.0
	@echo "golangci-lint v2.5.0 installed successfully."

# Run linter
lint: tools
	@echo "Running golangci-lint..."
	@golangci-lint run

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

# Build Docker image
build-container:
	@echo "Building Docker image $(IMAGE_NAME):$(IMAGE_TAG)..."
	@docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .
	@echo "Image build complete."

# Run Docker container
run-container:
	@echo "Running Docker container on port $(PORT)..."
	@docker run --rm -p $(PORT):8080 --name $(IMAGE_NAME) $(IMAGE_NAME):$(IMAGE_TAG)
