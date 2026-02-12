# Python Packaging

> **TL;DR:** Use **PDM** as the package manager with **PEP 621** metadata in `pyproject.toml`. PDM is lightweight, standards-compliant, and easily replaceable by any PEP 621-compatible tool.

## Overview

This document provides a step-by-step guide for creating a new Python package. While many packaging tools exist, this guide recommends PDM for the following reasons:

- **PEP 621 support:** All metadata lives in a single `pyproject.toml` file.
- **Portability:** PEP 621 is an official standard, so PDM can be replaced by any compliant tool.
- **Lightweight:** Unlike Poetry, PDM does not require GCC-compiled components.

### Tool Comparison

| Tool       | Description                                                                      |
|------------|----------------------------------------------------------------------------------|
| Conda      | Cross-platform binary package manager (primarily for Anaconda)                   |
| Setuptools | Fully-featured, legacy package builder                                           |
| Pipenv     | Virtualenv and dependency manager                                                |
| Poetry     | Dependency manager + builder + publisher (uses custom `pyproject.toml` sections) |
| **PDM**    | Dependency manager + builder + publisher with full PEP 621 support               |

For a more comprehensive comparison, see [Python Packaging Tools Comparison](https://chadsmith.dev/python-packaging/).

## Step-by-Step: Package Initialization with PDM

### 1. Install PDM

```shell
pip install -U pdm
```

### 2. Create and Initialize the Package

```shell
mkdir icarus && cd icarus
pdm init
```

Answer the interactive prompts according to your package's requirements. Use `Proprietary` as the license for non-open-source packages.

### 3. Create the Package Module

```shell
mkdir icarus
```

Create the initialization file:
```python
# icarus/__init__.py
__version__ = "1.0.0"
from .icarus import main
```

Create the `__main__` file (enables `python -m icarus`):
```python
# icarus/__main__.py
import sys
from .icarus import main

if __name__ == "__main__":
    sys.exit(main())
```

Create the entrypoint:
```python
# icarus/icarus.py
from . import __version__

def main() -> int:
    print(f"Icarus v{__version__}")
    return 0
```

### 4. Configure `pyproject.toml`

Enable dynamic versioning so the version is defined once in `__init__.py`:

```toml
[project]
name = "icarus"
description = "A demo package"
authors = [{ name = "Author", email = "author@example.com" }]
dependencies = []
requires-python = ">=3.10"
license = { text = "Proprietary" }
dynamic = ["version"]

[project.urls]
Homepage = "https://github.com/org/icarus"

[tool.pdm]
version = { from = "icarus/__init__.py" }

[build-system]
requires = ["pdm-pep517"]
build-backend = "pdm.pep517.api"
```

### 5. Create `README.md`

Provide usage instructions in the project README.

### 6. Final Structure

```
icarus/
├── icarus/
│   ├── __init__.py
│   ├── __main__.py
│   └── icarus.py
├── pyproject.toml
└── README.md
```

### Useful Commands

```bash
pdm add 'opencv-python>=4.5'     # Add a dependency
pip install .                      # Install the package locally
python3 -m build -s .             # Build source distribution
python3 -m build -w .             # Build wheel distribution
```

## Next Steps

- [Package Metadata Formats](Package-Metadata-Formats.md) -- Understand the evolution from `setup.py` to PEP 621.
- [Standard Package Layout](Standard-Package-Layout.md) -- Learn where to place non-code files.

## References

- [PDM Documentation](https://pdm-project.org/)
- [Pip Documentation](https://pip.pypa.io/)
- [PEP 621 -- Storing Project Metadata in pyproject.toml](https://peps.python.org/pep-0621/)
