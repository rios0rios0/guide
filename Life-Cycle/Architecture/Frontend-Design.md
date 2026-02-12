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

![](.assets/frontend_architecture_1.png)

### Applied to Feature Modules

![](.assets/frontend_architecture_2.png)

### Dependency Direction

Dependencies always point toward the **Domain** layer, which is the most abstract and stable:

![](.assets/frontend_architecture_3.png)

## State Management

Currently, state is managed locally within Presentation components. A dedicated State Layer may be introduced in the future as requirements evolve.

## References

- [A Different Approach to Frontend Architecture](https://dev.to/huytaquoc/a-different-approach-to-frontend-architecture-38d4)
- [Scalable Frontend #1 -- Architecture Fundamentals](https://blog.codeminer42.com/scalable-frontend-1-architecture-9b80a16b8ec7/)
- [Clean Architecture for the Rest of Us](https://pusher.com/tutorials/clean-architecture-introduction/)
