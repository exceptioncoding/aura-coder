# Aura - AI Coding Agent Build Configuration

.PHONY: build clean test run install help

# Build variables
BINARY_NAME=aura
BUILD_DIR=build
VERSION=1.0.0

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -ldflags "-X main.version=$(VERSION)" .

# Build for multiple platforms
build-all: build-linux build-darwin build-windows

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .

build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .

# Run the application
run:
	$(GOBUILD) -o $(BINARY_NAME) . && ./$(BINARY_NAME)

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)

# Test the application
test:
$(GOTEST) -v ./...

# Run unit tests only
test-unit:
$(GOTEST) -v ./config ./llm ./models

# Run integration tests
test-integration:
$(GOTEST) -v ./app

# Run end-to-end tests
test-e2e:
$(GOTEST) -v ./e2e

# Run tests with coverage
test-coverage:
$(GOTEST) -v -cover -coverprofile=coverage.out ./...
$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Run benchmarks
test-bench:
$(GOTEST) -bench=. -benchmem ./...

# Run tests with race detection
test-race:
$(GOTEST) -v -race ./...

# Run comprehensive test suite
test-all: test test-coverage test-bench test-race
@echo "All tests completed successfully!"

# Run tests with custom test runner
test-runner:
$(GOBUILD) -o test_runner test_runner.go && ./test_runner

# Quick test (unit tests only)
test-quick:
$(GOTEST) -short ./...

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Install to system PATH
install: build
	cp $(BINARY_NAME) /usr/local/bin/

# Create build directory
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# Prepare for multi-platform builds
build-all: $(BUILD_DIR) build-linux build-darwin build-windows

# Show help
help:
@echo "Available targets:"
@echo ""
@echo "Build targets:"
@echo "  build           Build the application for current platform"
@echo "  build-all       Build for all supported platforms"
@echo "  run             Build and run the application"
@echo "  clean           Remove build artifacts"
@echo "  install         Install to system PATH"
@echo ""
@echo "Test targets:"
@echo "  test            Run all tests"
@echo "  test-unit       Run unit tests only (config, llm, models)"
@echo "  test-integration Run integration tests (app)"
@echo "  test-e2e        Run end-to-end tests"
@echo "  test-coverage   Run tests with coverage report"
@echo "  test-bench      Run benchmark tests"
@echo "  test-race       Run tests with race detection"
@echo "  test-all        Run comprehensive test suite"
@echo "  test-runner     Run tests with custom test runner"
@echo "  test-quick      Run quick tests (short mode)"
@echo ""
@echo "Other targets:"
@echo "  deps            Download and tidy dependencies"
@echo "  help            Show this help message"
