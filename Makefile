# Cointracker - Ethereum Transaction Exporter Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=cointracker
BINARY_UNIX=$(BINARY_NAME)_unix

# Build target for local development
.PHONY: build
build:
	$(GOBUILD) -o $(BINARY_NAME) .

# Clean build artifacts
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Run tests
.PHONY: test
test:
	$(GOTEST) -v ./pkg/...

# Run all tests (including integration tests)
.PHONY: test-all
test-all:
	$(GOTEST) -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	$(GOTEST) -race -coverprofile=coverage.out -covermode=atomic ./pkg/...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Download dependencies
.PHONY: deps
deps:
	$(GOMOD) download
	$(GOMOD) verify

# Tidy up dependencies
.PHONY: tidy
tidy:
	$(GOMOD) tidy

# Build for Linux
.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) .

# Build for Windows
.PHONY: build-windows
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME).exe .

# Build for macOS
.PHONY: build-macos
build-macos:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)_darwin_amd64 .
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BINARY_NAME)_darwin_arm64 .

# Build for all platforms
.PHONY: build-all
build-all: build-linux build-windows build-macos

# Run the application (requires ETHERSCAN_API_KEY environment variable)
.PHONY: run
run:
	$(GOBUILD) -o $(BINARY_NAME) . && ./$(BINARY_NAME)

# Format code
.PHONY: fmt
fmt:
	$(GOCMD) fmt ./...

# Vet code for potential issues
.PHONY: vet
vet:
	$(GOCMD) vet ./...

# Run linter (requires golangci-lint to be installed)
.PHONY: lint
lint:
	golangci-lint run

# Install the binary to GOPATH/bin
.PHONY: install
install:
	$(GOCMD) install .

# Example usage with sample address
.PHONY: example
example: build
	@echo "Running example with sample address..."
	@if [ -z "$$ETHERSCAN_API_KEY" ]; then \
		echo "Error: ETHERSCAN_API_KEY environment variable is required"; \
		exit 1; \
	fi
	./$(BINARY_NAME) fetch --address 0xa39b189482f984388a34460636fea9eb181ad1a6 --output example_transactions.csv

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary for current platform"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run tests for pkg directory"
	@echo "  test-all      - Run all tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  deps          - Download and verify dependencies"
	@echo "  tidy          - Tidy up go.mod and go.sum"
	@echo "  build-linux   - Build for Linux"
	@echo "  build-windows - Build for Windows" 
	@echo "  build-macos   - Build for macOS (both amd64 and arm64)"
	@echo "  build-all     - Build for all platforms"
	@echo "  run           - Build and run the application"
	@echo "  fmt           - Format Go code"
	@echo "  vet           - Vet Go code"
	@echo "  lint          - Run golangci-lint"
	@echo "  install       - Install binary to GOPATH/bin"
	@echo "  example       - Run example with sample address"
	@echo "  help          - Show this help message"

# Default target
.DEFAULT_GOAL := build