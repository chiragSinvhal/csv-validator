# CSV Validator Makefile

.PHONY: help build run test clean docker-build docker-run deps fmt lint vet check coverage

# Default target
help:
	@echo "Available commands:"
	@echo "  build            - Build the application"
	@echo "  run              - Run the application"
	@echo "  test             - Run tests"
	@echo "  coverage         - Run tests with coverage"
	@echo "  deps             - Download dependencies"
	@echo "  fmt              - Format code"
	@echo "  lint             - Run linter"
	@echo "  vet              - Run go vet"
	@echo "  check            - Run all checks"
	@echo "  clean            - Clean build artifacts"
	@echo "  docker-build     - Build Docker image"
	@echo "  docker-run       - Run Docker container"

# Build the application
build:
	@echo "Building csv-validator..."
	go build -o bin/csv-validator .

# Run the application
run:
	@echo "Running csv-validator..."
	go run .

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests (verbose)..."
	go test -v ./...

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	go test -cover ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed"; \
	fi

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Run all checks
check: fmt vet test
	@echo "All checks passed!"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -rf uploads/ downloads/

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t csv-validator .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 csv-validator

# Development setup
setup:
	@echo "Setting up development environment..."
	go mod download
	cp .env.example .env
	mkdir -p uploads downloads
	@echo "Development environment ready!"
	@echo "Edit .env file with your configuration before running."