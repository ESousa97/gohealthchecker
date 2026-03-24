# Makefile for GoHealthChecker

# Variables
BINARY_NAME=gohealthchecker
CMD_PATH=./cmd/gohealthchecker
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "v0.1.0")

# Main Commands
.PHONY: build run test lint clean install help

help: ## Show this help menu
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the project binary
	@echo "Building $(BINARY_NAME)..."
	@go build -ldflags "-X main.Version=$(VERSION)" -o bin/$(BINARY_NAME) $(CMD_PATH)

run: ## Run the application directly from source
	@go run $(CMD_PATH)

test: ## Run all unit tests
	@go test -v ./...

lint: ## Run the linter (requires golangci-lint)
	@if command -v golangci-lint >/dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install at https://golangci-lint.run"; \
		exit 1; \
	fi

clean: ## Remove temporary files and binaries
	@rm -rf bin/
	@go clean

install: ## Install binary to system (GOPATH/bin)
	@go install $(CMD_PATH)
