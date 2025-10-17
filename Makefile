# Variables
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_TIME := $(shell date -u '+%Y-%m-%d %H:%M:%S')
LDFLAGS := -X 'github.com/ofkm/goobrew/internal/version.Version=$(VERSION)' \
           -X 'github.com/ofkm/goobrew/internal/version.Commit=$(COMMIT)' \
           -X 'github.com/ofkm/goobrew/internal/version.BuildTime=$(BUILD_TIME)'

# Build binary
.PHONY: build
build:
	@echo "Building goobrew..."
	go build -ldflags="$(LDFLAGS)" -o goobrew .

# Build for release (stripped binary)
.PHONY: build-release
build-release:
	@echo "Building goobrew for release..."
	go build -ldflags="$(LDFLAGS) -s -w" -trimpath -o goobrew .

# Install binary
.PHONY: install
install:
	@echo "Installing goobrew..."
	go install -ldflags="$(LDFLAGS)" .

# Run tests
.PHONY: test
test:
	go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linter
.PHONY: lint
lint:
	golangci-lint run --timeout=5m

# Clean build artifacts
.PHONY: clean
clean:
	rm -f goobrew coverage.out coverage.html

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build          - Build the binary"
	@echo "  build-release  - Build optimized binary for release"
	@echo "  install        - Install the binary to GOPATH/bin"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  lint           - Run golangci-lint"
	@echo "  clean          - Remove build artifacts"
	@echo "  help           - Show this help message"

.DEFAULT_GOAL := build
