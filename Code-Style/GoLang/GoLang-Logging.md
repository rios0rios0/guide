# Go Logging

> **TL;DR:** Use **[Logrus](https://github.com/sirupsen/logrus)** for all logging. Do not use Go's standard `log` package or `fmt.Println` for application logging. Always import with the alias `logger`. Use structured logging with `WithFields()` instead of string interpolation.

## Overview

Consistent, structured logging is essential for production observability. This page defines the mandatory logging library and patterns for all Go projects.

## Mandatory Library: Logrus

**Use [Logrus](https://github.com/sirupsen/logrus) for all logging.** Logrus provides structured logging, consistent log levels, JSON output support, and field-based contextual logging -- all of which are essential for production observability.

### Installation

```bash
go get github.com/sirupsen/logrus
```

**Important:** Always use the lowercase import path `github.com/sirupsen/logrus` (not the uppercase variant).

### Import Convention

Always import Logrus with the alias `logger` to keep usage concise and consistent across the codebase:

```go
import logger "github.com/sirupsen/logrus"
```

### Usage

```go
import logger "github.com/sirupsen/logrus"

func main() {
    logger.Info("application started")
    logger.WithFields(logger.Fields{
        "user_id": 42,
        "action":  "login",
    }).Info("user authenticated")
}
```

## Log Levels

| Level | Method           | When to Use                                                         |
|-------|------------------|---------------------------------------------------------------------|
| Trace | `logger.Trace()` | Very fine-grained diagnostic information                            |
| Debug | `logger.Debug()` | Diagnostic information useful during development                    |
| Info  | `logger.Info()`  | General operational events (application started, request processed) |
| Warn  | `logger.Warn()`  | Potential issues that do not prevent operation                      |
| Error | `logger.Error()` | Errors that prevent a specific operation but not the application    |
| Fatal | `logger.Fatal()` | Critical errors that require immediate application shutdown         |
| Panic | `logger.Panic()` | Critical errors that should panic after logging                     |

## Structured Logging

Always use `WithFields` to attach contextual data to log entries rather than interpolating values into the message string:

```go
// Correct -- structured fields
logger.WithFields(logger.Fields{
    "request_id": requestID,
    "status":     statusCode,
}).Info("request completed")

// Wrong -- string interpolation
logger.Infof("request %s completed with status %d", requestID, statusCode)
```

## Prohibited Patterns

```go
// Wrong -- standard library logger
import "log"
log.Println("something happened")

// Wrong -- fmt for application logging
fmt.Println("something happened")

// Wrong -- uppercase import path
import "github.com/Sirupsen/logrus"
```

## References

- [Logrus - Structured Logger for Go](https://github.com/sirupsen/logrus)
- [Logrus Logging Guide](https://betterstack.com/community/guides/logging/logrus/)
