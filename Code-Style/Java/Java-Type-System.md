# Java Type System

> **TL;DR:** Java is statically typed -- the compiler enforces type safety at build time. Use **records** for immutable data, **generics** with bounded type parameters for reusable abstractions, and **sealed classes** (Java 17+) for restricted type hierarchies. Never use raw types or `Object` as a catch-all parameter.

## Overview

Java's type system is static and checked at compile time. Combined with generics, records, and sealed types, it provides strong guarantees about program correctness. This page focuses on the patterns and principles that maximize type safety and code clarity in Java projects.

## Records

Use Java records (Java 16+) for immutable data carriers: DTOs, domain events, value objects, and listener contracts.

```java
// Correct -- immutable domain event
public record ItemEvent(
        Long organizationId,
        String referenceGuid,
        ItemCode code,
        ItemCategory category) {}

// Correct -- request DTO with validation
public record InsertItemRequest(
        @NotNull @Size(min = 1, max = 255) String name,
        @NotNull ItemCode code) {}

// Correct -- listener callbacks as a record
public record Listeners(Runnable onSuccess, Consumer<Exception> onError) {}
```

Records automatically provide:
- `final` fields with constructor, getters, `equals()`, `hashCode()`, and `toString()`
- Immutability by design
- Compact and readable declarations

## Generics

Use generics to create type-safe, reusable abstractions. Always use bounded type parameters when the generic type must satisfy constraints.

### When to Use

```java
// Correct -- generic entity that works with different related types
public final class Item<A> {
    private final ItemCode code;
    private final ItemSeverity severity;
    private List<A> affected;

    public Boolean hasData() {
        return !affected.isEmpty();
    }
}

// Correct -- generic repository interface
public interface QueryDslItemsRepository<T> extends ItemsRepository {
    Page<T> findAllWithFilters(Map<String, Object> filters, Pageable pageable);
}
```

### When Not to Use

```java
// Unnecessary -- only works with one type
public class ItemProcessor<T extends Item<?>> {
    public void process(T item) { ... }
}

// Better -- use the concrete type directly
public class ItemProcessor {
    public void process(Item<?> item) { ... }
}
```

### Bounded Type Parameters

Use bounded type parameters to constrain generic types:

```java
// Upper bound -- T must be Comparable
public static <T extends Comparable<T>> T max(List<T> items) {
    return items.stream().max(Comparator.naturalOrder()).orElseThrow();
}

// Multiple bounds
public static <T extends Serializable & Comparable<T>> void sort(List<T> items) {
    Collections.sort(items);
}
```

## Annotations

### Jakarta Validation

Use Jakarta Bean Validation annotations on request DTOs to validate input at the system boundary:

```java
public record InsertItemRequest(
        @NotNull @Size(min = 1, max = 255) String name,
        @NotNull ItemCode code,
        @Min(0) Long amount) {}
```

### JPA/Persistence

JPA annotations belong **only** on infrastructure models (`Jpa*` classes), never on domain entities:

```java
// Correct -- annotations on infrastructure model
@Entity
@Table(name = "item")
public class JpaItem {
    @Id
    @GeneratedValue(strategy = GenerationType.SEQUENCE)
    private Long id;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "category_id")
    private JpaCategory category;
}

// Wrong -- annotations on domain entity
public class Item {
    @Id
    private Long id;  // Domain entities must be framework-free
}
```

### Lombok

Use Lombok to reduce boilerplate while maintaining clarity:

| Annotation                 | Purpose                                                          |
|----------------------------|------------------------------------------------------------------|
| `@Getter` / `@Setter`      | Generate accessors                                               |
| `@NoArgsConstructor`       | JPA requirement for entity classes                               |
| `@RequiredArgsConstructor` | Constructor injection (generates constructor for `final` fields) |
| `@SuperBuilder`            | Fluent builder pattern with inheritance support                  |
| `@Slf4j`                   | Generate SLF4J logger field                                      |

### Spring

| Annotation        | Layer          | Purpose                    |
|-------------------|----------------|----------------------------|
| `@Component`      | Domain         | Commands, listeners        |
| `@RestController` | Infrastructure | HTTP controllers           |
| `@Service`        | Infrastructure | Service implementations    |
| `@Repository`     | Infrastructure | Repository implementations |

## Sealed Classes (Java 17+)

Use sealed classes and interfaces to restrict type hierarchies:

```java
// Only the listed classes can extend ItemResult
public sealed interface ItemResult
        permits ItemResult.Success, ItemResult.NotFound, ItemResult.Error {

    record Success(Item item) implements ItemResult {}
    record NotFound(Long id) implements ItemResult {}
    record Error(Exception cause) implements ItemResult {}
}
```

Sealed types work well with `switch` expressions (Java 21+ pattern matching):

```java
return switch (result) {
    case ItemResult.Success s -> ResponseEntity.ok(s.item());
    case ItemResult.NotFound n -> ResponseEntity.notFound().build();
    case ItemResult.Error e -> ResponseEntity.internalServerError().build();
};
```

## Prohibited Patterns

```java
// Wrong -- raw type (loses type safety)
List items = new ArrayList();

// Wrong -- Object as catch-all parameter
public void process(Object data) { ... }

// Wrong -- unchecked cast without type guard
Item item = (Item) someObject;

// Wrong -- using Optional as a field or parameter
public class ItemHolder {
    private Optional<Item> item;  // Optional is for return types only
}
```

## References

- [Java Records](https://docs.oracle.com/en/java/javase/21/language/records.html)
- [Java Generics](https://docs.oracle.com/javase/tutorial/java/generics/)
- [Sealed Classes](https://docs.oracle.com/en/java/javase/21/language/sealed-classes-and-interfaces.html)
- [Jakarta Bean Validation](https://beanvalidation.org/)
- [Effective Java -- Item 26: Don't use raw types](https://www.oreilly.com/library/view/effective-java/9780134686097/)
