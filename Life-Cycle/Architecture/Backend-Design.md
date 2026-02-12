# Backend Design

> **TL;DR:** Separate code into `domain` (contracts) and `infrastructure` (implementations). The five principal actors are Entities, Controllers, Commands, Services, and Repositories. Use Mappers to isolate layers and Dependency Injection to wire them together.

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

![](.assets/requests_flow.png)

### Mapping/Parsing Flow

Demonstrates how Mappers isolate layers from each other, preventing hard coupling between frameworks or external tools:

![](.assets/mapping_flow.png)

### Dependency Injection Flow

Shows how DI wires contracts to their implementations at runtime:

![](.assets/dependency_injection_flow.png)

## References

- [CQRS -- Martin Fowler](https://martinfowler.com/bliki/CQRS.html)
- [Mocks Aren't Stubs -- Martin Fowler](https://martinfowler.com/articles/mocksArentStubs.html)
- [Clean Architecture -- Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
