# Go

> **TL;DR:** Use `snake_case` for file names, a short abbreviation of the type as the method receiver (e.g., `c` for `Client`), [Dig](https://github.com/uber-go/dig) for dependency injection, [golangci-lint](https://golangci-lint.run/) for linting, [Logrus](https://github.com/sirupsen/logrus) for logging, and [testify](https://github.com/stretchr/testify) for testing. Entities must be framework-agnostic.

## Overview

This series of pages outlines Go-specific conventions, with detailed explanations of recommended approaches and the reasoning behind them. For the general baseline, refer to the Code Style guide. See the sub-pages for specific topics:

## Go Proverbs

The [Go Proverbs](https://go-proverbs.github.io/) capture the language's design philosophy:

- Don't communicate by sharing memory, share memory by communicating.
- Concurrency is not parallelism.
- The bigger the interface, the weaker the abstraction.
- Make the zero value useful.
- A little copying is better than a little dependency.
- Clear is better than clever.
- Errors are values.
- Don't just check errors, handle them gracefully.

---

# Go Conventions

> **TL;DR:** Use `snake_case` for file names, a short abbreviation of the type as the method receiver name (e.g., `c` for `Command`), and follow the strict naming patterns for Commands, Controllers, Repositories, and Mappers. Entities must be framework-agnostic. Use [Dig](https://github.com/uber-go/dig) for dependency injection.

## Overview

This document defines Go-specific naming conventions and component patterns. For the general baseline, refer to the Code Style guide. The architectural layers referenced here are defined in the Backend Design section.

## File Naming

All file names must use `snake_case`:

```
list_users_command.go     # Correct
listUsersCommand.go       # Wrong
ListUsersCommand.go       # Wrong
```

## General Conventions

1. Use a **one or two letter abbreviation** of the type name as the method receiver (e.g., `c` for `Command`, `r` for `Repository`, `m` for `Mapper`). Do not use generic names like `self`, `this`, or `me` -- this follows the [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments#receiver-names) convention and is enforced by revive's `receiver-naming` rule. The receiver name must be **consistent** across all methods of the same type.
2. Only attach a method to a struct when the method **needs to mutate** the struct's state.
3. For an introduction to the DTO pattern, refer to [this article](https://www.baeldung.com/java-dto-pattern).

## Entities

Entities are the core of the application. All business logic related to properties and fields belongs inside the entity.

**Entities must be free of any framework or external tool dependencies.** Do not use tags (e.g., `json`, `gorm`) inside entity structs.

## Commands

| Element     | Pattern                           | Example                                                                     |
|-------------|-----------------------------------|-----------------------------------------------------------------------------|
| File name   | `<operation>_<entity>_command.go` | `list_users_command.go`                                                     |
| Struct name | `<Operation><Entity>Command`      | `ListUsersCommand`                                                          |
| Method name | `Execute`                         | `func (c ListUsersCommand) Execute(listeners ListUsersCommandListeners)` |

**Notes:**
- Use plural entity names when the operation targets multiple entities.
- Use the standard [operations vocabulary](../../Code-Style.md#operations-vocabulary).
- Listeners must reflect all possible controller responses.

## Controllers

| Element     | Pattern                              | Example                                     |
|-------------|--------------------------------------|---------------------------------------------|
| File name   | `<operation>_<entity>_controller.go` | `list_users_controller.go`                  |
| Struct name | `<Operation><Entity>Controller`      | `ListUsersController`                       |
| Method name | `Execute`                            | `func (c ListUsersController) Execute()` |

## Services

This layer is **not used** in Go projects.

## Repositories

### Contract (Domain Layer)

| Element        | Pattern                  | Example               |
|----------------|--------------------------|-----------------------|
| File name      | `<entity>_repository.go` | `users_repository.go` |
| Interface name | `<Entity>Repository`     | `UsersRepository`     |

### Implementation (Infrastructure Layer)

| Element     | Pattern                            | Example                   |
|-------------|------------------------------------|---------------------------|
| File name   | `<library>_<entity>_repository.go` | `pgx_users_repository.go` |
| Struct name | `<Library><Entity>Repository`      | `PgxUsersRepository`      |

### Method Naming

Methods follow a logical sequence: find one, find all, filter, check existence, save one, save all, delete.

```go
// Find a single entity by a specific field
func (r UsersRepository) FindByTargetField(targetField any) entities.User

// Find multiple entities by a specific field
func (r UsersRepository) FindAllByTargetField(targetField any) []entities.User

// Check existence (returns boolean)
func (r UsersRepository) HasBooleanVerification(targetField any) bool

// Persist a single entity
func (r UsersRepository) Save(user entities.User)

// Persist multiple entities
func (r UsersRepository) SaveAll(users []entities.User)

// Remove a single entity by a specific field
func (r UsersRepository) DeleteByTargetField(targetField any)
```

**Notes:**
- `TargetField` is a placeholder (e.g., `Id`, `Name`, `Email`).
- `BooleanVerification` is a placeholder (e.g., `UserInGroup`, `UserPermission`).
- Implementations use the same signatures but attach to the concrete struct (e.g., `PgxUsersRepository`).

## Mappers

### Repository Mappers

| Element     | Pattern              | Example          |
|-------------|----------------------|------------------|
| File name   | `<entity>_mapper.go` | `user_mapper.go` |
| Struct name | `<Entity>Mapper`     | `UserMapper`     |

```go
// Infrastructure DTO -> Domain Entity
func (m UserMapper) MapToEntity(infra any) entities.User
func (m UserMapper) MapToEntities(infra []any) []entities.User

// Domain Entity -> Infrastructure DTO
func (m UserMapper) MapToExternal(user entities.User) models.External
func (m UserMapper) MapToExternals(users []entities.User) []models.External
```

### Controller Mappers

| Element              | Pattern                                   | Example                          |
|----------------------|-------------------------------------------|----------------------------------|
| File name (request)  | `<operation>_<entity>_request_mapper.go`  | `insert_user_request_mapper.go`  |
| File name (response) | `<operation>_<entity>_response_mapper.go` | `insert_user_response_mapper.go` |

```go
// Request -> Entity (no inverse mapping)
func (m InsertUserRequestMapper) MapToEntity(request InsertUserRequest) entities.User
func (m InsertUserRequestMapper) MapToEntities(requests []InsertUserRequest) []entities.User

// Entity -> Response (no inverse mapping)
func (m InsertUserResponseMapper) MapToResponse(user entities.User) responses.InsertUserResponse
func (m InsertUserResponseMapper) MapToResponses(users []entities.User) []responses.InsertUserResponse
```

**Important:** Do not use `json` tags outside the infrastructure layer. Tags are restricted to request and response DTOs.

## Models

Models reside exclusively in the infrastructure layer and represent DTOs for external data sources (databases, APIs, queues, etc.). They resemble entities but are **not** entities.

Each model is prefixed with the name of the external tool it communicates with:

| Example       | Source             |
|---------------|--------------------|n| `AwsFile`     | AWS S3             |
| `ApiDocument` | External API       |
| `PgxUser`     | PostgreSQL via pgx |

## Dependency Injection

Use [Uber Dig](https://github.com/uber-go/dig) for runtime dependency injection via constructor injection. Dig resolves dependencies automatically by matching types from registered provider functions, requiring no code generation or manual wiring.

### Container File Convention

Each architectural layer must have a `container.go` file that registers its own providers. A top-level orchestrator calls each layer's registration function in dependency order.

| File                                                | Purpose                                      |
|-----------------------------------------------------|----------------------------------------------|
| `cmd/<app>/dig.go`                                  | Creates the container and invokes root types |
| `internal/container.go`                             | Orchestrates registration across all layers  |
| `internal/domain/entities/container.go`             | Registers entity providers (or no-op)        |
| `internal/domain/commands/container.go`             | Registers command providers (or no-op)       |
| `internal/infrastructure/controllers/container.go`  | Registers controller providers               |
| `internal/infrastructure/repositories/container.go` | Registers repository providers               |

### Orchestrator Pattern

The top-level orchestrator registers providers in bottom-up dependency order:

```go
package internal

import "go.uber.org/dig"

func RegisterProviders(container *dig.Container) error {
    if err := repositories.RegisterProviders(container); err != nil {
        return err
    }
    if err := entities.RegisterProviders(container); err != nil {
        return err
    }
    if err := commands.RegisterProviders(container); err != nil {
        return err
    }
    if err := controllers.RegisterProviders(container); err != nil {
        return err
    }
    if err := container.Provide(NewAppInternal); err != nil {
        return err
    }
    return nil
}
```

### Layer Registration

Each layer registers its constructors. Dig resolves dependencies by matching constructor parameter types to previously registered providers:

```go
package controllers

import "go.uber.org/dig"

func RegisterProviders(container *dig.Container) error {
    if err := container.Provide(NewListUsersController); err != nil {
        return err
    }
    if err := container.Provide(NewDeleteUserController); err != nil {
        return err
    }
    return nil
}
```

For layers with no providers, maintain the function as a no-op for architectural consistency:

```go
package commands

import "go.uber.org/dig"

func RegisterProviders(_ *dig.Container) error {
    return nil
}
```

### Injection Functions

Create injection functions in `cmd/<app>/dig.go` that build the container and invoke the desired root type:

```go
package main

import (
    "go.uber.org/dig"
    "myapp/internal"
    "myapp/internal/infrastructure/controllers"
)

func injectController() *controllers.ListUsersController {
    container := dig.New()
    if err := internal.RegisterProviders(container); err != nil {
        panic(err)
    }

    var controller *controllers.ListUsersController
    if err := container.Invoke(func(c *controllers.ListUsersController) {
        controller = c
    }); err != nil {
        panic(err)
    }
    return controller
}
```

### Anonymous Providers for Complex Initialization

When a provider requires post-construction setup (e.g., registering adapters), use an anonymous function:

```go
if err := container.Provide(func() *ServiceRegistry {
    registry := NewServiceRegistry()
    registry.Register("github", github.NewAdapter())
    registry.Register("gitlab", gitlab.NewAdapter())
    return registry
}); err != nil {
    return err
}
```

### Type Aggregation

Collect multiple concrete types into a slice for bulk injection:

```go
func NewControllers(
    listController *ListUsersController,
    deleteController *DeleteUserController,
) *[]entities.Controller {
    return &[]entities.Controller{
        listController,
        deleteController,
    }
}
```

---

# Go Formatting and Linting

> **TL;DR:** Use **gofmt** (built-in) for code formatting, **goimports** for import ordering, and **golangci-lint** as the linter aggregator. These tools are non-negotiable and must be integrated into every project's CI pipeline.

## Overview

Go's toolchain includes a built-in formatter (`gofmt`) that eliminates all debates about code style. Combined with `goimports` for import management and `golangci-lint` for static analysis, this toolchain ensures consistent, high-quality code across all projects.

## Formatter: gofmt

`gofmt` is Go's official code formatter. It ships with the Go toolchain and produces a single canonical formatting for any Go source file. There are no configuration options -- this is by design.

```bash
# Format a single file
gofmt -w main.go

# Format all files in the current module
gofmt -w .
```

**Do not use alternative formatters.** `gofmt` is the universal standard in the Go ecosystem, and all Go code must be formatted with it.

## Import Ordering: goimports

[goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports) extends `gofmt` by automatically managing import statements -- adding missing imports, removing unused ones, and grouping them into sections:

1. Standard library
2. Third-party packages
3. Application packages

```bash
# Install
go install golang.org/x/tools/cmd/goimports@latest

# Run
goimports -w .
```

Example output:

```go
import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"myapp/domain/commands"
	"myapp/infrastructure/controllers"
)
```

## Linter: golangci-lint

[golangci-lint](https://golangci-lint.run/) aggregates dozens of Go linters into a single tool. It is fast, configurable, and must be used in all projects.

### Installation

```bash
# Binary installation (recommended for CI)
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Or via go install
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Usage

```bash
# Run all enabled linters
golangci-lint run

# Run on specific directories
golangci-lint run ./domain/... ./infrastructure/...
```

### Configuration

Place a `.golangci.yml` file in the project root. A recommended baseline configuration:

```yaml
run:
  timeout: 5m

linters:
  enable:
    - errcheck
    - govet
    - staticcheck
    - unused
    - gosimple
    - ineffassign
    - typecheck
    - misspell
    - gocyclo
    - revive
    - gocritic
    - nakedret
    - prealloc

linters-settings:
  gocyclo:
    min-complexity: 15
  revive:
    rules:
      - name: unexported-return
        disabled: true
      - name: receiver-naming

issues:
  exclude-use-default: false
```

## Editor Configuration

### Visual Studio Code

Install the [Go extension](https://marketplace.visualstudio.com/items?itemName=golang.go) and add the following settings:

```json
{
    "go.formatTool": "goimports",
    "go.lintTool": "golangci-lint",
    "go.lintOnSave": "package",
    "[go]": {
        "editor.formatOnSave": true,
        "editor.codeActionsOnSave": {
            "source.organizeImports": "explicit"
        }
    }
}
```

### NeoVim

Use [nvim-lspconfig](https://github.com/neovim/nvim-lspconfig) with `gopls` and configure format-on-save:

```lua
require("lspconfig").gopls.setup({
    settings = {
        gopls = {
            gofumpt = true,
            analyses = {
                unusedparams = true,
                shadow = true,
            },
            staticcheck = true,
        },
    },
})
```

---

# Go Type System

> **TL;DR:** Go is statically typed -- the compiler enforces type safety at build time. Design small interfaces, accept interfaces and return structs, and use generics (Go 1.18+) only when they reduce duplication without sacrificing readability.

## Overview

Go's type system is static and checked at compile time, so explicit type annotations on variables are rarely needed (the compiler infers them via `:=`). This page focuses on the patterns and principles that maximize type safety and code clarity in Go projects.

## Interface Design

### Accept Interfaces, Return Structs

Functions should accept interfaces as parameters and return concrete structs. This keeps the caller flexible while keeping the implementation explicit:

```go
// Correct -- accepts an interface
func ProcessItems(repo domain.ItemsRepository) error {
    items, err := repo.FindAll()
    // ...
}

// Wrong -- accepts a concrete type, limiting testability
func ProcessItems(repo *repositories.PgxItemsRepository) error {
    items, err := repo.FindAll()
    // ...
}
```

### Small Interfaces

The bigger the interface, the weaker the abstraction. Define interfaces with only the methods the consumer needs:

```go
// Correct -- small, focused interface
type ItemReader interface {
    FindByID(ctx context.Context, id int64) (*entities.Item, error)
}

// Avoid -- large interface that forces implementors to define unused methods
type ItemRepository interface {
    FindByID(ctx context.Context, id int64) (*entities.Item, error)
    FindAll(ctx context.Context) ([]entities.Item, error)
    Save(ctx context.Context, item *entities.Item) error
    Delete(ctx context.Context, id int64) error
    // ... many more methods
}
```

When the full interface is needed (e.g., in the domain layer contract), keep it. But when a function only reads items, accept `ItemReader` instead of the full `ItemRepository`.

### Interface Placement

Interfaces belong in the package that **uses** them, not the package that implements them. In the project's architecture, domain-layer interfaces are defined in `domain/repositories/` and implemented in `infrastructure/repositories/`.

## Generics (Go 1.18+)

Use generics when they eliminate genuine code duplication across multiple types. Do not use generics for single-type operations or when an interface would be clearer.

### When to Use

```go
// Correct -- generic function that works across multiple slice types
func Contains[T comparable](slice []T, target T) bool {
    for _, item := range slice {
        if item == target {
            return true
        }
    }
    return false
}
```

### When Not to Use

```go
// Unnecessary -- only works with one type, a regular function is clearer
func FindUserByID[T entities.User](users []T, id int64) *T { ... }

// Better
func FindUserByID(users []entities.User, id int64) *entities.User { ... }
```

## Type Assertions and Type Switches

Prefer type switches over chains of type assertions:

```go
// Correct -- type switch
switch v := value.(type) {
case string:
    logger.Info(v)
case int:
    logger.Info(strconv.Itoa(v))
default:
    logger.Warn("unexpected type")
}

// Avoid -- chained assertions
if s, ok := value.(string); ok {
    logger.Info(s)
} else if i, ok := value.(int); ok {
    logger.Info(strconv.Itoa(i))
}
```

## Prohibited Patterns

```go
// Wrong -- empty interface as a catch-all parameter
func Process(data interface{}) { ... }
func Process(data any) { ... }
```

Using `any` (`interface{}`) as a function parameter defeats the purpose of static typing. Define a proper interface or use generics with type constraints instead.

---

# Go Logging

> **TL;DR:** Use **[Logrus](https://github.com/sirupsen/logrus)** for all logging. Do not use Go's standard `log` package or `fmt.Println` for application logging. Always import with the alias `logger`. Use structured logging with `WithFields()` instead of string interpolation.

## Overview

Consistent, structured logging is essential for production observability. This page defines the mandatory logging library and patterns for all Go projects.

## Mandatory Library: Logrus

**Use [Logrus](https://github.com/sirupsen/logrus) for all logging.** Logrus provides structured logging, consistent log levels, JSON output support, and field-based contextual logging -- all of which are essential for production observability.

### Installation

```bash
go get github.com/sirupsen/logrus
```

**Important:** Always use the lowercase import path `github.com/sirupsen/logrus` (not the uppercase variant).

### Import Convention

Always import Logrus with the alias `logger` to keep usage concise and consistent across the codebase:

```go
import logger "github.com/sirupsen/logrus"
```

### Usage

```go
import logger "github.com/sirupsen/logrus"

func main() {
    logger.Info("application started")
    logger.WithFields(logger.Fields{
        "user_id": 42,
        "action":  "login",
    }).Info("user authenticated")
}
```

## Log Levels

| Level | Method           | When to Use                                                         |
|-------|------------------|---------------------------------------------------------------------|
| Trace | `logger.Trace()` | Very fine-grained diagnostic information                            |
| Debug | `logger.Debug()` | Diagnostic information useful during development                    |
| Info  | `logger.Info()`  | General operational events (application started, request processed) |
| Warn  | `logger.Warn()`  | Potential issues that do not prevent operation                      |
| Error | `logger.Error()` | Errors that prevent a specific operation but not the application    |
| Fatal | `logger.Fatal()` | Critical errors that require immediate application shutdown         |
| Panic | `logger.Panic()` | Critical errors that should panic after logging                     |

## Structured Logging

Always use `WithFields` to attach contextual data to log entries rather than interpolating values into the message string:

```go
// Correct -- structured fields
logger.WithFields(logger.Fields{
    "request_id": requestID,
    "status":     statusCode,
}).Info("request completed")

// Wrong -- string interpolation
logger.Infof("request %s completed with status %d", requestID, statusCode)
```

## Prohibited Patterns

```go
// Wrong -- standard library logger
import "log"
log.Println("something happened")

// Wrong -- fmt for application logging
fmt.Println("something happened")

// Wrong -- uppercase import path
import "github.com/Sirupsen/logrus"
```

---

# Go Testing Conventions

> **TL;DR:** Use build flags (`//go:build unit` or `//go:build integration`) on every test file. Place test files next to production code with the `_test.go` suffix. Use `stretchr/testify` for suites and assertions. Test packages must be **external** to the production package. All tests must follow the BDD pattern with `// given`, `// when`, `// then` comment blocks. Unit tests must run in **parallel** using `t.Parallel()` + `t.Run()`. Integration tests use **suites** with setup/teardown and are NOT parallel.

## Overview

Go discovers test files automatically by scanning for files ending in `_test.go`. This document defines the conventions for organizing and writing tests across all Go projects.

## File Structure

```
main/
  domain/
  infrastructure/
    repositories/
      sqlx_items_repository.go
      sqlx_items_repository_test.go       <-- placed next to production file
test/
  domain/
    builders/                              <-- test data builders
    doubles/
      repositories/                        <-- stubs, dummies, fakes
    helpers/
  infrastructure/
    doubles/
      repositories/
```

## General Conventions

1. **Build flags are mandatory.** Every test file must start with a build flag specifying the test type:
   ```go
   //go:build unit
   ```
   ```go
   //go:build integration
   ```
2. **External test packages.** The test package must be outside the production code package. For example, if the production code is in `package commands`, the test file must use `package commands_test`.
3. **Testing framework.** Use [`stretchr/testify`](https://github.com/stretchr/testify) for test suites and assertions.
4. **File naming.** Test files use the `_test` suffix (e.g., `sqlx_items_repository_test.go`).
5. **File placement.** Test files are placed next to the corresponding production file.
6. **BDD structure.** Every test must use `// given`, `// when`, `// then` comment blocks to separate preconditions, actions, and assertions.
7. **Parallel unit tests.** All unit tests must call `t.Parallel()` and use `t.Run()` sub-tests. All sub-tests within a function execute in parallel.
8. **Sequential integration tests.** Integration tests use `suite.Suite` with setup/teardown and must NOT be parallel, as they share database state.

## Unit Tests (Parallel with `t.Run`)

Unit tests must be structured for **parallel execution**. Each top-level test function calls `t.Parallel()`, and each scenario is a `t.Run()` sub-test. All sub-tests within the same function run concurrently.

### Command Tests

```go
package commands_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"myapp/domain/commands"
	"myapp/test/domain/doubles"
	domainErrors "myapp/domain/errors"
)

func TestDeleteItemCommand(t *testing.T) {
	t.Parallel()

	t.Run("should call OnSuccess when the item is deleted", func(t *testing.T) {
		// given
		repository := doubles.NewItemRepositoryStub()
		command := commands.NewDeleteItemCommand(repository)
		onSuccessCalled := false

		// when
		listeners := commands.DeleteItemListeners{
			OnSuccess: func() {
				onSuccessCalled = true
			},
		}
		command.Execute(context.TODO(), 1, listeners)

		// then
		assert.True(t, onSuccessCalled, "OnSuccess should have been called")
	})

	t.Run("should call OnNotFound when the item is not found", func(t *testing.T) {
		// given
		repository := doubles.NewItemRepositoryStub().WithOnError(domainErrors.ErrRecordNotFound)
		command := commands.NewDeleteItemCommand(repository)
		onNotFoundCalled := false

		// when
		listeners := commands.DeleteItemListeners{
			OnNotFound: func() {
				onNotFoundCalled = true
			},
		}
		command.Execute(context.TODO(), 1, listeners)

		// then
		require.True(t, onNotFoundCalled, "OnNotFound should have been called")
	})

	t.Run("should call OnError when there is an error processing the delete", func(t *testing.T) {
		// given
		dbProcessErr := errors.New("test error")
		repository := doubles.NewItemRepositoryStub().WithOnError(dbProcessErr)
		command := commands.NewDeleteItemCommand(repository)

		// when & then
		listeners := commands.DeleteItemListeners{
			OnError: func(err error) {
				assert.Error(t, err)
			},
		}
		command.Execute(context.TODO(), 1, listeners)
	})
}
```

**Key points:**
- `t.Parallel()` is called at the top of `TestDeleteItemCommand`, enabling all `t.Run()` sub-tests to execute concurrently.
- Each sub-test is self-contained -- it creates its own doubles, command, and listeners.
- Listeners pattern reflects all possible controller responses: `OnSuccess`, `OnNotFound`, `OnError`.

### Controller Tests

```go
package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"myapp/infrastructure/controllers"
	"myapp/test/domain/doubles"
)

func TestListItemsController(t *testing.T) {
	t.Parallel()

	t.Run("should respond 200 (OK) when items are listed successfully", func(t *testing.T) {
		// given
		command := doubles.NewListItemsCommandStub()
		ctrl := controllers.NewListItemsController(command)

		// when
		req, _ := http.NewRequest("GET", "/items", nil)
		w := httptest.NewRecorder()
		ctrl.Execute(w, req)

		// then
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should respond 500 (Internal Server Error) when command fails", func(t *testing.T) {
		// given
		command := doubles.NewListItemsCommandStub().WithOnError()
		ctrl := controllers.NewListItemsController(command)

		// when
		req, _ := http.NewRequest("GET", "/items", nil)
		w := httptest.NewRecorder()
		ctrl.Execute(w, req)

		// then
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
```

## Service Tests

This layer is **not used** in Go projects.

## Integration Tests (Suite with Setup/Teardown)

Integration tests use `suite.Suite` because they require shared infrastructure (databases, external services). They are **NOT parallel** due to shared mutable state. Use `SetupSuite`, `SetupTest`, and `TearDownTest` to manage the lifecycle.

Group sub-tests by outcome or feature using `suite.Run()`.

### Repository Tests

```go
package repositories_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/suite"

	"myapp/domain"
	"myapp/domain/entities"
	domainErrors "myapp/domain/errors"
	"myapp/infrastructure/repositories"
	"myapp/test/domain/builders"
	dbTest "myapp/test/infrastructure/helpers"
)

type ItemsRepositorySuite struct {
	dbTest.DatabaseSuite
}

func (s *ItemsRepositorySuite) getSeeders() []dbTest.TableSeed {
	return []dbTest.TableSeed{
		{
			MigrationFileRelativePathFromRoot: "db/seeders/test/items/item.sql",
			TableName:                         "items",
		},
	}
}

func (s *ItemsRepositorySuite) SetupSuite() {
	s.SetupDatabase()
}

func (s *ItemsRepositorySuite) SetupTest() {
	s.RunSeeders(s.DB, s.getSeeders())
}

func (s *ItemsRepositorySuite) TearDownTest() {
	sqlResp := s.DB.MustExec("DELETE FROM items")
	rows, err := sqlResp.RowsAffected()
	s.Require().NoError(err)
	s.Require().Positive(rows)
}

func (s *ItemsRepositorySuite) createItem() entities.Item {
	return builders.NewItemBuilder().
		WithID(0).
		WithExternalID(gofakeit.UUID()).
		WithPayload().
		Build()
}

func (s *ItemsRepositorySuite) getRandomID() int64 {
	var randomID int64
	query := "SELECT id FROM items LIMIT 1"
	err := s.DB.Get(&randomID, query)
	s.Require().NoError(err, "failed to get random row")
	return randomID
}

func (s *ItemsRepositorySuite) newRepository() *repositories.SQLXItemsRepository {
	return repositories.NewSQLXItemsRepository(s.DB)
}

func (s *ItemsRepositorySuite) newContext() context.Context {
	return context.Background()
}

func (s *ItemsRepositorySuite) newPagination(page, size int) domain.Pagination {
	return builders.NewPaginationBuilder().
		WithPage(page).
		WithSize(size).
		Build()
}

func (s *ItemsRepositorySuite) defaultPagination() domain.Pagination {
	return s.newPagination(1, 100)
}

func (s *ItemsRepositorySuite) newListParameters(
	pagination domain.Pagination,
	filters map[string]any,
	sorting []domain.SortField,
) domain.ListParameters {
	if filters == nil {
		filters = make(map[string]any)
	}
	if sorting == nil {
		sorting = []domain.SortField{}
	}
	return domain.ListParameters{
		Pagination: pagination,
		Filters:    filters,
		Sorting:    sorting,
	}
}

// --- Success Cases ---

func (s *ItemsRepositorySuite) TestItemsSuccess() {
	s.Run("should save a new item successfully", func() {
		// given
		item := s.createItem()
		repository := s.newRepository()
		ctx := s.newContext()

		// when
		created, err := repository.Save(ctx, &item)

		// then
		s.Require().NoError(err)
		s.Require().Positive(created.ID)
	})

	s.Run("should list all items successfully", func() {
		// given
		params := s.newListParameters(s.defaultPagination(), nil, nil)
		repository := s.newRepository()
		ctx := s.newContext()

		// when
		list, err := repository.ListAll(ctx, params)

		// then
		s.Greater(len(list.GetContent()), 1)
		s.NoError(err)
	})

	s.Run("should delete an item successfully", func() {
		// given
		itemID := s.getRandomID()
		repository := s.newRepository()
		ctx := s.newContext()

		// when
		err := repository.Delete(ctx, itemID)

		// then
		s.NoError(err)
	})

	s.Run("should get the target item by ID successfully", func() {
		// given
		itemID := s.getRandomID()
		repository := s.newRepository()
		ctx := s.newContext()

		// when
		item, err := repository.GetByID(ctx, itemID)

		// then
		s.Require().NoError(err)
		s.Require().NotNil(item)
		s.Require().Equal(item.ID, itemID)
	})

	s.Run("should update an item successfully", func() {
		// given
		itemID := s.getRandomID()
		repository := s.newRepository()
		ctx := s.newContext()

		existingItem, err := repository.GetByID(ctx, itemID)
		s.Require().NoError(err)

		originalTitle := existingItem.Title
		existingItem.Title = gofakeit.Word() + "_updated"

		// when
		updated, err := repository.Update(ctx, existingItem)

		// then
		s.Require().NoError(err)
		s.Equal(existingItem.ID, updated.ID)
		s.NotEqual(originalTitle, updated.Title)
	})
}

// --- Error Cases ---

func (s *ItemsRepositorySuite) TestItemsError() {
	s.Run("should return an error when saving with invalid data", func() {
		// given
		item := builders.NewItemBuilder().WithPayload().Build()
		item.ExternalID = "" // invalid
		repository := s.newRepository()

		// when
		_, err := repository.Save(s.newContext(), &item)

		// then
		s.Error(err)
	})

	s.Run("should return an error when listing with invalid pagination", func() {
		// given
		pagination := s.newPagination(-1, -1)
		params := s.newListParameters(pagination, nil, nil)
		repository := s.newRepository()
		ctx := s.newContext()

		// when
		list, err := repository.ListAll(ctx, params)

		// then
		s.Require().Error(err)
		s.Require().Nil(list)
	})

	s.Run("should return an error when deleting a non-existent item", func() {
		// given
		repository := s.newRepository()

		// when
		err := repository.Delete(s.newContext(), -9)

		// then
		s.Error(err)
	})

	s.Run("should return an error when getting an item that does not exist", func() {
		// given
		const itemID int64 = 99999
		repository := s.newRepository()
		ctx := s.newContext()

		// when
		item, err := repository.GetByID(ctx, itemID)

		// then
		s.Require().Error(err)
		s.Require().Nil(item)
	})

	s.Run("should return an error when updating a non-existent item", func() {
		// given
		const nonExistentID int64 = 99999
		item := s.createItem()
		item.ID = nonExistentID
		repository := s.newRepository()
		ctx := s.newContext()

		// when
		_, err := repository.Update(ctx, &item)

		// then
		s.Require().Error(err)
		s.Require().ErrorIs(err, domainErrors.ErrRecordNotFound)
	})
}

// --- Filter Cases ---

func (s *ItemsRepositorySuite) TestListAllWithFilters() {
	pagination := s.defaultPagination()
	repository := s.newRepository()
	ctx := s.newContext()

	s.Run("should filter items by category ID", func() {
		// given
		params := s.newListParameters(pagination, map[string]any{
			"category_id": 1,
		}, nil)

		// when
		list, _ := repository.ListAll(ctx, params)

		// then
		s.Require().NotEmpty(list.GetContent())
		for _, item := range list.GetContent() {
			s.Equal(uint(1), item.Category.ID)
		}
	})

	s.Run("should filter items by severity", func() {
		// given
		params := s.newListParameters(pagination, map[string]any{
			"severity": "high",
		}, nil)

		// when
		list, _ := repository.ListAll(ctx, params)

		// then
		s.Require().NotEmpty(list.GetContent())
		for _, item := range list.GetContent() {
			s.Equal("high", item.Level)
		}
	})

	s.Run("should return empty list when filter matches no items", func() {
		// given
		params := s.newListParameters(pagination, map[string]any{
			"category_id": 999,
		}, nil)

		// when
		list, err := repository.ListAll(ctx, params)

		// then
		s.Require().NoError(err)
		s.Require().Empty(list.GetContent())
	})
}

// --- Sorting Cases ---

func (s *ItemsRepositorySuite) TestListAllWithSorting() {
	pagination := s.defaultPagination()
	repository := s.newRepository()
	ctx := s.newContext()

	s.Run("should sort items by severity ascending", func() {
		// given
		params := s.newListParameters(pagination, nil, []domain.SortField{
			{Field: "severity", Direction: "asc"},
		})

		// when
		list, _ := repository.ListAll(ctx, params)

		// then
		s.Greater(len(list.GetContent()), 1)
		content := list.GetContent()
		for i := 1; i < len(content); i++ {
			s.GreaterOrEqual(content[i].Level, content[i-1].Level)
		}
	})

	s.Run("should sort items by severity descending", func() {
		// given
		params := s.newListParameters(pagination, nil, []domain.SortField{
			{Field: "severity", Direction: "desc"},
		})

		// when
		list, err := repository.ListAll(ctx, params)

		// then
		s.Require().NoError(err)
		content := list.GetContent()
		for i := 1; i < len(content); i++ {
			s.LessOrEqual(content[i].Level, content[i-1].Level)
		}
	})
}

// --- Pagination Cases ---

func (s *ItemsRepositorySuite) TestListAllPagination() {
	s.Run("should paginate items correctly", func() {
		// given
		pagination := s.newPagination(1, 2)
		params := s.newListParameters(pagination, nil, nil)
		repository := s.newRepository()
		ctx := s.newContext()

		// when
		list, _ := repository.ListAll(ctx, params)

		// then
		s.LessOrEqual(len(list.GetContent()), 2)
	})
}

// --- Entry Point ---

func TestSQLXItemsRepository(t *testing.T) {
	suite.Run(t, new(ItemsRepositorySuite))
}
```

**Key points:**
- **No `t.Parallel()` in suites.** Integration tests share a database through `SetupTest`/`TearDownTest`, so they must run sequentially.
- **`suite.Run()`** groups related sub-tests within a single test method (e.g., `TestItemsSuccess`, `TestItemsError`).
- **Helper methods** on the suite (e.g., `createItem()`, `getRandomID()`, `newRepository()`) reduce repetition and keep tests focused on the scenario.
- **Builders** construct test entities with fluent APIs (e.g., `NewItemBuilder().WithID(0).Build()`).
- **Seeders** populate the database before each test, and `TearDownTest` cleans up after each test.
- **Test methods are grouped by concern:** success, error, filters, sorting, pagination.

---

# Go Project Structure

> **TL;DR:** Follow the domain/infrastructure layer separation. Use `go.mod` for dependency management. Place test files next to production code. Use the `test/` directory for shared test helpers, builders, and doubles.

## Overview

This page defines the standard directory layout and dependency management practices for all Go projects. The architecture follows the Backend Design specification, separating code into `domain` (contracts) and `infrastructure` (implementations) layers.

## Directory Structure

```
cmd/
  <app>/
    main.go                   application entry point
    dig.go                    DI container creation and injection functions
internal/
  container.go              top-level DI provider orchestrator
  domain/                   (contracts)
    commands/
      container.go            DI registration for commands (or no-op)
    entities/
      container.go            DI registration for entities (or no-op)
    repositories/
  infrastructure/           (implementations)
    controllers/
      container.go            DI registration for controllers
      mappers/
      requests/
      responses/
    repositories/             prefixed with the tool name; returns database models
      container.go            DI registration for repositories
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

| Directory                               | Purpose                                                    |
|-----------------------------------------|------------------------------------------------------------|
| `cmd/<app>/`                            | Application entry point and DI injection functions         |
| `internal/domain/commands/`             | Business logic implementations                             |
| `internal/domain/entities/`             | Framework-agnostic domain entities                         |
| `internal/domain/repositories/`         | Repository interface contracts                             |
| `internal/infrastructure/controllers/`  | HTTP controllers (request/response handling)               |
| `internal/infrastructure/repositories/` | Repository implementations with library-specific code      |
| `test/domain/builders/`                 | Builder pattern implementations for constructing test data |
| `test/domain/doubles/`                  | Test doubles (stubs, dummies, fakers, in-memory)           |

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
    go.uber.org/dig v1.18.0
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

| File            | Purpose                                                  |
|-----------------|----------------------------------------------------------|
| `go.mod`        | Module path and dependency declarations                  |
| `go.sum`        | Dependency checksums (auto-generated, must be committed) |
| `.golangci.yml` | golangci-lint configuration                              |
| `container.go`  | Dig provider registration (one per architectural layer)  |
| `.editorconfig` | Editor standardization                                   |

---

# Python

> **TL;DR:** Follow the Zen of Python -- prioritize readability, simplicity, and explicitness. Use [Black](https://black.readthedocs.io/) for formatting, [isort](https://pycqa.github.io/isort/) for imports, [Flake8](https://flake8.pycqa.org/) for linting, type hints on all functions, [Loguru](https://loguru.readthedocs.io/) for logging, [pytest](https://docs.pytest.org/) for testing, and [PDM](https://pdm-project.org/) for packaging.

## Overview

This series of pages outlines best practices for Python development, with detailed explanations of recommended approaches and the reasoning behind them. For the general baseline, refer to the Code Style guide. See the sub-pages for specific topics:

## The Zen of Python

The Zen of Python defines the core philosophy for writing Pythonic code. Access it anytime by running `import this` in a Python REPL:

```
Beautiful is better than ugly.
Explicit is better than implicit.
Simple is better than complex.
Complex is better than complicated.
Flat is better than nested.
Sparse is better than dense.
Readability counts.
Special cases aren't special enough to break the rules.
Although practicality beats purity.
Errors should never pass silently.
Unless explicitly silenced.
In the face of ambiguity, refuse the temptation to guess.
There should be one-- and preferably only one --obvious way to do it.
Although that way may not be obvious at first unless you're Dutch.
Now is better than never.
Although never is often better than *right* now.
If the implementation is hard to explain, it's a bad idea.
If the implementation is easy to explain, it may be a good idea.
Namespaces are one honking great idea -- let's do more of those!
```

---

# Python Conventions

> **TL;DR:** Use `snake_case` for file and function names, `PascalCase` for classes. Follow the `<Operation><Entity>` naming pattern for Commands and Controllers with an `execute` method. Use descriptive names and write comments that explain **why**, not just **what**.

## Overview

This document defines Python-specific naming conventions and component patterns. For the general baseline, refer to the Code Style guide. The architectural layers referenced here are defined in the Backend Design section.

## File Naming

All file names must use `snake_case`:

```
list_users_command.py     # Correct
listUsersCommand.py       # Wrong
ListUsersCommand.py       # Wrong
```

Python package names cannot contain dashes, so `ion-cannon` becomes `ion_cannon`.

## General Conventions

### Use Descriptive Names

Avoid over-abbreviated variable and class names. While short names are acceptable for prototypes or mathematical algorithms, production code must be readable by any team member.

**Poorly named:**
```python
r = requests.get("https://example.com")
rj = r.json()
rc = r.status_code
```

**Well named:**
```python
response = requests.get("https://example.com")
response_body = response.json()
status_code = response.status_code
```

A more complex comparison:

```python
# Cryptic
q1 = queue.Queue()
q2 = queue.Queue()

while True:
    p = q1.get() * q2.get()
```

```python
# Self-documenting
customer_bill_queue = queue.Queue()
tax_rate_queue = queue.Queue()

while True:
    after_tax_amount = customer_bill_queue.get() * tax_rate_queue.get()
```

### Write Meaningful Comments

Write comments that capture the **intent** behind non-obvious logic:

```python
# No context
for row in [r for r in results if r is not None]:
    parsed.append(dict(zip(field_name, row)))
```

```python
# Clear intent
# Filter out empty rows from the database query results
for row in [r for r in results if r is not None]:
    # Combine field names with row values into a dictionary
    # (similar to using a MySQLCursorDict cursor)
    parsed.append(dict(zip(field_name, row)))
```

## Entities

Entities are the core of the application. All business logic related to properties and fields belongs inside the entity.

**Entities must be free of any framework or external tool dependencies.** Do not use ORM-specific base classes or decorators inside entity definitions.

```python
class User:
    def __init__(self, id: int, name: str, email: str) -> None:
        self.id = id
        self.name = name
        self.email = email

    def is_active(self) -> bool:
        return self.email is not None
```

## Commands

| Element     | Pattern                              | Example                    |
|-------------|--------------------------------------|----------------------------|
| File name   | `<operation>_<entity>_command.py`    | `list_users_command.py`    |
| Class name  | `<Operation><Entity>Command`         | `ListUsersCommand`         |
| Method name | `execute`                            | `def execute(self) -> ...` |

**Notes:**
- Use plural entity names when the operation targets multiple entities.
- Use the standard [operations vocabulary](../../Code-Style.md#operations-vocabulary).
- Callbacks must reflect all possible controller responses.

```python
from dataclasses import dataclass
from typing import Callable

@dataclass
class ListUsersCommandListeners:
    on_success: Callable[[list], None]
    on_error: Callable[[Exception], None]

class ListUsersCommand:
    def __init__(self, repository: "UsersRepository") -> None:
        self._repository = repository

    def execute(self, listeners: ListUsersCommandListeners) -> None:
        try:
            users = self._repository.find_all()
            listeners.on_success(users)
        except Exception as e:
            listeners.on_error(e)
```

## Controllers

| Element     | Pattern                                 | Example                       |
|-------------|-----------------------------------------|-------------------------------|
| File name   | `<operation>_<entity>_controller.py`    | `list_users_controller.py`    |
| Class name  | `<Operation><Entity>Controller`         | `ListUsersController`         |
| Method name | `execute`                               | `def execute(self) -> ...`    |

## Services

This layer usage depends on the framework. When used, services encapsulate reusable business logic that does not fit into a single command.

| Element     | Pattern                    | Example              |
|-------------|----------------------------|----------------------|
| File name   | `<entity>_service.py`      | `users_service.py`   |
| Class name  | `<Entity>Service`          | `UsersService`       |

## Repositories

### Contract (Domain Layer)

| Element     | Pattern                       | Example                |
|-------------|-------------------------------|------------------------|
| File name   | `<entity>_repository.py`      | `users_repository.py`  |
| Class name  | `<Entity>Repository` (ABC)    | `UsersRepository`      |

```python
from abc import ABC, abstractmethod

class UsersRepository(ABC):
    @abstractmethod
    def find_by_id(self, user_id: int) -> "User":
        ...

    @abstractmethod
    def find_all(self) -> list["User"]:
        ...

    @abstractmethod
    def save(self, user: "User") -> "User":
        ...

    @abstractmethod
    def delete(self, user_id: int) -> None:
        ...
```

### Implementation (Infrastructure Layer)

| Element     | Pattern                                   | Example                          |
|-------------|-------------------------------------------|----------------------------------|
| File name   | `<library>_<entity>_repository.py`        | `sqlalchemy_users_repository.py` |
| Class name  | `<Library><Entity>Repository`             | `SQLAlchemyUsersRepository`      |

## Mappers

### Repository Mappers

| Element     | Pattern                 | Example             |
|-------------|-------------------------|---------------------|
| File name   | `<entity>_mapper.py`    | `user_mapper.py`    |
| Class name  | `<Entity>Mapper`        | `UserMapper`        |

### Controller Mappers

| Element              | Pattern                                      | Example                         |
|----------------------|----------------------------------------------|---------------------------------|
| File name (request)  | `<operation>_<entity>_request_mapper.py`     | `insert_user_request_mapper.py` |
| File name (response) | `<operation>_<entity>_response_mapper.py`    | `insert_user_response_mapper.py`|

## Models

Models reside exclusively in the infrastructure layer and represent DTOs for external data sources (databases, APIs, queues, etc.). Each model is prefixed with the name of the external tool it communicates with:

| Example              | Source             |
|----------------------|--------------------|
| `SQLAlchemyUser`     | PostgreSQL via SQLAlchemy |
| `ApiDocument`        | External API       |
| `RedisSession`       | Redis cache        |

## Dependency Injection

Use constructor injection. For frameworks that support it (e.g., FastAPI), leverage the built-in dependency injection system. For standalone applications, wire dependencies manually in the composition root.

---

# Python Formatting and Linting

> **TL;DR:** Use **Black** for code formatting, **isort** for import sorting, and **Flake8** for linting. This combination ensures consistent, readable code across all Python projects.

## Overview

Different developers have different coding habits -- variable naming, indentation style, import ordering, and more. Python's [PEP 8](https://peps.python.org/pep-0008/) standard defines a comprehensive set of rules for consistent code styling. You do not need to memorize PEP 8; tooling handles enforcement automatically.

## Recommended Toolchain

| Tool       | Purpose                                          | Customizability                                        |
|------------|--------------------------------------------------|--------------------------------------------------------|
| **Black**  | PEP 8-compliant code formatter                   | Minimal (by design -- ensures all code looks the same) |
| **isort**  | Import sorter (stdlib, third-party, application) | Moderate (supports Black-compatible formatting)        |
| **Flake8** | Linter for PEP 8 compliance and code quality     | High                                                   |
| YAPF       | PEP 8-compliant formatter by Google              | Very high (multiple base styles)                       |
| autopep8   | PEP 8 auto-fixer                                 | Moderate                                               |

### Why Black?

Black's strict, opinionated formatting means all code formatted by Black looks identical. This consistency reduces cognitive overhead during code reviews and makes unfamiliar codebases immediately readable. Black is also trivial to set up -- there are no configuration knobs to tune.

**YAPF** (Google-style formatting):
```python
records = get_records(args.key, feed, args.output, args.range,
                      args.field, args.extra, args.expiration)
```

**Black** formatting:
```python
records = get_records(
    args.key,
    feed,
    args.output,
    args.range,
    args.field,
    args.extra,
    args.expiration,
)
```

### Import Ordering: isort

isort organizes imports into three sections (standard library, third-party, and application-specific), sorts them lexicographically and case-insensitively, and supports Black-compatible formatting:

```python
import argparse
import math
import multiprocessing
from typing import Type, Union

import cv2
import ffmpeg
from loguru import logger
from rich.progress import (
    BarColumn,
    Progress,
    ProgressColumn,
)
from rich.text import Text

from . import __version__
from .decoder import VideoDecoder
from .encoder import VideoEncoder
```

## Editor Configuration

### Visual Studio Code

Recommended extensions and settings:

- [Python](https://marketplace.visualstudio.com/items?itemName=ms-python.python) (includes [Pylance](https://marketplace.visualstudio.com/items?itemName=ms-python.vscode-pylance))
- [Python Indent](https://marketplace.visualstudio.com/items?itemName=KevinRose.vsc-python-indent)
- [autoDocstring](https://marketplace.visualstudio.com/items?itemName=njpwerner.autodocstring)

```json
{
    "python.formatting.provider": "black",
    "python.languageServer": "Pylance",
    "python.linting.flake8Enabled": true,
    "python.analysis.typeCheckingMode": "strict"
}
```

### NeoVim

Recommended plugins:

- [nvie/vim-flake8](https://github.com/nvie/vim-flake8) -- Flake8 linting
- [psf/black](https://github.com/psf/black) -- Black formatter integration
- [neoclide/coc.nvim](https://github.com/neoclide/coc.nvim) -- Extension loader (use `:CoCInstall coc-pyright` for Python support)
- [dense-analysis/ale](https://github.com/dense-analysis/ale) -- Asynchronous linting engine
- [kkoomen/vim-doge](https://github.com/kkoomen/vim-doge) -- Documentation generation

ALE configuration for Black:
```vim
let g:ale_fixers = {
\ '*': ['trim_whitespace', 'prettier'],
\ 'python': ['black']
\}
let g:ale_completion_enabled = 1
```

---

# Python Type System

> **TL;DR:** Add type hints to all function parameters, return types, and non-obvious variables. Type hints enable static analysis, catch bugs before runtime, and make code self-documenting.

## Overview

Python is dynamically typed, which means type errors often surface only at runtime -- or worse, in production. [PEP 484](https://peps.python.org/pep-0484/) introduced type hints, allowing developers to declare expected types for variables and function signatures. Modern linters and IDEs interpret these annotations to warn about type mismatches during development, significantly reducing debugging time.

## Type Annotation Requirements

**All function parameters and return types must have type hints.** Variables should be annotated when the type is not obvious from the assignment.

### Why Type Hints Matter

Consider the following function without type hints:

```python
def collatz(n):
    if n % 2 == 0:
        return n / 2
    else:
        return 3 * n + 1

print(collatz(5).denominator)   # Works (int has .denominator)
print(collatz(6).denominator)   # Fails (float has no .denominator)
```

Without type annotations, this bug is invisible until execution. With type hints, the linter catches the mismatch immediately:

```python
def collatz(n: int) -> int:
    if n % 2 == 0:
        return n / 2    # Linter warns: float is incompatible with int
    else:
        return 3 * n + 1
```

## Syntax

### Functions

```python
def function(param: <type>) -> <return_type>:
    ...
```

Example:

```python
import pathlib

def delete_file(path: pathlib.Path) -> int:
    """Delete the file at the given path.

    :param path: Path to the file to delete.
    :return: 0 if deleted successfully, 1 otherwise.
    """
    if path.is_file():
        path.unlink()
        return 0
    return 1
```

### Variables

```python
<variable>: <type> = <value>
```

Example:

```python
import queue

item: int = 1
some_queue: queue.Queue[int] = queue.Queue()
some_queue.put(item)
```

## Key Type Patterns

### The `typing` Module

The `typing` module supports advanced type constructs such as `Union`, `Optional`, `Type`, and generics:

```python
from typing import Type, Union

def func(param: Type[Union[ClassA, ClassB]]) -> None:
    ...
```

Alternative syntax:

```python
from typing import Type, Union

def func(param: Union[Type[ClassA], Type[ClassB]]) -> None:
    ...
```

### Python 3.10+ Union Syntax

For Python 3.10+, use the built-in `|` operator instead of `Union`:

```python
def func(param: type[ClassA | ClassB]) -> None:
    ...
```

## Prohibited Patterns

```python
# Wrong -- missing type hints on function signatures
def process(data):
    return data.transform()

# Wrong -- using Any as a catch-all
from typing import Any
def process(data: Any) -> Any:
    return data.transform()
```

Define proper types or use `Protocol` for structural typing instead.

---

# Python Logging

> **TL;DR:** Use **[Loguru](https://loguru.readthedocs.io/)** for all application logging. Do not use the standard `logging` module or `print()` for logging in production code. Always import as `from loguru import logger`. Use STDOUT for normal output and STDERR for warnings and errors.

## Overview

Programs frequently need to output messages for progress tracking, state reporting, and error diagnostics. This page defines the mandatory logging library and patterns for all Python projects.

## Mandatory Library: Loguru

**Use [Loguru](https://loguru.readthedocs.io/) for all application logging.** Do not use Python's built-in `logging` module. Loguru provides a simpler API, colorized output, structured exception formatting, and sensible defaults with minimal configuration.

### Installation

```bash
pip install loguru
# or via PDM
pdm add loguru
```

### Import Convention

```python
from loguru import logger
```

### Usage

```python
from loguru import logger

logger.info("application started")
logger.warning("disk usage above threshold")
logger.error("failed to connect to database")
```

## Log Levels

| Level    | Severity | Method              | When to Use                                               |
|----------|----------|---------------------|-----------------------------------------------------------|
| TRACE    | 5        | `logger.trace()`    | Very fine-grained diagnostic information                  |
| DEBUG    | 10       | `logger.debug()`    | Diagnostic information useful during development          |
| INFO     | 20       | `logger.info()`     | General operational events (application started, request processed) |
| SUCCESS  | 25       | `logger.success()`  | Successful completion of a significant operation          |
| WARNING  | 30       | `logger.warning()`  | Potential issues that do not prevent operation             |
| ERROR    | 40       | `logger.error()`    | Errors that prevent a specific operation but not the application |
| CRITICAL | 50       | `logger.critical()` | Critical errors that require immediate attention          |

## Structured Logging

### Custom Format

The default format can be verbose. A recommended concise format:

```python
import sys
from loguru import logger

LOGURU_FORMAT = (
    "<green>{time:HH:mm:ss.SSSSSS!UTC}</green> "
    "<level>{level: <8}</level> "
    "<level>{message}</level>"
)
logger.remove()
logger.add(sys.stderr, colorize=True, format=LOGURU_FORMAT)
```

### Exception Formatting

Loguru's `logger.exception()` produces colorized, structured exception output with variable values labeled inline, making debugging significantly faster:

```python
try:
    risky_operation()
except Exception:
    logger.exception("operation failed with unexpected error")
```

## Printing

The `print()` function writes to STDOUT by default. It is acceptable for simple scripts and CLI output, but **must not be used as a substitute for logging** in application code. In production, ensure messages are directed to the appropriate output stream:

```python
import sys
print("Fetching data from Redis")
print("Error: Could not connect to Redis", file=sys.stderr)
```

### String Formatting

Use f-strings (preferred) or `.format()` for variable interpolation:

```python
number = 42
print(f"The number is: {number}")
print("The number is: {}".format(number))
```

## Prohibited Patterns

The following patterns must **not** be used for application logging:

```python
# Wrong -- standard library logging module
import logging
logging.info("something happened")
logging.getLogger(__name__).warning("issue detected")

# Wrong -- print for application logging
print("Processing request...")

# Wrong -- f-string logging without Loguru
print(f"Error: {error_message}")
```

**Correct equivalent:**

```python
from loguru import logger

logger.info("something happened")
logger.warning("issue detected")
logger.info("processing request")
logger.error(f"error: {error_message}")
```

## Why Not the Standard `logging` Module?

The standard `logging` module is part of Python's standard library and provides basic functionality. However, it requires significant boilerplate to configure properly (formatters, handlers, log levels per module), its default output is unstructured, and it lacks features like colorized output and automatic exception variable labeling. Loguru provides all of these out of the box with a single import.

For reference, the standard library supports five levels: DEBUG, INFO, WARNING, ERROR, and CRITICAL. Full details are in the [logging HOWTO](https://docs.python.org/3/howto/logging.html). However, **this is provided for awareness only -- always use Loguru in our projects.**

---

# Python Testing

> **TL;DR:** Use **pytest** for all testing. All tests must follow the BDD pattern with `# given`, `# when`, `# then` comment blocks. Place tests in a `/tests/` directory mirroring the source structure. Use fixtures for setup/teardown and factories for test data construction.

## Overview

This document defines the testing conventions for all Python projects. For the cross-language testing standards (BDD structure, test doubles, builders), refer to the Tests page.

## Testing Framework

| Tool                                                 | Purpose                                 |
|------------------------------------------------------|-----------------------------------------|
| [pytest](https://docs.pytest.org/)                   | Testing framework and runner            |
| [pytest-cov](https://pytest-cov.readthedocs.io/)    | Coverage reporting                      |
| [Faker](https://faker.readthedocs.io/)               | Fake data generation library            |

### Installation

```bash
pdm add -dG test pytest pytest-cov faker
```

## BDD Structure (Given / When / Then)

**Every test must use `# given`, `# when`, `# then` comment blocks** to clearly separate preconditions, actions, and assertions. This is mandatory across all test files.

```python
def test_should_return_user_when_found():
    # given
    repository = InMemoryUsersRepository()
    repository.save(User(id=1, name="John", email="john@test.com"))
    command = GetUserCommand(repository)

    # when
    result = command.execute(user_id=1)

    # then
    assert result is not None
    assert result.name == "John"
```

## File Placement

Tests live in a `/tests/` directory at the project root, mirroring the source code structure:

```
project/
  src/
    domain/
      commands/
        list_users_command.py
      entities/
        user.py
      repositories/
        users_repository.py
    infrastructure/
      controllers/
        list_users_controller.py
      repositories/
        sqlalchemy_users_repository.py
  tests/
    domain/
      commands/
        test_list_users_command.py
      builders/
        user_builder.py
      doubles/
        in_memory_users_repository.py
    infrastructure/
      controllers/
        test_list_users_controller.py
      repositories/
        test_sqlalchemy_users_repository.py
    conftest.py
```

## File Naming

All test files must use the `test_` prefix:

```
test_list_users_command.py     # Correct
list_users_command_test.py     # Wrong
test_list_users.py             # Wrong (must match the source file name)
```

## Unit Tests

Unit tests validate individual components in isolation. Use test doubles (stubs, dummies, fakers, in-memory implementations) to replace external dependencies.

### Command Tests

```python
import pytest
from domain.commands import DeleteUserCommand
from tests.domain.doubles import InMemoryUsersRepository, UsersRepositoryStub

class TestDeleteUserCommand:
    def test_should_call_on_success_when_user_is_deleted(self):
        # given
        repository = InMemoryUsersRepository()
        repository.save(User(id=1, name="John", email="john@test.com"))
        command = DeleteUserCommand(repository)
        success_called = False

        def on_success():
            nonlocal success_called
            success_called = True

        # when
        command.execute(user_id=1, on_success=on_success, on_error=lambda e: None)

        # then
        assert success_called is True

    def test_should_call_on_error_when_user_not_found(self):
        # given
        repository = InMemoryUsersRepository()
        command = DeleteUserCommand(repository)
        error_received = None

        def on_error(err):
            nonlocal error_received
            error_received = err

        # when
        command.execute(user_id=999, on_success=lambda: None, on_error=on_error)

        # then
        assert error_received is not None
```

### Controller Tests

```python
from fastapi.testclient import TestClient

class TestListUsersController:
    def test_should_respond_200_when_users_listed_successfully(self, client: TestClient):
        # when
        response = client.get("/users")

        # then
        assert response.status_code == 200

    def test_should_respond_500_when_command_fails(self, client_with_error: TestClient):
        # when
        response = client_with_error.get("/users")

        # then
        assert response.status_code == 500
```

## Integration Tests

Integration tests validate components working together with real infrastructure (databases, external services). Use pytest fixtures for setup and teardown.

### Repository Tests

```python
import pytest
from faker import Faker
from domain.entities import Item
from infrastructure.repositories import SQLAlchemyItemsRepository
from tests.domain.builders import ItemBuilder

fake = Faker()

@pytest.fixture
def repository(db_session):
    return SQLAlchemyItemsRepository(db_session)

@pytest.fixture(autouse=True)
def seed_data(db_session):
    """Seed the database before each test and clean up after."""
    # setup
    db_session.execute("INSERT INTO items (name, external_id) VALUES ('test', 'ext-1')")
    db_session.commit()
    yield
    # teardown
    db_session.execute("DELETE FROM items")
    db_session.commit()

class TestSQLAlchemyItemsRepository:
    # --- Success Cases ---

    def test_should_save_a_new_item_successfully(self, repository):
        # given
        item = ItemBuilder().with_external_id(fake.uuid4()).build()

        # when
        created = repository.save(item)

        # then
        assert created.id is not None
        assert created.id > 0

    def test_should_list_all_items_successfully(self, repository):
        # given / when
        items = repository.find_all()

        # then
        assert len(items) >= 1

    def test_should_get_item_by_id_successfully(self, repository, seed_data):
        # given
        items = repository.find_all()
        item_id = items[0].id

        # when
        item = repository.find_by_id(item_id)

        # then
        assert item is not None
        assert item.id == item_id

    # --- Error Cases ---

    def test_should_raise_when_item_not_found(self, repository):
        # given
        non_existent_id = 99999

        # when / then
        with pytest.raises(ItemNotFoundError):
            repository.find_by_id(non_existent_id)
```

## Test Doubles

Use the following test double types, in order of preference:

| Type       | Description                                       | When to Use                                    |
|------------|---------------------------------------------------|------------------------------------------------|
| Stub       | Returns pre-configured (canned) answers           | When you need controlled return values         |
| Dummy      | Fills required parameters, never actually used    | When a parameter is required but irrelevant    |
| In-memory  | Implements logic without external dependencies    | When you need realistic behavior without I/O   |
| Faker      | Generates realistic fake data                     | When you need varied, realistic test inputs    |
| **Mock**   | Records and verifies calls                        | **Avoid when possible** -- prefer stubs        |

## Builders

Use the builder pattern for constructing complex test objects:

```python
class UserBuilder:
    def __init__(self):
        self._id = 1
        self._name = "Default User"
        self._email = "default@test.com"

    def with_id(self, id: int) -> "UserBuilder":
        self._id = id
        return self

    def with_name(self, name: str) -> "UserBuilder":
        self._name = name
        return self

    def with_email(self, email: str) -> "UserBuilder":
        self._email = email
        return self

    def build(self) -> User:
        return User(id=self._id, name=self._name, email=self._email)
```

Usage:

```python
user = UserBuilder().with_id(0).with_name("John").with_email("john@test.com").build()
```

---

# Python Project Structure

> **TL;DR:** Use **PDM** as the package manager with **PEP 621** metadata in `pyproject.toml`. Follow the conventional directory layout: `/<package_name>/` for source code, `/tests/` for tests, `/examples/` for usage examples. Avoid legacy `setup.py` and `setup.cfg` formats.

## Overview

This page defines the standard directory layout, package management, and distribution practices for all Python projects. It consolidates directory structure, packaging with PDM, and metadata format guidelines.

## Directory Structure

```
project/
  <package_name>/           source code
    __init__.py
    __main__.py             enables `python -m <package_name>`
  examples/                 example scripts
  tests/                    external tests (pytest)
    conftest.py
  .editorconfig
  Dockerfile
  MANIFEST.in               additional files for source distributions
  pyproject.toml            PEP 518 + PEP 621 metadata (single source of truth)
  README.md
```

### Key Directories and Files

| Path                   | Purpose                                                            |
|------------------------|--------------------------------------------------------------------|
| `/<package_name>/`     | Main source code. Replace with the actual package name (e.g., `icarus`). Python names cannot contain dashes, so `ion-cannon` becomes `ion_cannon`. |
| `/examples/`           | Example scripts demonstrating how to use the library or application. |
| `/tests/`              | External tests and test data. Use [pytest](https://docs.pytest.org/) to run all tests. |
| `/Dockerfile`          | Container definition. For multiple variants: `Dockerfile.alpine`, `Dockerfile.ubuntu`, `Dockerfile.slim-debian`. |
| `/MANIFEST.in`         | Specifies additional files to include in source distributions (sdist). |
| `/pyproject.toml`      | PEP 518 + PEP 621-compliant metadata. Single source of truth for project configuration. |
| `/README.md`           | Project description and usage instructions. Acceptable formats: Markdown, reStructuredText, or plain text. |

## Package Manager: PDM

**Use [PDM](https://pdm-project.org/) as the package manager.** PDM supports PEP 621 natively, is lightweight (unlike Poetry, it does not require GCC-compiled components), and can be replaced by any PEP 621-compatible tool.

### Tool Comparison

| Tool       | Description                                                                      |
|------------|----------------------------------------------------------------------------------|
| Conda      | Cross-platform binary package manager (primarily for Anaconda)                   |
| Setuptools | Fully-featured, legacy package builder                                           |
| Pipenv     | Virtualenv and dependency manager                                                |
| Poetry     | Dependency manager + builder + publisher (uses custom `pyproject.toml` sections) |
| **PDM**    | Dependency manager + builder + publisher with full PEP 621 support               |

### Installation

```bash
pip install -U pdm
```

### Initialization

```bash
mkdir icarus && cd icarus
pdm init
```

Answer the interactive prompts according to your package's requirements. Use `Proprietary` as the license for non-open-source packages.

## Dependency Management

```bash
# Add a dependency
pdm add 'opencv-python>=4.5'

# Add a development dependency
pdm add -dG test pytest pytest-cov

# Install all dependencies
pdm install

# Update all dependencies
pdm update

# Remove a dependency
pdm remove opencv-python
```

## Build & Distribution

```bash
# Install the package locally
pip install .

# Build source distribution
python3 -m build -s .

# Build wheel distribution
python3 -m build -w .
```

### Docker

```bash
# Build a specific Dockerfile variant
docker build -f Dockerfile.alpine -t app:1.0.0-alpine .
```

## Key Configuration Files

### pyproject.toml (PEP 621)

**PEP 621 is the recommended standard** because it consolidates all metadata into a single, readable file. Avoid legacy `setup.py` and `setup.cfg` formats.

```toml
[project]
name = "icarus"
description = "A demo package"
authors = [{ name = "Author", email = "author@example.com" }]
dependencies = []
requires-python = ">=3.10"
license = { text = "Proprietary" }
dynamic = ["version"]

[project.urls]
Homepage = "https://github.com/org/icarus"

[tool.pdm]
version = { from = "icarus/__init__.py" }

[build-system]
requires = ["pdm-pep517"]
build-backend = "pdm.pep517.api"
```

Enable dynamic versioning so the version is defined once in `__init__.py`:

```python
# icarus/__init__.py
__version__ = "1.0.0"
```

### Package Entry Points

Create the `__main__` file to enable `python -m icarus`:

```python
# icarus/__main__.py
import sys
from .icarus import main

if __name__ == "__main__":
    sys.exit(main())
```

### MANIFEST.in

Specifies additional files to include in source distributions:

```text
recursive-include models *
include somedir/README.md
```

### Metadata Format Evolution

For historical context, Python package metadata has evolved through several formats:

| Format       | Standard  | Description                                                  |
|--------------|-----------|--------------------------------------------------------------|
| `setup.py`   | Legacy    | Python script executed during build; not declarative         |
| `setup.cfg`  | PEP 517   | Declarative config, but separate from build system config    |
| `pyproject.toml` | PEP 621 + PEP 518 | Consolidated, readable, single-file metadata (recommended) |

---

# Java

> **TL;DR:** Use [Spring Boot](https://spring.io/projects/spring-boot) as the application framework, [Google Java Format](https://github.com/google/google-java-format) (via [Spotless](https://github.com/diffplug/spotless)) for formatting, [Checkstyle](https://checkstyle.sourceforge.io/) + [PMD](https://pmd.github.io/) for linting, [MapStruct](https://mapstruct.org/) for object mapping, [SLF4J](https://www.slf4j.org/) with [Logback](https://logback.qos.ch/) for logging, [JUnit 5](https://junit.org/junit5/) + [Mockito](https://site.mockito.org/) for testing, and [Gradle](https://gradle.org/) for builds.

## Overview

This series of pages outlines Java-specific conventions, with detailed explanations of recommended approaches and the reasoning behind them. For the general baseline, refer to the Code Style guide. See the sub-pages for specific topics:

## Liquibase

All database versioning is managed through [Liquibase](https://www.liquibase.org/). Follow these mandatory conventions:

| Rule                 | Description                                                                |
|----------------------|----------------------------------------------------------------------------|
| File format          | Always use `.yaml` (not `.yml`)                                            |
| Quoting              | Always use single quotes                                                   |
| Naming               | Always use `snake_case`                                                    |
| Rollbacks            | Every changeset must include a rollback statement                          |
| Pre-conditions       | Every changeset must include pre-conditions to validate the database state |
| Boolean/Time columns | Use descriptive names: `created_at`, `enabled`, `updated_at`, `activated`  |

### Column Definition Order

Constraints must appear **before** the column name in column definitions:

```yaml
# Correct
- column:
    constraints:
        nullable: false
    name: 'email'
    type: 'VARCHAR(255)'

# Wrong
- column:
    name: 'email'
    type: 'VARCHAR(255)'
    constraints:
        nullable: false
```

### Constraint Naming

| Type                   | Pattern                                   | Example                               |
|------------------------|-------------------------------------------|---------------------------------------|
| Unique Index (UIX)     | `<table>_<column>_uix`                    | `user_name_uix`                       |
| Foreign Key (FK)       | `<origin_table>_<destination_table>_fkey` | `user_user_contract_information_fkey` |
| Primary Key (PK)       | `<table>_pkey`                            | `user_service_pkey`                   |
| Materialized View (MV) | `<view>_uix`                              | `rich_user_ix`                        |

### Database Versioning

Follow [Semantic Versioning](https://semver.org/) for database changelogs:

| Change Type                | Version Update |
|----------------------------|----------------|
| New table                  | `X.Y+1.0`      |
| Alter existing table       | `X.Y.Z+1`      |
| Breaking structural change | `X+1.0.0`      |

Packages published for testing must use the `-SNAPSHOT` suffix (e.g., `3.0.0-SNAPSHOT`). Snapshot packages must **never** be used in production.

### Complete Liquibase Example

```yaml
databaseChangeLog:
  - changeSet:
      id: '00-create-table-example'
      author: 'author.name'
      context: 'prod, test'
      preConditions:
        - not:
            tableExists:
              tableName: 'example'
      changes:
        - createTable:
            tableName: 'example'
            columns:
              - column:
                  autoIncrement: true
                  constraints:
                    nullable: false
                    primaryKey: true
                    primaryKeyName: 'example_pkey'
                  name: 'id'
                  type: 'BIGINT'
              - column:
                  constraints:
                    nullable: false
                    unique: true
                    uniqueConstraintName: 'example_name_uix'
                  name: 'name'
                  type: 'VARCHAR(50)'
              - column:
                  name: 'active'
                  type: 'BOOLEAN'
              - column:
                  constraints:
                    nullable: false
                  name: 'created_at'
                  type: 'TIMESTAMP'
      rollback:
        - dropTable:
            tableName: 'example'
```

## Java Design Principles

The following principles guide Java development across all projects:

- **Program to interfaces, not implementations.** Depend on abstractions (interfaces) rather than concrete classes. This enables testability and loose coupling.
- **Favor composition over inheritance.** Build complex behavior by composing objects rather than through deep class hierarchies.
- **Immutability by default.** Use `final` fields, records, and immutable collections wherever possible. Mutable state should be the exception, not the rule.
- **Fail fast.** Validate inputs early and throw meaningful exceptions rather than propagating invalid state through the system.
- **Convention over configuration.** Leverage Spring Boot's auto-configuration and opinionated defaults rather than writing boilerplate.

---

# Java Conventions

> **TL;DR:** Use `PascalCase` for class names and follow the strict `<Operation><Entity>` naming patterns for Commands, Controllers, Services, Repositories, and Mappers. Entities must be framework-agnostic. Use Spring constructor injection for dependency injection and [MapStruct](https://mapstruct.org/) for object mapping.

## Overview

This document defines Java-specific naming conventions and component patterns. For the general baseline, refer to the Code Style guide. The architectural layers referenced here are defined in the Backend Design section.

## File Naming

Java enforces that the file name matches the public class name:

```
ListItemsCommand.java       # Correct -- matches class ListItemsCommand
list_items_command.java      # Wrong -- Java uses PascalCase file names
```

## General Conventions

1. Use **constructor injection** for all dependencies -- never use field injection (`@Autowired` on fields).
2. Mark injected fields as `final` to enforce immutability.
3. Use [Lombok](https://projectlombok.org/) to reduce boilerplate (`@Getter`, `@Setter`, `@NoArgsConstructor`, `@RequiredArgsConstructor`, `@SuperBuilder`).
4. Use Java **records** for immutable DTOs, domain events, and value objects.
5. Use `@Component` for commands and listeners, `@RestController` for controllers, `@Repository` for repository implementations, and `@Service` for infrastructure service implementations.
6. For an introduction to the DTO pattern, refer to [this article](https://www.baeldung.com/java-dto-pattern).

## Entities

Entities are the core of the application. All business logic related to properties and fields belongs inside the entity.

**Entities must be free of any persistence or framework annotations.** Do not use `@Entity`, `@Table`, `@Column`, or any Jakarta/JPA annotations inside domain entities. Persistence concerns belong in infrastructure models.

```java
public final class Item<A> {
    private final ItemCode code;
    private final ItemSeverity severity;
    private Boolean isCountable;
    private Long amountAffected;
    private List<A> affected;

    public Boolean hasData() {
        return amountAffected > 0 || !isCountable;
    }
}
```

Use Java records for simple domain events and value objects:

```java
public record ItemEvent(
        Long organizationId,
        String referenceGuid,
        ItemCode code,
        ItemCategory category) {}
```

## Commands

| Element     | Pattern                           | Example                    |
|-------------|-----------------------------------|----------------------------|
| File name   | `<Operation><Entity>Command.java` | `InsertItemCommand.java`   |
| Class name  | `<Operation><Entity>Command`      | `InsertItemCommand`        |
| Method name | `execute`                         | `public void execute(...)` |

**Notes:**
- Use plural entity names when the operation targets multiple entities.
- Use the standard [operations vocabulary](../../Code-Style.md#operations-vocabulary).
- Commands must define a `Listeners` record for all possible outcomes (callback pattern).

```java
@Component
@RequiredArgsConstructor
public class InsertItemCommand {
    private final SendItemService service;

    public void execute(final ItemEvent event, final Listeners listeners) {
        try {
            service.send(event);
            listeners.onSuccess().run();
        } catch (Exception e) {
            listeners.onError().accept(e);
        }
    }

    public record Listeners(Runnable onSuccess, Consumer<Exception> onError) {}
}
```

## Controllers

| Element     | Pattern                              | Example                                 |
|-------------|--------------------------------------|-----------------------------------------|
| File name   | `<Operation><Entity>Controller.java` | `ListItemsController.java`              |
| Class name  | `<Operation><Entity>Controller`      | `ListItemsController`                   |
| Method name | `execute`                            | `public ResponseEntity<?> execute(...)` |

```java
@RestController
@RequiredArgsConstructor
@RequestMapping("${api.v1-prefix}")
public class InsertItemController {
    private final InsertItemCommand command;

    @PostMapping("/items")
    public ResponseEntity<?> execute(
            @RequestHeader("organization-id") final Long organizationId,
            @Valid @RequestBody final InsertItemRequest request) {
        final var entity = InsertItemRequestMapper.INSTANCE.mapToEntity(organizationId, request);
        final var ref = new AtomicReference<ResponseEntity<?>>();

        final var listeners = new InsertItemCommand.Listeners(
                () -> ref.set(ResponseEntity.status(HttpStatus.CREATED).build()),
                (error) -> ref.set(ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).build()));
        command.execute(entity, listeners);
        return ref.get();
    }
}
```

## Services

The Services layer provides an abstraction between Commands and Repositories. The **domain layer** defines the service contract (interface), and the **infrastructure layer** provides the implementation.

### Contract (Domain Layer)

| Element        | Pattern                           | Example                |
|----------------|-----------------------------------|------------------------|
| File name      | `<Operation><Entity>Service.java` | `SendItemService.java` |
| Interface name | `<Operation><Entity>Service`      | `SendItemService`      |

```java
public interface SendItemService {
    void send(ItemEvent event);
}
```

### Implementation (Infrastructure Layer)

| Element    | Pattern                              | Example                   |
|------------|--------------------------------------|---------------------------|
| File name  | `Jpa<Operation><Entity>Service.java` | `JpaSendItemService.java` |
| Class name | `Jpa<Operation><Entity>Service`      | `JpaSendItemService`      |

```java
@Service
@RequiredArgsConstructor
public class JpaSendItemService implements SendItemService {
    private final ItemsRepository repository;
    private final ItemMapper mapper;

    @Override
    public void send(final ItemEvent event) {
        final var model = mapper.toModel(event);
        repository.save(model);
    }
}
```

**Note:** The `Jpa` prefix indicates the infrastructure tool used. If using a different technology (e.g., REST client, gRPC), use the corresponding prefix (e.g., `RestSendItemService`, `GrpcSendItemService`).

## Repositories

### Contract (Domain Layer)

| Element        | Pattern                   | Example                |
|----------------|---------------------------|------------------------|
| File name      | `<Entity>Repository.java` | `ItemsRepository.java` |
| Interface name | `<Entity>Repository`      | `ItemsRepository`      |

```java
public interface ItemsRepository {
    Boolean hasItems(Long categoryId);
    Page<JpaItem> findAllByCategoryId(Long categoryId, Pageable pageable);
}
```

### Implementation (Infrastructure Layer)

| Element    | Pattern                      | Example                   |
|------------|------------------------------|---------------------------|
| File name  | `Jpa<Entity>Repository.java` | `JpaItemsRepository.java` |
| Class name | `Jpa<Entity>Repository`      | `JpaItemsRepository`      |

For Spring Data JPA repositories:

```java
@Repository
public interface JpaItemsRepository
        extends JpaRepository<JpaItem, Long>, ItemsRepository {}
```

### QueryDSL Repositories

For complex queries that go beyond Spring Data JPA:

| Element    | Pattern                           | Example                        |
|------------|-----------------------------------|--------------------------------|
| File name  | `QueryDsl<Entity>Repository.java` | `QueryDslItemsRepository.java` |
| Class name | `QueryDsl<Entity>Repository`      | `QueryDslItemsRepository`      |

## Mappers

Use [MapStruct](https://mapstruct.org/) for all object mapping. Mappers isolate layers from each other and prevent hard coupling to frameworks.

### Repository Mappers

| Element   | Pattern               | Example           |
|-----------|-----------------------|-------------------|
| File name | `<Entity>Mapper.java` | `ItemMapper.java` |
| Interface | `<Entity>Mapper`      | `ItemMapper`      |

```java
@Mapper(unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface ItemMapper {
    ItemMapper INSTANCE = Mappers.getMapper(ItemMapper.class);

    Item toEntity(JpaItem model);
    JpaItem toModel(Item entity);
}
```

### Controller Mappers

| Element              | Pattern                                        | Example                           |
|----------------------|------------------------------------------------|-----------------------------------|
| File name (request)  | `<Operation><Entity>RequestMapper.java`        | `InsertItemRequestMapper.java`    |
| File name (response) | `<Operation><Entity>ResponseMapper.java`       | `InsertItemResponseMapper.java`   |

```java
@Mapper(unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface InsertItemRequestMapper {
    InsertItemRequestMapper INSTANCE = Mappers.getMapper(InsertItemRequestMapper.class);

    default ItemEvent mapToEntity(final Long organizationId, final InsertItemRequest request) {
        return new ItemEvent(
                organizationId,
                request.referenceGuid(),
                request.code(),
                request.category());
    }
}
```

## Models

Models reside exclusively in the infrastructure layer and represent JPA entities for database persistence. They resemble domain entities but are **not** domain entities.

Each model is prefixed with `Jpa`:

| Element    | Pattern            | Example        |
|------------|--------------------|----------------|
| File name  | `Jpa<Entity>.java` | `JpaItem.java` |
| Class name | `Jpa<Entity>`      | `JpaItem`      |

```java
@Entity
@Getter
@Setter
@SuperBuilder
@NoArgsConstructor
@Table(name = "item")
public class JpaItem {
    @Id
    @GeneratedValue(strategy = GenerationType.SEQUENCE)
    private Long id;

    @CreationTimestamp
    @Column(nullable = false)
    private LocalDateTime createdAt;

    @Column(nullable = false)
    private String name;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "category_id")
    private JpaCategory category;
}
```

## Listeners

Listeners handle asynchronous messages from queues (e.g., Kafka). They act as infrastructure-layer entry points, similar to controllers but for message-driven operations.

| Element     | Pattern                            | Example                      |
|-------------|------------------------------------|------------------------------|
| File name   | `Process<Entity>Listener.java`     | `ProcessItemListener.java`   |
| Class name  | `Process<Entity>Listener`          | `ProcessItemListener`        |

```java
@Component
@RequiredArgsConstructor
public class ProcessItemListener {
    private final InsertItemCommand command;

    @KafkaListener(
            topics = "items-topic",
            groupId = "${queue.items.groupId}")
    @RetryableTopic(attempts = "1")
    public void queueListener(
            @Payload final ItemMessage message,
            final Acknowledgment acknowledgment) {
        acknowledgment.acknowledge();
        final var event = message.toDomain();
        command.execute(event, listeners);
    }
}
```

### Messages

| Element     | Pattern                  | Example            |
|-------------|--------------------------|--------------------|
| File name   | `<Entity>Message.java`   | `ItemMessage.java` |
| Class name  | `<Entity>Message`        | `ItemMessage`      |

```java
public record ItemMessage(
        @JsonProperty("organization_id") Long organizationId,
        @JsonProperty("reference_guid") String referenceGuid,
        @JsonProperty("code") String code) {

    public ItemEvent toDomain() {
        return new ItemEvent(organizationId, referenceGuid, ItemCode.valueOf(code), null);
    }
}
```

## Dependency Injection

Use **Spring constructor injection** for all dependency wiring. Lombok's `@RequiredArgsConstructor` generates the constructor automatically from `final` fields.

```java
// Correct -- constructor injection via Lombok
@Component
@RequiredArgsConstructor
public class InsertItemCommand {
    private final SendItemService service;
    private final ItemsRepository repository;
}

// Wrong -- field injection
@Component
public class InsertItemCommand {
    @Autowired
    private SendItemService service;
}
```

Constructor injection ensures:
- Dependencies are explicit and visible.
- Objects are fully initialized upon creation.
- Fields can be declared `final`, enforcing immutability.
- Unit tests can inject dependencies without Spring context.

---

# Java Formatting and Linting

> **TL;DR:** Use **[Google Java Format](https://github.com/google/google-java-format)** (via the [Spotless](https://github.com/diffplug/spotless) Gradle plugin) for code formatting, **[Checkstyle](https://checkstyle.sourceforge.io/)** for style enforcement, and **[PMD](https://pmd.github.io/)** for static code analysis. These tools are non-negotiable and must be integrated into every project's CI pipeline.

## Overview

Java's ecosystem does not include a built-in formatter like Go, so the team standardizes on Google Java Format. Combined with Checkstyle for style rules and PMD for bug detection, this toolchain ensures consistent, high-quality code across all projects.

## Formatter: Google Java Format (via Spotless)

[Google Java Format](https://github.com/google/google-java-format) produces a single canonical formatting for any Java source file. It is applied via the [Spotless](https://github.com/diffplug/spotless) Gradle plugin, which also handles import ordering and license header enforcement.

### Gradle Configuration

```groovy
plugins {
    id 'com.diffplug.spotless' version '6.25.0'
}

spotless {
    java {
        googleJavaFormat('1.22.0')
        removeUnusedImports()
        trimTrailingWhitespace()
        endWithNewline()
    }
}
```

### Usage

```bash
# Auto-format all Java files
gradle spotlessApply

# Check formatting without modifying files (CI mode)
gradle spotlessCheck
```

**Do not use alternative formatters.** Google Java Format is the team standard and all Java code must be formatted with it.

## Import Ordering

Spotless handles import ordering automatically when using Google Java Format. The standard grouping is:

1. `com.*` -- third-party packages
2. `java.*` -- standard library
3. `javax.*` / `jakarta.*` -- extension libraries
4. Static imports (last)

Each group is separated by a blank line. Unused imports are removed automatically.

## Linter: Checkstyle

[Checkstyle](https://checkstyle.sourceforge.io/) enforces coding standards and style rules. The project uses the Google ruleset as a baseline with project-specific customizations.

### Gradle Configuration

```groovy
plugins {
    id 'checkstyle'
}

checkstyle {
    toolVersion = '10.12.1'
    configFile = file('src/main/resources/app/quality/checkstyle-google-ruleset.xml')
    maxWarnings = 0
    maxErrors = 0
}
```

### Usage

```bash
# Run Checkstyle on main sources
gradle checkstyleMain

# Run Checkstyle on test sources
gradle checkstyleTest
```

### Key Rules Enforced

| Rule               | Description                                   |
|--------------------|-----------------------------------------------|
| `IndentationCheck` | Enforces consistent indentation (2 spaces)    |
| `LineLength`       | Maximum line length (100 characters)          |
| `NeedBraces`       | All control structures must use braces        |
| `JavadocMethod`    | Public methods require Javadoc (configurable) |
| `UnusedImports`    | No unused imports allowed                     |

## Static Analysis: PMD

[PMD](https://pmd.github.io/) detects potential bugs, dead code, suboptimal code, and overly complicated expressions.

### Gradle Configuration

```groovy
plugins {
    id 'pmd'
}

pmd {
    toolVersion = '7.1.0'
    ruleSetFiles = files('src/main/resources/app/quality/pmd-custom-ruleset.xml')
    consoleOutput = true
}
```

### Usage

```bash
# Run PMD analysis
gradle pmdMain
gradle pmdTest
```

## Security Analysis: SpotBugs

[SpotBugs](https://spotbugs.github.io/) detects potential security vulnerabilities and bug patterns in compiled bytecode.

### Gradle Configuration

```groovy
plugins {
    id 'com.github.spotbugs' version '6.0.22'
}

spotbugs {
    toolVersion = '4.8.4'
    excludeFilter = file('src/main/resources/app/security/spotbugs-security-exclude.xml')
}
```

### Usage

```bash
# Run SpotBugs on main sources
gradle spotbugsMain

# Run SpotBugs on test sources
gradle spotbugsTest
```

## Editor Configuration

### IntelliJ IDEA

1. Install the [google-java-format plugin](https://plugins.jetbrains.com/plugin/8527-google-java-format).
2. Enable it in **Settings > google-java-format Settings > Enable google-java-format**.
3. Enable **Settings > Editor > General > Auto Import > Optimize imports on the fly**.
4. Install the [Checkstyle-IDEA plugin](https://plugins.jetbrains.com/plugin/1065-checkstyle-idea) and point it to the project's ruleset.

### Visual Studio Code

Install the [Language Support for Java](https://marketplace.visualstudio.com/items?itemName=redhat.java) extension and add the following settings:

```json
{
    "java.format.settings.url": "https://raw.githubusercontent.com/google/styleguide/gh-pages/eclipse-java-google-style.xml",
    "java.format.enabled": true,
    "editor.formatOnSave": true,
    "[java]": {
        "editor.defaultFormatter": "redhat.java"
    }
}
```

---

# Java Type System

> **TL;DR:** Java is statically typed -- the compiler enforces type safety at build time. Use **records** for immutable data, **generics** with bounded type parameters for reusable abstractions, and **sealed classes** (Java 17+) for restricted type hierarchies. Never use raw types or `Object` as a catch-all parameter.

## Overview

Java's type system is static and checked at compile time. Combined with generics, records, and sealed types, it provides strong guarantees about program correctness. This page focuses on the patterns and principles that maximize type safety and code clarity in Java projects.

## Records

Use Java records (Java 16+) for immutable data carriers: DTOs, domain events, value objects, and listener contracts.

```java
// Correct -- immutable domain event
public record ItemEvent(
        Long organizationId,
        String referenceGuid,
        ItemCode code,
        ItemCategory category) {}

// Correct -- request DTO with validation
public record InsertItemRequest(
        @NotNull @Size(min = 1, max = 255) String name,
        @NotNull ItemCode code) {}

// Correct -- listener callbacks as a record
public record Listeners(Runnable onSuccess, Consumer<Exception> onError) {}
```

Records automatically provide:
- `final` fields with constructor, getters, `equals()`, `hashCode()`, and `toString()`
- Immutability by design
- Compact and readable declarations

## Generics

Use generics to create type-safe, reusable abstractions. Always use bounded type parameters when the generic type must satisfy constraints.

### When to Use

```java
// Correct -- generic entity that works with different related types
public final class Item<A> {
    private final ItemCode code;
    private final ItemSeverity severity;
    private List<A> affected;

    public Boolean hasData() {
        return !affected.isEmpty();
    }
}

// Correct -- generic repository interface
public interface QueryDslItemsRepository<T> extends ItemsRepository {
    Page<T> findAllWithFilters(Map<String, Object> filters, Pageable pageable);
}
```

### When Not to Use

```java
// Unnecessary -- only works with one type
public class ItemProcessor<T extends Item<?>> {
    public void process(T item) { ... }
}

// Better -- use the concrete type directly
public class ItemProcessor {
    public void process(Item<?> item) { ... }
}
```

### Bounded Type Parameters

Use bounded type parameters to constrain generic types:

```java
// Upper bound -- T must be Comparable
public static <T extends Comparable<T>> T max(List<T> items) {
    return items.stream().max(Comparator.naturalOrder()).orElseThrow();
}

// Multiple bounds
public static <T extends Serializable & Comparable<T>> void sort(List<T> items) {
    Collections.sort(items);
}
```

## Annotations

### Jakarta Validation

Use Jakarta Bean Validation annotations on request DTOs to validate input at the system boundary:

```java
public record InsertItemRequest(
        @NotNull @Size(min = 1, max = 255) String name,
        @NotNull ItemCode code,
        @Min(0) Long amount) {}
```

### JPA/Persistence

JPA annotations belong **only** on infrastructure models (`Jpa*` classes), never on domain entities:

```java
// Correct -- annotations on infrastructure model
@Entity
@Table(name = "item")
public class JpaItem {
    @Id
    @GeneratedValue(strategy = GenerationType.SEQUENCE)
    private Long id;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "category_id")
    private JpaCategory category;
}

// Wrong -- annotations on domain entity
public class Item {
    @Id
    private Long id;  // Domain entities must be framework-free
}
```

### Lombok

Use Lombok to reduce boilerplate while maintaining clarity:

| Annotation                 | Purpose                                                          |
|----------------------------|------------------------------------------------------------------|
| `@Getter` / `@Setter`      | Generate accessors                                               |
| `@NoArgsConstructor`       | JPA requirement for entity classes                               |
| `@RequiredArgsConstructor` | Constructor injection (generates constructor for `final` fields) |
| `@SuperBuilder`            | Fluent builder pattern with inheritance support                  |
| `@Slf4j`                   | Generate SLF4J logger field                                      |

### Spring

| Annotation        | Layer          | Purpose                    |
|-------------------|----------------|----------------------------|
| `@Component`      | Domain         | Commands, listeners        |
| `@RestController` | Infrastructure | HTTP controllers           |
| `@Service`        | Infrastructure | Service implementations    |
| `@Repository`     | Infrastructure | Repository implementations |

## Sealed Classes (Java 17+)

Use sealed classes and interfaces to restrict type hierarchies:

```java
// Only the listed classes can extend ItemResult
public sealed interface ItemResult
        permits ItemResult.Success, ItemResult.NotFound, ItemResult.Error {

    record Success(Item item) implements ItemResult {}
    record NotFound(Long id) implements ItemResult {}
    record Error(Exception cause) implements ItemResult {}
}
```

Sealed types work well with `switch` expressions (Java 21+ pattern matching):

```java
return switch (result) {
    case ItemResult.Success s -> ResponseEntity.ok(s.item());
    case ItemResult.NotFound n -> ResponseEntity.notFound().build();
    case ItemResult.Error e -> ResponseEntity.internalServerError().build();
};
```

## Prohibited Patterns

```java
// Wrong -- raw type (loses type safety)
List items = new ArrayList();

// Wrong -- Object as catch-all parameter
public void process(Object data) { ... }

// Wrong -- unchecked cast without type guard
Item item = (Item) someObject;

// Wrong -- using Optional as a field or parameter
public class ItemHolder {
    private Optional<Item> item;  // Optional is for return types only
}
```

---

# Java Logging

> **TL;DR:** Use **[SLF4J](https://www.slf4j.org/)** with **[Logback](https://logback.qos.ch/)** (the Spring Boot default) for all logging. Use Lombok's `@Slf4j` annotation to generate the logger field. Use `{}` placeholders for structured logging instead of string concatenation. Never use `System.out.println` or `java.util.logging`.

## Overview

Consistent, structured logging is essential for production observability. This page defines the mandatory logging library and patterns for all Java projects.

## Mandatory Library: SLF4J with Logback

**Use [SLF4J](https://www.slf4j.org/) with [Logback](https://logback.qos.ch/) for all logging.** This is the default logging framework in Spring Boot. SLF4J provides the API (facade), and Logback provides the implementation -- together they offer structured logging, configurable levels, and JSON output for production environments.

### Installation

Spring Boot includes SLF4J and Logback by default. No additional dependencies are required:

```groovy
// Already included via spring-boot-starter
implementation 'org.springframework.boot:spring-boot-starter'
```

### Import Convention

Use Lombok's `@Slf4j` annotation to generate the logger field automatically:

```java
import lombok.extern.slf4j.Slf4j;

@Slf4j
@Component
public class InsertItemCommand {
    public void execute(final ItemEvent event) {
        log.info("Processing item event for organization: {}", event.organizationId());
    }
}
```

When Lombok is not available, create the logger manually:

```java
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class InsertItemCommand {
    private static final Logger log = LoggerFactory.getLogger(InsertItemCommand.class);
}
```

## Log Levels

| Level | Method        | When to Use                                                                                      |
|-------|---------------|--------------------------------------------------------------------------------------------------|
| TRACE | `log.trace()` | Very fine-grained diagnostic information (method entry/exit, loop iterations)                    |
| DEBUG | `log.debug()` | Diagnostic information useful during development (variable values, flow decisions)               |
| INFO  | `log.info()`  | General operational events (application started, request processed, message sent)                |
| WARN  | `log.warn()`  | Potential issues that do not prevent operation (deprecated API usage, retry attempts)            |
| ERROR | `log.error()` | Errors that prevent a specific operation but not the application (failed request, invalid input) |

**Note:** Unlike some logging frameworks, SLF4J does not have a `FATAL` level. Use `ERROR` for critical errors and handle application shutdown via Spring's shutdown hooks.

## Structured Logging

Always use `{}` placeholders to attach contextual data to log entries rather than concatenating values into the message string:

```java
// Correct -- placeholder substitution
log.info("Message sent successfully to topic: {}", topicName);
log.error("Error sending message {} to topic: {}", messageId, topicName, throwable);

// Wrong -- string concatenation
log.info("Message sent successfully to topic: " + topicName);

// Wrong -- String.format
log.info(String.format("Message sent to topic: %s", topicName));
```

### Exception Logging

When logging exceptions, pass the throwable as the **last argument** -- SLF4J automatically appends the stack trace:

```java
try {
    service.send(event);
} catch (Exception e) {
    // Correct -- throwable as last argument, stack trace printed automatically
    log.error("Failed to send item event for organization: {}", event.organizationId(), e);
}
```

### Structured Fields (MDC)

Use the Mapped Diagnostic Context (MDC) for contextual fields that span multiple log entries:

```java
import org.slf4j.MDC;

MDC.put("organizationId", String.valueOf(organizationId));
MDC.put("requestId", requestId);
try {
    log.info("Processing request");
    // ... business logic ...
    log.info("Request completed");
} finally {
    MDC.clear();
}
```

## Configuration

Configure logging in `application.yaml`:

```yaml
logging:
  level:
    root: 'warn'
    com.example.myapp: 'info'
  pattern:
    console: '%d{yyyy-MM-dd HH:mm:ss.SSS} %5p %.20logger{39} : %m%n%wEx'
```

## Prohibited Patterns

```java
// Wrong -- System.out for application logging
System.out.println("Something happened");
System.err.println("Error occurred");

// Wrong -- java.util.logging
import java.util.logging.Logger;
Logger.getLogger("MyClass").info("message");

// Wrong -- string concatenation in log messages
log.info("Processing item " + item.getId() + " for org " + orgId);

// Wrong -- unnecessary isEnabled check with placeholder syntax
if (log.isDebugEnabled()) {
    log.debug("Processing item: {}", item.getId());  // SLF4J already skips evaluation
}
```

**Note:** The `isEnabled` guard is only necessary when the argument computation itself is expensive (e.g., serializing a large object). For simple field access, SLF4J's lazy evaluation via `{}` placeholders is sufficient.

---

# Java Testing

> **TL;DR:** Use **[JUnit 5](https://junit.org/junit5/)** as the testing framework with **[Mockito](https://site.mockito.org/)** for mocking and **[Java Faker](https://github.com/DiUS/java-faker)** for test data generation. Tag tests with `@Tag("unit")` or `@Tag("integration")`. All tests must follow the BDD pattern with `// given`, `// when`, `// then` comment blocks. Use `@DisplayName` for human-readable test descriptions. Unit tests run in **parallel**. Integration tests use **Spring Boot test slices** with setup/teardown and are NOT parallel.

## Overview

JUnit 5 discovers test files automatically via classpath scanning. This document defines the conventions for organizing and writing tests across all Java projects.

## File Structure

```
src/test/java/.../
  <module>/
    domain/
      builders/                         test data builders
      commands/
        InsertItemCommandTest.java      command unit tests
    infrastructure/
      builders/                         JPA entity builders
      controllers/
        InsertItemControllerTest.java   controller unit tests
      repositories/
        JpaItemsRepositoryTest.java     repository integration tests
      services/
        JpaSendItemServiceTest.java     service unit tests
      listeners/
        ProcessItemListenerTest.java    listener tests
      doubles/
        SendItemServiceStub.java        stubs
        DummyItemService.java           dummies
        InMemoryItemsRepository.java    in-memory implementations

src/test/resources/
  db/                                   test database migrations
  http/                                 HTTP test data (REST requests)
  queue/                                message test data
```

## General Conventions

1. **Test tagging is mandatory.** Every test class must have `@Tag("unit")` or `@Tag("integration")`.
2. **BDD structure.** Every test must use `// given`, `// when`, `// then` comment blocks to separate preconditions, actions, and assertions.
3. **Display names.** Use `@DisplayName` on every test method to provide a human-readable description of the scenario.
4. **Testing framework.** Use [JUnit 5](https://junit.org/junit5/) with [Mockito](https://site.mockito.org/) for mocking and [AssertJ](https://assertj.github.io/doc/) or JUnit assertions for verification.
5. **File naming.** Test classes use the `Test` suffix (e.g., `InsertItemCommandTest.java`).
6. **File placement.** Test classes mirror the source structure under `src/test/java/`.
7. **Parallel unit tests.** Unit tests are configured for parallel class-level execution via JUnit properties.
8. **Sequential integration tests.** Integration tests share database state and must NOT run in parallel.

## Unit Tests (Parallel with @Tag)

Unit tests must be lightweight, fast, and isolated. They are configured for **parallel execution** at the class level.

### Command Tests

```java
@Tag("unit")
@NoArgsConstructor(access = AccessLevel.PRIVATE)
class InsertItemCommandTest {

    @Test
    @DisplayName("should call onSuccess when the item is inserted")
    void shouldCallOnSuccess() {
        // given
        final var event = new ItemEventBuilder().build();
        final var service = new SendItemServiceStub().withOnSuccess();
        final var command = new InsertItemCommand(service);
        final var onSuccess = mock(Runnable.class);

        // when
        final var listeners = new InsertItemCommand.Listeners(
                onSuccess,
                (e) -> fail("onError should not be called"));
        command.execute(event, listeners);

        // then
        verify(onSuccess, times(1)).run();
    }

    @Test
    @DisplayName("should call onError when the service throws an exception")
    void shouldCallOnError() {
        // given
        final var event = new ItemEventBuilder().build();
        final var service = new SendItemServiceStub().withOnError(new RuntimeException("test"));
        final var command = new InsertItemCommand(service);
        final var errorRef = new AtomicReference<Exception>();

        // when
        final var listeners = new InsertItemCommand.Listeners(
                () -> fail("onSuccess should not be called"),
                errorRef::set);
        command.execute(event, listeners);

        // then
        assertNotNull(errorRef.get());
        assertEquals("test", errorRef.get().getMessage());
    }
}
```

**Key points:**
- `@Tag("unit")` marks the class for unit test execution.
- `@NoArgsConstructor(access = AccessLevel.PRIVATE)` prevents instantiation outside JUnit.
- Each test is self-contained -- it creates its own doubles, command, and listeners.
- The Listeners pattern reflects all possible outcomes: `onSuccess`, `onError`.

### Controller Tests

```java
@Tag("unit")
@NoArgsConstructor(access = AccessLevel.PRIVATE)
class ListItemsControllerTest {

    @Test
    @DisplayName("should respond 200 (OK) when items are listed successfully")
    void shouldRespondOk() {
        // given
        final var command = new ListItemsCommandStub().withOnSuccess();
        final var controller = new ListItemsController(command);

        // when
        final var response = controller.execute(1L, PageRequest.of(0, 10));

        // then
        assertEquals(HttpStatus.OK.value(), response.getStatusCode().value());
    }

    @Test
    @DisplayName("should respond 500 (Internal Server Error) when command fails")
    void shouldRespondInternalServerError() {
        // given
        final var command = new ListItemsCommandStub().withOnError();
        final var controller = new ListItemsController(command);

        // when
        final var response = controller.execute(1L, PageRequest.of(0, 10));

        // then
        assertEquals(HttpStatus.INTERNAL_SERVER_ERROR.value(), response.getStatusCode().value());
    }
}
```

## Service Tests

```java
@Tag("unit")
@NoArgsConstructor(access = AccessLevel.PRIVATE)
class JpaSendItemServiceTest {

    @Test
    @DisplayName("should save the mapped item when event is valid")
    void shouldSaveItem() {
        // given
        final var event = new ItemEventBuilder().build();
        final var repository = new InMemoryItemsRepository();
        final var mapper = ItemMapper.INSTANCE;
        final var service = new JpaSendItemService(repository, mapper);

        // when
        service.send(event);

        // then
        assertEquals(1, repository.count());
    }

    @Test
    @DisplayName("should throw when repository fails to save")
    void shouldThrowWhenRepositoryFails() {
        // given
        final var event = new ItemEventBuilder().build();
        final var repository = new InMemoryItemsRepository().withOnError(new RuntimeException("db error"));
        final var mapper = ItemMapper.INSTANCE;
        final var service = new JpaSendItemService(repository, mapper);

        // when & then
        assertThrows(RuntimeException.class, () -> service.send(event));
    }
}
```

## Integration Tests (Spring Boot + TestContainers)

Integration tests verify the full stack (database, HTTP, messaging) using real infrastructure. They are **NOT parallel** due to shared mutable state.

### Repository Tests

```java
@Tag("integration")
@SpringBootTest
@ActiveProfiles("test")
class JpaItemsRepositoryTest {

    @Autowired
    private JpaItemsRepository repository;

    @Autowired
    private JdbcTemplate jdbcTemplate;

    @BeforeEach
    void setUp() {
        // Seed test data
        jdbcTemplate.execute("INSERT INTO item (id, name, created_at) VALUES (1, 'test-item', NOW())");
    }

    @AfterEach
    void tearDown() {
        jdbcTemplate.execute("DELETE FROM item");
    }

    @Test
    @DisplayName("should find item by ID successfully")
    void shouldFindById() {
        // given
        final var itemId = 1L;

        // when
        final var result = repository.findById(itemId);

        // then
        assertTrue(result.isPresent());
        assertEquals("test-item", result.get().getName());
    }

    @Test
    @DisplayName("should return empty when item does not exist")
    void shouldReturnEmpty() {
        // given
        final var nonExistentId = 99999L;

        // when
        final var result = repository.findById(nonExistentId);

        // then
        assertTrue(result.isEmpty());
    }

    @Test
    @DisplayName("should save a new item successfully")
    void shouldSaveItem() {
        // given
        final var item = JpaItem.builder()
                .name("new-item")
                .createdAt(LocalDateTime.now())
                .build();

        // when
        final var saved = repository.save(item);

        // then
        assertNotNull(saved.getId());
        assertEquals("new-item", saved.getName());
    }

    @Test
    @DisplayName("should delete an item successfully")
    void shouldDeleteItem() {
        // given
        final var itemId = 1L;

        // when
        repository.deleteById(itemId);

        // then
        assertTrue(repository.findById(itemId).isEmpty());
    }
}
```

**Key points:**
- `@Tag("integration")` marks the class for integration test execution.
- `@SpringBootTest` loads the full application context.
- `@ActiveProfiles("test")` activates the test profile (`application-test.yaml`).
- `@BeforeEach` / `@AfterEach` manage test data lifecycle.
- Tests are grouped by outcome: success, error, edge cases.

## Test Doubles

Follow the [Martin Fowler taxonomy](https://martinfowler.com/articles/mocksArentStubs.html) for test doubles:

| Type          | Purpose                                    | Example                                     |
|---------------|--------------------------------------------|---------------------------------------------|
| **Stub**      | Returns canned answers, no logic           | `SendItemServiceStub`                       |
| **Dummy**     | Minimal implementation, ready-made answers | `DummyItemService`                          |
| **In-Memory** | In-memory logic without external modules   | `InMemoryItemsRepository`                   |
| **Faker**     | External library generating realistic data | `ItemEventBuilder` (using Java Faker)       |
| **Mock**      | Mimics and verifies method calls           | Mockito `mock()` -- **avoid when possible** |

### Stub Example

```java
public class SendItemServiceStub implements SendItemService {
    private Exception error;

    public SendItemServiceStub withOnSuccess() {
        this.error = null;
        return this;
    }

    public SendItemServiceStub withOnError(final Exception error) {
        this.error = error;
        return this;
    }

    @Override
    public void send(final ItemEvent event) {
        if (error != null) {
            throw new RuntimeException(error);
        }
    }
}
```

### In-Memory Repository Example

```java
public class InMemoryItemsRepository implements ItemsRepository {
    private final List<JpaItem> items = new ArrayList<>();
    private Exception error;

    public InMemoryItemsRepository withOnError(final Exception error) {
        this.error = error;
        return this;
    }

    public long count() {
        return items.size();
    }

    @Override
    public JpaItem save(final JpaItem item) {
        if (error != null) {
            throw new RuntimeException(error);
        }
        items.add(item);
        return item;
    }
}
```

## Builders

Use the Builder Design Pattern to construct complex test objects step by step. Builders keep test setup readable and reusable across test suites.

```java
@NoArgsConstructor
public final class ItemEventBuilder {
    private Long organizationId;
    private String referenceGuid;
    private ItemCode code;
    private ItemCategory category;

    public ItemEventBuilder withOrganizationId(final Long organizationId) {
        this.organizationId = organizationId;
        return this;
    }

    public ItemEventBuilder withCode(final ItemCode code) {
        this.code = code;
        return this;
    }

    public ItemEvent build() {
        final var faker = new Faker();
        return new ItemEvent(
                organizationId != null ? organizationId : faker.number().randomNumber(),
                referenceGuid != null ? referenceGuid : faker.internet().uuid(),
                code != null ? code : ItemCode.values()[faker.number().numberBetween(0, ItemCode.values().length)],
                category != null ? category : ItemCategory.DEFAULT);
    }
}
```

Usage in tests:

```java
// Default values (randomized via Faker)
final var event = new ItemEventBuilder().build();

// Custom values
final var event = new ItemEventBuilder()
        .withOrganizationId(42L)
        .withCode(ItemCode.CRITICAL)
        .build();
```

## Parallel Execution

Configure parallel execution in `build.gradle`:

```groovy
test {
    useJUnitPlatform()
    systemProperty 'junit.jupiter.execution.parallel.enabled', 'true'
    systemProperty 'junit.jupiter.execution.parallel.mode.classes.default', 'concurrent'
}
```

This runs test **classes** concurrently while methods within the same class run in the same thread, preventing shared-state conflicts within a test class.

## Seeds

Seed files populate test databases with known data:

- The seed file name must match the **table name** (e.g., `item.sql`).
- For multiple seeds targeting the same table, use numbered suffixes: `item_01.sql`, `item_02.sql`.
- Constants must be named according to the file name.

---

# Java Project Structure

> **TL;DR:** Follow the domain/infrastructure layer separation within each business module. Use [Gradle](https://gradle.org/) for builds and dependency management. Place test files under `src/test/java/` mirroring the source structure. Use module-based organization for multi-domain applications.

## Overview

This page defines the standard directory layout and dependency management practices for all Java projects. The architecture follows the Backend Design specification, separating code into `domain` (contracts) and `infrastructure` (implementations) layers within each business module.

## Directory Structure

```
src/main/java/.../
  Startup.java                           application entry point (@SpringBootApplication)
  global/                                cross-cutting concerns
    errors/                                global error handlers
    helpers/                               shared utilities
  items/                                 business module (one per domain concept)
    domain/
      commands/
        InsertItemCommand.java             business logic (write operation)
        ListItemsCommand.java              business logic (read operation)
      entities/
        Item.java                          pure domain entity (no annotations)
        ItemEvent.java                     domain event (record)
        ItemCode.java                      domain enum
      services/
        SendItemService.java               service contract (interface)
      helpers/
        ItemHelper.java                    domain utilities
    infrastructure/
      controllers/
        InsertItemController.java          HTTP endpoint (write)
        ListItemsController.java           HTTP endpoint (read)
        requests/
          InsertItemRequest.java           request DTO (record)
        responses/
          ItemResponse.java                response DTO (record)
        mappers/
          InsertItemRequestMapper.java     request-to-entity mapper
          ItemResponseMapper.java          entity-to-response mapper
      listeners/
        ProcessItemListener.java           Kafka listener
        messages/
          ItemMessage.java                 message DTO (record)
        mappers/
          ItemMessageMapper.java           message-to-entity mapper
      repositories/
        contracts/
          ItemsRepository.java             repository interface
          QueryDslItemsRepository.java     QueryDSL extension
        JpaItemsRepository.java            Spring Data JPA implementation
        models/
          JpaItem.java                     JPA entity model
        helpers/
          ItemQueryBuilder.java            QueryDSL query construction
      services/
        JpaSendItemService.java            service implementation
        mappers/
          ItemMapper.java                  entity-to-model mapper

src/main/resources/
  application.yaml                       main configuration
  application-prod.yaml                  production profile
  application-test.yaml                  test profile
  messages.properties                    i18n (English)
  messages_pt_BR.properties              i18n (Portuguese)
  db/
    changelog/                           Liquibase migration scripts (YAML)
    seeds/                               test data seeders (SQL)
  app/
    quality/
      checkstyle-google-ruleset.xml      Checkstyle configuration
      pmd-custom-ruleset.xml             PMD configuration
    security/
      spotbugs-security-exclude.xml      SpotBugs exclusions
      dependency-check-suppress.xml      OWASP dependency suppression

src/test/java/.../
  items/                                 mirrors source module structure
    domain/
      builders/
        ItemEventBuilder.java            test data builder
      commands/
        InsertItemCommandTest.java       command unit tests
    infrastructure/
      builders/
        JpaItemBuilder.java              JPA entity builder
      controllers/
        InsertItemControllerTest.java    controller unit tests
      repositories/
        JpaItemsRepositoryTest.java      repository integration tests
      services/
        JpaSendItemServiceTest.java      service unit tests
      listeners/
        ProcessItemListenerTest.java     listener tests
      doubles/
        SendItemServiceStub.java         stubs
        DummyItemService.java            dummies
        InMemoryItemsRepository.java     in-memory implementations

src/test/resources/
  db/                                    test-specific migrations
  http/                                  HTTP test request data
  queue/                                 Kafka test message data
```

### Key Directories

| Directory                               | Purpose                                                 |
|-----------------------------------------|---------------------------------------------------------|
| `global/`                               | Cross-cutting concerns (error handling, shared helpers) |
| `<module>/domain/commands/`             | Business logic implementations (CQRS write side)        |
| `<module>/domain/entities/`             | Framework-agnostic domain entities and events           |
| `<module>/domain/services/`             | Service contracts (interfaces)                          |
| `<module>/infrastructure/controllers/`  | HTTP controllers (CQRS read side + write endpoints)     |
| `<module>/infrastructure/listeners/`    | Kafka/queue message listeners                           |
| `<module>/infrastructure/repositories/` | Repository implementations with JPA/QueryDSL            |
| `<module>/infrastructure/services/`     | Service implementations with infrastructure tools       |
| `src/main/resources/db/changelog/`      | Liquibase database migration scripts                    |
| `src/main/resources/app/quality/`       | Code quality tool configurations                        |

## Package Manager: Gradle

All Java projects use [Gradle](https://gradle.org/) for build automation and dependency management.

### Project Setup

```groovy
// settings.gradle
rootProject.name = 'my-application'

// build.gradle
plugins {
    id 'java'
    id 'org.springframework.boot' version '3.3.4'
    id 'io.spring.dependency-management' version '1.1.6'
}

group = 'com.example'
version = '1.0.0'

java {
    sourceCompatibility = JavaVersion.VERSION_21
    targetCompatibility = JavaVersion.VERSION_21
}
```

### Managing Dependencies

```groovy
dependencies {
    // Spring Boot starters
    implementation 'org.springframework.boot:spring-boot-starter-web'
    implementation 'org.springframework.boot:spring-boot-starter-data-jpa'
    implementation 'org.springframework.boot:spring-boot-starter-validation'

    // Database
    runtimeOnly 'org.postgresql:postgresql'
    implementation 'org.liquibase:liquibase-core'

    // Messaging
    implementation 'org.springframework.kafka:spring-kafka'

    // Mapping
    implementation 'org.mapstruct:mapstruct:1.6.2'
    annotationProcessor 'org.mapstruct:mapstruct-processor:1.6.2'

    // Utilities
    compileOnly 'org.projectlombok:lombok'
    annotationProcessor 'org.projectlombok:lombok'

    // Testing
    testImplementation 'org.springframework.boot:spring-boot-starter-test'
    testImplementation 'com.github.javafaker:javafaker:1.0.2'
    testImplementation 'org.awaitility:awaitility:4.2.2'
    testImplementation 'org.testcontainers:junit-jupiter'
}
```

### Gradle Wrapper

Always use the Gradle Wrapper (`gradlew`) to ensure consistent build tool versions across environments:

```bash
# Build the project
./gradlew build

# Run tests
./gradlew test

# Format code
./gradlew spotlessApply

# Run quality checks
./gradlew checkstyleMain pmdMain spotbugsMain
```

The `gradle/wrapper/` directory and `gradlew` / `gradlew.bat` scripts must be committed to version control.

## Build & Distribution

### Building

```bash
# Build the JAR
./gradlew bootJar

# Build without running tests
./gradlew bootJar -x test
```

### Running

```bash
# Run via Gradle
./gradlew bootRun

# Run the compiled JAR
java -jar build/libs/my-application-1.0.0.jar
```

### Docker

Use multi-stage builds to produce minimal container images:

```dockerfile
FROM eclipse-temurin:21-jdk-alpine AS builder
WORKDIR /app
COPY gradlew build.gradle settings.gradle ./
COPY gradle/ gradle/
RUN ./gradlew dependencies --no-daemon
COPY src/ src/
RUN ./gradlew bootJar --no-daemon -x test

FROM eclipse-temurin:21-jre-alpine
COPY --from=builder /app/build/libs/*.jar /app/app.jar
ENTRYPOINT ["java", "-jar", "/app/app.jar"]
```

## Key Configuration Files

| File                | Purpose                                                            |
|---------------------|--------------------------------------------------------------------|
| `build.gradle`      | Build configuration, plugins, dependencies, and quality tool setup |
| `settings.gradle`   | Project name and multi-module settings                             |
| `gradle.properties` | Gradle build properties (JVM args, versions)                       |
| `application.yaml`  | Spring Boot configuration (profiles, datasource, kafka, logging)   |
| `lombok.config`     | Lombok behavior configuration                                      |
| `compose.yaml`      | Docker Compose for production                                      |
| `compose.dev.yaml`  | Docker Compose for local development                               |
| `.editorconfig`     | Editor standardization                                             |

---

# JavaScript & TypeScript Conventions

> **TL;DR:** Use `snake_case` for file names and `camelCase` for HTML attribute values. Prefer modern syntax (arrow functions, template strings, optional chaining, nullish coalescing). Never use `any` -- use `unknown` instead. Favor immutability and destructuring.

## General

### File Names

All file names must use `snake_case`.

```
src/my_class.js    ✅ Correct
src/myClass.js       ❌ Wrong
```

### Formatting and Linting

Use **Prettier** for code formatting and **ESLint** for linting. Every Prettier and ESLint rule is considered part of this guide. Disable rules only in exceptional cases and always leave a comment explaining why.

## HTML

### Camel Case for `id` and `data-test-subj`

Use camelCase for the values of `id` and `data-test-subj` attributes:

```html
<button id="veryImportantButton" data-test-subj="clickMeButton">Click me</button>
```

The only exception is when dynamically generating values, where hyphens may be used as delimiters:

```jsx
buttons.map(btn => (
  <button
    id={`veryImportantButton-${btn.id}`}
    data-test-subj={`clickMeButton-${btn.id}`}
  >
    {btn.label}
  </button>
))
```

## TypeScript / JavaScript

### Prefer Modern Syntax

- Prefer **arrow functions** over function expressions.
- Prefer **template strings** over string concatenation.
- Prefer the **spread operator** (`[...arr]`) over `arr.slice()` for copying arrays.
- Use **optional chaining** (`?.`) and **nullish coalescing** (`??`) over `lodash.get` and similar utilities.

### Avoid Mutability

Do not reassign variables, modify object properties, or push values to arrays. Instead, create new variables and shallow copies:

```js
// ✅ Good
function addBar(foos, foo) {
  const newFoo = { ...foo, name: 'bar' };
  return [...foos, newFoo];
}

// ❌ Bad
function addBar(foos, foo) {
  foo.name = 'bar';
  foos.push(foo);
}
```

### Avoid `any`

Since TypeScript 3.0 introduced the [`unknown` type](https://mariusschulz.com/blog/the-unknown-type-in-typescript), there is rarely a valid reason to use `any`. Replace `any` with either a generic type parameter or `unknown`, combined with type narrowing.

### Use Object Destructuring

Destructuring reduces temporary references and prevents typo-related bugs:

```js
// ✅ Best
function fullName({ first, last }) {
  return `${first} ${last}`;
}

// ❌ Bad
function fullName(user) {
  const first = user.first;
  const last = user.last;
  return `${first} ${last}`;
}
```

### Use Array Destructuring

Avoid accessing array values by index. When direct access is necessary, use array destructuring:

```js
const arr = [1, 2, 3];

// ✅ Good
const [first, second] = arr;

// ❌ Bad
const first = arr[0];
const second = arr[1];
```

---

# JavaScript & TypeScript Testing

> **TL;DR:** Use Jest for testing, Enzyme for React component rendering, and Faker for generating test data. All tests must follow the BDD pattern with `// given`, `// when`, `// then` comment blocks. Test business rules and user flows rather than framework internals. Use `mount` for integration tests (triggers `useEffect`) and `shallow` for lightweight unit tests.

## Overview

Front-end testing can be challenging. This document provides clear guidelines on **what** to test and **how** to test it, ensuring consistency across the codebase. It is a living document and should be updated as new modules and features introduce new testing requirements.

## Tools

| Tool                                         | Purpose                                 |
|----------------------------------------------|-----------------------------------------|
| [Jest](https://jestjs.io/)                   | JavaScript testing framework and runner |
| [Enzyme](https://enzymejs.github.io/enzyme/) | React component testing utility         |
| [Faker](https://github.com/marak/Faker.js/)  | Fake data generation library            |

## BDD Structure (Given / When / Then)

**Every test must use `// given`, `// when`, `// then` comment blocks** to clearly separate preconditions, actions, and assertions. This is mandatory across all test files.

```ts
it('should render Content when not loading', () => {
  // given
  const props = { loading: false };

  // when
  const component = shallow(<ExampleComponent {...props} />);

  // then
  expect(component.find(Content).exists()).toBeTruthy();
  expect(component.find(Loading).exists()).toBeFalsy();
});
```

## Unit Tests

### What to Test

Focus on testing **your own** business rules and code behavior, not the framework's rendering internals:

**a. Props are passed correctly:**
```ts
it('should render an ExampleChildren with the right props', () => {
  // given
  const component = mount(<ExampleComponent prop={'prop_value'} />);

  // then
  expect(component.find(ExampleChildren).prop('prop')).toBe('prop_value');
});
```

**b. Component output matches expected state:**
```ts
it('should render the Content component', () => {
  // when
  const component = shallow(<ExampleComponent loading={false} />);

  // then
  expect(component.find(Content).exists()).toBeTruthy();
  expect(component.find(Loading).exists()).toBeFalsy();
});

it('should render Loading component', () => {
  // when
  const component = shallow(<ExampleComponent loading={true} />);

  // then
  expect(component.find(Content).exists()).toBeFalsy();
  expect(component.find(Loading).exists()).toBeTruthy();
});
```

**c. User interactions mutate state correctly:**
```ts
it('should present the Congrats component when button clicked', async () => {
  // given
  const component = mount(<ExampleComponent />);

  // when
  component.find(ExampleButton).simulate('click');
  component.update();

  // then
  expect(component.find(Congrats).exists()).toBeTruthy();
});
```

**d. Async actions execute correctly:**
```ts
import { flushComponent } from "tests/public/utils/flush_component.ts";

it('should render the Children with the loaded data', async () => {
  // given
  const exampleService = new MockExampleService();
  const exampleData = faker.random.alphaNumeric();
  exampleService.withSuccess(exampleData);
  const component = mount(
    <ExampleAsyncComponent />,
    { services: { exampleService } }
  );

  // when
  await flushComponent(component);

  // then
  expect(component.find(Children).exists()).toBeTruthy();
  expect(component.find(Children).prop('loadedData')).toBeTruthy();
});
```

**e. Logic returns expected output:**

Avoid duplicating business rules in tests. Hard-code inputs and expected outputs:
```ts
it('should return the sum multiplied by two', () => {
  // given
  const numbersCollection = new NumbersCollection([1, 2, 3]);

  // when
  const result = numbersCollection.sumAndDuplicate(); // (1 + 2 + 3) * 2

  // then
  expect(result).toBe(12);
});
```

**f. Router navigation (`useHistory`):**
```tsx
const mockHistoryPush = jest.fn();

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useHistory: () => ({
    push: mockHistoryPush,
  }),
}));

it('should call push', async () => {
  // given
  const component = mount(<Something />);

  // when
  await flushComponent(component);

  // then
  expect(mockHistoryPush).toHaveBeenCalled();
});
```

### How to Test

| Utility         | Purpose                                                                                  |
|-----------------|------------------------------------------------------------------------------------------|
| `act`           | Required for any action that triggers React state updates                                |
| `flushPromises` | Forces Jest to resolve all pending promises (e.g., API requests)                         |
| `mount`         | Renders full DOM including `useEffect` hooks; suitable for integration tests             |
| `shallow`       | Renders a lightweight representation without `useEffect`; more performant for unit tests |

**Tip:** Use `console.log(component.debug())` to inspect the rendered output when debugging test failures.

---

# YAML Conventions

> **TL;DR:** Always use the `.yaml` extension (not `.yml`). Always quote strings with **single quotes**. Use **double quotes** only when the string contains variable interpolation or escape sequences. Never leave string values unquoted. Do not quote booleans or numbers. These rules apply to all YAML files: pipeline configurations, Kubernetes manifests, infrastructure-as-code, and YAML code blocks inside Markdown.

## Overview

[YAML](https://yaml.org/) (YAML Ain't Markup Language) is a human-readable data serialization format used extensively for configuration files, CI/CD pipelines, Kubernetes manifests, and infrastructure-as-code. The current specification is [YAML 1.2.2](https://yaml.org/spec/1.2.2/), published on October 1, 2021.

This document defines the team's conventions for writing consistent, unambiguous YAML across all projects and documentation.

## File Extension

**Always use `.yaml` as the file extension.**

The [official YAML FAQ](https://yaml.org/faq.html) recommends `.yaml` as the preferred extension. The `.yml` variant exists for historical reasons (legacy Windows three-character extension limits) and must be avoided.

```
# Correct
docker-compose.yaml
config.yaml
deployment.yaml
.golangci.yaml

# Wrong
docker-compose.yml
config.yml
deployment.yml
.golangci.yml
```

### Exceptions

Some tools enforce a specific filename that uses `.yml` and do not accept alternatives. In these rare cases, the tool's requirement takes precedence:

| Tool                    | Required Filename     | Reason                                                             |
|-------------------------|-----------------------|--------------------------------------------------------------------|
| Azure DevOps            | `azure-pipelines.yml` | Only recognizes this exact filename                                |
| Docker Compose (legacy) | `docker-compose.yml`  | Older versions required this name (modern versions accept `.yaml`) |

If a tool accepts both extensions, always use `.yaml`.

## String Quoting

### Rule: Always Quote Strings with Single Quotes

All string values must be explicitly quoted with **single quotes** (`'...'`). Unquoted strings are technically valid YAML, but they introduce ambiguity -- values like `yes`, `no`, `on`, `off`, `null`, or numeric-looking strings can be silently misinterpreted as booleans, nulls, or numbers by the YAML parser.

```yaml
# Correct
name: 'my-service'
image: 'nginx:1.25-alpine'
region: 'us-east-1'
environment: 'production'

# Wrong -- unquoted strings
name: my-service
image: nginx:1.25-alpine
region: us-east-1
environment: production
```

### Rule: Use Double Quotes Only for Interpolation or Escape Sequences

Use **double quotes** (`"..."`) only when the string contains variable interpolation (e.g., environment variable expansion) or escape sequences (e.g., `\n`, `\t`):

```yaml
# Correct -- double quotes for interpolation
connection_string: "${DATABASE_HOST}:${DATABASE_PORT}"
greeting: "Hello,\tWorld\n"

# Correct -- single quotes for everything else
host: 'localhost'
port_label: '8080'
```

### Rule: Do Not Quote Booleans or Numbers

Booleans and numbers are native YAML types and must **not** be quoted:

```yaml
# Correct
enabled: true
replicas: 3
timeout: 30.5
debug: false

# Wrong -- quoted booleans and numbers
enabled: 'true'
replicas: '3'
timeout: '30.5'
```

### Summary Table

| Value Type                   | Quoting Style | Example                   |
|------------------------------|---------------|---------------------------|
| Plain string                 | Single quotes | `name: 'my-app'`          |
| String with variables        | Double quotes | `url: "${API_HOST}/v1"`   |
| String with escape sequences | Double quotes | `message: "line1\nline2"` |
| Boolean                      | No quotes     | `enabled: true`           |
| Number (integer)             | No quotes     | `replicas: 3`             |
| Number (float)               | No quotes     | `ratio: 0.75`             |
| Null                         | No quotes     | `value: null`             |

## YAML in Markdown Code Blocks

When embedding YAML snippets inside Markdown files (READMEs, guides, documentation), the same conventions apply. Code examples must be consistent with production YAML:

````markdown
```yaml
# Correct -- follows all conventions
apiVersion: 'apps/v1'
kind: 'Deployment'
metadata:
  name: 'my-service'
  labels:
    app: 'my-service'
    environment: 'production'
spec:
  replicas: 3
  selector:
    matchLabels:
      app: 'my-service'
```
````

This includes YAML examples in:

- `README.md` files
- Wiki pages and guides
- Pull request descriptions
- Code review comments
- Any other documentation containing YAML snippets

## Scope of Application

These conventions apply to **all** YAML files across the organization:

| Domain                        | Examples                                                    |
|-------------------------------|-------------------------------------------------------------|
| **CI/CD pipelines**           | GitHub Actions workflows, GitLab CI, Azure DevOps pipelines |
| **Kubernetes**                | Deployments, Services, ConfigMaps, Ingresses, Helm values   |
| **Infrastructure-as-Code**    | Docker Compose, Ansible playbooks, CloudFormation templates |
| **Application configuration** | Spring Boot `application.yaml`, service configs             |
| **Tooling**                   | `.golangci.yaml`, `.hadolint.yaml`, `autoupdate.yaml`       |
| **Documentation**             | YAML code blocks inside Markdown files                      |

---

# Code Style

> **TL;DR:** Follow the naming conventions and file structures defined here as the baseline. Language-specific guides override these defaults where applicable. Use the standard operations vocabulary (`List`, `Get`, `Insert`, `Update`, `Delete`, and their batch variants) for naming files and classes.

## Overview

This document establishes the baseline coding standards for file creation, function naming, and variable naming across all projects. If a language-specific guide exists (see the sub-pages for Go, JavaScript, Java, and Python), its conventions take precedence over this document.

The architectural layers referenced throughout are defined in the Backend Design section. All language-specific examples follow **Hexagonal Architecture (Ports & Adapters)** with **Domain-Driven Design (DDD)** and **CQRS** -- see Backend Design for details.

## File Structure

The default file structure follows the Backend Design specification, which separates code into `domain` (contracts) and `infrastructure` (implementations) layers.

## Service Naming

When creating service files, the naming must reflect both the **entity** being operated on and whether the operation targets a **single record or a batch**.

### Operations Vocabulary

Use the following standard operation prefixes consistently across all projects:

| Operation     | Description               |
|---------------|---------------------------|
| `List`        | Retrieve multiple records |
| `Get`         | Retrieve a single record  |
| `Insert`      | Create a single record    |
| `Update`      | Modify a single record    |
| `Delete`      | Remove a single record    |
| `BatchInsert` | Create multiple records   |
| `BatchUpdate` | Modify multiple records   |
| `BatchDelete` | Remove multiple records   |
| `DeleteAll`   | Remove all records        |

---

# Git Flow

> **TL;DR:** Use a feature-branch model with `main` always deployable. Synchronize branches with `git rebase` (never merge). Follow the `type(TASK-ID): message` commit format in simple past tense. Use Semantic Versioning (MAJOR.MINOR.PATCH) for releases. Flag breaking changes in commits, CHANGELOG.md, and PRs.

## Overview

This document defines the Git workflow, naming conventions, commit message format, versioning strategy, and branch management practices for all projects.

## Feature Branch Model

The workflow follows a feature-branch model with these principles:

1. The `main` branch is **always in a deployable state** and ready for production.
2. All changes reach `main` through feature branches and pull requests.
3. Branch synchronization and conflict resolution must use `git rebase`. See [this tutorial](https://www.atlassian.com/git/tutorials/rewriting-history/git-rebase) for details.

### Development Workflow

1. Create a **feature branch**: `feat/TASK-ID`
2. Develop and commit incrementally.
3. Optionally create development builds for testing.
4. **Rebase with `main`** to incorporate upstream changes and resolve conflicts.
5. **Create a Pull Request** targeting `main`, adding reviewers.
6. After code review approval, **merge** the feature branch into `main`.
7. The CI pipeline builds `main` and deploys to the QA environment.
8. After QA approval, follow the `CHANGELOG.md` instructions to generate a production release.
9. The CI pipeline deploys the release to production.

## Naming Conventions

### Branches

**With a task/ticket ID:**
```
feat/TICKET-000
fix/TICKET-990
```

**Without a task/ticket ID:**
```
feat/add-logs
fix/input-mask
```

### Branch Types

| Type       | Purpose                                    |
|------------|--------------------------------------------|
| `feat`     | New feature implementation                 |
| `fix`      | Bug fix for an existing issue              |
| `refactor` | Code restructuring without behavior change |
| `chore`    | Infrastructure or tooling improvement      |
| `test`     | New test scenario                          |
| `docs`     | Documentation change                       |

## Commit Messages

### Format

```
type(SCOPE): message
```

Where **SCOPE** is the task ID (if available) or a short descriptive scope.

### Rules

- **No period** at the end of the subject line.
- **Do not capitalize** the first letter.
- Use **simple past tense**: `changed`, `fixed`, `removed`, `added` (not present continuous).
- Wrap code references in backticks: `setAnyThing`.

### Short Commit Message

```
fix(TICKET-000): changed the button colors
```

### Long Commit Message

```
fix(TICKET-000): changed the button colors

- created a new CSS file for the buttons
- changed the color of the cancel button from blue to red
```

### Message Footer -- Breaking Changes

```
**BREAKING CHANGE:** isolate scope bindings definition has changed and the inject option for the directive controller injection was removed
```

### Message Footer -- Referencing Issues

```
Closes TICKET-567
Closes TICKET-568
```

### Complete Example

```
fix(TICKET-567+TICKET-568+TICKET-569): changed the button colors

- created a new CSS file for the buttons
- changed the color of the cancel button from blue to red

**BREAKING CHANGE:** isolate scope bindings definition has changed and the inject option for the directive controller injection was removed

Closes TICKET-567
Closes TICKET-568
Closes TICKET-569
```

### Automatic Changelog Generation

```bash
git log <last tag> HEAD --pretty=format:%s
git log <last release> HEAD --grep feature
```

## Semantic Versioning

Follow [Semantic Versioning (SemVer)](https://semver.org/) for all releases:

| Release Type | Version Change     | When to Use                                                        |
|--------------|--------------------|--------------------------------------------------------------------|
| **MAJOR**    | `1.0.0` -> `2.0.0` | Breaking changes, new native modules, drastic structural changes   |
| **MINOR**    | `1.0.0` -> `1.1.0` | Incremental features, no new native modules, no structural changes |
| **PATCH**    | `1.0.0` -> `1.0.1` | Bug fixes in production only                                       |

## Working with Branches

To update a feature branch that is behind `main`:

1. Ensure all changes are committed (nothing in stash).
2. Push to the remote: `git push origin <your-branch>`
3. Switch to main: `git checkout main`
4. Update main: `git pull origin main`
5. Switch back: `git checkout <your-branch>`
6. Rebase: `git rebase main`
7. Resolve any conflicts, then `git add <file>` and `git rebase --continue`.
8. Force-push: `git push -f` (necessary because rebase rewrites commit hashes).
9. Create the Pull Request.

## Breaking Changes

When code changes alter public interfaces, flag the breaking change in **three places**:

1. **Commit footer:**
   ```
   refactor(input): used props for state management
   **BREAKING CHANGE:** input behavior now must be implemented by the peer, including value and handleChange
   ```

2. **CHANGELOG.md:**
   ```
   - **BREAKING CHANGE:** updated input to use props for state management
   ```

3. **Pull Request description:**
   ```
   **BREAKING CHANGE:** updated input to use props for state management
   ```

---

# Merge Guide

> **TL;DR:** For dependent (chained) branches, merge from the outermost branch inward before merging into `main`. For independent branches, merge one at a time, rebasing each subsequent branch on the updated `main`.

## Overview

This guide covers the correct merge order for two common scenarios: dependent branches and independent branches. Following these procedures ensures a clean, linear commit history.

## Case 1: Dependent Branches

When branches are chained (e.g., `test/2` depends on `test/1`), merge from the **outermost** branch inward before merging into `main`.

### Procedure

1. Merge `test/2` into `test/1`.
2. Merge `test/1` into `main`.

For deeper chains (N branches):

1. Merge `test/N` into `test/N-1`.
2. Merge `test/N-1` into `test/N-2`.
3. Continue until `test/1`.
4. Merge `test/1` into `main`.

### Result

---

## Case 2: Independent Branches

When multiple branches originate independently from `main`, merge them **one at a time**, rebasing each subsequent branch on the updated `main`.

### Procedure

1. Merge `test/1` into `main`.
2. **Update `main` locally** after the merge.
3. Rebase `test/2` with the updated `main`.
4. Merge `test/2` into `main`.
5. Repeat for additional branches.

This produces independent "triangles" in the commit graph, keeping the history clean and traceable.

### Handling Rebase Conflicts

If conflicts occur during the rebase:

1. Resolve the conflicts in your IDE.
2. Stage the resolved files: `git add <resolved-file>`
3. Continue the rebase: `git rebase --continue`
4. Force-push: `git push -f`

Then update `main` locally and proceed with the merge. The resulting graph:

---

# Testing Standards

> **TL;DR:** All tests must follow the BDD (Behavior-Driven Development) pattern with `// given`, `// when`, `// then` comment blocks. Write descriptive test names per layer: Commands (`"should call <LISTENER> when ..."`), Controllers (`"should respond <HTTP_STATUS_CODE> when ..."`), Services/Repos (`"should ... when ..."` with success + failure pairs). Prefer Stubs, Dummies, and In-memory doubles over Mocks. Use Builders for test data construction.

## Overview

This document defines the testing standards applicable across all languages and projects. The architectural layers referenced here are defined in the Backend Design section.

## BDD Structure (Given / When / Then)

**All tests in all languages must follow the BDD pattern** with three clearly separated blocks using comments. This structure makes every test readable, self-documenting, and consistent across the entire codebase.

| Block     | Purpose                                                                               | Also Known As |
|-----------|---------------------------------------------------------------------------------------|---------------|
| **given** | Set up preconditions -- initialize objects, configure doubles, prepare input data     | Arrange       |
| **when**  | Execute the action under test -- call the method, trigger the event, send the request | Act           |
| **then**  | Assert expected outcomes -- verify return values, check side effects, validate state  | Assert        |

### Comment Syntax by Language

| Language                | Given      | When      | Then      |
|-------------------------|------------|-----------|-----------|
| Go                      | `// given` | `// when` | `// then` |
| JavaScript / TypeScript | `// given` | `// when` | `// then` |
| Java                    | `// given` | `// when` | `// then` |
| Python                  | `# given`  | `# when`  | `# then`  |

### Example (JavaScript)

```ts
it('should render Content when not loading', () => {
  // given
  const props = { loading: false };

  // when
  const component = shallow(<ExampleComponent {...props} />);

  // then
  expect(component.find(Content).exists()).toBeTruthy();
  expect(component.find(Loading).exists()).toBeFalsy();
});
```

### Example (Go)

```go
func (s *CommandSuite) TestCreateUserCommand() {
    // given
    s.repository.On("Save", mock.Anything).Return(nil)
    user := &cmd.User{
        Name:  "John Doe",
        Email: "johndoe@example.com",
    }

    // when
    err := s.command.Run(user)

    // then
    s.repository.AssertExpectations(s.T())
    assert.Nil(s.T(), err)
}
```

### Example (Python)

```python
def test_create_user_successfully(self):
    # given
    user_data = {"name": "John Doe", "email": "john@example.com"}

    # when
    result = self.service.create_user(user_data)

    # then
    assert result.name == "John Doe"
    assert result.email == "john@example.com"
```

## Test Description Patterns

### Commands

Create one test per listener in the command. Use this description format:

```
"should call <LISTENER> when ..."
```

### Controllers

Create one test per HTTP status code. Verify both the status code and the response body:

```
"should respond <HTTP_STATUS_CODE> (HTTP_STATUS) when ..."
```

### Services

Create **at least** one success and one failure test per public method:

```
"should ... when ..."
```

### Repositories

Follow the same conventions as [Services](#services).

## Test Doubles

Based on [Martin Fowler's taxonomy](https://martinfowler.com/bliki/TestDouble.html), a **test double** is any object that replaces a real dependency during testing. The team uses the following types:

| Type          | Purpose                                               | Guidelines                                                           |
|---------------|-------------------------------------------------------|----------------------------------------------------------------------|
| **Stub**      | Returns pre-configured (canned) answers               | No in-memory logic; return static values only                        |
| **Dummy**     | Fills required parameters that are never used         | Return minimal values (empty lists, `null`)                          |
| **In-memory** | Implements logic in memory without external modules   | Use for lightweight simulations of repositories or services          |
| **Faker**     | Generates realistic fake data via an external library | Use libraries like Faker.js or Go Faker                              |
| **Mock**      | Records and verifies method calls                     | **Avoid when possible.** Use only when no other double type suffices |

## Builders

The [Builder Design Pattern](https://refactoring.guru/design-patterns/builder) enables step-by-step construction of complex test objects. Builders separate object construction from representation, making test setup readable and reusable:

```
UserBuilder.new()
    .withName("John Doe")
    .withEmail("john@example.com")
    .build()
```

Use builders to construct entities, DTOs, and other complex objects needed for automated testing.

---

# Architecture

> **TL;DR:** All applications follow **Clean Architecture** principles, separating business logic from infrastructure concerns. Dependencies always point inward toward the domain. Adhere to SOLID principles for maintainable, testable code.

## Overview

Modern application development is fundamentally a task of **managing dependencies**. Poor architectural design leads to rigid, fragile codebases that are costly to change. Since requirements evolve, libraries get deprecated, and external services change, the development team must actively maintain a clean dependency graph.

|        Polluted Architecture        |       Clean Architecture        |
|:---------------------------------------:|:-----------------------------------:|

Clean Architecture enforces a clear separation of concerns where **business rules are independent of frameworks, databases, and delivery mechanisms**. The inner layers define policies; the outer layers implement mechanisms.

## Sub-Pages

---

# Backend Design

> **TL;DR:** Separate code into `domain` (contracts) and `infrastructure` (implementations). The five principal actors are Entities, Controllers, Commands, Services, and Repositories. Use Mappers to isolate layers and Dependency Injection to wire them together.

## Architectural Foundations

The backend architecture draws from three complementary patterns. All language-specific examples in the Code Style guides follow these patterns:

- **[Hexagonal Architecture (Ports & Adapters)](https://alistaiccockburn.com/Hexagonal+architecture):** The `domain`/`infrastructure` split maps directly to the hexagonal model. The domain layer defines **ports** (interfaces for repositories, services) that represent what the application needs. The infrastructure layer provides **adapters** (implementations) that connect the application to databases, HTTP, message queues, and other external systems. Dependencies always point inward -- infrastructure depends on domain, never the reverse.

- **[Domain-Driven Design (DDD)](https://www.domainlanguage.com/ddd/):** Entities encapsulate business logic and are free of framework dependencies. Repositories abstract persistence behind domain-defined contracts. Commands orchestrate domain operations. The ubiquitous language of the business domain drives naming conventions across all layers.

- **[CQRS (Command Query Responsibility Segregation)](https://martinfowler.com/bliki/CQRS.html):** Write operations flow through Commands (which mutate state), while read operations flow through Controllers that query repositories directly or via services. This separation enables independent optimization of reads and writes.

## Concepts

### Rule Types

1. **Domain rules (enterprise rules):** Internal business rules that are independent of external services and do not change when external systems change.
2. **Application rules:** Rules that define how the application behaves based on the state of domain entities. They orchestrate domain rules.
3. **Adapter rules:** Conversions and adaptations between external APIs, databases, and domain entity classes.

### Principal Actors

| Actor            | Responsibility                                                                                                    |
|------------------|-------------------------------------------------------------------------------------------------------------------|
| **Entities**     | Core business objects. Contain domain rules about what the modeled objects are and how they behave.               |
| **Controllers**  | Bridge between the view (client) and the business layer. Receive requests and delegate to commands.               |
| **Commands**     | Implement business/feature logic and apply domain rules.                                                          |
| **Services**     | Handle application-level concerns: parsing, conversions, transformations.                                         |
| **Repositories** | Abstract all data access. Changing the database technology requires modifying only the repository implementation. |

## File Structure

The structure scales with complexity. Below are three tiers, from simplest to most complete:

### Minimal (3 layers)

```
domain/                   (contracts)
  commands/                 business logic (no separate contract)
  entities/
  repositories/
infrastructure/           (implementations)
  repositories/             prefixed with tool name
```

### Intermediate (with Services)

```
domain/                   (contracts)
  commands/
  entities/
  repositories/             contract used by services
  services/                 contract used by commands
infrastructure/           (implementations)
  repositories/             prefixed with tool name
  services/                 prefixed with tool name
```

### Complete (with API layer)

```
domain/                   (contracts)
  commands/
  entities/
  repositories/             contract used by services
  services/                 contract used by commands
infrastructure/           (implementations)
  controllers/
    requests/
    responses/
    mappers/
  repositories/             prefixed with tool name
  services/                 prefixed with tool name
    mappers/
```

## Architectural Flows

### Request Flow

Illustrates how an HTTP request from the browser traverses the application layers:

### Mapping/Parsing Flow

Demonstrates how Mappers isolate layers from each other, preventing hard coupling between frameworks or external tools:

### Dependency Injection Flow

Shows how DI wires contracts to their implementations at runtime:

---

# Frontend Design

> **TL;DR:** The frontend follows a lightweight Clean Architecture with 5 layers: Domain, Service, Infrastructure, Presentation, and Main. Dependencies always point toward the Domain layer. The Main layer handles dependency injection.

## Layers

As a lightweight adaptation of Clean Architecture, the frontend is divided into five layers:

| Layer              | Responsibility                                                                         |
|--------------------|----------------------------------------------------------------------------------------|
| **Domain**         | The most abstract layer. Defines Entities and Contracts for all other layers.          |
| **Service**        | Provides communication with APIs and external services.                                |
| **Infrastructure** | Implements technology-specific interfaces.                                             |
| **Presentation**   | Contains all view-rendering code (components, pages, hooks).                           |
| **Main**           | The "wiring" layer. Instantiates Providers and Services, handles Dependency Injection. |

## Architecture Diagrams

### Layer Dependencies

### Applied to Feature Modules

### Dependency Direction

Dependencies always point toward the **Domain** layer, which is the most abstract and stable:

## State Management

Currently, state is managed locally within Presentation components. A dedicated State Layer may be introduced in the future as requirements evolve.

---

# Security Practices

> **TL;DR:** Follow OWASP Top 10 and MITRE ATT&CK guidelines. Run `make sast` before every push to detect secrets (Gitleaks), code vulnerabilities (CodeQL, Semgrep), IaC misconfigurations (Trivy), and Dockerfile issues (Hadolint). Never hard-code secrets -- use environment variables or secret managers.

## Overview

This document defines the security best practices that must be considered when designing, developing, and deploying applications. These practices are grounded in industry standards from OWASP and MITRE, and enforced through the team's SAST toolchain.

## Security Checklist

### 1. Input Validation and Sanitization (OWASP A1)

Ensure all user input is validated and sanitized to prevent injection attacks (SQL injection, XSS, command injection):

- Validate data types and input length.
- Filter or encode special characters.
- Use parameterized queries for database access.

### 2. Authentication and Authorization (OWASP A2)

Implement robust identity and access management:

- Enforce **multi-factor authentication (MFA)**.
- Implement **role-based access controls (RBAC)**.
- Secure session management with proper token expiration and rotation.

### 3. Encryption (OWASP A3)

Protect sensitive data both in transit and at rest:

- Use **HTTPS/TLS** for all data in transit.
- Encrypt sensitive data stored on disk or in databases.
- Manage encryption keys securely (never hard-code secrets).

### 4. Access Controls (OWASP A4)

Apply the **principle of least privilege** and **need-to-know** across all systems:

- Implement fine-grained access controls.
- Use access control lists (ACLs) where appropriate.
- Regularly review and revoke unnecessary privileges.

### 5. Security Testing (OWASP A5)

Integrate security testing into the development pipeline:

- Perform **penetration testing** regularly.
- Use **Dynamic Application Security Testing (DAST)** in CI/CD pipelines.
- Conduct vulnerability assessments on dependencies (e.g., `npm audit`, `snyk`).

### 6. Incident Response (MITRE ATT&CK)

Maintain a documented incident response plan:

- Define scope identification procedures.
- Establish damage containment protocols.
- Document recovery and post-mortem processes.

### 7. Logging and Monitoring (MITRE ATT&CK)

Implement comprehensive observability:

- Use **centralized logging** and event management systems.
- Deploy **intrusion detection and prevention systems (IDS/IPS)**.
- Set up alerting for anomalous behavior.

### 8. Network Security (MITRE ATT&CK)

Protect the network perimeter and internal communications:

- Implement **firewalls** and **network segmentation**.
- Use **VPNs** for remote access.
- Apply **network access controls** to limit lateral movement.

## SAST Pipeline (Static Application Security Testing)

All projects enforce security through a standardized SAST pipeline provided by the [Pipelines repository](https://github.com/rios0rios0/pipelines). The full list of available security tools, their configuration, and usage is documented in the [Available Tools & Scripts](https://github.com/rios0rios0/pipelines?tab=readme-ov-file#-available-tools--scripts) section of the pipelines README.

The pipeline is invoked locally via `make sast`, which runs the following tools in sequence:

| Tool                                                 | Category  | What It Detects                                                                                     |
|------------------------------------------------------|-----------|-----------------------------------------------------------------------------------------------------|
| **[CodeQL](https://codeql.github.com/)**             | SAST      | SQL injection, XSS, path traversal, insecure deserialization, and other vulnerability patterns      |
| **[Semgrep](https://semgrep.dev/)**                  | SAST      | OWASP Top 10 patterns, hardcoded secrets, language-specific anti-patterns, best practice violations |
| **[Trivy](https://trivy.dev/)**                      | IaC / SCA | Infrastructure-as-Code misconfigurations and dependency vulnerabilities                             |
| **[Hadolint](https://github.com/hadolint/hadolint)** | Linting   | Dockerfile best practice violations (unpinned base images, running as root, missing health checks)  |
| **[Gitleaks](https://gitleaks.io/)**                 | Secrets   | API keys, tokens, passwords, private keys, and other secrets committed to Git history               |

Beyond SAST, the pipelines repository also provides **SCA (Software Composition Analysis)** tools such as `govulncheck` (Go), `Safety` (Python), `OWASP Dependency-Check` (Java), and `yarn npm audit` (JavaScript). See the [full tool reference](https://github.com/rios0rios0/pipelines?tab=readme-ov-file#-available-tools--scripts) for details.

### Running Locally

```bash
make setup   # Clone/update the pipelines repository (first time only)
make sast    # Run the full SAST suite
```

All reports are generated in `build/reports/`. If any tool reports findings, fix them before pushing code.

### Secret Hygiene

Leaked secrets are one of the most common and dangerous security issues. Follow these rules:

- **Never** hard-code secrets, API keys, tokens, or passwords in source code.
- Use **environment variables** or **secret managers** (1Password, HashiCorp Vault, AWS Secrets Manager).
- If Gitleaks detects a leaked secret, **rotate the credential immediately** and remove it from the codebase.
- Use `.env` files for local development, but **never commit them** (add `.env` to `.gitignore`).

## Pre-Commit Hooks

While a formal pre-commit hook standard has not yet been adopted, the team strongly recommends running `make lint && make sast` before every commit. This serves as a manual pre-commit gate.

A future standard may use [pre-commit](https://pre-commit.com/) to automate the following hooks:

| Hook             | Tool                                                   | Purpose                                      |
|------------------|--------------------------------------------------------|----------------------------------------------|
| Secret detection | Gitleaks                                               | Prevent secrets from entering the repository |
| Linting          | Language-specific (golangci-lint, Black, ESLint, etc.) | Enforce code style                           |
| Static analysis  | Semgrep (targeted rules)                               | Catch common security anti-patterns          |
| Commit message   | Custom script                                          | Enforce `type(SCOPE): message` format        |

Until formalized, treat `make lint && make sast` as a **mandatory pre-push step**.

---

# CI/CD Pipeline

> **TL;DR:** A complete CI/CD pipeline consists of 10 stages, from source control to continuous feedback. All projects use a shared Makefile with `make lint`, `make test`, and `make sast` targets. The SAST suite runs CodeQL, Semgrep, Trivy, Hadolint, and Gitleaks. Always run these locally before pushing code.

## Overview

This document defines the standards for building Continuous Integration and Continuous Deployment (CI/CD) pipelines. The pipeline described below represents the ideal end-to-end flow; individual projects may implement a subset based on their requirements.

## Pipeline Stages

| Stage                            | Description                                                                                                   |
|----------------------------------|---------------------------------------------------------------------------------------------------------------|
| **1. Source Control Management** | Code is stored in a VCS (Git) and committed to the repository by developers.                                  |
| **2. Build**                     | Code is compiled and packaged into a deployable artifact by the build system (e.g., Jenkins, GitHub Actions). |
| **3. Unit Testing**              | Automated unit tests validate that individual components function correctly, catching bugs early.             |
| **4. Integration Testing**       | The code is integrated with other system components and tested for correct inter-module behavior.             |
| **5. Static Analysis (SAST)**    | Code quality, security, and secret detection checks are performed using the SAST toolchain (see below).       |
| **6. Deployment (Staging)**      | The artifact is deployed to a staging environment for further testing and validation.                         |
| **7. Acceptance Testing**        | Users and stakeholders verify that the system meets business requirements in the staging environment.         |
| **8. Release (Production)**      | After passing all acceptance tests, the artifact is released to production.                                   |
| **9. Monitoring**                | The production system is monitored for correct operation and early detection of issues.                       |
| **10. Continuous Feedback**      | Feedback from users and stakeholders drives improvements to the system and the pipeline itself.               |

## Local Quality Gates

All projects use a shared `Makefile` that imports targets from the [pipelines repository](https://github.com/rios0rios0/pipelines). Developers must run these targets locally before pushing code:

```bash
make lint    # Fix all linter issues
make test    # Run the full test suite
make sast    # Run the complete SAST security suite
```

Individual SAST tools can also be run separately when you only need to check a specific category:

```bash
make semgrep   # Run Semgrep pattern-based analysis only
make trivy     # Run Trivy IaC/dependency scanning only
make hadolint  # Run Hadolint Dockerfile linting only
make gitleaks  # Run Gitleaks secret detection only
```

**Never call tool binaries (e.g., `golangci-lint`, `pytest`, `eslint`, `semgrep`) directly.** Always use the Makefile targets (`make lint`, `make test`, `make sast`, or the per-tool SAST targets above), which invoke the [pipelines repository](https://github.com/rios0rios0/pipelines) scripts to load the correct configuration before running each tool.

## SAST Toolchain

All SAST, SCA, and quality tools are provided by the [Pipelines repository](https://github.com/rios0rios0/pipelines). For the full list of available tools, configurations, and usage instructions, see the [Available Tools & Scripts](https://github.com/rios0rios0/pipelines?tab=readme-ov-file#-available-tools--scripts) section of the pipelines README.

The `make sast` target orchestrates five security tools that run both locally and in CI pipelines:

| Tool                                                 | Purpose                                                                                                     | Output              |
|------------------------------------------------------|-------------------------------------------------------------------------------------------------------------|---------------------|
| **[CodeQL](https://codeql.github.com/)**             | Static Application Security Testing -- detects SQL injection, XSS, path traversal, insecure deserialization | SARIF report        |
| **[Semgrep](https://semgrep.dev/)**                  | Pattern-based static analysis with OWASP Top 10, secrets, and best-practice rule sets                       | JSON report         |
| **[Trivy](https://trivy.dev/)**                      | Infrastructure-as-Code misconfiguration scanning and dependency vulnerability scanning                      | SARIF / JSON report |
| **[Hadolint](https://github.com/hadolint/hadolint)** | Dockerfile linting and best practice enforcement                                                            | SARIF report        |
| **[Gitleaks](https://gitleaks.io/)**                 | Secret and credential detection across the entire Git history                                               | JSON report         |

All reports are generated in the `build/reports/` directory.

### False Positive Management

Each tool supports project-level configuration for suppressing false positives:

| Tool     | Configuration File                       |
|----------|------------------------------------------|
| CodeQL   | `.codeql-false-positives`                |
| Semgrep  | `.semgrepignore`, `.semgrepexcluderules` |
| Trivy    | `.trivyignore`                           |
| Hadolint | `.hadolint.yaml`                         |
| Gitleaks | `.gitleaks.toml`                         |

## Pre-Commit Hooks

While no formal pre-commit hook standard has been defined yet, developers are strongly encouraged to run `make lint` and `make sast` before every commit. These targets serve as the equivalent of a pre-commit gate, catching issues before they enter the repository.

A future standard may formalize this using tools like [pre-commit](https://pre-commit.com/) to automatically run the following checks on every commit:

- **Linting** (language-specific formatters and linters)
- **Secret detection** (Gitleaks)
- **Static analysis** (Semgrep with targeted rule sets)

Until a formal standard is adopted, treat `make lint && make sast` as a mandatory manual pre-push step.

## Key Principles

- **Automate everything.** Every stage that can be automated should be automated.
- **Fail fast.** Run the fastest and cheapest checks (linting, unit tests) first.
- **Shift left on security.** Run SAST tools locally and in CI, not just in production.
- **Immutable artifacts.** Build once, deploy the same artifact to every environment.
- **Rollback capability.** Every deployment must support quick rollback to the previous version.

---

# Documentation & Change Control

> **TL;DR:** Every project must maintain a `CHANGELOG.md` (following [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)) and a `README.md`. Every code change must be accompanied by the corresponding documentation update -- changelog **always**, README and other docs (e.g., `.github/copilot-instructions.md`) **whenever behavior, configuration, or architecture changes**.

## Overview

Documentation and change control are integral parts of the engineering workflow, not afterthoughts. A well-maintained changelog and README provide traceability for stakeholders, reduce onboarding friction for new team members, and ensure that the state of the project is always understandable from its documentation alone.

**Every change introduced to a project must include updates to the relevant documentation files.** This is enforced as part of the development workflow, not as a separate task.

## Required Documentation Files

Every project must contain at minimum:

| File                                  | Purpose                                                      | Update Frequency                                             |
|---------------------------------------|--------------------------------------------------------------|--------------------------------------------------------------|
| **`README.md`**                       | Describes the project, its usage, setup, and architecture    | When behavior, configuration, CLI, or setup changes          |
| **`CONTRIBUTING.md`**                 | Guides contributors on prerequisites, workflow, and standards | When prerequisites, workflow, or project structure changes    |
| **`CHANGELOG.md`**                    | Records all notable changes, organized by version            | **Every change** (always)                                    |
| **`.github/copilot-instructions.md`** | AI assistant context for the project structure and workflows | When architecture, commands, or development workflow changes |

**Templates are available for standardized project setup:**

## Changelog Standard

All changelogs must follow the [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) format, combined with [Semantic Versioning](https://semver.org/).

### Format

```markdown
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- added new feature X that does Y

### Changed
- changed behavior of Z to handle edge case W

### Fixed
- fixed a bug where A caused B

## [1.0.0] - 2026-01-15

### Added
- added initial project setup with Clean Architecture
```

### Change Categories

| Category       | When to Use                                       |
|----------------|---------------------------------------------------|
| **Added**      | New features, new files, new capabilities         |
| **Changed**    | Modifications to existing functionality           |
| **Deprecated** | Features that will be removed in a future version |
| **Removed**    | Features that were removed                        |
| **Fixed**      | Bug fixes                                         |
| **Security**   | Vulnerability fixes                               |

### Writing Rules

- **Write for humans, not machines.** Describe *what changed and why*, not implementation details.
- **Use simple past tense.** "added", "changed", "fixed", "removed" -- consistent with the commit message standard.
- **Start each entry with a lowercase verb.** Example: `- added automatic Dockerfile image tag update`.
- **Be specific.** Bad: `- updated dependencies`. Good: `- added JavaScript updater supporting npm, yarn, and pnpm projects`.
- **Link to issues or PRs** when the change is non-trivial.
- **Group related changes** in a single entry rather than listing every file touched.

### The `[Unreleased]` Section

All in-progress changes go under `[Unreleased]`. When a release is cut:

1. Create a branch `bump/x.x.x`.
2. Move entries from `[Unreleased]` to a new version heading with the release date.
3. Open a Pull Request targeting `main`.
4. After merge, create a Git tag for the version.

## README Standard

The `README.md` must accurately describe the current state of the project. Update it whenever:

- A new feature changes how users interact with the project.
- CLI commands, flags, or configuration options are added, changed, or removed.
- Setup instructions, prerequisites, or environment requirements change.
- The project structure or architecture changes significantly.
- New dependencies or integrations are introduced.

### Recommended Sections

| Section                              | Purpose                                       |
|--------------------------------------|-----------------------------------------------|
| **Title and description**            | One-line summary of what the project does     |
| **Quick start / Installation**       | How to get running in under 5 minutes         |
| **Usage**                            | Commands, configuration, and examples         |
| **Architecture / Project structure** | High-level overview of directories and layers |
| **Development**                      | How to build, test, and contribute            |
| **References**                       | Links to external documentation               |

## AI Assistant Instructions

Projects that use AI-assisted development (GitHub Copilot, Cursor, etc.) should maintain a `.github/copilot-instructions.md` file. This file provides the AI with project-specific context about:

- Project purpose and architecture
- Build, test, and lint commands with expected timings
- Repository structure and key files
- Development workflow and validation steps
- Testing infrastructure and conventions

Update this file whenever the development workflow, architecture, or key commands change.

## Workflow Integration

Documentation updates must be part of the same commit or PR that introduces the change:

1. **Write the code change.**
2. **Update `CHANGELOG.md`** -- add an entry under `[Unreleased]` describing what changed.
3. **Update `README.md`** -- if the change affects usage, setup, or architecture.
4. **Update `.github/copilot-instructions.md`** -- if the change affects build commands, project structure, or development workflow.
5. **Commit everything together.** Documentation and code ship as one unit.

**Never merge a PR that introduces user-facing or architectural changes without the corresponding documentation update.**

## Automation with AutoBump

[AutoBump](https://github.com/rios0rios0/autobump) is a CLI tool that automates the release step of the changelog workflow. When the `[Unreleased]` section is ready to ship, AutoBump detects the project language, moves unreleased entries into a new versioned section with the current date, updates language-specific version files (e.g., `go.mod`, `package.json`, `pyproject.toml`, `build.gradle`), commits, pushes, and opens a merge/pull request -- all in a single command.

It supports Go, Java, Python, TypeScript, and C# projects, with automatic language detection, and works across GitHub, GitLab, and Azure DevOps.

**AutoBump does not replace the discipline of writing changelog entries.** The `CHANGELOG.md`, `README.md`, and other documentation files must already exist and be maintained by the team as part of every change. AutoBump only automates the versioning and release ceremony -- not the content creation.

---

# README Template

> **TL;DR:** Copy this template when creating a new project. Replace all `{PLACEHOLDERS}` with project-specific values. Remove optional sections that do not apply. Every project must have a README that follows this structure.

## Overview

This template defines the standard structure for all project README files. It ensures consistency across repositories and covers the sections described in the Documentation & Change Control guide. Copy the raw template below, replace placeholders, and remove optional sections wrapped in HTML comments.

## Placeholder Reference

| Placeholder      | Description                        | Example                |
|------------------|------------------------------------|------------------------|
| `{project-name}` | Human-readable project name        | `AutoBump`             |
| `{ORG}`          | GitHub organization or user        | `rios0rios0`           |
| `{REPO}`         | Repository name                    | `autobump`             |
| `{PACKAGE}`      | Published package name (npm, PyPI) | `@rios0rios0/autobump` |

## Template

````markdown
<h1 align="center">{project-name}</h1>
<p align="center">
    <a href="https://github.com/{ORG}/{REPO}/releases/latest">
        <img src="https://img.shields.io/github/release/{ORG}/{REPO}.svg?style=for-the-badge&logo=github" alt="Latest Release"/></a>
    <a href="https://github.com/{ORG}/{REPO}/blob/main/LICENSE">
        <img src="https://img.shields.io/github/license/{ORG}/{REPO}.svg?style=for-the-badge&logo=github" alt="License"/></a>
    <a href="https://github.com/{ORG}/{REPO}/actions/workflows/default.yaml">
        <img src="https://img.shields.io/github/actions/workflow/status/{ORG}/{REPO}/default.yaml?branch=main&style=for-the-badge&logo=github" alt="Build Status"/></a>
    <!-- Add SonarCloud badges when configured -->
    <!--
    <a href="https://sonarcloud.io/summary/overall?id={ORG}_{REPO}">
        <img src="https://img.shields.io/sonar/coverage/{ORG}_{REPO}?server=https%3A%2F%2Fsonarcloud.io&style=for-the-badge&logo=sonarqubecloud" alt="Coverage"/></a>
    <a href="https://sonarcloud.io/summary/overall?id={ORG}_{REPO}">
        <img src="https://img.shields.io/sonar/quality_gate/{ORG}_{REPO}?server=https%3A%2F%2Fsonarcloud.io&style=for-the-badge&logo=sonarqubecloud" alt="Quality Gate"/></a>
    -->
    <!-- Language-specific badges (uncomment the one that applies) -->
    <!-- Go: -->
    <!-- <a href="https://pkg.go.dev/github.com/{ORG}/{REPO}"><img src="https://img.shields.io/badge/go-reference-007d9c?style=for-the-badge&logo=go" alt="Go Reference"/></a> -->
    <!-- Java: -->
    <!-- <a href="https://central.sonatype.com/artifact/{GROUP_ID}/{ARTIFACT_ID}"><img src="https://img.shields.io/maven-central/v/{GROUP_ID}/{ARTIFACT_ID}?style=for-the-badge&logo=apachemaven" alt="Maven Central"/></a> -->
    <!-- Python: -->
    <!-- <a href="https://pypi.org/project/{PACKAGE}"><img src="https://img.shields.io/pypi/v/{PACKAGE}?style=for-the-badge&logo=pypi" alt="PyPI"/></a> -->
    <!-- JavaScript/TypeScript: -->
    <!-- <a href="https://www.npmjs.com/package/{PACKAGE}"><img src="https://img.shields.io/npm/v/{PACKAGE}?style=for-the-badge&logo=npm" alt="npm"/></a> -->
    <!-- OpenSSF: -->
    <!-- <a href="https://www.bestpractices.dev/projects/{ID}"><img src="https://img.shields.io/cii/level/{ID}?style=for-the-badge&logo=opensourceinitiative" alt="OpenSSF Best Practices"/></a> -->
</p>

One to two sentence description of what this project does and why it exists.

## Features

- **Feature One**: brief explanation
- **Feature Two**: brief explanation
- **Feature Three**: brief explanation

<!-- OPTIONAL: only for multi-ecosystem tools -->
<!--
## Supported Ecosystems

| Ecosystem | Detection | Version File |
|-----------|-----------|--------------|
| Go        | `go.mod`  | `go.mod`     |
-->

## Installation

<!-- Keep only the section that applies to your project's language -->

<!-- Go CLI or library -->
```bash
go install github.com/{ORG}/{REPO}@latest
# or for libraries:
go get github.com/{ORG}/{REPO}
```

<!-- Java (Gradle) -->
```groovy
implementation '{GROUP_ID}:{ARTIFACT_ID}:{VERSION}'
```

<!-- Python -->
```bash
pdm add {PACKAGE}
# or:
pip install {PACKAGE}
```

<!-- JavaScript/TypeScript -->
```bash
npm install {PACKAGE}
```

Download pre-built binaries from the [releases page](https://github.com/{ORG}/{REPO}/releases).

<!-- OPTIONAL: only for CLIs with config files -->
<!--
## Configuration

Create `~/.config/{REPO}.yaml`:

```yaml
key: 'value'
```
-->

## Usage

```bash
{REPO} [flags] [arguments]
```

Brief explanation of the primary workflow.

<!-- OPTIONAL: only for libraries and complex CLIs -->
<!--
## Architecture

```
{REPO}/
├── domain/           # core business objects and contracts
├── infrastructure/   # implementations
└── ...
```
-->

<!-- OPTIONAL: only for libraries exposing interfaces -->
<!--
## API Reference

- **`InterfaceName`**: what it does
-->

## Contributing

Contributions are welcome. See CONTRIBUTING.md for guidelines.

## License

See [LICENSE](LICENSE) file for details.
````

---

# CONTRIBUTING Template

> **TL;DR:** Copy this template when creating a new project. Replace all `{PLACEHOLDERS}` with project-specific values. Remove optional sections that do not apply. Every project must have a CONTRIBUTING.md that follows this structure.

## Overview

This template defines the standard structure for all project CONTRIBUTING files. It ensures every repository provides clear, consistent onboarding for contributors. Copy the raw template below, replace placeholders, and remove optional sections wrapped in HTML comments.

## Placeholder Reference

| Placeholder         | Description                          | Example                                                                   |
|---------------------|--------------------------------------|---------------------------------------------------------------------------|
| `{LANGUAGE}`        | Primary programming language         | `Go 1.23`, `Java 21`, `Python 3.12`, `Node.js 20`                         |
| `{INSTALL_COMMAND}` | Dependency install command           | `go mod download`, `./gradlew dependencies`, `pdm install`, `npm install` |
| `{EXTENSION_TYPE}`  | Plugin/provider type (if applicable) | `Provider`, `Updater`, `Plugin`                                           |

## Language-Specific Prerequisites

Use the prerequisite block that matches your project:

| Language                  | Prerequisites                                                                          |
|---------------------------|----------------------------------------------------------------------------------------|
| **Go**                    | Go 1.23+, Make                                                                         |
| **Java**                  | Java 21+ (Eclipse Temurin), Gradle (via wrapper), Make, Docker (for integration tests) |
| **Python**                | Python 3.12+, PDM, Make                                                                |
| **JavaScript/TypeScript** | Node.js 20+, npm, Make                                                                 |

## Template

````markdown
# Contributing

Contributions are welcome. By participating, you agree to maintain a respectful and constructive environment.

For coding standards, testing patterns, architecture guidelines, commit conventions, and all
development practices, refer to the **[Development Guide](https://github.com/rios0rios0/guide/wiki)**.

## Prerequisites

- {LANGUAGE}
- [Make](https://www.gnu.org/software/make/)
<!-- Add any other tools required by your project -->
<!-- Java projects: -->
<!-- - Docker (for integration tests with TestContainers) -->
<!-- Python projects: -->
<!-- - [PDM](https://pdm-project.org/) -->

## Development Workflow

1. Fork and clone the repository
2. Create a branch: `git checkout -b feat/my-change`
3. Install dependencies:
   ```bash
   {INSTALL_COMMAND}
   ```
4. Make your changes
5. Validate:
   ```bash
   make lint
   make test
   make sast
   ```
6. Update `CHANGELOG.md` under `[Unreleased]`
7. Commit following the [commit conventions](https://github.com/rios0rios0/guide/wiki/Git-Flow)
8. Open a pull request against `main`

<!-- OPTIONAL: only when the project requires environment variables or local services -->
<!--
## Local Environment

Copy `.env.example` to `.env` and fill in the required values:

```bash
cp .env.example .env
```

Start local services:

```bash
docker compose -f compose.dev.yaml up -d
```

| Variable | Description | Required |
|----------|-------------|----------|
| `DB_HOST` | Database hostname | Yes |
| `DB_PORT` | Database port | Yes |
-->

<!-- OPTIONAL: only for projects where contributors extend functionality -->
<!--
## Adding a New {EXTENSION_TYPE}

1. Create the implementation file following the [naming conventions](https://github.com/rios0rios0/guide/wiki/Code-Style)
2. Implement the required interface/contract
3. Register it in the dependency injection wiring
4. Add tests following the [testing guide](https://github.com/rios0rios0/guide/wiki/Tests)
5. Update `CHANGELOG.md` with an entry under `[Unreleased] > Added`
-->
````

---

# CHANGELOG Formatting

> **TL;DR:** Proper nouns are capitalized, code identifiers use backticks, acronyms and protocol names use their official casing, versions and library names use backticks.

## Rules

### 1. Proper Nouns -- Capitalize Correctly

Product names, company names, and personal names must use their official capitalization.

| Wrong      | Correct    |
|------------|------------|
| github     | GitHub     |
| kubernetes | Kubernetes |
| docker     | Docker     |
| terraform  | Terraform  |
| john doe   | John Doe   |
| elasticsearch | Elasticsearch |
| postgresql | PostgreSQL |
| javascript | JavaScript |
| typescript | TypeScript |

### 2. Code Identifiers -- Use Backticks

Class names, function names, variable names, file names, and any code reference must be wrapped in backticks.

| Wrong                                      | Correct                                        |
|--------------------------------------------|------------------------------------------------|
| added CreateUser command                   | added `CreateUser` command                     |
| fixed bug in handleRequest function        | fixed bug in `handleRequest` function          |
| renamed user_id to userId                  | renamed `user_id` to `userId`                  |
| updated Dockerfile                         | updated `Dockerfile`                           |
| changed settings.json configuration        | changed `settings.json` configuration          |

### 3. Acronyms and Protocol Names -- Use Official Casing

Technical acronyms must be uppercased. Branded protocol and technology names must use their official casing.

| Wrong  | Correct |
|--------|---------|
| http   | HTTP    |
| api    | API     |
| sql    | SQL     |
| css    | CSS     |
| html   | HTML    |
| jwt    | JWT     |
| grpc   | gRPC    |
| graphql | GraphQL |
| rest   | REST    |
| oauth  | OAuth   |

### 4. Versions -- Use Backticks

Version numbers (with or without the `v` prefix) must be wrapped in backticks.

| Wrong                          | Correct                            |
|--------------------------------|------------------------------------|
| bumped to v1.2.0               | bumped to `v1.2.0`                 |
| upgraded from 2.0.0 to 3.0.0  | upgraded from `2.0.0` to `3.0.0`  |

### 5. Library and Package Names -- Use Backticks

Dependency names, package names, and module names must be wrapped in backticks.

| Wrong                              | Correct                                |
|------------------------------------|----------------------------------------|
| upgraded lodash to 4.17.21         | upgraded `lodash` to `4.17.21`         |
| added gin-gonic as HTTP framework  | added `gin-gonic` as HTTP framework    |
| replaced moment with dayjs         | replaced `moment` with `dayjs`         |

## Examples

**Bad:**

```markdown
- added createUser endpoint using express framework
- upgraded golang to v1.23
- fixed sql injection in handleLogin
- integrated with github actions for ci/cd
```

**Good:**

```markdown
- added `CreateUser` endpoint using `express` framework
- upgraded Go to `v1.23`
- fixed SQL injection in `handleLogin`
- integrated with GitHub Actions for CI/CD
```

---

# Mapper Design Pattern (Replacing Switch/Case)

> **TL;DR:** Replace `switch`/`case` statements with a map (dictionary) that associates each key with a handler function. This eliminates branching, simplifies adding new cases, and follows the Open/Closed Principle.

## Overview

The `switch`/`case` construct is one of the most common sources of code that violates the **Open/Closed Principle** -- every time a new case is needed the existing function must be modified. The Mapper Design Pattern replaces the branching logic with a **lookup map** that associates each key with its handler. Adding a new case becomes adding a new entry to the map, without touching the dispatch logic.

### Benefits

| Benefit              | Explanation                                                              |
|----------------------|--------------------------------------------------------------------------|
| **Open/Closed**      | New cases are added by inserting a map entry, not by editing a function. |
| **Readability**      | The mapping is declarative -- you see all cases at a glance.             |
| **Testability**      | Each handler can be tested in isolation.                                 |
| **Single dispatch**  | The lookup + call replaces an unbounded chain of comparisons.            |

## Problem -- Switch/Case

The examples below use a notification system that sends messages through different channels. The `switch` version hard-codes every channel:

### JavaScript (Switch)

```javascript
function notify(channel, message) {
  switch (channel) {
    case "email":
      sendEmail(message);
      break;
    case "sms":
      sendSms(message);
      break;
    case "slack":
      sendSlack(message);
      break;
    default:
      throw new Error(`Unknown channel: ${channel}`);
  }
}
```

### Java (Switch)

```java
public class NotificationService {
    public void notify(final String channel, final String message) {
        switch (channel) {
            case "email" -> sendEmail(message);
            case "sms" -> sendSms(message);
            case "slack" -> sendSlack(message);
            default -> throw new IllegalArgumentException("Unknown channel: " + channel);
        }
    }
}
```

### Go (Switch)

```go
func notify(channel, message string) error {
	switch channel {
	case "email":
		return sendEmail(message)
	case "sms":
		return sendSms(message)
	case "slack":
		return sendSlack(message)
	default:
		return fmt.Errorf("unknown channel: %s", channel)
	}
}
```

### Python (Switch)

```python
def notify(channel: str, message: str) -> None:
    if channel == "email":
        send_email(message)
    elif channel == "sms":
        send_sms(message)
    elif channel == "slack":
        send_slack(message)
    else:
        raise ValueError(f"Unknown channel: {channel}")
```

## Solution -- Mapper Design Pattern

Replace the branching with a map from key to handler. The dispatch function looks up the key and calls the handler.

### JavaScript

```javascript
const notificationHandlers = new Map([
  ["email", sendEmail],
  ["sms", sendSms],
  ["slack", sendSlack],
]);

function notify(channel, message) {
  const handler = notificationHandlers.get(channel);
  if (!handler) {
    throw new Error(`Unknown channel: ${channel}`);
  }
  handler(message);
}
```

Adding a new channel (e.g., push notifications) requires only one line:

```javascript
notificationHandlers.set("push", sendPush);
```

### Java

```java
public class NotificationService {
    private static final Map<String, Consumer<String>> HANDLERS = Map.of(
            "email", NotificationService::sendEmail,
            "sms", NotificationService::sendSms,
            "slack", NotificationService::sendSlack
    );

    public void notify(final String channel, final String message) {
        final var handler = HANDLERS.get(channel);
        if (handler == null) {
            throw new IllegalArgumentException("Unknown channel: " + channel);
        }
        handler.accept(message);
    }
}
```

For richer handler contracts, define a functional interface or use a strategy class:

```java
@FunctionalInterface
public interface NotificationHandler {
    void send(String message);
}

public class NotificationService {
    private final Map<String, NotificationHandler> handlers;

    public NotificationService(final Map<String, NotificationHandler> handlers) {
        this.handlers = handlers;
    }

    public void notify(final String channel, final String message) {
        final var handler = handlers.get(channel);
        if (handler == null) {
            throw new IllegalArgumentException("Unknown channel: " + channel);
        }
        handler.send(message);
    }
}
```

This second form also supports **dependency injection** -- the handler map can be wired by Spring:

```java
@Configuration
public class NotificationConfig {
    @Bean
    public Map<String, NotificationHandler> notificationHandlers(
            final EmailHandler email,
            final SmsHandler sms,
            final SlackHandler slack) {
        return Map.of("email", email, "sms", sms, "slack", slack);
    }
}
```

### Go

```go
type NotificationHandler func(message string) error

var notificationHandlers = map[string]NotificationHandler{
	"email": sendEmail,
	"sms":   sendSms,
	"slack": sendSlack,
}

func notify(channel, message string) error {
	handler, ok := notificationHandlers[channel]
	if !ok {
		return fmt.Errorf("unknown channel: %s", channel)
	}
	return handler(message)
}
```

For handlers that require dependencies, use a struct with a method:

```go
type NotificationHandler interface {
	Send(message string) error
}

type NotificationService struct {
	handlers map[string]NotificationHandler
}

func NewNotificationService(handlers map[string]NotificationHandler) *NotificationService {
	return &NotificationService{handlers: handlers}
}

func (s *NotificationService) Notify(channel, message string) error {
	handler, ok := s.handlers[channel]
	if !ok {
		return fmt.Errorf("unknown channel: %s", channel)
	}
	return handler.Send(message)
}
```

### Python

```python
from typing import Callable

notification_handlers: dict[str, Callable[[str], None]] = {
    "email": send_email,
    "sms": send_sms,
    "slack": send_slack,
}

def notify(channel: str, message: str) -> None:
    handler = notification_handlers.get(channel)
    if handler is None:
        raise ValueError(f"Unknown channel: {channel}")
    handler(message)
```

For handlers that need dependencies, use a class-based approach:

```python
from abc import ABC, abstractmethod

class NotificationHandler(ABC):
    @abstractmethod
    def send(self, message: str) -> None: ...

class NotificationService:
    def __init__(self, handlers: dict[str, NotificationHandler]) -> None:
        self._handlers = handlers

    def notify(self, channel: str, message: str) -> None:
        handler = self._handlers.get(channel)
        if handler is None:
            raise ValueError(f"Unknown channel: {channel}")
        handler.send(message)
```

## When to Use

| Scenario                                               | Recommendation   |
|--------------------------------------------------------|------------------|
| 2-3 simple cases unlikely to grow                      | `switch` is fine |
| 4+ cases or cases that change often                    | Use the Mapper   |
| Handlers with different dependencies or configurations | Use the Mapper   |
| Cases determined at runtime (e.g., plugins)            | Use the Mapper   |

---

# Forking Technique

> **TL;DR:** Reserve `main` for the upstream community version and use `custom` as the working branch. Rebase `custom` onto updated `main` when upstream releases new versions. Tag fork releases with an incremental fourth digit (e.g., `1.0.0.1`).

## Overview

When forking an open-source project, it is common for upstream maintainers to take significant time reviewing and merging community contributions. This standard defines a forking strategy that maintains compatibility with upstream releases while preserving custom modifications.

## Strategy

### Branch Convention

- The **`main`** branch mirrors the upstream community version. It is updated by rebasing on upstream releases.
- The **`custom`** branch serves as the team's working branch (equivalent to `main` for internal purposes).

### Synchronization

To stay current with the upstream project:

1. Update `main` by rebasing on the latest upstream release.
2. Rebase `custom` onto the updated `main`, placing custom modifications on top of the newest version.

### Versioning

Fork versions use an **incremental fourth digit** appended to the upstream version:

| Scenario                                  | Version   |
|-------------------------------------------|-----------|
| Upstream at `1.0.0`, fork synced          | `1.0.0.0` |
| New fork release                          | `1.0.0.1` |
| Upstream updates to `1.0.1`, fork rebased | `1.0.1.0` |
| Next fork release                         | `1.0.1.1` |

The fourth digit resets to `0` each time the fork is rebased on a new upstream version.

## Caveats

If your CI/CD or release tooling does not support four-segment version numbers (`X.Y.Z.N`), use a dash separator instead: `X.Y.Z-N`.

---

# Bulk Operations Across Multiple Repositories

> **TL;DR:** Use a 4-phase workflow (Discover, Apply, Git Ops, Create PRs) to apply changes across all repositories under a workspace root. Auto-detect the hosting vendor from `git remote` to create PRs with the correct CLI. Always stash, fetch, rebase, and restore to preserve local state.

## Overview

When managing many repositories within an organization, it is common to apply the same change across all of them -- updating configuration files, fixing security findings, bumping dependencies, or standardizing tooling. This cookbook defines a repeatable, vendor-agnostic workflow for bulk operations that works with GitHub, Azure DevOps, GitLab, or any Git hosting provider.

## Workspace Layout

Set `ORG_ROOT` to the directory that contains the repositories. The structure can be flat or nested -- the discovery phase scans for `.git` directories at any depth.

**Flat structure (e.g., GitHub organizations):**

```
<ORG_ROOT>/
├── repo-alpha/
├── repo-beta/
└── repo-gamma/
```

**Nested structure (e.g., Azure DevOps projects):**

```
<ORG_ROOT>/
├── project-a/
│   ├── repo-one/
│   └── repo-two/
├── project-b/
│   └── repo-three/
```

**Key rule:** repos are identified by the presence of a `.git` directory, regardless of nesting depth.

## Phase 1: Discovery -- Find All Repos

Use a Python script (executed with `/usr/bin/python3`) to discover all git repositories under `ORG_ROOT`. The `max_depth` parameter controls how deep to scan (set to `None` for unlimited depth):

```python
import os

ORG_ROOT = "/path/to/your/workspace"  # Set this to your workspace root
MAX_DEPTH = 3  # Maximum directory depth to scan (None for unlimited)

def discover_repos(root, max_depth=None):
    repos = []
    root_depth = root.rstrip(os.sep).count(os.sep)
    for dirpath, dirnames, _ in os.walk(root):
        current_depth = dirpath.rstrip(os.sep).count(os.sep) - root_depth
        if max_depth is not None and current_depth >= max_depth:
            dirnames.clear()
            continue
        if ".git" in dirnames:
            repos.append(dirpath)
            dirnames.remove(".git")  # Don't descend into .git
            dirnames.clear()  # Don't descend into subdirectories of a repo
    return sorted(repos)

repos = discover_repos(ORG_ROOT, MAX_DEPTH)
```

## Phase 2: Apply Changes

For each repository, apply the required file modifications using Claude's Read, Write, and Edit tools.

**Important:** Track which repos were actually modified. Only proceed to Phase 3 for repos with real changes.

## Phase 3: Git Operations -- Branch, Commit, Push

### Critical Environment Setup

Git commands MUST use the full path and SSH must be configured for non-interactive mode:

```python
import subprocess, os

GIT = "/usr/bin/git"
SSH_CMD = "ssh -o BatchMode=yes -o ConnectTimeout=15"

def git(args, cwd, timeout=120):
    env = os.environ.copy()
    env["GIT_SSH_COMMAND"] = SSH_CMD
    try:
        r = subprocess.run(
            [GIT] + args, cwd=cwd,
            capture_output=True, text=True, timeout=timeout, env=env
        )
        return r.returncode, r.stdout.strip(), r.stderr.strip()
    except subprocess.TimeoutExpired:
        return -1, "", "TIMEOUT"
```

### Per-Repository Git Workflow

The workflow MUST preserve pre-existing local work and ensure the new branch is always created from an **up-to-date** default branch.

```python
BRANCH = "chore/your-branch-name"  # Follow branch naming from the Git Flow rule

def restore(repo_path, original_branch, has_stash):
    git(["checkout", original_branch], repo_path)
    if has_stash:
        git(["stash", "pop"], repo_path)

# 1. Detect the default branch
rc, out, _ = git(["symbolic-ref", "refs/remotes/origin/HEAD"], repo_path)
default_branch = out.replace("refs/remotes/origin/", "") if rc == 0 and out else "main"

# 2. Save the current branch
rc, original_branch, _ = git(["branch", "--show-current"], repo_path)
if not original_branch:
    original_branch = default_branch

# 3. Stash any pre-existing uncommitted changes
rc, stash_out, _ = git(["stash", "push", "-m", "bulk-op-auto-stash"], repo_path)
has_stash = "No local changes" not in stash_out

# 4. Switch to the default branch
git(["checkout", default_branch], repo_path)

# 5. Fetch ALL remotes and rebase onto the latest remote default branch
#    CRITICAL: without this the branch is created from stale local state
git(["fetch", "--all"], repo_path, timeout=120)
git(["pull", "--rebase"], repo_path, timeout=120)

# 6. Delete old feature branch if it exists (idempotency for reruns)
git(["branch", "-D", BRANCH], repo_path)

# 7. Create the feature branch from the now up-to-date default branch
git(["checkout", "-b", BRANCH], repo_path)

# 8. === Apply your changes here ===

# 9. Stage all changes
git(["add", "-A"], repo_path)

# 10. Verify there are staged changes
rc, diff, _ = git(["diff", "--cached", "--name-only"], repo_path)
if not diff:
    restore(repo_path, original_branch, has_stash)
    continue

# 11. Commit (follow commit message standards from the Git Flow rule)
msg = "chore(maintenance): your commit message in simple past tense"
rc, _, err = git(["commit", "-m", msg], repo_path)
if rc != 0:
    rc, _, err = git(["commit", "--no-verify", "-m", msg], repo_path)
    if rc != 0:
        print(f"FAIL(commit: {err})")
        restore(repo_path, original_branch, has_stash)
        continue

# 12. Push the branch (force is safe because this is our own new branch)
rc, _, err = git(["push", "-u", "origin", BRANCH, "--force"], repo_path, timeout=120)
if rc != 0:
    print(f"FAIL(push: {err})")
    restore(repo_path, original_branch, has_stash)
    continue

# 13. Restore the original branch and any stashed work
restore(repo_path, original_branch, has_stash)
```

### Why `fetch --all && pull --rebase` Is Mandatory

Without fetching and rebasing, the local default branch may be weeks or months behind the remote. This causes PR diffs with unrelated old changes, merge conflicts, and reviewer confusion. **NEVER skip the fetch + rebase step.**

### Handling Pre-commit Hooks

Some repos have pre-commit hooks that may block commits. The workflow above handles this with an automatic `--no-verify` retry.

## Phase 4: Create Pull Requests

### Detecting Vendor from Remote URL

Parse the remote URL to determine which CLI to use for PR creation:

```python
import re

def detect_vendor(repo_path):
    """Detect the Git hosting vendor from the remote URL."""
    rc, url, _ = git(["remote", "get-url", "origin"], repo_path)
    if rc != 0:
        return "unknown", url

    if "github.com" in url:
        return "github", url
    elif "dev.azure.com" in url or "ssh.dev.azure.com" in url or "visualstudio.com" in url:
        return "azure-devops", url
    elif "gitlab.com" in url or "gitlab" in url.lower():
        return "gitlab", url
    else:
        return "unknown", url
```

### Extracting Owner and Repo Name

```python
def extract_repo_info(url):
    """Extract owner/org and repo name from a remote URL."""
    # SSH format: git@github.com:owner/repo.git
    ssh_match = re.match(r"git@[^:]+:(.+?)(?:\.git)?$", url)
    if ssh_match:
        parts = ssh_match.group(1).split("/")
        return "/".join(parts[:-1]), parts[-1]

    # HTTPS format: https://github.com/owner/repo.git
    https_match = re.match(r"https?://[^/]+/(.+?)(?:\.git)?$", url)
    if https_match:
        parts = https_match.group(1).split("/")
        return "/".join(parts[:-1]), parts[-1]

    return None, None
```

### Creating PRs by Vendor

**GitHub** -- use `gh pr create`:

```bash
gh pr create \
  --repo "<owner>/<repo>" \
  --head "<branch>" \
  --base "<default-branch>" \
  --title "<pr-title>" \
  --body "<pr-body>"
```

**Azure DevOps** -- use `az repos pr create`:

```bash
az repos pr create \
  --organization "https://dev.azure.com/<org>" \
  --project "<project>" \
  --repository "<repo>" \
  --source-branch "<branch>" \
  --target-branch "<default-branch>" \
  --title "<pr-title>" \
  --description "<pr-body>"
```

**GitLab** -- use `glab mr create`:

```bash
glab mr create \
  --repo "<owner>/<repo>" \
  --source-branch "<branch>" \
  --target-branch "<default-branch>" \
  --title "<pr-title>" \
  --description "<pr-body>"
```

**Unknown vendor** -- skip PR creation and report the branch:

```
PR skipped for <repo>: branch <branch> pushed. Create PR manually.
```

### Batch PR Creation

Create PRs in parallel batches of up to 10 at a time.

## Complete End-to-End Script Template

```python
#!/usr/bin/env python3
"""Bulk operation: <describe what this does>"""
import subprocess, os, sys, json, re

GIT = "/usr/bin/git"
ORG_ROOT = "/path/to/your/workspace"  # Set this to your workspace root
BRANCH = "chore/your-branch-name"
MAX_DEPTH = 3

def git(args, cwd, timeout=120):
    env = os.environ.copy()
    env["GIT_SSH_COMMAND"] = "ssh -o BatchMode=yes -o ConnectTimeout=15"
    try:
        r = subprocess.run([GIT] + args, cwd=cwd, capture_output=True, text=True, timeout=timeout, env=env)
        return r.returncode, r.stdout.strip(), r.stderr.strip()
    except subprocess.TimeoutExpired:
        return -1, "", "TIMEOUT"

def restore(repo_path, original_branch, has_stash):
    git(["checkout", original_branch], repo_path)
    if has_stash:
        git(["stash", "pop"], repo_path)

def discover_repos(root, max_depth=None):
    repos = []
    root_depth = root.rstrip(os.sep).count(os.sep)
    for dirpath, dirnames, _ in os.walk(root):
        current_depth = dirpath.rstrip(os.sep).count(os.sep) - root_depth
        if max_depth is not None and current_depth >= max_depth:
            dirnames.clear()
            continue
        if ".git" in dirnames:
            repos.append(dirpath)
            dirnames.remove(".git")
            dirnames.clear()
    return sorted(repos)

def detect_vendor(repo_path):
    rc, url, _ = git(["remote", "get-url", "origin"], repo_path)
    if rc != 0:
        return "unknown", url
    if "github.com" in url:
        return "github", url
    elif "dev.azure.com" in url or "ssh.dev.azure.com" in url:
        return "azure-devops", url
    elif "gitlab" in url.lower():
        return "gitlab", url
    return "unknown", url

# --- Phase 1: Discover repos ---
repos = discover_repos(ORG_ROOT, MAX_DEPTH)

# --- Phase 2 & 3: Apply changes, branch, commit, push ---
pushed = []
for repo_path in repos:
    name = repo_path.replace(ORG_ROOT + "/", "")
    sys.stdout.write(f"{name}: ")
    sys.stdout.flush()

    rc, out, _ = git(["symbolic-ref", "refs/remotes/origin/HEAD"], repo_path)
    defb = out.replace("refs/remotes/origin/", "") if rc == 0 and out else "main"

    rc, original_branch, _ = git(["branch", "--show-current"], repo_path)
    if not original_branch:
        original_branch = defb

    rc, stash_out, _ = git(["stash", "push", "-m", "bulk-op-auto-stash"], repo_path)
    has_stash = "No local changes" not in stash_out

    git(["checkout", defb], repo_path)
    git(["fetch", "--all"], repo_path, timeout=120)
    rc, _, err = git(["pull", "--rebase"], repo_path, timeout=120)
    if rc != 0:
        print(f"FAIL(pull --rebase: {err})")
        restore(repo_path, original_branch, has_stash)
        continue

    git(["branch", "-D", BRANCH], repo_path)
    git(["checkout", "-b", BRANCH], repo_path)

    # TODO: === Apply your changes here ===

    rc, st, _ = git(["status", "--porcelain"], repo_path)
    if not st:
        print("SKIP(clean)")
        restore(repo_path, original_branch, has_stash)
        continue

    git(["add", "-A"], repo_path)

    rc, diff, _ = git(["diff", "--cached", "--name-only"], repo_path)
    if not diff:
        print("SKIP(nothing staged)")
        restore(repo_path, original_branch, has_stash)
        continue

    msg = "chore(maintenance): <your commit message>"
    rc, _, err = git(["commit", "-m", msg], repo_path)
    if rc != 0:
        rc, _, err = git(["commit", "--no-verify", "-m", msg], repo_path)
        if rc != 0:
            print(f"FAIL(commit: {err})")
            restore(repo_path, original_branch, has_stash)
            continue

    rc, _, err = git(["push", "-u", "origin", BRANCH, "--force"], repo_path, timeout=120)
    if rc != 0:
        print(f"FAIL(push: {err})")
        restore(repo_path, original_branch, has_stash)
        continue

    vendor, url = detect_vendor(repo_path)
    pushed.append({
        "path": name,
        "vendor": vendor,
        "url": url,
        "default": defb,
        "title": msg,
    })
    print("OK")
    restore(repo_path, original_branch, has_stash)

with open("/tmp/pushed_repos.json", "w") as f:
    json.dump(pushed, f, indent=2)
print(f"\nDone: {len(pushed)} repos pushed. Saved to /tmp/pushed_repos.json")
```

After running the script, read `/tmp/pushed_repos.json` and create PRs using the appropriate CLI for each vendor.

## Common Pitfalls & Solutions

| Problem                             | Cause                                             | Solution                                                                             |
|-------------------------------------|---------------------------------------------------|--------------------------------------------------------------------------------------|
| PR has merge conflicts / stale diff | Branch created from outdated local default branch | Always run `git fetch --all && git pull --rebase` before creating the feature branch |
| Pre-existing user work is lost      | Didn't stash before switching branches            | Always `git stash push` before checkout, `git stash pop` after restoring             |
| Repo left on wrong branch           | Didn't restore original branch                    | Always save `original_branch` and checkout back to it on all code paths              |
| `git push` hangs forever            | SSH waiting for interactive auth                  | Set `GIT_SSH_COMMAND="ssh -o BatchMode=yes -o ConnectTimeout=15"`                    |
| `git` or `python3` not found        | Not in PATH in subprocess                         | Use full paths: `/usr/bin/git`, `/usr/bin/python3`                                   |
| Shell loop silently fails           | zsh variable conflicts (`status` is read-only)    | Use Python instead of shell loops                                                    |
| Commit fails with hook error        | Repo has git hooks                                | Retry with `--no-verify` flag                                                        |
| Wrong target branch                 | Not all repos use `main`                          | Always detect via `git symbolic-ref refs/remotes/origin/HEAD`                        |

## Checklist for Every Bulk Operation

1. **Discover** all repos under the workspace root
2. **Save state** -- record current branch and stash uncommitted changes
3. **Sync** -- switch to default branch, `fetch --all`, `pull --rebase`
4. **Branch** -- create feature branch from the now-current default HEAD
5. **Apply** file changes
6. **Track** which repos have actual changes (`git status --porcelain`)
7. **Commit** with proper message format (see the Git Flow rule)
8. **Push** with SSH batch mode and timeout
9. **Restore** -- switch back to original branch and pop the stash
10. **Print progress** for every repo (OK / SKIP / FAIL + reason)
11. **Create PRs** using the detected vendor CLI in batches of 10
12. **Report** final summary with PR numbers
