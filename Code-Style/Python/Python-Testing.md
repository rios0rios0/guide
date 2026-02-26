# Python Testing

> **TL;DR:** Use **pytest** for all testing. All tests must follow the BDD pattern with `# given`, `# when`, `# then` comment blocks. Place tests in a `/tests/` directory mirroring the source structure. Use fixtures for setup/teardown and factories for test data construction.

## Overview

This document defines the testing conventions for all Python projects. For the cross-language testing standards (BDD structure, test doubles, builders), refer to the [Tests](../../Life-Cycle/Tests.md) page.

## Testing Framework

| Tool                                                 | Purpose                                 |
|------------------------------------------------------|-----------------------------------------|
| [pytest](https://docs.pytest.org/)                   | Testing framework and runner            |
| [pytest-cov](https://pytest-cov.readthedocs.io/)    | Coverage reporting                      |
| [Faker](https://faker.readthedocs.io/)               | Fake data generation library            |

### Installation

```bash
pdm add -dG test pytest pytest-cov faker
```

## BDD Structure (Given / When / Then)

**Every test must use `# given`, `# when`, `# then` comment blocks** to clearly separate preconditions, actions, and assertions. This is mandatory across all test files.

```python
def test_should_return_user_when_found():
    # given
    repository = InMemoryUsersRepository()
    repository.save(User(id=1, name="John", email="john@test.com"))
    command = GetUserCommand(repository)

    # when
    result = command.execute(user_id=1)

    # then
    assert result is not None
    assert result.name == "John"
```

## File Placement

Tests live in a `/tests/` directory at the project root, mirroring the source code structure:

```
project/
  src/
    domain/
      commands/
        list_users_command.py
      entities/
        user.py
      repositories/
        users_repository.py
    infrastructure/
      controllers/
        list_users_controller.py
      repositories/
        sqlalchemy_users_repository.py
  tests/
    domain/
      commands/
        test_list_users_command.py
      builders/
        user_builder.py
      doubles/
        in_memory_users_repository.py
    infrastructure/
      controllers/
        test_list_users_controller.py
      repositories/
        test_sqlalchemy_users_repository.py
    conftest.py
```

## File Naming

All test files must use the `test_` prefix:

```
test_list_users_command.py     # Correct
list_users_command_test.py     # Wrong
test_list_users.py             # Wrong (must match the source file name)
```

## Unit Tests

Unit tests validate individual components in isolation. Use test doubles (stubs, dummies, fakers, in-memory implementations) to replace external dependencies.

### Command Tests

```python
import pytest
from domain.commands import DeleteUserCommand
from tests.domain.doubles import InMemoryUsersRepository, UsersRepositoryStub

class TestDeleteUserCommand:
    def test_should_call_on_success_when_user_is_deleted(self):
        # given
        repository = InMemoryUsersRepository()
        repository.save(User(id=1, name="John", email="john@test.com"))
        command = DeleteUserCommand(repository)
        success_called = False

        def on_success():
            nonlocal success_called
            success_called = True

        # when
        command.execute(user_id=1, on_success=on_success, on_error=lambda e: None)

        # then
        assert success_called is True

    def test_should_call_on_error_when_user_not_found(self):
        # given
        repository = InMemoryUsersRepository()
        command = DeleteUserCommand(repository)
        error_received = None

        def on_error(err):
            nonlocal error_received
            error_received = err

        # when
        command.execute(user_id=999, on_success=lambda: None, on_error=on_error)

        # then
        assert error_received is not None
```

### Controller Tests

```python
from fastapi.testclient import TestClient

class TestListUsersController:
    def test_should_respond_200_when_users_listed_successfully(self, client: TestClient):
        # when
        response = client.get("/users")

        # then
        assert response.status_code == 200

    def test_should_respond_500_when_command_fails(self, client_with_error: TestClient):
        # when
        response = client_with_error.get("/users")

        # then
        assert response.status_code == 500
```

## Integration Tests

Integration tests validate components working together with real infrastructure (databases, external services). Use pytest fixtures for setup and teardown.

### Repository Tests

```python
import pytest
from faker import Faker
from domain.entities import Item
from infrastructure.repositories import SQLAlchemyItemsRepository
from tests.domain.builders import ItemBuilder

fake = Faker()

@pytest.fixture
def repository(db_session):
    return SQLAlchemyItemsRepository(db_session)

@pytest.fixture(autouse=True)
def seed_data(db_session):
    """Seed the database before each test and clean up after."""
    # setup
    db_session.execute("INSERT INTO items (name, external_id) VALUES ('test', 'ext-1')")
    db_session.commit()
    yield
    # teardown
    db_session.execute("DELETE FROM items")
    db_session.commit()

class TestSQLAlchemyItemsRepository:
    # --- Success Cases ---

    def test_should_save_a_new_item_successfully(self, repository):
        # given
        item = ItemBuilder().with_external_id(fake.uuid4()).build()

        # when
        created = repository.save(item)

        # then
        assert created.id is not None
        assert created.id > 0

    def test_should_list_all_items_successfully(self, repository):
        # given / when
        items = repository.find_all()

        # then
        assert len(items) >= 1

    def test_should_get_item_by_id_successfully(self, repository, seed_data):
        # given
        items = repository.find_all()
        item_id = items[0].id

        # when
        item = repository.find_by_id(item_id)

        # then
        assert item is not None
        assert item.id == item_id

    # --- Error Cases ---

    def test_should_raise_when_item_not_found(self, repository):
        # given
        non_existent_id = 99999

        # when / then
        with pytest.raises(ItemNotFoundError):
            repository.find_by_id(non_existent_id)
```

## Test Doubles

Use the following test double types, in order of preference:

| Type       | Description                                       | When to Use                                    |
|------------|---------------------------------------------------|------------------------------------------------|
| Stub       | Returns pre-configured (canned) answers           | When you need controlled return values         |
| Dummy      | Fills required parameters, never actually used    | When a parameter is required but irrelevant    |
| In-memory  | Implements logic without external dependencies    | When you need realistic behavior without I/O   |
| Faker      | Generates realistic fake data                     | When you need varied, realistic test inputs    |
| **Mock**   | Records and verifies calls                        | **Avoid when possible** -- prefer stubs        |

## Builders

Use the builder pattern for constructing complex test objects:

```python
class UserBuilder:
    def __init__(self):
        self._id = 1
        self._name = "Default User"
        self._email = "default@test.com"

    def with_id(self, id: int) -> "UserBuilder":
        self._id = id
        return self

    def with_name(self, name: str) -> "UserBuilder":
        self._name = name
        return self

    def with_email(self, email: str) -> "UserBuilder":
        self._email = email
        return self

    def build(self) -> User:
        return User(id=self._id, name=self._name, email=self._email)
```

Usage:

```python
user = UserBuilder().with_id(0).with_name("John").with_email("john@test.com").build()
```

## References

- [pytest Documentation](https://docs.pytest.org/)
- [pytest Fixtures](https://docs.pytest.org/en/stable/how-to/fixtures.html)
- [Faker Documentation](https://faker.readthedocs.io/)
- [Given-When-Then -- Martin Fowler](https://martinfowler.com/bliki/GivenWhenThen.html)
