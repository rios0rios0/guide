# Java Logging

> **TL;DR:** Use **[SLF4J](https://www.slf4j.org/)** with **[Logback](https://logback.qos.ch/)** (the Spring Boot default) for all logging. Use Lombok's `@Slf4j` annotation to generate the logger field. Use `{}` placeholders for structured logging instead of string concatenation. Never use `System.out.println` or `java.util.logging`.

## Overview

Consistent, structured logging is essential for production observability. This page defines the mandatory logging library and patterns for all Java projects.

## Mandatory Library: SLF4J with Logback

**Use [SLF4J](https://www.slf4j.org/) with [Logback](https://logback.qos.ch/) for all logging.** This is the default logging framework in Spring Boot. SLF4J provides the API (facade), and Logback provides the implementation -- together they offer structured logging, configurable levels, and JSON output for production environments.

### Installation

Spring Boot includes SLF4J and Logback by default. No additional dependencies are required:

```groovy
// Already included via spring-boot-starter
implementation 'org.springframework.boot:spring-boot-starter'
```

### Import Convention

Use Lombok's `@Slf4j` annotation to generate the logger field automatically:

```java
import lombok.extern.slf4j.Slf4j;

@Slf4j
@Component
public class InsertItemCommand {
    public void execute(final ItemEvent event) {
        log.info("Processing item event for organization: {}", event.organizationId());
    }
}
```

When Lombok is not available, create the logger manually:

```java
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class InsertItemCommand {
    private static final Logger log = LoggerFactory.getLogger(InsertItemCommand.class);
}
```

## Log Levels

| Level | Method        | When to Use                                                                                      |
|-------|---------------|--------------------------------------------------------------------------------------------------|
| TRACE | `log.trace()` | Very fine-grained diagnostic information (method entry/exit, loop iterations)                    |
| DEBUG | `log.debug()` | Diagnostic information useful during development (variable values, flow decisions)               |
| INFO  | `log.info()`  | General operational events (application started, request processed, message sent)                |
| WARN  | `log.warn()`  | Potential issues that do not prevent operation (deprecated API usage, retry attempts)            |
| ERROR | `log.error()` | Errors that prevent a specific operation but not the application (failed request, invalid input) |

**Note:** Unlike some logging frameworks, SLF4J does not have a `FATAL` level. Use `ERROR` for critical errors and handle application shutdown via Spring's shutdown hooks.

## Structured Logging

Always use `{}` placeholders to attach contextual data to log entries rather than concatenating values into the message string:

```java
// Correct -- placeholder substitution
log.info("Message sent successfully to topic: {}", topicName);
log.error("Error sending message {} to topic: {}", messageId, topicName, throwable);

// Wrong -- string concatenation
log.info("Message sent successfully to topic: " + topicName);

// Wrong -- String.format
log.info(String.format("Message sent to topic: %s", topicName));
```

### Exception Logging

When logging exceptions, pass the throwable as the **last argument** -- SLF4J automatically appends the stack trace:

```java
try {
    service.send(event);
} catch (Exception e) {
    // Correct -- throwable as last argument, stack trace printed automatically
    log.error("Failed to send item event for organization: {}", event.organizationId(), e);
}
```

### Structured Fields (MDC)

Use the Mapped Diagnostic Context (MDC) for contextual fields that span multiple log entries:

```java
import org.slf4j.MDC;

MDC.put("organizationId", String.valueOf(organizationId));
MDC.put("requestId", requestId);
try {
    log.info("Processing request");
    // ... business logic ...
    log.info("Request completed");
} finally {
    MDC.clear();
}
```

## Configuration

Configure logging in `application.yaml`:

```yaml
logging:
  level:
    root: 'warn'
    com.example.myapp: 'info'
  pattern:
    console: '%d{yyyy-MM-dd HH:mm:ss.SSS} %5p %.20logger{39} : %m%n%wEx'
```

## Prohibited Patterns

```java
// Wrong -- System.out for application logging
System.out.println("Something happened");
System.err.println("Error occurred");

// Wrong -- java.util.logging
import java.util.logging.Logger;
Logger.getLogger("MyClass").info("message");

// Wrong -- string concatenation in log messages
log.info("Processing item " + item.getId() + " for org " + orgId);

// Wrong -- unnecessary isEnabled check with placeholder syntax
if (log.isDebugEnabled()) {
    log.debug("Processing item: {}", item.getId());  // SLF4J already skips evaluation
}
```

**Note:** The `isEnabled` guard is only necessary when the argument computation itself is expensive (e.g., serializing a large object). For simple field access, SLF4J's lazy evaluation via `{}` placeholders is sufficient.

## References

- [SLF4J Manual](https://www.slf4j.org/manual.html)
- [Logback Documentation](https://logback.qos.ch/documentation.html)
- [Spring Boot Logging](https://docs.spring.io/spring-boot/reference/features/logging.html)
- [Lombok @Slf4j](https://projectlombok.org/features/log)
