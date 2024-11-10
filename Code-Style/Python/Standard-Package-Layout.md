# Standard Package Layout
This page provides guidelines on where to place files that aren’t the module’s code. Python doesn’t have an official standard for package repository layout, but there are some generally accepted conventions.
Below are common directory/file names and their purposes. If unsure about naming a directory/file not listed here, consider referring to [Golang’s project layout standard](https://github.com/golang-standards/project-layout). Directories end with a `/`.

## Table of Contents
- [`/[package_name]/`](#package_name)
- [`/examples/`](#examples)
- [`/tests/`](#tests)
- [`/Dockerfile`](#dockerfile)
- [`/MANIFEST.in`](#manifestin)
- [`/pyproject.toml`](#pyprojecttoml)
- [`/README.md`](#readmemd)

## `/[package_name]/`
The module’s main code. Replace `[package_name]` with the name of the package, e.g., if the package name is `icarus`, name this directory `icarus`. Python file/directory names can’t contain dashes (`-`), so `ion-cannon` should be `ion_cannon`.

## `/examples/`
Example codes demonstrating how to run the code or call the library.

## `/tests/`
External tests and test data. Consider using [pytest](https://docs.pytest.org/en/7.0.x/) to execute all tests with a single `pytest` command for detailed test results. Example content for `tests/test_api.py`:
```python
import requests

def test_api():
    response = requests.get("https://some.api.local/health")
    assert response.status_code == 200
```

Run all tests in the `tests` folder with `pytest`.

## `/Dockerfile`
The project’s Dockerfile. For multiple Dockerfiles (e.g., for different architectures), name them as follows:
- `/Dockerfile.alpine`
- `/Dockerfile.ubuntu`
- `/Dockerfile.slim-debian`

Use `docker build -f Dockerfile.alpine -t app:1.0.0-alpine` to specify which Dockerfile to build from.

## `/MANIFEST.in`
Specifies additional data to include in the sdist package. Example content:
```text
recursive-include models *
include somedir/README.md
```

## `/pyproject.toml`
PEP 518 and PEP 621-compliant package metadata.

## `/README.md`
The project’s README file. It can be in Markdown (`README.md`), reStructuredText (`README.rst`), or plain text (`README` or `README.txt`). This file should describe the program’s purpose and usage instructions.
Example content:
```markdown
# Project Name

A brief description.

## Usage
Install and use the program as follows:
```bash
pip install program
python -m program
```
