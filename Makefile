# Makefile for home-ip-updater service
# Provides targets for building, testing, linting, and code quality checks

# Project configuration
PROJECT_NAME := "home-ip-updater"
PKG := "github.com/a-castellano/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

# Available targets
.PHONY: all build clean test test_integration test_messagebroker coverage coverhtml lint race msan

# Default target - builds the project
all: build

# Lint the Go source files for common issues and style violations
lint: ## Lint the GO_FILES
	@go vet ./...

# Run unit tests with short timeout
test: ## Run unit tests
	@go test --tags=unit_tests -short ./...

# Run integration tests with short timeout
test_integration: ## Run integration tests
	@go test --tags=integration_tests -short ./...

# Run tests with race condition detection
race: ## Run data race detector
	@go test -race -short ./...

# Run tests with memory sanitizer for memory safety issues
msan: ## Run memory sanitizer
	@go test -msan -short ./...

# Generate code coverage report
coverage: ## Generate global code coverage report
	./development/coverage.sh;

# Render the report produced by coverage as HTML, without running the tests again
coverhtml: ## Generate global code coverage report in HTML
	go tool cover -html=cover/coverage.report -o coverage.html;

# Build the binary executable
build: ## Build the binary file
	@go build -v $(PKG)

# Remove previous build artifacts
clean: ## Remove previous build
	@rm -f $(PROJECT_NAME)

# Display help information for all available targets
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
