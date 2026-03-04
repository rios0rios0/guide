# Mapper Design Pattern (Replacing Switch/Case)

> **TL;DR:** Replace `switch`/`case` statements with a map (dictionary) that associates each key with a handler function. This eliminates branching, simplifies adding new cases, and follows the Open/Closed Principle.

## Overview

The `switch`/`case` construct is one of the most common sources of code that violates the **Open/Closed Principle** -- every time a new case is needed the existing function must be modified. The Mapper Design Pattern replaces the branching logic with a **lookup map** that associates each key with its handler. Adding a new case becomes adding a new entry to the map, without touching the dispatch logic.

### Benefits

| Benefit              | Explanation                                                              |
|----------------------|--------------------------------------------------------------------------|
| **Open/Closed**      | New cases are added by inserting a map entry, not by editing a function. |
| **Readability**      | The mapping is declarative -- you see all cases at a glance.             |
| **Testability**      | Each handler can be tested in isolation.                                 |
| **Single dispatch**  | The lookup + call replaces an unbounded chain of comparisons.            |

## Problem -- Switch/Case

The examples below use a notification system that sends messages through different channels. The `switch` version hard-codes every channel:

### JavaScript (Switch)

```javascript
function notify(channel, message) {
  switch (channel) {
    case "email":
      sendEmail(message);
      break;
    case "sms":
      sendSms(message);
      break;
    case "slack":
      sendSlack(message);
      break;
    default:
      throw new Error(`Unknown channel: ${channel}`);
  }
}
```

### Java (Switch)

```java
public class NotificationService {
    public void notify(final String channel, final String message) {
        switch (channel) {
            case "email" -> sendEmail(message);
            case "sms" -> sendSms(message);
            case "slack" -> sendSlack(message);
            default -> throw new IllegalArgumentException("Unknown channel: " + channel);
        }
    }
}
```

### Go (Switch)

```go
func notify(channel, message string) error {
	switch channel {
	case "email":
		return sendEmail(message)
	case "sms":
		return sendSms(message)
	case "slack":
		return sendSlack(message)
	default:
		return fmt.Errorf("unknown channel: %s", channel)
	}
}
```

### Python (Switch)

```python
def notify(channel: str, message: str) -> None:
    if channel == "email":
        send_email(message)
    elif channel == "sms":
        send_sms(message)
    elif channel == "slack":
        send_slack(message)
    else:
        raise ValueError(f"Unknown channel: {channel}")
```

## Solution -- Mapper Design Pattern

Replace the branching with a map from key to handler. The dispatch function looks up the key and calls the handler.

### JavaScript

```javascript
const notificationHandlers = new Map([
  ["email", sendEmail],
  ["sms", sendSms],
  ["slack", sendSlack],
]);

function notify(channel, message) {
  const handler = notificationHandlers.get(channel);
  if (!handler) {
    throw new Error(`Unknown channel: ${channel}`);
  }
  handler(message);
}
```

Adding a new channel (e.g., push notifications) requires only one line:

```javascript
notificationHandlers.set("push", sendPush);
```

### Java

```java
public class NotificationService {
    private static final Map<String, Consumer<String>> HANDLERS = Map.of(
            "email", NotificationService::sendEmail,
            "sms", NotificationService::sendSms,
            "slack", NotificationService::sendSlack
    );

    public void notify(final String channel, final String message) {
        final var handler = HANDLERS.get(channel);
        if (handler == null) {
            throw new IllegalArgumentException("Unknown channel: " + channel);
        }
        handler.accept(message);
    }
}
```

For richer handler contracts, define a functional interface or use a strategy class:

```java
@FunctionalInterface
public interface NotificationHandler {
    void send(String message);
}

public class NotificationService {
    private final Map<String, NotificationHandler> handlers;

    public NotificationService(final Map<String, NotificationHandler> handlers) {
        this.handlers = handlers;
    }

    public void notify(final String channel, final String message) {
        final var handler = handlers.get(channel);
        if (handler == null) {
            throw new IllegalArgumentException("Unknown channel: " + channel);
        }
        handler.send(message);
    }
}
```

This second form also supports **dependency injection** -- the handler map can be wired by Spring:

```java
@Configuration
public class NotificationConfig {
    @Bean
    public Map<String, NotificationHandler> notificationHandlers(
            final EmailHandler email,
            final SmsHandler sms,
            final SlackHandler slack) {
        return Map.of("email", email, "sms", sms, "slack", slack);
    }
}
```

### Go

```go
type NotificationHandler func(message string) error

var notificationHandlers = map[string]NotificationHandler{
	"email": sendEmail,
	"sms":   sendSms,
	"slack": sendSlack,
}

func notify(channel, message string) error {
	handler, ok := notificationHandlers[channel]
	if !ok {
		return fmt.Errorf("unknown channel: %s", channel)
	}
	return handler(message)
}
```

For handlers that require dependencies, use a struct with a method:

```go
type NotificationHandler interface {
	Send(message string) error
}

type NotificationService struct {
	handlers map[string]NotificationHandler
}

func NewNotificationService(handlers map[string]NotificationHandler) *NotificationService {
	return &NotificationService{handlers: handlers}
}

func (s *NotificationService) Notify(channel, message string) error {
	handler, ok := s.handlers[channel]
	if !ok {
		return fmt.Errorf("unknown channel: %s", channel)
	}
	return handler.Send(message)
}
```

### Python

```python
from typing import Callable

notification_handlers: dict[str, Callable[[str], None]] = {
    "email": send_email,
    "sms": send_sms,
    "slack": send_slack,
}


def notify(channel: str, message: str) -> None:
    handler = notification_handlers.get(channel)
    if handler is None:
        raise ValueError(f"Unknown channel: {channel}")
    handler(message)
```

For handlers that need dependencies, use a class-based approach:

```python
from abc import ABC, abstractmethod


class NotificationHandler(ABC):
    @abstractmethod
    def send(self, message: str) -> None: ...


class NotificationService:
    def __init__(self, handlers: dict[str, NotificationHandler]) -> None:
        self._handlers = handlers

    def notify(self, channel: str, message: str) -> None:
        handler = self._handlers.get(channel)
        if handler is None:
            raise ValueError(f"Unknown channel: {channel}")
        handler.send(message)
```

## When to Use

| Scenario                                               | Recommendation   |
|--------------------------------------------------------|------------------|
| 2-3 simple cases unlikely to grow                      | `switch` is fine |
| 4+ cases or cases that change often                    | Use the Mapper   |
| Handlers with different dependencies or configurations | Use the Mapper   |
| Cases determined at runtime (e.g., plugins)            | Use the Mapper   |

## References

- [Open/Closed Principle -- Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2014/05/12/TheOpenClosedPrinciple.html)
- [Replace Conditional with Polymorphism -- Refactoring Guru](https://refactoring.guru/replace-conditional-with-polymorphism)
- [Strategy Pattern -- Refactoring Guru](https://refactoring.guru/design-patterns/strategy)
