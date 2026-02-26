# Go Conventions

> **TL;DR:** Use `snake_case` for file names, `self` as the method receiver name, and follow the strict naming patterns for Commands, Controllers, Repositories, and Mappers. Entities must be framework-agnostic. Use [Wire](https://github.com/google/wire) for dependency injection.

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

Use [Google Wire](https://github.com/google/wire) for compile-time dependency injection. Wire generates code at build time, avoiding the runtime overhead and complexity of reflection-based DI containers.

## References

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Google Wire - Dependency Injection](https://github.com/google/wire)
- [DTO Pattern](https://www.baeldung.com/java-dto-pattern)
