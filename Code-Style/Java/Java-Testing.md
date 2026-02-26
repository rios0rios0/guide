# Java Testing

> **TL;DR:** Use **[JUnit 5](https://junit.org/junit5/)** as the testing framework with **[Mockito](https://site.mockito.org/)** for mocking and **[Java Faker](https://github.com/DiUS/java-faker)** for test data generation. Tag tests with `@Tag("unit")` or `@Tag("integration")`. All tests must follow the BDD pattern with `// given`, `// when`, `// then` comment blocks. Use `@DisplayName` for human-readable test descriptions. Unit tests run in **parallel**. Integration tests use **Spring Boot test slices** with setup/teardown and are NOT parallel.

## Overview

JUnit 5 discovers test files automatically via classpath scanning. This document defines the conventions for organizing and writing tests across all Java projects.

## File Structure

```
src/test/java/.../
  <module>/
    domain/
      builders/                         test data builders
      commands/
        InsertItemCommandTest.java      command unit tests
    infrastructure/
      builders/                         JPA entity builders
      controllers/
        InsertItemControllerTest.java   controller unit tests
      repositories/
        JpaItemsRepositoryTest.java     repository integration tests
      services/
        JpaSendItemServiceTest.java     service unit tests
      listeners/
        ProcessItemListenerTest.java    listener tests
      doubles/
        SendItemServiceStub.java        stubs
        DummyItemService.java           dummies
        InMemoryItemsRepository.java    in-memory implementations

src/test/resources/
  db/                                   test database migrations
  http/                                 HTTP test data (REST requests)
  queue/                                message test data
```

## General Conventions

1. **Test tagging is mandatory.** Every test class must have `@Tag("unit")` or `@Tag("integration")`.
2. **BDD structure.** Every test must use `// given`, `// when`, `// then` comment blocks to separate preconditions, actions, and assertions.
3. **Display names.** Use `@DisplayName` on every test method to provide a human-readable description of the scenario.
4. **Testing framework.** Use [JUnit 5](https://junit.org/junit5/) with [Mockito](https://site.mockito.org/) for mocking and [AssertJ](https://assertj.github.io/doc/) or JUnit assertions for verification.
5. **File naming.** Test classes use the `Test` suffix (e.g., `InsertItemCommandTest.java`).
6. **File placement.** Test classes mirror the source structure under `src/test/java/`.
7. **Parallel unit tests.** Unit tests are configured for parallel class-level execution via JUnit properties.
8. **Sequential integration tests.** Integration tests share database state and must NOT run in parallel.

## Unit Tests (Parallel with @Tag)

Unit tests must be lightweight, fast, and isolated. They are configured for **parallel execution** at the class level.

### Command Tests

```java
@Tag("unit")
@NoArgsConstructor(access = AccessLevel.PRIVATE)
class InsertItemCommandTest {

    @Test
    @DisplayName("should call onSuccess when the item is inserted")
    void shouldCallOnSuccess() {
        // given
        final var event = new ItemEventBuilder().build();
        final var service = new SendItemServiceStub().withOnSuccess();
        final var command = new InsertItemCommand(service);
        final var onSuccess = mock(Runnable.class);

        // when
        final var listeners = new InsertItemCommand.Listeners(
                onSuccess,
                (e) -> fail("onError should not be called"));
        command.execute(event, listeners);

        // then
        verify(onSuccess, times(1)).run();
    }

    @Test
    @DisplayName("should call onError when the service throws an exception")
    void shouldCallOnError() {
        // given
        final var event = new ItemEventBuilder().build();
        final var service = new SendItemServiceStub().withOnError(new RuntimeException("test"));
        final var command = new InsertItemCommand(service);
        final var errorRef = new AtomicReference<Exception>();

        // when
        final var listeners = new InsertItemCommand.Listeners(
                () -> fail("onSuccess should not be called"),
                errorRef::set);
        command.execute(event, listeners);

        // then
        assertNotNull(errorRef.get());
        assertEquals("test", errorRef.get().getMessage());
    }
}
```

**Key points:**
- `@Tag("unit")` marks the class for unit test execution.
- `@NoArgsConstructor(access = AccessLevel.PRIVATE)` prevents instantiation outside JUnit.
- Each test is self-contained -- it creates its own doubles, command, and listeners.
- The Listeners pattern reflects all possible outcomes: `onSuccess`, `onError`.

### Controller Tests

```java
@Tag("unit")
@NoArgsConstructor(access = AccessLevel.PRIVATE)
class ListItemsControllerTest {

    @Test
    @DisplayName("should respond 200 (OK) when items are listed successfully")
    void shouldRespondOk() {
        // given
        final var command = new ListItemsCommandStub().withOnSuccess();
        final var controller = new ListItemsController(command);

        // when
        final var response = controller.execute(1L, PageRequest.of(0, 10));

        // then
        assertEquals(HttpStatus.OK.value(), response.getStatusCode().value());
    }

    @Test
    @DisplayName("should respond 500 (Internal Server Error) when command fails")
    void shouldRespondInternalServerError() {
        // given
        final var command = new ListItemsCommandStub().withOnError();
        final var controller = new ListItemsController(command);

        // when
        final var response = controller.execute(1L, PageRequest.of(0, 10));

        // then
        assertEquals(HttpStatus.INTERNAL_SERVER_ERROR.value(), response.getStatusCode().value());
    }
}
```

## Service Tests

```java
@Tag("unit")
@NoArgsConstructor(access = AccessLevel.PRIVATE)
class JpaSendItemServiceTest {

    @Test
    @DisplayName("should save the mapped item when event is valid")
    void shouldSaveItem() {
        // given
        final var event = new ItemEventBuilder().build();
        final var repository = new InMemoryItemsRepository();
        final var mapper = ItemMapper.INSTANCE;
        final var service = new JpaSendItemService(repository, mapper);

        // when
        service.send(event);

        // then
        assertEquals(1, repository.count());
    }

    @Test
    @DisplayName("should throw when repository fails to save")
    void shouldThrowWhenRepositoryFails() {
        // given
        final var event = new ItemEventBuilder().build();
        final var repository = new InMemoryItemsRepository().withOnError(new RuntimeException("db error"));
        final var mapper = ItemMapper.INSTANCE;
        final var service = new JpaSendItemService(repository, mapper);

        // when & then
        assertThrows(RuntimeException.class, () -> service.send(event));
    }
}
```

## Integration Tests (Spring Boot + TestContainers)

Integration tests verify the full stack (database, HTTP, messaging) using real infrastructure. They are **NOT parallel** due to shared mutable state.

### Repository Tests

```java
@Tag("integration")
@SpringBootTest
@ActiveProfiles("test")
class JpaItemsRepositoryTest {

    @Autowired
    private JpaItemsRepository repository;

    @Autowired
    private JdbcTemplate jdbcTemplate;

    @BeforeEach
    void setUp() {
        // Seed test data
        jdbcTemplate.execute("INSERT INTO item (id, name, created_at) VALUES (1, 'test-item', NOW())");
    }

    @AfterEach
    void tearDown() {
        jdbcTemplate.execute("DELETE FROM item");
    }

    @Test
    @DisplayName("should find item by ID successfully")
    void shouldFindById() {
        // given
        final var itemId = 1L;

        // when
        final var result = repository.findById(itemId);

        // then
        assertTrue(result.isPresent());
        assertEquals("test-item", result.get().getName());
    }

    @Test
    @DisplayName("should return empty when item does not exist")
    void shouldReturnEmpty() {
        // given
        final var nonExistentId = 99999L;

        // when
        final var result = repository.findById(nonExistentId);

        // then
        assertTrue(result.isEmpty());
    }

    @Test
    @DisplayName("should save a new item successfully")
    void shouldSaveItem() {
        // given
        final var item = JpaItem.builder()
                .name("new-item")
                .createdAt(LocalDateTime.now())
                .build();

        // when
        final var saved = repository.save(item);

        // then
        assertNotNull(saved.getId());
        assertEquals("new-item", saved.getName());
    }

    @Test
    @DisplayName("should delete an item successfully")
    void shouldDeleteItem() {
        // given
        final var itemId = 1L;

        // when
        repository.deleteById(itemId);

        // then
        assertTrue(repository.findById(itemId).isEmpty());
    }
}
```

**Key points:**
- `@Tag("integration")` marks the class for integration test execution.
- `@SpringBootTest` loads the full application context.
- `@ActiveProfiles("test")` activates the test profile (`application-test.yaml`).
- `@BeforeEach` / `@AfterEach` manage test data lifecycle.
- Tests are grouped by outcome: success, error, edge cases.

## Test Doubles

Follow the [Martin Fowler taxonomy](https://martinfowler.com/articles/mocksArentStubs.html) for test doubles:

| Type          | Purpose                                    | Example                                     |
|---------------|--------------------------------------------|---------------------------------------------|
| **Stub**      | Returns canned answers, no logic           | `SendItemServiceStub`                       |
| **Dummy**     | Minimal implementation, ready-made answers | `DummyItemService`                          |
| **In-Memory** | In-memory logic without external modules   | `InMemoryItemsRepository`                   |
| **Faker**     | External library generating realistic data | `ItemEventBuilder` (using Java Faker)       |
| **Mock**      | Mimics and verifies method calls           | Mockito `mock()` -- **avoid when possible** |

### Stub Example

```java
public class SendItemServiceStub implements SendItemService {
    private Exception error;

    public SendItemServiceStub withOnSuccess() {
        this.error = null;
        return this;
    }

    public SendItemServiceStub withOnError(final Exception error) {
        this.error = error;
        return this;
    }

    @Override
    public void send(final ItemEvent event) {
        if (error != null) {
            throw new RuntimeException(error);
        }
    }
}
```

### In-Memory Repository Example

```java
public class InMemoryItemsRepository implements ItemsRepository {
    private final List<JpaItem> items = new ArrayList<>();
    private Exception error;

    public InMemoryItemsRepository withOnError(final Exception error) {
        this.error = error;
        return this;
    }

    public long count() {
        return items.size();
    }

    @Override
    public JpaItem save(final JpaItem item) {
        if (error != null) {
            throw new RuntimeException(error);
        }
        items.add(item);
        return item;
    }
}
```

## Builders

Use the Builder Design Pattern to construct complex test objects step by step. Builders keep test setup readable and reusable across test suites.

```java
@NoArgsConstructor
public final class ItemEventBuilder {
    private Long organizationId;
    private String referenceGuid;
    private ItemCode code;
    private ItemCategory category;

    public ItemEventBuilder withOrganizationId(final Long organizationId) {
        this.organizationId = organizationId;
        return this;
    }

    public ItemEventBuilder withCode(final ItemCode code) {
        this.code = code;
        return this;
    }

    public ItemEvent build() {
        final var faker = new Faker();
        return new ItemEvent(
                organizationId != null ? organizationId : faker.number().randomNumber(),
                referenceGuid != null ? referenceGuid : faker.internet().uuid(),
                code != null ? code : ItemCode.values()[faker.number().numberBetween(0, ItemCode.values().length)],
                category != null ? category : ItemCategory.DEFAULT);
    }
}
```

Usage in tests:

```java
// Default values (randomized via Faker)
final var event = new ItemEventBuilder().build();

// Custom values
final var event = new ItemEventBuilder()
        .withOrganizationId(42L)
        .withCode(ItemCode.CRITICAL)
        .build();
```

## Parallel Execution

Configure parallel execution in `build.gradle`:

```groovy
test {
    useJUnitPlatform()
    systemProperty 'junit.jupiter.execution.parallel.enabled', 'true'
    systemProperty 'junit.jupiter.execution.parallel.mode.classes.default', 'concurrent'
}
```

This runs test **classes** concurrently while methods within the same class run in the same thread, preventing shared-state conflicts within a test class.

## Seeds

Seed files populate test databases with known data:

- The seed file name must match the **table name** (e.g., `item.sql`).
- For multiple seeds targeting the same table, use numbered suffixes: `item_01.sql`, `item_02.sql`.
- Constants must be named according to the file name.

## References

- [JUnit 5 User Guide](https://junit.org/junit5/docs/current/user-guide/)
- [Mockito Documentation](https://site.mockito.org/)
- [Java Faker](https://github.com/DiUS/java-faker)
- [Awaitility](https://github.com/awaitility/awaitility)
- [TestContainers](https://testcontainers.com/)
- [Given-When-Then -- Martin Fowler](https://martinfowler.com/bliki/GivenWhenThen.html)
- [Mocks Aren't Stubs -- Martin Fowler](https://martinfowler.com/articles/mocksArentStubs.html)
