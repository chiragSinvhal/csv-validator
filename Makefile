# CSV Validator Makefile

.PHONY: help build run test clean docker-build docker-run deps fmt lint vet check coverage integration-test

# Default target
help:
	@echo "Available commands:"
	@echo "  build            - Build the application"
	@echo "  run              - Run the application"
	@echo "  test             - Run tests"
	@echo "  test-v           - Run tests with verbose output"
	@echo "  coverage         - Run tests with coverage"
	@echo "  integration-test - Run API integration tests"
	@echo "  deps             - Download dependencies"
	@echo "  fmt              - Format code"
	@echo "  lint             - Run linter"
	@echo "  vet              - Run go vet"
	@echo "  check            - Run all checks (fmt, vet, lint, test)"
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
test-v:
	@echo "Running tests (verbose)..."
	go test -v ./...

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	go test -v -cover ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run integration tests
integration-test:
	@echo "Running integration tests..."
	./scripts/integration-tests.sh

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
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
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
	rm -rf uploads/

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
	go mod tidy
	cp .env.example .env
	mkdir -p uploads
	@echo "Development environment ready!"
	@echo "Edit .env file with your configuration before running."

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Development tools installed!"