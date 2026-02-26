# Python Project Structure

> **TL;DR:** Use **PDM** as the package manager with **PEP 621** metadata in `pyproject.toml`. Follow the conventional directory layout: `/<package_name>/` for source code, `/tests/` for tests, `/examples/` for usage examples. Avoid legacy `setup.py` and `setup.cfg` formats.

## Overview

This page defines the standard directory layout, package management, and distribution practices for all Python projects. It consolidates directory structure, packaging with PDM, and metadata format guidelines.

## Directory Structure

```
project/
  <package_name>/           source code
    __init__.py
    __main__.py             enables `python -m <package_name>`
  examples/                 example scripts
  tests/                    external tests (pytest)
    conftest.py
  .editorconfig
  Dockerfile
  MANIFEST.in               additional files for source distributions
  pyproject.toml            PEP 518 + PEP 621 metadata (single source of truth)
  README.md
```

### Key Directories and Files

| Path                   | Purpose                                                            |
|------------------------|--------------------------------------------------------------------|
| `/<package_name>/`     | Main source code. Replace with the actual package name (e.g., `icarus`). Python names cannot contain dashes, so `ion-cannon` becomes `ion_cannon`. |
| `/examples/`           | Example scripts demonstrating how to use the library or application. |
| `/tests/`              | External tests and test data. Use [pytest](https://docs.pytest.org/) to run all tests. |
| `/Dockerfile`          | Container definition. For multiple variants: `Dockerfile.alpine`, `Dockerfile.ubuntu`, `Dockerfile.slim-debian`. |
| `/MANIFEST.in`         | Specifies additional files to include in source distributions (sdist). |
| `/pyproject.toml`      | PEP 518 + PEP 621-compliant metadata. Single source of truth for project configuration. |
| `/README.md`           | Project description and usage instructions. Acceptable formats: Markdown, reStructuredText, or plain text. |

## Package Manager: PDM

**Use [PDM](https://pdm-project.org/) as the package manager.** PDM supports PEP 621 natively, is lightweight (unlike Poetry, it does not require GCC-compiled components), and can be replaced by any PEP 621-compatible tool.

### Tool Comparison

| Tool       | Description                                                                      |
|------------|----------------------------------------------------------------------------------|
| Conda      | Cross-platform binary package manager (primarily for Anaconda)                   |
| Setuptools | Fully-featured, legacy package builder                                           |
| Pipenv     | Virtualenv and dependency manager                                                |
| Poetry     | Dependency manager + builder + publisher (uses custom `pyproject.toml` sections) |
| **PDM**    | Dependency manager + builder + publisher with full PEP 621 support               |

### Installation

```bash
pip install -U pdm
```

### Initialization

```bash
mkdir icarus && cd icarus
pdm init
```

Answer the interactive prompts according to your package's requirements. Use `Proprietary` as the license for non-open-source packages.

## Dependency Management

```bash
# Add a dependency
pdm add 'opencv-python>=4.5'

# Add a development dependency
pdm add -dG test pytest pytest-cov

# Install all dependencies
pdm install

# Update all dependencies
pdm update

# Remove a dependency
pdm remove opencv-python
```

## Build & Distribution

```bash
# Install the package locally
pip install .

# Build source distribution
python3 -m build -s .

# Build wheel distribution
python3 -m build -w .
```

### Docker

```bash
# Build a specific Dockerfile variant
docker build -f Dockerfile.alpine -t app:1.0.0-alpine .
```

## Key Configuration Files

### pyproject.toml (PEP 621)

**PEP 621 is the recommended standard** because it consolidates all metadata into a single, readable file. Avoid legacy `setup.py` and `setup.cfg` formats.

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

Enable dynamic versioning so the version is defined once in `__init__.py`:

```python
# icarus/__init__.py
__version__ = "1.0.0"
```

### Package Entry Points

Create the `__main__` file to enable `python -m icarus`:

```python
# icarus/__main__.py
import sys
from .icarus import main

if __name__ == "__main__":
    sys.exit(main())
```

### MANIFEST.in

Specifies additional files to include in source distributions:

```text
recursive-include models *
include somedir/README.md
```

### Metadata Format Evolution

For historical context, Python package metadata has evolved through several formats:

| Format       | Standard  | Description                                                  |
|--------------|-----------|--------------------------------------------------------------|
| `setup.py`   | Legacy    | Python script executed during build; not declarative         |
| `setup.cfg`  | PEP 517   | Declarative config, but separate from build system config    |
| `pyproject.toml` | PEP 621 + PEP 518 | Consolidated, readable, single-file metadata (recommended) |

## References

- [PDM Documentation](https://pdm-project.org/)
- [PEP 518 -- Build System Requirements](https://peps.python.org/pep-0518/)
- [PEP 621 -- Storing Project Metadata in pyproject.toml](https://peps.python.org/pep-0621/)
- [Python Packaging User Guide](https://packaging.python.org/)
- [pytest Documentation](https://docs.pytest.org/)
- [Dockerfile Best Practices](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)
