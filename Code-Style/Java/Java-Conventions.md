# Java Conventions

> **TL;DR:** Use `PascalCase` for class names and follow the strict `<Operation><Entity>` naming patterns for Commands, Controllers, Services, Repositories, and Mappers. Entities must be framework-agnostic. Use Spring constructor injection for dependency injection and [MapStruct](https://mapstruct.org/) for object mapping.

## Overview

This document defines Java-specific naming conventions and component patterns. For the general baseline, refer to the [Code Style](../../Code-Style.md) guide. The architectural layers referenced here are defined in the [Backend Design](../../Life-Cycle/Architecture/Backend-Design.md) section.

## File Naming

Java enforces that the file name matches the public class name:

```
ListItemsCommand.java       # Correct -- matches class ListItemsCommand
list_items_command.java      # Wrong -- Java uses PascalCase file names
```

## General Conventions

1. Use **constructor injection** for all dependencies -- never use field injection (`@Autowired` on fields).
2. Mark injected fields as `final` to enforce immutability.
3. Use [Lombok](https://projectlombok.org/) to reduce boilerplate (`@Getter`, `@Setter`, `@NoArgsConstructor`, `@RequiredArgsConstructor`, `@SuperBuilder`).
4. Use Java **records** for immutable DTOs, domain events, and value objects.
5. Use `@Component` for commands and listeners, `@RestController` for controllers, `@Repository` for repository implementations, and `@Service` for infrastructure service implementations.
6. For an introduction to the DTO pattern, refer to [this article](https://www.baeldung.com/java-dto-pattern).

## Entities

Entities are the core of the application. All business logic related to properties and fields belongs inside the entity.

**Entities must be free of any persistence or framework annotations.** Do not use `@Entity`, `@Table`, `@Column`, or any Jakarta/JPA annotations inside domain entities. Persistence concerns belong in infrastructure models.

```java
public final class Item<A> {
    private final ItemCode code;
    private final ItemSeverity severity;
    private Boolean isCountable;
    private Long amountAffected;
    private List<A> affected;

    public Boolean hasData() {
        return amountAffected > 0 || !isCountable;
    }
}
```

Use Java records for simple domain events and value objects:

```java
public record ItemEvent(
        Long organizationId,
        String referenceGuid,
        ItemCode code,
        ItemCategory category) {}
```

## Commands

| Element     | Pattern                           | Example                    |
|-------------|-----------------------------------|----------------------------|
| File name   | `<Operation><Entity>Command.java` | `InsertItemCommand.java`   |
| Class name  | `<Operation><Entity>Command`      | `InsertItemCommand`        |
| Method name | `execute`                         | `public void execute(...)` |

**Notes:**
- Use plural entity names when the operation targets multiple entities.
- Use the standard [operations vocabulary](../../Code-Style.md#operations-vocabulary).
- Commands must define a `Listeners` record for all possible outcomes (callback pattern).

```java
@Component
@RequiredArgsConstructor
public class InsertItemCommand {
    private final SendItemService service;

    public void execute(final ItemEvent event, final Listeners listeners) {
        try {
            service.send(event);
            listeners.onSuccess().run();
        } catch (Exception e) {
            listeners.onError().accept(e);
        }
    }

    public record Listeners(Runnable onSuccess, Consumer<Exception> onError) {}
}
```

## Controllers

| Element     | Pattern                              | Example                                 |
|-------------|--------------------------------------|-----------------------------------------|
| File name   | `<Operation><Entity>Controller.java` | `ListItemsController.java`              |
| Class name  | `<Operation><Entity>Controller`      | `ListItemsController`                   |
| Method name | `execute`                            | `public ResponseEntity<?> execute(...)` |

```java
@RestController
@RequiredArgsConstructor
@RequestMapping("${api.v1-prefix}")
public class InsertItemController {
    private final InsertItemCommand command;

    @PostMapping("/items")
    public ResponseEntity<?> execute(
            @RequestHeader("organization-id") final Long organizationId,
            @Valid @RequestBody final InsertItemRequest request) {
        final var entity = InsertItemRequestMapper.INSTANCE.mapToEntity(organizationId, request);
        final var ref = new AtomicReference<ResponseEntity<?>>();

        final var listeners = new InsertItemCommand.Listeners(
                () -> ref.set(ResponseEntity.status(HttpStatus.CREATED).build()),
                (error) -> ref.set(ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).build()));
        command.execute(entity, listeners);
        return ref.get();
    }
}
```

## Services

The Services layer provides an abstraction between Commands and Repositories. The **domain layer** defines the service contract (interface), and the **infrastructure layer** provides the implementation.

### Contract (Domain Layer)

| Element        | Pattern                           | Example                |
|----------------|-----------------------------------|------------------------|
| File name      | `<Operation><Entity>Service.java` | `SendItemService.java` |
| Interface name | `<Operation><Entity>Service`      | `SendItemService`      |

```java
public interface SendItemService {
    void send(ItemEvent event);
}
```

### Implementation (Infrastructure Layer)

| Element    | Pattern                              | Example                   |
|------------|--------------------------------------|---------------------------|
| File name  | `Jpa<Operation><Entity>Service.java` | `JpaSendItemService.java` |
| Class name | `Jpa<Operation><Entity>Service`      | `JpaSendItemService`      |

```java
@Service
@RequiredArgsConstructor
public class JpaSendItemService implements SendItemService {
    private final ItemsRepository repository;
    private final ItemMapper mapper;

    @Override
    public void send(final ItemEvent event) {
        final var model = mapper.toModel(event);
        repository.save(model);
    }
}
```

**Note:** The `Jpa` prefix indicates the infrastructure tool used. If using a different technology (e.g., REST client, gRPC), use the corresponding prefix (e.g., `RestSendItemService`, `GrpcSendItemService`).

## Repositories

### Contract (Domain Layer)

| Element        | Pattern                   | Example                |
|----------------|---------------------------|------------------------|
| File name      | `<Entity>Repository.java` | `ItemsRepository.java` |
| Interface name | `<Entity>Repository`      | `ItemsRepository`      |

```java
public interface ItemsRepository {
    Boolean hasItems(Long categoryId);
    Page<JpaItem> findAllByCategoryId(Long categoryId, Pageable pageable);
}
```

### Implementation (Infrastructure Layer)

| Element    | Pattern                      | Example                   |
|------------|------------------------------|---------------------------|
| File name  | `Jpa<Entity>Repository.java` | `JpaItemsRepository.java` |
| Class name | `Jpa<Entity>Repository`      | `JpaItemsRepository`      |

For Spring Data JPA repositories:

```java
@Repository
public interface JpaItemsRepository
        extends JpaRepository<JpaItem, Long>, ItemsRepository {}
```

### QueryDSL Repositories

For complex queries that go beyond Spring Data JPA:

| Element    | Pattern                           | Example                        |
|------------|-----------------------------------|--------------------------------|
| File name  | `QueryDsl<Entity>Repository.java` | `QueryDslItemsRepository.java` |
| Class name | `QueryDsl<Entity>Repository`      | `QueryDslItemsRepository`      |

## Mappers

Use [MapStruct](https://mapstruct.org/) for all object mapping. Mappers isolate layers from each other and prevent hard coupling to frameworks.

### Repository Mappers

| Element   | Pattern               | Example           |
|-----------|-----------------------|-------------------|
| File name | `<Entity>Mapper.java` | `ItemMapper.java` |
| Interface | `<Entity>Mapper`      | `ItemMapper`      |

```java
@Mapper(unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface ItemMapper {
    ItemMapper INSTANCE = Mappers.getMapper(ItemMapper.class);

    Item toEntity(JpaItem model);
    JpaItem toModel(Item entity);
}
```

### Controller Mappers

| Element              | Pattern                                        | Example                           |
|----------------------|------------------------------------------------|-----------------------------------|
| File name (request)  | `<Operation><Entity>RequestMapper.java`        | `InsertItemRequestMapper.java`    |
| File name (response) | `<Operation><Entity>ResponseMapper.java`       | `InsertItemResponseMapper.java`   |

```java
@Mapper(unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface InsertItemRequestMapper {
    InsertItemRequestMapper INSTANCE = Mappers.getMapper(InsertItemRequestMapper.class);

    default ItemEvent mapToEntity(final Long organizationId, final InsertItemRequest request) {
        return new ItemEvent(
                organizationId,
                request.referenceGuid(),
                request.code(),
                request.category());
    }
}
```

## Models

Models reside exclusively in the infrastructure layer and represent JPA entities for database persistence. They resemble domain entities but are **not** domain entities.

Each model is prefixed with `Jpa`:

| Element    | Pattern            | Example        |
|------------|--------------------|----------------|
| File name  | `Jpa<Entity>.java` | `JpaItem.java` |
| Class name | `Jpa<Entity>`      | `JpaItem`      |

```java
@Entity
@Getter
@Setter
@SuperBuilder
@NoArgsConstructor
@Table(name = "item")
public class JpaItem {
    @Id
    @GeneratedValue(strategy = GenerationType.SEQUENCE)
    private Long id;

    @CreationTimestamp
    @Column(nullable = false)
    private LocalDateTime createdAt;

    @Column(nullable = false)
    private String name;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "category_id")
    private JpaCategory category;
}
```

## Listeners

Listeners handle asynchronous messages from queues (e.g., Kafka). They act as infrastructure-layer entry points, similar to controllers but for message-driven operations.

| Element     | Pattern                            | Example                      |
|-------------|------------------------------------|------------------------------|
| File name   | `Process<Entity>Listener.java`     | `ProcessItemListener.java`   |
| Class name  | `Process<Entity>Listener`          | `ProcessItemListener`        |

```java
@Component
@RequiredArgsConstructor
public class ProcessItemListener {
    private final InsertItemCommand command;

    @KafkaListener(
            topics = "items-topic",
            groupId = "${queue.items.groupId}")
    @RetryableTopic(attempts = "1")
    public void queueListener(
            @Payload final ItemMessage message,
            final Acknowledgment acknowledgment) {
        acknowledgment.acknowledge();
        final var event = message.toDomain();
        command.execute(event, listeners);
    }
}
```

### Messages

| Element     | Pattern                  | Example            |
|-------------|--------------------------|--------------------|
| File name   | `<Entity>Message.java`   | `ItemMessage.java` |
| Class name  | `<Entity>Message`        | `ItemMessage`      |

```java
public record ItemMessage(
        @JsonProperty("organization_id") Long organizationId,
        @JsonProperty("reference_guid") String referenceGuid,
        @JsonProperty("code") String code) {

    public ItemEvent toDomain() {
        return new ItemEvent(organizationId, referenceGuid, ItemCode.valueOf(code), null);
    }
}
```

## Dependency Injection

Use **Spring constructor injection** for all dependency wiring. Lombok's `@RequiredArgsConstructor` generates the constructor automatically from `final` fields.

```java
// Correct -- constructor injection via Lombok
@Component
@RequiredArgsConstructor
public class InsertItemCommand {
    private final SendItemService service;
    private final ItemsRepository repository;
}

// Wrong -- field injection
@Component
public class InsertItemCommand {
    @Autowired
    private SendItemService service;
}
```

Constructor injection ensures:
- Dependencies are explicit and visible.
- Objects are fully initialized upon creation.
- Fields can be declared `final`, enforcing immutability.
- Unit tests can inject dependencies without Spring context.

## References

- [Spring Framework Reference](https://docs.spring.io/spring-framework/reference/)
- [MapStruct Reference Guide](https://mapstruct.org/documentation/stable/reference/html/)
- [Lombok Features](https://projectlombok.org/features/)
- [DTO Pattern](https://www.baeldung.com/java-dto-pattern)
