Scaffold a new Go backend project following Clean Architecture with domain/infrastructure separation, manual constructor DI, testify testing, and standard naming conventions.

For detailed Go conventions, refer to the Go rule. For architecture patterns, refer to the Architecture rule. For testing standards, refer to the Testing rule. For Makefile setup, refer to the CI/CD rule.

## Directory Structure

Create the following layout:

```
<project>/
├── cmd/
│   └── <app>/
│       ├── main.go
│       └── container.go
├── internal/
│   ├── container.go
│   ├── domain/
│   │   ├── commands/
│   │   │   └── container.go
│   │   ├── entities/
│   │   │   └── container.go
│   │   └── repositories/
│   └── infrastructure/
│       ├── controllers/
│       │   ├── container.go
│       │   ├── mappers/
│       │   ├── requests/
│       │   └── responses/
│       └── repositories/
│           ├── container.go
│           ├── mappers/
│           └── models/
├── test/
│   ├── domain/
│   │   ├── builders/
│   │   ├── doubles/
│   │   │   └── repositories/
│   │   └── helpers/
│   └── infrastructure/
│       └── doubles/
│           └── repositories/
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Step-by-Step

### 1. Initialize the module

```bash
mkdir <project> && cd <project>
go mod init <module-path>
```

### 2. Create the domain layer -- entities (pure, no framework tags)

```go
// internal/domain/entities/user.go
package entities

type User struct {
    ID    string
    Name  string
    Email string
}
```

### 3. Create repository contracts

```go
// internal/domain/repositories/users_repository.go
package repositories

import "module/internal/domain/entities"

type UsersRepository interface {
    FindAll() ([]entities.User, error)
}
```

### 4. Create commands (business logic)

```go
// internal/domain/commands/list_users_command.go
package commands

import (
    "module/internal/domain/entities"
    "module/internal/domain/repositories"
)

type ListUsersCommand struct {
    repo repositories.UsersRepository
}

func NewListUsersCommand(repo repositories.UsersRepository) *ListUsersCommand {
    return &ListUsersCommand{repo: repo}
}

func (c *ListUsersCommand) Execute() ([]entities.User, error) {
    return c.repo.FindAll()
}
```

### 5. Create infrastructure implementations (prefixed with library name)

```go
// internal/infrastructure/repositories/in_memory_users_repository.go
package repositories

import "module/internal/domain/entities"

type InMemoryUsersRepository struct {
    users []entities.User
}

func NewInMemoryUsersRepository() *InMemoryUsersRepository {
    return &InMemoryUsersRepository{}
}

func (r *InMemoryUsersRepository) FindAll() ([]entities.User, error) {
    return r.users, nil
}
```

### 6. Create controllers

```go
// internal/infrastructure/controllers/list_users_controller.go
package controllers

import (
    "net/http"

    "module/internal/domain/commands"
)

type ListUsersController struct {
    command *commands.ListUsersCommand
}

func NewListUsersController(command *commands.ListUsersCommand) *ListUsersController {
    return &ListUsersController{command: command}
}

func (c *ListUsersController) Execute(w http.ResponseWriter, r *http.Request) {
    // call command, map response, write HTTP response
}
```

### 7. Assemble dependencies manually

Use explicit constructor calls for compile-time checked dependency injection. Keep the composition root in `cmd/<app>/container.go`; use `container.go` for module assembly where useful. These are ordinary typed functions, not a runtime container or generated provider registry.

Wire is no longer maintained. New projects must not introduce Wire or Dig. Migrate existing Wire services in a dedicated PR per service; encourage the same migration for Dig services. Preserve their current build during unrelated changes. Manual wiring removes runtime dependency resolution and catches missing arguments and incompatible types at compilation; do not claim a performance improvement without measuring the application.

### Composition Root

```go
// cmd/app/container.go
package main

import (
    "module/internal/domain/commands"
    "module/internal/infrastructure/controllers"
    "module/internal/infrastructure/repositories"
)

