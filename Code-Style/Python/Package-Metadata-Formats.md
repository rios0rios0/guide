# Package Metadata Formats
Python package metadata can be written in a few formats, primarily in `setup.py`, `setup.cfg`, and `pyproject.toml`.
This page will introduce the differences between the different formats and how to write them. **TL;DR: use PEP 621.**

## Table of Contents
- [setup.py](#PackageMetadataFormats-setup.py)
- [PEP 517](#PackageMetadataFormats-PEP517)
- [PEP 518](#PackageMetadataFormats-PEP518)
- [PEP 621](#PackageMetadataFormats-PEP621)

All the formats/standards listed below can be recognized by Pip.
While the wheels built with all three standards won't be that different from each other, PEP 621 is generally preferred since it is the least complex to read/write/manage.

## [setup.py](https://pip.pypa.io/en/latest/reference/build-system/setup-py/)
`setup.py` was the only metadata definition format Pip recognized before PEP 517 and PEP 518 were introduced.
It is essentially a Python script that is executed then the package needs to be built/installed.
You should avoid using this legacy format since it is only readable by `Setuptools`, and it's not written in a declarative manner.
Below are some useful resources if you have to use it:

- <https://docs.python.org/3/distutils/setupscript.html>
- <https://setuptools.pypa.io/en/latest/userguide/keywords.html>

```python
#!/usr/bin/python3
from distutils.core import setup

setup(
    name="icarus",
    version="1.0.0",
    description="A demo package",
    author="Leon Lin",
    author_email="i@k4yt3x.com",
    url="https://gitlab.com/demo-group/demo-subgroup/icarus",
    packages=["icarus"],
)
```

## [PEP 517](https://www.python.org/dev/peps/pep-0517/)
PEP 517 allows users to write the project's metadata to a build-system independent CFG file.
You can refer to these pages for advanced usages of the `setup.cfg` file:

* <https://setuptools.pypa.io/en/latest/userguide/declarative_config.html>
* <https://setuptools.pypa.io/en/latest/userguide/keywords.html>

```ini
[metadata]
name = icarus
version = 1.0.0
author = Leon Lin
author_email = i@k4yt3x.com
license = Proprietary
description = A demo package
url = "https://gitlab.com/demo-group/demo-subgroup/icarus"
long_description = file: README.md
long_description_content_type = text/markdown

[options]
packages = find:
install_requires =
    opencv-python>=4.5
python_requires = >=3.6
```

## [PEP 518](https://www.python.org/dev/peps/pep-0518/)
PEP 518 allows users to define the minimum build system requirements for a package. I.e., what packages are required to build this package.
These build dependencies are not required/installed for runtime. The build requirements can only be defined in `pyproject.toml`.
Therefore, it is less cumbersome to use PEP 621 + PEP 518 than PEP 517 + PEP 518 since the former option requires only one `pyproject.toml` file instead of two files.
Here are some useful resources on how to write the `pyproject.toml` file for PEP 518:

* <https://www.python.org/dev/peps/pep-0518/>
* <https://pip.pypa.io/en/stable/reference/build-system/pyproject-toml/>

```toml
[build-system]
requires = ["pdm-pep517"]
build-backend = "pdm.pep517.api"
```

## [PEP 621](https://www.python.org/dev/peps/pep-0621/)
PEP 621 allows all the project's metadata to be written into a single `pyproject.toml` file.
Here are some resources you can look at for writing advanced `pyproject.toml` files:

* <https://www.python.org/dev/peps/pep-0621/>
* <https://pdm.fming.dev/pyproject/pep621/>
* <https://python-poetry.org/docs/pyproject/>

```toml
[project]
name = "icarus"
description = "A demo package"
authors = [{ name = "Leon Lin", email = "i@k4yt3x.com" }]
dependencies = []
requires-python = ">=3.10"
license = { text = "Proprietary" }
dynamic = ["version"]

[project.urls]
Homepage = "https://gitlab.com/demo-group/demo-subgroup/icarus"

[tool.pdm]
version = { from = "icarus/__init__.py" }

[build-system]
requires = ["pdm-pep517"]
build-backend = "pdm.pep517.api"
```
