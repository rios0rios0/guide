# Python Conventions

> **TL;DR:** Use `snake_case` for file and function names, `PascalCase` for classes. Follow the `<Operation><Entity>` naming pattern for Commands and Controllers with an `execute` method. Use descriptive names and write comments that explain **why**, not just **what**.

## Overview

This document defines Python-specific naming conventions and component patterns. For the general baseline, refer to the [Code Style](../../Code-Style.md) guide. The architectural layers referenced here are defined in the [Backend Design](../../Life-Cycle/Architecture/Backend-Design.md) section.

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

## References

- [PEP 8 -- Naming Conventions](https://peps.python.org/pep-0008/#naming-conventions)
- [PEP 257 -- Docstring Conventions](https://peps.python.org/pep-0257/)
- [ABC Module Documentation](https://docs.python.org/3/library/abc.html)
