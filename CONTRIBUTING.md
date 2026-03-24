# Contributing Guide

Thank you for your interest in contributing to `gohealthchecker`! This document guides you on how to set up your environment and the process for submitting improvements.

## Prerequisites

- **Go**: >= 1.25
- **Make**: For automated command execution.
- **Git**: For version control.

## Environment Setup

1. Fork the repository.
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR-USERNAME/gohealthchecker.git
   cd gohealthchecker
   ```
3. Install dependencies:
   ```bash
   go mod download
   ```

## Development Lifecycle

Use the `Makefile` commands to facilitate the workflow:

- **Tests**: `make test`
- **Lint**: `make lint` (requires golangci-lint)
- **Build**: `make build`

## Code Standards

- Follow the conventions of [Effective Go](https://golang.org/doc/effective_go.html).
- Ensure all exported code is documented (standard `godoc`).
- Use `go fmt` before each commit.

## Pull Request Process

1. Create a branch for your change: `git checkout -b feature/my-improvement`.
2. Make your commits following clear patterns.
3. Ensure tests pass (`make test`).
4. Submit the PR to the `main` branch.

## Where to Contribute?

- TUI improvements (new visual components).
- New notification types (Telegram, Email, etc.).
- Concurrency optimizations.
- Documentation improvements.
