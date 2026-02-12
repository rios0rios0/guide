# Architecture

> **TL;DR:** All applications follow **Clean Architecture** principles, separating business logic from infrastructure concerns. Dependencies always point inward toward the domain. Adhere to SOLID principles for maintainable, testable code.

## Overview

Modern application development is fundamentally a task of **managing dependencies**. Poor architectural design leads to rigid, fragile codebases that are costly to change. Since requirements evolve, libraries get deprecated, and external services change, the development team must actively maintain a clean dependency graph.

|        Polluted Architecture[^1]        |       Clean Architecture[^1]        |
|:---------------------------------------:|:-----------------------------------:|
| ![](.assets/not-clean-architecture.png) | ![](.assets/clean-architecture.png) |

[^1]: [Clean Architecture Introduction](https://pusher.com/tutorials/clean-architecture-introduction)

Clean Architecture enforces a clear separation of concerns where **business rules are independent of frameworks, databases, and delivery mechanisms**. The inner layers define policies; the outer layers implement mechanisms.

## Sub-Pages

- [Backend Design](Architecture/Backend-Design.md) -- Layers, actors, file structure, and flow diagrams
- [Frontend Design](Architecture/Frontend-Design.md) -- 5-layer frontend architecture

## References

- [SOLID Principles](https://en.wikipedia.org/wiki/SOLID)
- [Uncle Bob's Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Clean Architecture Book](https://www.amazon.com/Clean-Architecture-Craftsmans-Software-Structure/dp/0134494164)
- [Clean Architecture Introduction](https://pusher.com/tutorials/clean-architecture-introduction)
