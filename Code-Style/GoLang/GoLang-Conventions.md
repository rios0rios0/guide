# Go Conventions

> **TL;DR:** Use `snake_case` for file names, a short abbreviation of the type as the method receiver name (e.g., `c` for `Command`), and follow the strict naming patterns for Commands, Controllers, Repositories, and Mappers. Entities must be framework-agnostic. Use manual constructor injection with compile-time checks.

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

Use explicit constructor calls for compile-time checked dependency injection. Keep the composition root in `cmd/<app>/container.go`; use `container.go` for module assembly where useful. These are ordinary typed functions, not a runtime container or generated provider registry.

Wire is no longer maintained. New projects must not introduce Wire or Dig. Migrate existing Wire services in a dedicated PR per service; encourage the same migration for Dig services. Preserve their current build during unrelated changes. Manual wiring removes runtime dependency resolution and catches missing arguments and incompatible types at compilation; do not claim a performance improvement without measuring the application.

### Composition Root

```go
// cmd/app/container.go
package main

import (
    "example.com/app/internal/domain/commands"
    "example.com/app/internal/infrastructure/controllers"
    "example.com/app/internal/infrastructure/repositories"
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

## References

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Wire maintenance status](https://github.com/google/wire)
- [DTO Pattern](https://www.baeldung.com/java-dto-pattern)
