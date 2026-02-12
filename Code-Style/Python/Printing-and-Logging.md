# Printing and Logging

> **TL;DR:** Prefer **Loguru** over the standard `logging` library for its simpler API, colorized output, and structured exception formatting. Use STDOUT for normal output and STDERR for warnings and errors.

## Overview

Programs frequently need to output messages for progress tracking, state reporting, and error diagnostics. This page covers the fundamentals of printing and logging in Python, from basic `print()` usage to production-grade logging with Loguru.

## Printing

The `print()` function writes to STDOUT by default. In production code, ensure messages are directed to the appropriate output stream:

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

Advanced formatting (refer to [the documentation](https://docs.python.org/3/tutorial/inputoutput.html) for full details):

```python
>>> print(f"{10.12345:x^20.2f}")
xxxxxxx10.12xxxxxxxx
```

## Standard Library: `logging`

Python's built-in [`logging`](https://docs.python.org/3/library/logging.html) module provides basic logging capabilities:

```python
import logging
logging.basicConfig(format="%(asctime)s %(message)s", level=logging.DEBUG)
logging.info("Application started")
```

| Level    | Method               |
|----------|----------------------|
| DEBUG    | `logging.debug()`    |
| INFO     | `logging.info()`     |
| WARNING  | `logging.warning()`  |
| ERROR    | `logging.error()`    |
| CRITICAL | `logging.critical()` |

For advanced usage, refer to the [logging HOWTO](https://docs.python.org/3/howto/logging.html).

## Third-Party Library: Loguru (Recommended)

[Loguru](https://loguru.readthedocs.io/) is the recommended logging library for all projects. It provides a richer feature set with minimal configuration:

```python
from loguru import logger
logger.info("Application started")
```

### Log Levels

| Level    | Severity | Method              |
|----------|----------|---------------------|
| TRACE    | 5        | `logger.trace()`    |
| DEBUG    | 10       | `logger.debug()`    |
| INFO     | 20       | `logger.info()`     |
| SUCCESS  | 25       | `logger.success()`  |
| WARNING  | 30       | `logger.warning()`  |
| ERROR    | 40       | `logger.error()`    |
| CRITICAL | 50       | `logger.critical()` |

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

Loguru's `logger.exception()` produces colorized, structured exception output with variable values labeled inline, making debugging significantly faster.

## References

- [Python `logging` Documentation](https://docs.python.org/3/library/logging.html)
- [Loguru Documentation](https://loguru.readthedocs.io/)
- [Python String Formatting](https://docs.python.org/3/tutorial/inputoutput.html)
