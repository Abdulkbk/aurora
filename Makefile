# Aurora CLI Makefile

# Binary name
BINARY_NAME=aurora

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# Version info
VERSION?=0.1.0
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X github.com/abdulkbk/aurora/cmd.version=$(VERSION) -X github.com/abdulkbk/aurora/cmd.commit=$(COMMIT)"

.PHONY: all build clean test fmt tidy install help

# Default target
all: build

## build: Build the aurora binary
build:
	@echo ".. building $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "[ok] built $(BINARY_NAME)"

## install: Install aurora to GOPATH/bin
install:
	@echo ".. installing $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(GOPATH)/bin/$(BINARY_NAME) .
	@echo "[ok] installed to $(GOPATH)/bin/$(BINARY_NAME)"

## test: Run all tests
test:
	@echo ".. running tests..."
	$(GOTEST) -v ./...

## fmt: Format Go code
fmt:
	@echo ".. formatting code..."
	$(GOFMT) ./...

## tidy: Tidy and verify module dependencies
tidy:
	@echo ".. tidying modules..."
	$(GOMOD) tidy
	$(GOMOD) verify

## clean: Remove build artifacts
clean:
	@echo ".. cleaning..."
	@rm -f $(BINARY_NAME)
	@echo "[ok] cleaned"

## run: Build and run with sample args (for testing)
run: build
	./$(BINARY_NAME) --help

## help: Show this help message
help:
	@echo "Aurora CLI - Build custom Docker images from GitHub PRs for Polar"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## /  /'
