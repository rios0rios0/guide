# Standard Package Layout

> **TL;DR:** Follow the conventional directory structure: `/<package_name>/` for source code, `/tests/` for tests (pytest), `/examples/` for usage examples, `pyproject.toml` for metadata, and `Dockerfile` for containerization.

## Overview

Python does not enforce an official package repository layout, but the community has established widely accepted conventions. This page describes the standard directory structure and the purpose of each file. For items not listed here, consider referencing the [Go project layout standard](https://github.com/golang-standards/project-layout) as a general guide.

## Directory Structure

### `/[package_name]/`

The module's main source code. Replace `[package_name]` with the actual package name (e.g., `icarus`). Python names cannot contain dashes, so `ion-cannon` becomes `ion_cannon`.

### `/examples/`

Example scripts demonstrating how to use the library or run the application.

### `/tests/`

External tests and test data. Use [pytest](https://docs.pytest.org/) to execute all tests with a single command:

```python
# tests/test_api.py
import requests

def test_api():
    response = requests.get("https://some.api.local/health")
    assert response.status_code == 200
```

Run all tests:
```bash
pytest
```

### `/Dockerfile`

Container definition for the project. For multiple Dockerfiles targeting different base images:

- `Dockerfile.alpine`
- `Dockerfile.ubuntu`
- `Dockerfile.slim-debian`

Build a specific variant:
```bash
docker build -f Dockerfile.alpine -t app:1.0.0-alpine .
```

### `/MANIFEST.in`

Specifies additional files to include in source distributions (sdist):

```text
recursive-include models *
include somedir/README.md
```

### `/pyproject.toml`

PEP 518 and PEP 621-compliant package metadata. This is the single source of truth for project configuration.

### `/README.md`

Project description and usage instructions. Acceptable formats: Markdown (`.md`), reStructuredText (`.rst`), or plain text.

## References

- [Python Packaging User Guide](https://packaging.python.org/)
- [pytest Documentation](https://docs.pytest.org/)
- [Dockerfile Best Practices](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)
