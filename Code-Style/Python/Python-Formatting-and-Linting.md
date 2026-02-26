# Python Formatting and Linting

> **TL;DR:** Use **Black** for code formatting, **isort** for import sorting, and **Flake8** for linting. This combination ensures consistent, readable code across all Python projects.

## Overview

Different developers have different coding habits -- variable naming, indentation style, import ordering, and more. Python's [PEP 8](https://peps.python.org/pep-0008/) standard defines a comprehensive set of rules for consistent code styling. You do not need to memorize PEP 8; tooling handles enforcement automatically.

## Recommended Toolchain

| Tool       | Purpose                                          | Customizability                                        |
|------------|--------------------------------------------------|--------------------------------------------------------|
| **Black**  | PEP 8-compliant code formatter                   | Minimal (by design -- ensures all code looks the same) |
| **isort**  | Import sorter (stdlib, third-party, application) | Moderate (supports Black-compatible formatting)        |
| **Flake8** | Linter for PEP 8 compliance and code quality     | High                                                   |
| YAPF       | PEP 8-compliant formatter by Google              | Very high (multiple base styles)                       |
| autopep8   | PEP 8 auto-fixer                                 | Moderate                                               |

### Why Black?

Black's strict, opinionated formatting means all code formatted by Black looks identical. This consistency reduces cognitive overhead during code reviews and makes unfamiliar codebases immediately readable. Black is also trivial to set up -- there are no configuration knobs to tune.

**YAPF** (Google-style formatting):
```python
records = get_records(args.key, feed, args.output, args.range,
                      args.field, args.extra, args.expiration)
```

**Black** formatting:
```python
records = get_records(
    args.key,
    feed,
    args.output,
    args.range,
    args.field,
    args.extra,
    args.expiration,
)
```

### Import Ordering: isort

isort organizes imports into three sections (standard library, third-party, and application-specific), sorts them lexicographically and case-insensitively, and supports Black-compatible formatting:

```python
import argparse
import math
import multiprocessing
from typing import Type, Union

import cv2
import ffmpeg
from loguru import logger
from rich.progress import (
    BarColumn,
    Progress,
    ProgressColumn,
)
from rich.text import Text

from . import __version__
from .decoder import VideoDecoder
from .encoder import VideoEncoder
```

## Editor Configuration

### Visual Studio Code

Recommended extensions and settings:

- [Python](https://marketplace.visualstudio.com/items?itemName=ms-python.python) (includes [Pylance](https://marketplace.visualstudio.com/items?itemName=ms-python.vscode-pylance))
- [Python Indent](https://marketplace.visualstudio.com/items?itemName=KevinRose.vsc-python-indent)
- [autoDocstring](https://marketplace.visualstudio.com/items?itemName=njpwerner.autodocstring)

```json
{
    "python.formatting.provider": "black",
    "python.languageServer": "Pylance",
    "python.linting.flake8Enabled": true,
    "python.analysis.typeCheckingMode": "strict"
}
```

### NeoVim

Recommended plugins:

- [nvie/vim-flake8](https://github.com/nvie/vim-flake8) -- Flake8 linting
- [psf/black](https://github.com/psf/black) -- Black formatter integration
- [neoclide/coc.nvim](https://github.com/neoclide/coc.nvim) -- Extension loader (use `:CoCInstall coc-pyright` for Python support)
- [dense-analysis/ale](https://github.com/dense-analysis/ale) -- Asynchronous linting engine
- [kkoomen/vim-doge](https://github.com/kkoomen/vim-doge) -- Documentation generation

ALE configuration for Black:
```vim
let g:ale_fixers = {
\ '*': ['trim_whitespace', 'prettier'],
\ 'python': ['black']
\}
let g:ale_completion_enabled = 1
```

## References

- [PEP 8 -- Style Guide for Python Code](https://peps.python.org/pep-0008/)
- [Black Documentation](https://black.readthedocs.io/)
- [isort Documentation](https://pycqa.github.io/isort/)
- [Flake8 Documentation](https://flake8.pycqa.org/)
