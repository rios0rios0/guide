# Package Metadata Formats

> **TL;DR:** Use **PEP 621** to store all package metadata in a single `pyproject.toml` file. Avoid `setup.py` (legacy) and prefer PEP 621 + PEP 518 over PEP 517 + PEP 518 to minimize the number of configuration files.

## Overview

Python package metadata can be written in several formats: `setup.py`, `setup.cfg`, and `pyproject.toml`. All formats are recognized by Pip, and the resulting wheels are nearly identical. However, **PEP 621 is the recommended standard** because it consolidates all metadata into a single, readable `pyproject.toml` file.

## Format Comparison

### `setup.py` (Legacy)

The original metadata format -- a Python script executed during build/install. Avoid this format because it is only readable by Setuptools and is not declarative.

```python
from distutils.core import setup

setup(
    name="icarus",
    version="1.0.0",
    description="A demo package",
    author="Author",
    author_email="author@example.com",
    packages=["icarus"],
)
```

- [setup.py Reference](https://docs.python.org/3/distutils/setupscript.html)

### PEP 517 (`setup.cfg`)

Introduced build-system-independent declarative configuration via `setup.cfg`:

```ini
[metadata]
name = icarus
version = 1.0.0
author = Author
author_email = author@example.com
description = A demo package

[options]
packages = find:
install_requires =
    opencv-python>=4.5
python_requires = >=3.6
```

- [PEP 517](https://peps.python.org/pep-0517/)
- [Setuptools Declarative Config](https://setuptools.pypa.io/en/latest/userguide/declarative_config.html)

### PEP 518 (Build System Requirements)

Defines minimum build system requirements in `pyproject.toml`. Used in combination with either PEP 517 or PEP 621:

```toml
[build-system]
requires = ["pdm-pep517"]
build-backend = "pdm.pep517.api"
```

- [PEP 518](https://peps.python.org/pep-0518/)

### PEP 621 (Recommended)

Consolidates all metadata into a single `pyproject.toml` file -- the simplest and most modern approach:

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

- [PEP 621](https://peps.python.org/pep-0621/)
- [PDM pyproject.toml Reference](https://pdm-project.org/en/latest/reference/pyproject/)

## References

- [PEP 517](https://peps.python.org/pep-0517/)
- [PEP 518](https://peps.python.org/pep-0518/)
- [PEP 621](https://peps.python.org/pep-0621/)
- [Python Packaging User Guide](https://packaging.python.org/)
