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

---

# Forking Technique

> **TL;DR:** Reserve `main` for the upstream community version and use `custom` as the working branch. Rebase `custom` onto updated `main` when upstream releases new versions. Tag fork releases with an incremental fourth digit (e.g., `1.0.0.1`).

## Overview

When forking an open-source project, it is common for upstream maintainers to take significant time reviewing and merging community contributions. This standard defines a forking strategy that maintains compatibility with upstream releases while preserving custom modifications.

## Strategy

### Branch Convention

- The **`main`** branch mirrors the upstream community version. It is updated by rebasing on upstream releases.
- The **`custom`** branch serves as the team's working branch (equivalent to `main` for internal purposes).

### Synchronization

To stay current with the upstream project:

1. Update `main` by rebasing on the latest upstream release.
2. Rebase `custom` onto the updated `main`, placing custom modifications on top of the newest version.

### Versioning

Fork versions use an **incremental fourth digit** appended to the upstream version:

| Scenario                                  | Version   |
|-------------------------------------------|-----------|
| Upstream at `1.0.0`, fork synced          | `1.0.0.0` |
| New fork release                          | `1.0.0.1` |
| Upstream updates to `1.0.1`, fork rebased | `1.0.1.0` |
| Next fork release                         | `1.0.1.1` |

The fourth digit resets to `0` each time the fork is rebased on a new upstream version.

## Caveats

If your CI/CD or release tooling does not support four-segment version numbers (`X.Y.Z.N`), use a dash separator instead: `X.Y.Z-N`.
