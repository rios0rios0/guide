# Python Logging

> **TL;DR:** Use **[Loguru](https://loguru.readthedocs.io/)** for all application logging. Do not use the standard `logging` module or `print()` for logging in production code. Always import as `from loguru import logger`. Use STDOUT for normal output and STDERR for warnings and errors.

## Overview

Programs frequently need to output messages for progress tracking, state reporting, and error diagnostics. This page defines the mandatory logging library and patterns for all Python projects.

## Mandatory Library: Loguru

**Use [Loguru](https://loguru.readthedocs.io/) for all application logging.** Do not use Python's built-in `logging` module. Loguru provides a simpler API, colorized output, structured exception formatting, and sensible defaults with minimal configuration.

### Installation

```bash
pip install loguru
# or via PDM
pdm add loguru
```

### Import Convention

```python
from loguru import logger
```

### Usage

```python
from loguru import logger

logger.info("application started")
logger.warning("disk usage above threshold")
logger.error("failed to connect to database")
```

## Log Levels

| Level    | Severity | Method              | When to Use                                               |
|----------|----------|---------------------|-----------------------------------------------------------|
| TRACE    | 5        | `logger.trace()`    | Very fine-grained diagnostic information                  |
| DEBUG    | 10       | `logger.debug()`    | Diagnostic information useful during development          |
| INFO     | 20       | `logger.info()`     | General operational events (application started, request processed) |
| SUCCESS  | 25       | `logger.success()`  | Successful completion of a significant operation          |
| WARNING  | 30       | `logger.warning()`  | Potential issues that do not prevent operation             |
| ERROR    | 40       | `logger.error()`    | Errors that prevent a specific operation but not the application |
| CRITICAL | 50       | `logger.critical()` | Critical errors that require immediate attention          |

## Structured Logging

### Custom Format

The default format can be verbose. A recommended concise format:

```python
import sys
from loguru import logger

LOGURU_FORMAT = (
    "<green>{time:HH:mm:ss.SSSSSS!UTC}</green> "
    "<level>{level: <8}</level> "
    "<level>{message}</level>"
)
logger.remove()
logger.add(sys.stderr, colorize=True, format=LOGURU_FORMAT)
```

### Exception Formatting

Loguru's `logger.exception()` produces colorized, structured exception output with variable values labeled inline, making debugging significantly faster:

```python
try:
    risky_operation()
except Exception:
    logger.exception("operation failed with unexpected error")
```

## Printing

The `print()` function writes to STDOUT by default. It is acceptable for simple scripts and CLI output, but **must not be used as a substitute for logging** in application code. In production, ensure messages are directed to the appropriate output stream:

```python
import sys
print("Fetching data from Redis")
print("Error: Could not connect to Redis", file=sys.stderr)
```

### String Formatting

Use f-strings (preferred) or `.format()` for variable interpolation:

```python
number = 42
print(f"The number is: {number}")
print("The number is: {}".format(number))
```

## Prohibited Patterns

The following patterns must **not** be used for application logging:

```python
# Wrong -- standard library logging module
import logging
logging.info("something happened")
logging.getLogger(__name__).warning("issue detected")

# Wrong -- print for application logging
print("Processing request...")

# Wrong -- f-string logging without Loguru
print(f"Error: {error_message}")
```

**Correct equivalent:**

```python
from loguru import logger

logger.info("something happened")
logger.warning("issue detected")
logger.info("processing request")
logger.error(f"error: {error_message}")
```

## Why Not the Standard `logging` Module?

The standard `logging` module is part of Python's standard library and provides basic functionality. However, it requires significant boilerplate to configure properly (formatters, handlers, log levels per module), its default output is unstructured, and it lacks features like colorized output and automatic exception variable labeling. Loguru provides all of these out of the box with a single import.

For reference, the standard library supports five levels: DEBUG, INFO, WARNING, ERROR, and CRITICAL. Full details are in the [logging HOWTO](https://docs.python.org/3/howto/logging.html). However, **this is provided for awareness only -- always use Loguru in our projects.**

## References

- [Loguru Documentation](https://loguru.readthedocs.io/)
- [Loguru GitHub Repository](https://github.com/Delgan/loguru)
- [Python `logging` Documentation](https://docs.python.org/3/library/logging.html) (reference only)
- [Python String Formatting](https://docs.python.org/3/tutorial/inputoutput.html)
