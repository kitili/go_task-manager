# Go Task Manager Makefile

# Variables
BINARY_NAME=task-manager
API_BINARY_NAME=api
TEST_COVERAGE_FILE=coverage.out
TEST_COVERAGE_HTML=coverage.html

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet
GOLINT=golangci-lint

# Build flags
BUILD_FLAGS=-ldflags "-s -w"
TEST_FLAGS=-v -race -coverprofile=$(TEST_COVERAGE_FILE)
BENCHMARK_FLAGS=-bench=. -benchmem

# Default target
.PHONY: all
all: clean build test

# Build the main application
.PHONY: build
build:
	@echo "Building Go Task Manager..."
	$(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_NAME) .
	@echo "Building API server..."
	$(GOBUILD) $(BUILD_FLAGS) -o $(API_BINARY_NAME) ./cmd/api
	@echo "Build complete!"

# Build for production (optimized)
.PHONY: build-prod
build-prod:
	@echo "Building for production..."
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_NAME) .
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(API_BINARY_NAME) ./cmd/api
	@echo "Production build complete!"

# Run the application
.PHONY: run
run: build
	@echo "Running Go Task Manager..."
	./$(BINARY_NAME)

# Run the API server
.PHONY: run-api
run-api: build
	@echo "Running API server..."
	./$(API_BINARY_NAME)

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) $(TEST_FLAGS) ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage: test
	@echo "Generating coverage report..."
	$(GOTEST) -coverprofile=$(TEST_COVERAGE_FILE) ./...
	$(GOCMD) tool cover -html=$(TEST_COVERAGE_FILE) -o $(TEST_COVERAGE_HTML)
	@echo "Coverage report generated: $(TEST_COVERAGE_HTML)"

# Run unit tests only
.PHONY: test-unit
test-unit:
	@echo "Running unit tests..."
	$(GOTEST) $(TEST_FLAGS) -short ./...

# Run integration tests only
.PHONY: test-integration
test-integration:
	@echo "Running integration tests..."
	$(GOTEST) $(TEST_FLAGS) -run "Integration" ./...

# Run API tests only
.PHONY: test-api
test-api:
	@echo "Running API tests..."
	$(GOTEST) $(TEST_FLAGS) -run "TestAPI" ./internal/api/...

# Run benchmarks
.PHONY: benchmark
benchmark:
	@echo "Running benchmarks..."
	$(GOTEST) $(BENCHMARK_FLAGS) ./...

# Run performance tests
.PHONY: test-performance
test-performance:
	@echo "Running performance tests..."
	$(GOTEST) $(TEST_FLAGS) -run "TestMemoryUsage\|TestConcurrentAccess" ./...

# Run all tests including performance
.PHONY: test-all
test-all: test test-performance benchmark

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(API_BINARY_NAME)
	rm -f $(TEST_COVERAGE_FILE)
	rm -f $(TEST_COVERAGE_HTML)
	@echo "Clean complete!"

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

# Run linter
.PHONY: lint
lint:
	@echo "Running linter..."
	$(GOLINT) run ./...

# Run vet
.PHONY: vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...

# Download dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Install dependencies
.PHONY: install
install: deps
	@echo "Installing dependencies..."
	$(GOGET) -u ./...

# Generate Swagger documentation
.PHONY: swagger
swagger:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/api/main.go -o docs --exclude web

# Run the application in development mode
.PHONY: dev
dev: build
	@echo "Starting development server..."
	./$(BINARY_NAME) --mode=dev

# Run the API server in development mode
.PHONY: dev-api
dev-api: build
	@echo "Starting development API server..."
	./$(API_BINARY_NAME) --mode=dev

# Docker build
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t go-task-manager .

# Docker run
.PHONY: docker-run
docker-run: docker-build
	@echo "Running Docker container..."
	docker run -p 8080:8080 go-task-manager

# Database setup
.PHONY: db-setup
db-setup:
	@echo "Setting up database..."
	$(GOCMD) run . --setup-db

# Database migration
.PHONY: db-migrate
db-migrate:
	@echo "Running database migrations..."
	$(GOCMD) run . --migrate

# Database reset
.PHONY: db-reset
db-reset:
	@echo "Resetting database..."
	rm -f task_manager.db
	$(GOCMD) run . --setup-db

# Security scan
.PHONY: security
security:
	@echo "Running security scan..."
	gosec ./...

# Generate mocks
.PHONY: mocks
mocks:
	@echo "Generating mocks..."
	mockgen -source=internal/task/interface.go -destination=internal/task/mocks/mock_task_manager.go
	mockgen -source=internal/database/repository.go -destination=internal/database/mocks/mock_repository.go

# Pre-commit checks
.PHONY: pre-commit
pre-commit: fmt vet lint test
	@echo "Pre-commit checks complete!"

# CI/CD pipeline
.PHONY: ci
ci: deps fmt vet lint test test-coverage
	@echo "CI pipeline complete!"

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  build-prod     - Build for production"
	@echo "  run            - Run the application"
	@echo "  run-api        - Run the API server"
	@echo "  test           - Run all tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-unit      - Run unit tests only"
	@echo "  test-integration - Run integration tests only"
	@echo "  test-api       - Run API tests only"
	@echo "  benchmark      - Run benchmarks"
	@echo "  test-performance - Run performance tests"
	@echo "  test-all       - Run all tests including performance"
	@echo "  clean          - Clean build artifacts"
	@echo "  fmt            - Format code"
	@echo "  lint           - Run linter"
	@echo "  vet            - Run go vet"
	@echo "  deps           - Download dependencies"
	@echo "  install        - Install dependencies"
	@echo "  swagger        - Generate Swagger documentation"
	@echo "  dev            - Run in development mode"
	@echo "  dev-api        - Run API in development mode"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  db-setup       - Setup database"
	@echo "  db-migrate     - Run database migrations"
	@echo "  db-reset       - Reset database"
	@echo "  security       - Run security scan"
	@echo "  mocks          - Generate mocks"
	@echo "  pre-commit     - Run pre-commit checks"
	@echo "  ci             - Run CI pipeline"
	@echo "  help           - Show this help message"
