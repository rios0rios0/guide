Scaffold a new Go backend project following Clean Architecture with domain/infrastructure separation, Dig DI, testify testing, and standard naming conventions.

For detailed Go conventions, refer to the Go rule. For architecture patterns, refer to the Architecture rule. For testing standards, refer to the Testing rule. For Makefile setup, refer to the CI/CD rule.

## Directory Structure

Create the following layout:

```
<project>/
├── cmd/
│   └── <app>/
│       ├── main.go
│       └── dig.go
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
    FindByID(id string) (*entities.User, error)
    Insert(entity *entities.User) error
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
// internal/infrastructure/repositories/pgx_users_repository.go
package repositories

import "module/internal/domain/entities"

type PgxUsersRepository struct {
    // db connection
}

func NewPgxUsersRepository() *PgxUsersRepository {
    return &PgxUsersRepository{}
}

func (r *PgxUsersRepository) FindAll() ([]entities.User, error) {
    // implementation
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

### 7. Set up Dig for dependency injection

Create `container.go` files per layer. Each layer registers its own providers:

```go
// internal/infrastructure/repositories/container.go
package repositories

import "go.uber.org/dig"

func RegisterProviders(container *dig.Container) error {
    if err := container.Provide(NewPgxUsersRepository); err != nil {
        return err
    }
    return nil
}
```

```go
// internal/domain/commands/container.go
package commands

import "go.uber.org/dig"

func RegisterProviders(container *dig.Container) error {
    if err := container.Provide(NewListUsersCommand); err != nil {
        return err
    }
    return nil
}
```

```go
// internal/container.go
package internal

import (
    "go.uber.org/dig"

    "module/internal/domain/commands"
    "module/internal/infrastructure/controllers"
    "module/internal/infrastructure/repositories"
)

func RegisterProviders(container *dig.Container) error {
    if err := repositories.RegisterProviders(container); err != nil {
        return err
    }
    if err := commands.RegisterProviders(container); err != nil {
        return err
    }
    if err := controllers.RegisterProviders(container); err != nil {
        return err
    }
    return nil
}
```

```go
// cmd/<app>/dig.go
package main

import (
    "go.uber.org/dig"
    "module/internal"
    "module/internal/infrastructure/controllers"
)

func injectController() *controllers.ListUsersController {
    container := dig.New()
    if err := internal.RegisterProviders(container); err != nil {
        panic(err)
    }

    var ctrl *controllers.ListUsersController
    if err := container.Invoke(func(c *controllers.ListUsersController) {
        ctrl = c
    }); err != nil {
        panic(err)
    }
    return ctrl
}
```

### 8. Create test structure

- Place test files next to production files with `_test.go` suffix
- Add build flags: `//go:build unit` or `//go:build integration`
- Use `stretchr/testify` suites and `assert`/`require` for assertions
- Create builders in `test/domain/builders/` for test data
- Create doubles (stubs/dummies) in `test/domain/doubles/repositories/`
- Every test must use `// given`, `// when`, `// then` comment blocks
- Unit tests must call `t.Parallel()` at the top

Refer to the Go rule (Testing section) and the Testing rule for full conventions.

### 9. Create Makefile

The project Makefile must import from the shared [pipelines repository](https://github.com/rios0rios0/pipelines) and expose `lint`, `test`, and `sast` targets. Refer to the CI/CD rule for details.

## Naming Quick Reference

| Component             | File                           | Struct                    | Method                            |
|-----------------------|--------------------------------|---------------------------|-----------------------------------|
| Command               | `<op>_<entity>_command.go`     | `<Op><Entity>Command`     | `Execute`                         |
| Controller            | `<op>_<entity>_controller.go`  | `<Op><Entity>Controller`  | `Execute`                         |
| Repository (contract) | `<entity>_repository.go`       | `<Entity>Repository`      | varies                            |
| Repository (impl)     | `<lib>_<entity>_repository.go` | `<Lib><Entity>Repository` | varies                            |
| Mapper                | `<entity>_mapper.go`           | `<Entity>Mapper`          | `ToEntity`/`ToModel`/`ToResponse` |
| Container             | `container.go`                 | --                        | `RegisterProviders`               |

## Key Rules

- Use a **short abbreviation** of the type as the method receiver (e.g., `c` for Command, `r` for Repository)
- Only attach methods to a struct when the method mutates state
- Entities MUST be free of framework tags -- pure business logic only
- No Services layer in Go projects
- Use **Logrus** (`github.com/sirupsen/logrus`) for ALL logging with alias `logger`
- All file names use **snake_case**
