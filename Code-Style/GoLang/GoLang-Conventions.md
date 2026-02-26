# Go Conventions

> **TL;DR:** Use `snake_case` for file names, `self` as the method receiver name, and follow the strict naming patterns for Commands, Controllers, Repositories, and Mappers. Entities must be framework-agnostic. Use [Dig](https://github.com/uber-go/dig) for dependency injection.

## Overview

This document defines Go-specific naming conventions and component patterns. For the general baseline, refer to the [Code Style](../../Code-Style.md) guide. The architectural layers referenced here are defined in the [Backend Design](../../Life-Cycle/Architecture/Backend-Design.md) section.

## File Naming

All file names must use `snake_case`:

```
list_users_command.go     # Correct
listUsersCommand.go       # Wrong
ListUsersCommand.go       # Wrong
```

## General Conventions

1. Use `self` as the method receiver name (analogous to `this` in other languages).
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
| Method name | `Execute`                         | `func (self ListUsersCommand) Execute(listeners ListUsersCommandListeners)` |

**Notes:**
- Use plural entity names when the operation targets multiple entities.
- Use the standard [operations vocabulary](../../Code-Style.md#operations-vocabulary).
- Listeners must reflect all possible controller responses.

## Controllers

| Element     | Pattern                              | Example                                     |
|-------------|--------------------------------------|---------------------------------------------|
| File name   | `<operation>_<entity>_controller.go` | `list_users_controller.go`                  |
| Struct name | `<Operation><Entity>Controller`      | `ListUsersController`                       |
| Method name | `Execute`                            | `func (self ListUsersController) Execute()` |

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
func (self UsersRepository) FindByTargetField(targetField any) entities.User

// Find multiple entities by a specific field
func (self UsersRepository) FindAllByTargetField(targetField any) []entities.User

// Check existence (returns boolean)
func (self UsersRepository) HasBooleanVerification(targetField any) bool

// Persist a single entity
func (self UsersRepository) Save(user entities.User)

// Persist multiple entities
func (self UsersRepository) SaveAll(users []entities.User)

// Remove a single entity by a specific field
func (self UsersRepository) DeleteByTargetField(targetField any)
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
func (self UserMapper) MapToEntity(infra any) entities.User
func (self UserMapper) MapToEntities(infra []any) []entities.User

// Domain Entity -> Infrastructure DTO
func (self UserMapper) MapToExternal(user entities.User) models.External
func (self UserMapper) MapToExternals(users []entities.User) []models.External
```

### Controller Mappers

| Element              | Pattern                                   | Example                          |
|----------------------|-------------------------------------------|----------------------------------|
| File name (request)  | `<operation>_<entity>_request_mapper.go`  | `insert_user_request_mapper.go`  |
| File name (response) | `<operation>_<entity>_response_mapper.go` | `insert_user_response_mapper.go` |

```go
// Request -> Entity (no inverse mapping)
func (self InsertUserRequestMapper) MapToEntity(request InsertUserRequest) entities.User
func (self InsertUserRequestMapper) MapToEntities(requests []InsertUserRequest) []entities.User

// Entity -> Response (no inverse mapping)
func (self InsertUserResponseMapper) MapToResponse(user entities.User) responses.InsertUserResponse
func (self InsertUserResponseMapper) MapToResponses(users []entities.User) []responses.InsertUserResponse
```

**Important:** Do not use `json` tags outside the infrastructure layer. Tags are restricted to request and response DTOs.

## Models

Models reside exclusively in the infrastructure layer and represent DTOs for external data sources (databases, APIs, queues, etc.). They resemble entities but are **not** entities.

Each model is prefixed with the name of the external tool it communicates with:

| Example       | Source             |
|---------------|--------------------|
| `AwsFile`     | AWS S3             |
| `ApiDocument` | External API       |
| `PgxUser`     | PostgreSQL via pgx |

## Dependency Injection

Use [Uber Dig](https://github.com/uber-go/dig) for runtime dependency injection via constructor injection. Dig resolves dependencies automatically by matching types from registered provider functions, requiring no code generation or manual wiring.

### Container File Convention

Each architectural layer must have a `container.go` file that registers its own providers. A top-level orchestrator calls each layer's registration function in dependency order.

| File                                          | Purpose                                      |
|-----------------------------------------------|----------------------------------------------|
| `cmd/<app>/dig.go`                            | Creates the container and invokes root types |
| `internal/container.go`                       | Orchestrates registration across all layers  |
| `internal/domain/entities/container.go`       | Registers entity providers (or no-op)        |
| `internal/domain/commands/container.go`       | Registers command providers (or no-op)       |
| `internal/infrastructure/controllers/container.go` | Registers controller providers          |
| `internal/infrastructure/repositories/container.go` | Registers repository providers         |

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

## References

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Dig - Dependency Injection](https://github.com/uber-go/dig)
- [DTO Pattern](https://www.baeldung.com/java-dto-pattern)
