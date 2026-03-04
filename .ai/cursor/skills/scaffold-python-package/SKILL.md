---
name: scaffold-python-package
description: Scaffold a new Python package with PEP 621, PDM, Black/isort/Flake8, type hints, and Loguru logging. Use when the user asks to create, scaffold, or bootstrap a new Python project or package.
---

Scaffold a new Python package following PEP 621, PDM, Black/isort/Flake8, type hints, and Loguru logging.

For detailed Python conventions, refer to the Python rule. For testing patterns, refer to the Testing rule. For Makefile setup, refer to the CI/CD rule.

## Directory Structure

```
<project>/
├── <package_name>/
│   ├── __init__.py
│   └── main.py
├── tests/
│   ├── __init__.py
│   └── test_main.py
├── examples/
│   └── example_usage.py
├── Dockerfile
├── pyproject.toml
└── README.md
```

## Step-by-Step

### 1. Initialize with PDM

```bash
mkdir <project> && cd <project>
pdm init
```

### 2. Configure pyproject.toml (PEP 621)

```toml
[project]
name = "package-name"
version = "0.1.0"
description = "Brief description of the package"
readme = "README.md"
requires-python = ">=3.10"
license = {text = "MIT"}
authors = [
    {name = "Author Name", email = "author@example.com"},
]
dependencies = []

[project.optional-dependencies]
dev = [
    "black",
    "isort",
    "flake8",
    "pytest",
    "loguru",
]

[build-system]
requires = ["pdm-backend"]
build-backend = "pdm.backend"

[tool.black]
line-length = 120

[tool.isort]
profile = "black"
line_length = 120

[tool.flake8]
max-line-length = 120
```

### 3. Create the package module

```python
# <package_name>/__init__.py
"""Package description."""

__version__ = "0.1.0"
```

```python
# <package_name>/main.py
from loguru import logger


def main() -> None:
    """Entry point for the package."""
    logger.info("Starting application")
```

### 4. Add type hints everywhere

All functions must have type hints on parameters and return types:

```python
def process_data(input_data: list[str], max_items: int = 10) -> dict[str, int]:
    """Process the input data and return counts."""
    result: dict[str, int] = {}
    for item in input_data[:max_items]:
        result[item] = result.get(item, 0) + 1
    return result
```

### 5. Set up logging with Loguru

Always use Loguru -- NEVER use the standard `logging` module or `print()`:

```python
from loguru import logger

logger.info("Processing started")
logger.error("Something failed")

try:
    risky_operation()
except Exception:
    logger.exception("Formatted exception with traceback")
```

### 6. Create tests following BDD structure

```python
# tests/test_main.py
from package_name.main import main


def test_main_runs_without_error() -> None:
    # given
    # (no preconditions needed)

    # when / then
    main()  # should not raise
```

Every test must use `# given`, `# when`, `# then` comment blocks. Refer to the Testing rule for full conventions.

### 7. Install dev dependencies

```bash
pdm add -dG dev black isort flake8 pytest loguru
```

### 8. Create Makefile

The project Makefile must import from the shared [pipelines repository](https://github.com/rios0rios0/pipelines) and expose `lint`, `test`, and `sast` targets. Refer to the CI/CD rule for details. Always use `make lint` and `make test`, never invoke tools directly.

## Key Rules

- Use **Black** for formatting, **isort** for imports, **Flake8** for linting
- Use **Loguru** (`from loguru import logger`) -- never `logging` or `print()`
- Add type hints to ALL function parameters and return types
- Avoid `Any` entirely -- use `unknown` narrowing patterns instead
- Avoid `setup.py` and `setup.cfg` -- use `pyproject.toml` exclusively
- All file names use **snake_case**

## Formatting Tools Summary

| Tool       | Purpose         | Config                                  |
|------------|-----------------|-----------------------------------------|
| **Black**  | Code formatting | `[tool.black]` in pyproject.toml        |
| **isort**  | Import sorting  | `[tool.isort]` with `profile = "black"` |
| **Flake8** | Linting         | `[tool.flake8]` in pyproject.toml       |
