# Testing Standards

> **TL;DR:** All tests must follow the BDD (Behavior-Driven Development) pattern with `// given`, `// when`, `// then` comment blocks. Write descriptive test names per layer: Commands (`"should call <LISTENER> when ..."`), Controllers (`"should respond <HTTP_STATUS_CODE> when ..."`), Services/Repos (`"should ... when ..."` with success + failure pairs). Prefer Stubs, Dummies, and In-memory doubles over Mocks. Use Builders for test data construction.

## Overview

This document defines the testing standards applicable across all languages and projects. The architectural layers referenced here are defined in the [Backend Design](../Life-Cycle/Architecture/Backend-Design.md) section.

## BDD Structure (Given / When / Then)

**All tests in all languages must follow the BDD pattern** with three clearly separated blocks using comments. This structure makes every test readable, self-documenting, and consistent across the entire codebase.

| Block     | Purpose                                                                               | Also Known As |
|-----------|---------------------------------------------------------------------------------------|---------------|
| **given** | Set up preconditions -- initialize objects, configure doubles, prepare input data     | Arrange       |
| **when**  | Execute the action under test -- call the method, trigger the event, send the request | Act           |
| **then**  | Assert expected outcomes -- verify return values, check side effects, validate state  | Assert        |

### Comment Syntax by Language

| Language                | Given      | When      | Then      |
|-------------------------|------------|-----------|-----------|
| Go                      | `// given` | `// when` | `// then` |
| JavaScript / TypeScript | `// given` | `// when` | `// then` |
| Java                    | `// given` | `// when` | `// then` |
| Python                  | `# given`  | `# when`  | `# then`  |

### Example (JavaScript)

```ts
it('should render Content when not loading', () => {
  // given
  const props = { loading: false };

  // when
  const component = shallow(<ExampleComponent {...props} />);

  // then
  expect(component.find(Content).exists()).toBeTruthy();
  expect(component.find(Loading).exists()).toBeFalsy();
});
```

### Example (Go)

```go
func (s *CommandSuite) TestCreateUserCommand() {
    // given
    repository := doubles.NewInMemoryUsersRepository()
    command := commands.NewCreateUserCommand(repository)
    user := &cmd.User{
        Name:  "John Doe",
        Email: "johndoe@example.com",
    }

    // when
    err := command.Run(user)

    // then
    assert.Equal(s.T(), user, repository.FindByID(user.ID))
    assert.Nil(s.T(), err)
}
```

### Example (Python)

```python
def test_create_user_successfully(self):
    # given
    user_data = {"name": "John Doe", "email": "john@example.com"}

    # when
    result = self.service.create_user(user_data)

    # then
    assert result.name == "John Doe"
    assert result.email == "john@example.com"
```

## Test Description Patterns

### Commands

Create one test per listener in the command. Use this description format:

```
"should call <LISTENER> when ..."
```

### Controllers

Create one test per HTTP status code. Verify both the status code and the response body:

```
"should respond <HTTP_STATUS_CODE> (HTTP_STATUS) when ..."
```

### Services

Create **at least** one success and one failure test per public method:

```
"should ... when ..."
```

### Repositories

Follow the same conventions as [Services](#services).

## Test Doubles

Based on [Martin Fowler's taxonomy](https://martinfowler.com/bliki/TestDouble.html), a **test double** is any object that replaces a real dependency during testing. The team uses the following types:

| Type          | Purpose                                               | Guidelines                                                           |
|---------------|-------------------------------------------------------|----------------------------------------------------------------------|
| **Stub**      | Returns pre-configured (canned) answers               | No in-memory logic; return static values only                        |
| **Dummy**     | Fills required parameters that are never used         | Return minimal values (empty lists, `null`)                          |
| **In-memory** | Implements logic in memory without external modules   | Use for lightweight simulations of repositories or services          |
| **Faker**     | Generates realistic fake data via an external library | Use libraries like Faker.js or Go Faker                              |
| **Mock**      | Records and verifies method calls                     | **Avoid when possible.** Use only when no other double type suffices |

## Mocking Libraries

Write lightweight stubs, dummies, in-memory implementations, and recording doubles by hand. Use composition or interface embedding in Go (Go has no class inheritance). Test observable behavior. Assertions and suites from `testify/assert`, `testify/require`, and `testify/suite` are allowed; `testify/mock` is a mocking library.

Mocking libraries are prohibited by default. The only exception is an externally owned abstraction that cannot be tested with a handwritten double, consumer-owned interface, adapter, or interface extension/composition. Document the external constraint, the alternatives attempted, and why each fails beside the test or in its review. Convenience, a large interface, or call verification alone is insufficient. Keep any exception narrowly scoped and reviewable.

Examples requiring review include `testify/mock`, `golang/mock`, `go.uber.org/mock`, `mockery`, `go-sqlmock`, `httpmock`, and `gock`. The restriction covers the category, not only these packages. A checker reports a candidate violation; an exception marker does not prove the exception is justified.

Unit-test consumers against their own ports. Integration-test database adapters against a real database and migrations. Use `httptest.NewServer` for HTTP behavior; it is a real server under test control. A handwritten repository double is valid for command tests and does not replace integration coverage of SQL or transactions.

## Go Test Build Tags

Unit tests and helpers they import must be untagged. Both `//go:build unit` and legacy `// +build unit`, including combined or negated forms, are prohibited. Ordinary `go test ./...` must discover unit tests.

Keep `//go:build integration` for tests and helpers requiring external infrastructure. Dedicated end-to-end tags may gate genuine end-to-end tests. Do not require a unit tag in Makefiles, generators, or CI. Use `go test ./...` for units and `go test -tags=integration ./...` for the suite including integration tests.

## Builders

The [Builder Design Pattern](https://refactoring.guru/design-patterns/builder) enables step-by-step construction of complex test objects. Builders separate object construction from representation, making test setup readable and reusable:

```
UserBuilder.new()
    .withName("John Doe")
    .withEmail("john@example.com")
    .build()
```

Use builders to construct entities, DTOs, and other complex objects needed for automated testing.

## References

- [Test Double -- Martin Fowler](https://martinfowler.com/bliki/TestDouble.html)
- [Mocks Aren't Stubs -- Martin Fowler](https://martinfowler.com/articles/mocksArentStubs.html)
- [Builder Pattern -- Refactoring Guru](https://refactoring.guru/design-patterns/builder)
- [BDD -- Behavior-Driven Development](https://cucumber.io/docs/bdd/)
- [Given-When-Then](https://martinfowler.com/bliki/GivenWhenThen.html)
