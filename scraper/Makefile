.PHONY: build test integration-test clean setup-integration teardown-integration clean-containers

# Variables
APP_NAME := macrochain-scraper
BUILD_DIR := ./build
DOCKER_PREFIX := macrochain-scraper-test
POSTGRES_PORT := 5433
REDIS_PORT := 6380

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) .
	@echo "Build complete!"

# Run unit tests (excluding integration tests)
test:
	@echo "Running unit tests..."
	@go test -v ./...
	@echo "Unit tests complete!"

# Clean up any existing containers
clean-containers:
	@echo "Cleaning up any existing test containers..."
	@docker compose -f docker-compose.test.yml down -v 2>/dev/null || true

# Set up integration test environment
setup-integration: clean-containers
	@echo "Setting up integration test environment..."
	@POSTGRES_PORT=$(POSTGRES_PORT) REDIS_PORT=$(REDIS_PORT) docker compose -f docker-compose.test.yml up -d
	@echo "Waiting for containers to be ready..."
	@sleep 5

# Tear down integration test environment
teardown-integration:
	@echo "Tearing down integration test environment..."
	@docker compose -f docker-compose.test.yml down -v

# Run integration tests with proper error handling
integration-test: setup-integration
	@echo "Running integration tests..."
	@(set -o pipefail; \
	REDIS_HOST=localhost \
	REDIS_PORT=$(REDIS_PORT) \
	DB_HOST=localhost \
	DB_PORT=$(POSTGRES_PORT) \
	DB_USER=postgres \
	DB_PASSWORD=postgres \
	DB_NAME=macrochain_test \
	go test -v -tags=integration ./... || { \
		echo "Integration tests failed!"; \
		$(MAKE) teardown-integration; \
		exit 1; \
	})
	@echo "Integration tests complete!"
	@$(MAKE) teardown-integration

# Run integration tests with race detection
integration-test-race: setup-integration
	@echo "Running integration tests with race detection..."
	@(set -o pipefail; \
	REDIS_HOST=localhost \
	REDIS_PORT=$(REDIS_PORT) \
	DB_HOST=localhost \
	DB_PORT=$(POSTGRES_PORT) \
	DB_USER=postgres \
	DB_PASSWORD=postgres \
	DB_NAME=macrochain_test \
	go test -v -race -tags=integration ./... || { \
		echo "Integration tests with race detection failed!"; \
		$(MAKE) teardown-integration; \
		exit 1; \
	})
	@echo "Integration tests with race detection complete!"
	@$(MAKE) teardown-integration

# Clean build artifacts
clean: clean-containers
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete!"

# Default target
all: build 