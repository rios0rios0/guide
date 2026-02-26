# Go Project Structure

> **TL;DR:** Follow the domain/infrastructure layer separation. Use `go.mod` for dependency management. Place test files next to production code. Use the `test/` directory for shared test helpers, builders, and doubles.

## Overview

This page defines the standard directory layout and dependency management practices for all Go projects. The architecture follows the [Backend Design](../../Life-Cycle/Architecture/Backend-Design.md) specification, separating code into `domain` (contracts) and `infrastructure` (implementations) layers.

## Directory Structure

```
main/
  domain/                   (contracts)
    commands/                 business logic; the only layer without a contract
    entities/
    repositories/
    app.go
    main.go
    wire.go
    wire_gen.go
  infrastructure/           (implementations)
    controllers/
      mappers/
      requests/
      responses/
    repositories/             prefixed with the tool name; returns database models
      mappers/
      models/
test/
  domain/
    builders/                 test data builders
    doubles/
      repositories/           stubs, dummies, fakes
    helpers/
  infrastructure/
    doubles/
      repositories/
```

### Key Directories

| Directory                          | Purpose                                                      |
|------------------------------------|--------------------------------------------------------------|
| `main/domain/commands/`            | Business logic implementations                               |
| `main/domain/entities/`            | Framework-agnostic domain entities                           |
| `main/domain/repositories/`        | Repository interface contracts                               |
| `main/infrastructure/controllers/` | HTTP controllers (request/response handling)                 |
| `main/infrastructure/repositories/`| Repository implementations with library-specific code        |
| `test/domain/builders/`            | Builder pattern implementations for constructing test data   |
| `test/domain/doubles/`             | Test doubles (stubs, dummies, fakers, in-memory)             |

## Package Manager: Go Modules

All Go projects use [Go Modules](https://go.dev/ref/mod) for dependency management. The module is defined in `go.mod` at the project root.

### Initializing a Module

```bash
go mod init github.com/org/project-name
```

### Managing Dependencies

```bash
# Add a dependency
go get github.com/sirupsen/logrus

# Add a specific version
go get github.com/sirupsen/logrus@v1.9.3

# Update all dependencies
go get -u ./...

# Remove unused dependencies
go mod tidy
```

### go.mod Example

```go
module github.com/org/project-name

go 1.23

require (
    github.com/gorilla/mux v1.8.1
    github.com/sirupsen/logrus v1.9.3
    github.com/stretchr/testify v1.9.0
    github.com/google/wire v0.6.0
)
```

### go.sum

The `go.sum` file contains cryptographic checksums for all dependencies and must be committed to version control. Do not edit it manually.

## Build & Distribution

### Building

```bash
# Build the binary
go build -o bin/app ./main

# Build with version information
go build -ldflags "-X main.version=1.0.0" -o bin/app ./main
```

### Running

```bash
# Run directly
go run ./main

# Run the compiled binary
./bin/app
```

### Docker

Use multi-stage builds to produce minimal container images:

```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bin/app ./main

FROM alpine:3.19
COPY --from=builder /bin/app /bin/app
ENTRYPOINT ["/bin/app"]
```

## Key Configuration Files

| File             | Purpose                                              |
|------------------|------------------------------------------------------|
| `go.mod`         | Module path and dependency declarations              |
| `go.sum`         | Dependency checksums (auto-generated, must be committed) |
| `.golangci.yml`  | golangci-lint configuration                          |
| `wire.go`        | Wire dependency injection declarations               |
| `wire_gen.go`    | Wire-generated DI code (auto-generated, must be committed) |
| `.editorconfig`  | Editor standardization                               |

## References

- [Go Modules Reference](https://go.dev/ref/mod)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Docker Multi-stage Builds](https://docs.docker.com/build/building/multi-stage/)
