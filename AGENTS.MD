# Coding Agent Instructions

This document provides instructions for AI coding agents working on this repository.

## Project Overview

## Technology Stack

- **Go**: 1.25+ required
- **CLI**: [spf13/cobra](https://github.com/spf13/cobra)

## Directory Structure

``` bash
.
├── cmd/
├── internal/
├── tasks/
└── main.go
```

## Development Commands

```bash
# Install dependencies
go mod download

# Build
go build -o bin/diagram-gen

# Run tests
go test ./...

# Run tests with coverage
go test -coverprofile=c.out ./...

# Run tests with coverage and tail
go test -coverprofile=c.out ./... 2>&1 && go tool cover -func=c.out 2>&1 | tail -5

# Run tests with coverage verbose and tail
go test -v -coverprofile=c.out ./... 2>&1 && go tool cover -func=c.out 2>&1 | tail -5

# Format code
gofmt -w ./cmd ./internal ./main.go

# Run linter
golangci-lint-v2 run

# Run CLI
go run . --help
go run . generate path/to/file.go -o diagram.drawio
```

## Code Style Guidelines

1. **Error Handling**: Always handle errors explicitly. Do not ignore errors with `_`.
2. **CLI Commands**: Keep Cobra commands in `cmd/` and core logic in `internal/`.

## Testing

- Tests are in `*_test.go` files alongside source code
- Run `go test ./...` before committing changes

## Security Considerations

- Always validate user input using validator tags
