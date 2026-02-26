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

## References

- [Effective Go -- Interfaces](https://go.dev/doc/effective_go#interfaces)
- [Go Generics Tutorial](https://go.dev/doc/tutorial/generics)
- [Go Type Assertions](https://go.dev/tour/methods/15)
- [Go Proverbs -- The bigger the interface, the weaker the abstraction](https://go-proverbs.github.io/)
