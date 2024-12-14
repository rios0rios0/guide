# Python Packaging
This document outlines the best Python development practices and a step-by-step tutorial on how to create a new Python package. **TL;DR: Use Pip + PDM.**

## Table of Contents
- [Package Initialization with PDM](#PythonPackaging-PackageInitializationwi)
  - [Initialize Package](#PythonPackaging-InitializePackage)
  - [Create Package Base Code](#PythonPackaging-CreatePackageBaseCode)
  - [Update Package Metadata and README](#PythonPackaging-UpdatePackageMetadataan)
  - [Next Steps](#PythonPackaging-NextSteps)

There are a lot of Python packaging and distribution tools available. Below are some common ones:

| Name       | Description                                                                            |
|------------|----------------------------------------------------------------------------------------|
| Conda      | A cross-platform binary package manager mainly used by Anaconda                        |
| Setuptools | A fully-featured and popular package builder                                           |
| Pipenv     | A virtualenv and dependency manager                                                    |
| Poetry     | A dependency manager, package builder, and package publisher                           |
| PDM        | A dependency manager, package builder, and package publisher with full PEP 621 support |

The tools listed above are merely a fraction of all the tools available.
A more comprehensive comparison between all Python packaging and distribution tools can be found [here](https://chadsmith.dev/python-packaging/).
In practice, you can choose any tool you like to develop your Python package as long as Pip can recognize the format and install the package.

This guide **recommends users to use PDM** due to the following reasons:
- PDM supports **PEP 621**, which allows all the project's metadata to be included in one `pyproject.toml` file instead of spread across `setup.py`, `setup.cfg`, and `pyproject.toml`.
- PDM can easily be replaced in the future since **PEP 621** is an official standard. Poetry can also store all files in `pyproject.toml`, but it uses custom sections that cannot be recognized by other tools.
- PDM is relatively light-weight when compared to tools like Poetry, which requires certain components to be compiled by GCC.

## Package Initialization with PDM
This section will use PDM as the build tool to demonstrate how to create a standard Python package from scratch.

### Initialize Package
Install PDM using Pip:
```shell
pip install -U pdm
```

Create a directory for the package:
```shell
mkdir icarus
cd icarus
```

Initialize the package using PDM. Answer the questions according to your package's requirements.
You can use "Proprietary" as the license name for a non-open-source package.
```shell
$ pdm init
Creating a pyproject.toml for PDM...
Please enter the Python interpreter to use
0. /usr/bin/python3 (3.10)
1. /usr/bin/python (3.10)
Please select: [0]:
Using Python interpreter: /usr/bin/python3 (3.10)
Is the project a library that will be uploaded to PyPI? [y/N]:
License(SPDX name) [MIT]: Proprietary
Author name [leon.lin]: Leon Lin
Author email [i@k4yt3x.com]:
Python requires('*' to allow any) [>=3.10]: >=3.6
Changes are written to pyproject.toml.
```

### Create Package Base Code

Create the module's directory. This directory's name should be the same as the Python package's name.
```shell
mkdir icarus
```

Create the project initialization file:
```shell
cat <<EOF> icarus/__init__.py
#!/usr/bin/python3
# -*- coding: utf-8 -*-
__version__ = "1.0.0"
from .icarus import main
EOF
```

Create a package `__main__` file. This file will be executed when the module is executed directly like `python -m icarus`.
```shell
cat <<EOF> icarus/__main__.py
#!/usr/bin/python3
# -*- coding: utf-8 -*-
# local imports
from .icarus import main
# built-in imports
import sys
if __name__ == "__main__":
 sys.exit(main())
EOF
```

Create the module's entrypoint file:
```shell
cat <<EOF> icarus/icarus.py
#!/usr/bin/python3
# -*- coding: utf-8 -*-
from . import __version__
def main():
 print("Don't fly too close to the sun.")
 print(f"The package's version is: {__version__}")
EOF
```

Change the Python files' mode to executable:
```shell
chmod +x icarus/*.py
```

### Update Package Metadata and README

Edit the `pyproject.toml` file to instruct it to read the package's version from the package's `__init__` file dynamically.
This way the version only needs to be specified once in the package but is readable by both the packaging tool and the module's main code.
Also update the other relevant fields in the file like the package's name and description.
```diff
--- pyproject.toml.orig 2022-02-28 19:33:29.705268473 +0000
+++ pyproject.toml 2022-02-28 19:24:41.563891539 +0000
@@ -1,16 +1,17 @@
 [project]
-name = ""
-version = ""
-description = ""
+name = "icarus"
+description = "A demo package"
 authors = [{ name = "Leon Lin", email = "i@k4yt3x.com" }]
 dependencies = []
 requires-python = ">=3.10"
 license = { text = "Proprietary" }
+dynamic = ["version"]
 [project.urls]
-Homepage = ""
+Homepage = "https://gitlab.com/demo-group/demo-subgroup/icarus"
 [tool.pdm]
+version = { from = "icarus/__init__.py" }
 [build-system]
 requires = ["pdm-pep517"]
```

Create a `README.md` file for this repository:
```markdown
# Icarus
This is a demo package for our Python standardization document.

## Usages

Install the package:

@```shell
pip install .
@```
```

Run the package:
```bash
python3 -m icarus
```

### Next Steps
When you are done, your repository's structure should look something like this:
```
icarus
├── icarus
│   ├── __init__.py
│   ├── __main__.py
│   └── icarus.py
├── pyproject.toml
└── README.md
```

Your package is now initialized. You can start writing your code.
Below are some useful commands. For more information, refer to [PDM's documentation](https://pdm.fming.dev/) and [Pip's documentation](https://pip.pypa.io/en/stable/).
```bash
# add a new dependency
pdm add 'opencv-python>=4.5'

# install the package
pip install .

# build the package's source distribution (sdist)
pip install -U build
python3 -m build -s .

# build the package's binary distribution (bdist wheel)
pip install -U build wheel
python3 -m build -w .
```

Now you know how to build a simple Python package.
Here are some more articles helpful for you to better understand Python packaging:
- [Package-Metadata-Formats](Package-Metadata-Formats.md)
- [Standard Package Layout](Standard-Package-Layout.md)
