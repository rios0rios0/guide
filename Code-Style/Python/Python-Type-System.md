# Python Type System

> **TL;DR:** Add type hints to all function parameters, return types, and non-obvious variables. Type hints enable static analysis, catch bugs before runtime, and make code self-documenting.

## Overview

Python is dynamically typed, which means type errors often surface only at runtime -- or worse, in production. [PEP 484](https://peps.python.org/pep-0484/) introduced type hints, allowing developers to declare expected types for variables and function signatures. Modern linters and IDEs interpret these annotations to warn about type mismatches during development, significantly reducing debugging time.

## Type Annotation Requirements

**All function parameters and return types must have type hints.** Variables should be annotated when the type is not obvious from the assignment.

### Why Type Hints Matter

Consider the following function without type hints:

```python
def collatz(n):
    if n % 2 == 0:
        return n / 2
    else:
        return 3 * n + 1

print(collatz(5).denominator)   # Works (int has .denominator)
print(collatz(6).denominator)   # Fails (float has no .denominator)
```

Without type annotations, this bug is invisible until execution. With type hints, the linter catches the mismatch immediately:

```python
def collatz(n: int) -> int:
    if n % 2 == 0:
        return n / 2    # Linter warns: float is incompatible with int
    else:
        return 3 * n + 1
```

## Syntax

### Functions

```python
def function(param: <type>) -> <return_type>:
    ...
```

Example:

```python
import pathlib

def delete_file(path: pathlib.Path) -> int:
    """Delete the file at the given path.

    :param path: Path to the file to delete.
    :return: 0 if deleted successfully, 1 otherwise.
    """
    if path.is_file():
        path.unlink()
        return 0
    return 1
```

### Variables

```python
<variable>: <type> = <value>
```

Example:

```python
import queue

item: int = 1
some_queue: queue.Queue[int] = queue.Queue()
some_queue.put(item)
```

## Key Type Patterns

### The `typing` Module

The `typing` module supports advanced type constructs such as `Union`, `Optional`, `Type`, and generics:

```python
from typing import Type, Union

def func(param: Type[Union[ClassA, ClassB]]) -> None:
    ...
```

Alternative syntax:

```python
from typing import Type, Union

def func(param: Union[Type[ClassA], Type[ClassB]]) -> None:
    ...
```

### Python 3.10+ Union Syntax

For Python 3.10+, use the built-in `|` operator instead of `Union`:

```python
def func(param: type[ClassA | ClassB]) -> None:
    ...
```

## Prohibited Patterns

```python
# Wrong -- missing type hints on function signatures
def process(data):
    return data.transform()

# Wrong -- using Any as a catch-all
from typing import Any
def process(data: Any) -> Any:
    return data.transform()
```

Define proper types or use `Protocol` for structural typing instead.

## References

- [PEP 484 -- Type Hints](https://peps.python.org/pep-0484/)
- [typing Module Documentation](https://docs.python.org/3/library/typing.html)
- [mypy -- Static Type Checker](https://mypy.readthedocs.io/)
