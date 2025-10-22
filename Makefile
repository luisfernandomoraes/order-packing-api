.PHONY: help run build test test-verbose test-coverage lint lint-install clean dev swagger build-container run-container

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
	@echo "  make lint          - Run linter (golangci-lint)"
	@echo "  make lint-install  - Install golangci-lint"
	@echo "  make fmt           - Format code"
	@echo "  make swagger       - Generate Swagger documentation"
	@echo "  make swagger-install - Install swag CLI tool"
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
fmt:
	@echo "Formatting code..."
	@go fmt ./...


# Install swag CLI tool
swagger-install:
	@echo "Installing swag CLI tool..."
	@which swag > /dev/null || go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@which swag > /dev/null || go install github.com/swaggo/swag/cmd/swag@latest
	@swag init -g cmd/api/main.go -o docs
	@echo "Swagger documentation generated at docs/"

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

# Build Docker image
build-container:
	@echo "Building Docker image $(IMAGE_NAME):$(IMAGE_TAG)..."
	@docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .
	@echo "Image build complete."

# Run Docker container
run-container:
	@echo "Running Docker container on port $(PORT)..."
	@docker run --rm -p $(PORT):8080 --name $(IMAGE_NAME) $(IMAGE_NAME):$(IMAGE_TAG)