func initializeController() *controllers.ListUsersController {
    repository := repositories.NewInMemoryUsersRepository()
    command := commands.NewListUsersCommand(repository)
    return controllers.NewListUsersController(command)
}
```

The compiler checks that the repository satisfies the command's interface. Constructors return concrete types or the project's established interfaces. Aggregate dependencies with typed structs and explicit field assignments. Avoid reflection, service locators, global registries, no-op registration functions, and generated wiring.

Keep infrastructure imports in the outer composition root. Domain constructors accept domain interfaces and must not import infrastructure. Create shared clients once, propagate construction errors, and release acquired resources in reverse order on startup failure and shutdown.

Logging remains a separate choice: general projects use Logrus; projects with an established shared logging abstraction retain it. Selecting manual DI does not require changing the logger.

### Entry point

```go
// cmd/app/main.go
package main

import (
    "net/http"
    "time"

    logger "github.com/sirupsen/logrus"
)

func main() {
    controller := initializeController()
    router := http.NewServeMux()
    router.HandleFunc("/users", controller.Execute)
    server := &http.Server{Addr: ":8080", Handler: router, ReadHeaderTimeout: 5 * time.Second}
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        logger.WithField("error", err).Error("server stopped")
    }
}
```

### 8. Create test structure

- Place test files next to production files with `_test.go` suffix
- Unit tests do not need build tags — Go runs `_test.go` files by default. Add build flags (`//go:build integration`, `//go:build e2e`, etc.) only for tests that require external infrastructure
- Use `stretchr/testify` suites and `assert`/`require` for assertions
- Create builders in `test/domain/builders/` for test data
- Create doubles (stubs/dummies) in `test/domain/doubles/repositories/`
- Every test must use `// given`, `// when`, `// then` comment blocks
- Unit tests must call `t.Parallel()` at the top

Mocking libraries are prohibited unless an externally owned abstraction cannot be wrapped or doubled manually; document the constraint and rejected alternatives.

Provide working `test` and `test-unit` targets running `go test ./...`, and `test-integration` running `go test -tags=integration ./...`. Compile the complete scaffold and run both test targets before handing it off.

Refer to the Go rule (Testing section) and the Testing rule for full conventions.

### Example unit test

```go
// internal/domain/commands/list_users_command_test.go
package commands_test

import (
    "testing"
    "module/internal/domain/commands"
    "module/internal/domain/entities"
)

type usersStub struct { users []entities.User }
func (s usersStub) FindAll() ([]entities.User, error) { return s.users, nil }

func TestListUsersCommand_ReturnsUsers_WhenRepositorySucceeds(t *testing.T) {
    t.Parallel()
    // given
    command := commands.NewListUsersCommand(usersStub{users: []entities.User{{ID: "fixture-user"}}})
    // when
    users, err := command.Execute()
    // then
    if err != nil || len(users) != 1 || users[0].ID != "fixture-user" {
        t.Fatalf("unexpected users: %v, error: %v", users, err)
    }
}
```

### 9. Create Makefile

Provide these local targets, keeping integration tests gated:

```makefile
test:
	go test ./...

test-unit: test

test-integration:
	go test -tags=integration ./...

.PHONY: test test-unit test-integration
```

The project Makefile must import from the shared [pipelines repository](https://github.com/rios0rios0/pipelines) and expose `lint`, `test`, and `sast` targets. Refer to the CI/CD rule for details.

## Naming Quick Reference

| Component             | File                           | Struct                    | Method                            |
|-----------------------|--------------------------------|---------------------------|-----------------------------------|
| Command               | `<op>_<entity>_command.go`     | `<Op><Entity>Command`     | `Execute`                         |
| Controller            | `<op>_<entity>_controller.go`  | `<Op><Entity>Controller`  | `Execute`                         |
| Repository (contract) | `<entity>_repository.go`       | `<Entity>Repository`      | varies                            |
| Repository (impl)     | `<lib>_<entity>_repository.go` | `<Lib><Entity>Repository` | varies                            |
| Mapper                | `<entity>_mapper.go`           | `<Entity>Mapper`          | `ToEntity`/`ToModel`/`ToResponse` |
| Container             | `container.go`                 | --                        | typed assembly functions               |

## Key Rules

- Use a **short abbreviation** of the type as the method receiver (e.g., `c` for Command, `r` for Repository)
- Only attach methods to a struct when the method mutates state
- Entities MUST be free of framework tags -- pure business logic only
- No Services layer in Go projects
- Use **Logrus** (`github.com/sirupsen/logrus`) with alias `logger` in general projects; preserve a project-specific shared logger
- All file names use **snake_case**
