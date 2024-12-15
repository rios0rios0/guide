# Styling and Formatting (PEP 8)

**TL;DR: Use Black, isort, and Flake8.**

## Table of Contents
- [YAPF and Black](#StylingandFormatting\(PEP8\)-YAPFandBlack)
- [isort](#StylingandFormatting\(PEP8\)-isort)
- [Visual Studio Code](#StylingandFormatting\(PEP8\)-VisualStudio)
- [NeoVim](#StylingandFormatting\(PEP8\)-NeoVim)

Different developers tend to have different habits when writing code – how they name variables, choose whether to use tabs or spaces for indentation, etc.
If unregulated, it may be difficult for developers to read each other's code. Therefore, Python's PEP 8 standard defines a wide variety of rules on how to style your Python code so people's coding styles do not vary too much from each other.

For instance, it defines that an in-line comment should have exactly two spaces preceding the `#` and exactly one space between the `#` and the comment's content.
It also defines that you should import different modules on separate lines instead of the same line.

You do not need to fully memorize the PEP 8 standards. Code editors and IDEs often have extensions or provide internal support for code linting and formatting.
The most commonly used tools for automatically checking and highlighting (linting) PEP 8-noncompliance are Flake8 and Pylint. `autopep8` can then automatically format (fix) these noncompliance.

## YAPF and Black

Although PEP 8 already covers a lot of details on how the code should be written, it still doesn't cover enough for all code compliant with PEP 8 to share the same style.
YAPF and Black are PEP 8-compliant formatters with a lot more detailed custom formatting rules. YAPF is created by Google, comes with several base styles (pep8, google, yapf, and facebook), and is highly customizable.
Black's standards are much stricter and it is much less customizable.

Less customizability isn't a bad thing. It means that all code formatted by Black will look very similar.
This guide recommends users to use Black since once you get used to Black's style, you can easily read all other code that's formatted with Black.
Even if a piece of code isn't already formatted with Black, you can quickly format it to a style that you're familiar with. Black is also much easier to set up since there are no knobs to tune.

To give you a taste of how these formatters format code differently, here is how YAPF (`--style=google`) formats a long function call:

```python
records = get_records(args.key, feed, args.output, args.range,
                      args.field, args.extra, args.expiration)
```

Here is how Black formats the same function call:

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

## isort

PEP 8 has some definitions on [how Python modules should be organized](https://www.python.org/dev/peps/pep-0008/#id19), but as described in the previous section, the definitions aren't comprehensive.
In practice, organizations like the CIA has developed their own stricter standards on [how to format the imports](https://wikileaks.org/ciav7p1/cms/page_26607631.html) (warning: it is illegal for you to view this link if you're affiliated with the CIA since it's a leaked document).

`isort` is like YAPF and Black, but instead of sorting code, it only sorts imports, which YAPF and Black don't care enough about.
It splits imports into three sections (standard library, third-party, and application-specific), sorts imports lexicographically and case-insensitively, then moves from imports after normal imports.
It also supports formatting imports in the Black style.

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

Below is a summarization of all the formatters we've talked about. **Bold** items are recommended for you to use during your development.

| Name      | Description                                                                                                                                                                                                  |
|-----------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| autopep8  | `autopep8` checks the file against the official PEP 8 standard and automatically tries to fix all non-compliant code                                                                                         |
| YAPF      | YAPF is a formatter developed by Google. It allows the user to format the code with a stricter standard and is highly configurable. Officially included base standards are pep8, google, yapf, and facebook. |
| **Black** | Black is a formatter with a very strict style and very few "tweakable" options.                                                                                                                              |
| **isort** | `isort` automatically sorts and formats all imports, then splits them into sections per PEP 8, and optionally Black standards.                                                                               |

## Visual Studio Code

If you use VSCode for developing Python, here are some useful extensions and their recommended configurations:

- [Python](https://marketplace.visualstudio.com/items?itemName=ms-python.python): Basic Python support like auto-completes, includes [Pylance](https://marketplace.visualstudio.com/items?itemName=ms-python.vscode-pylance)
  ```json
  {
      "python.formatting.provider": "black",
      "python.languageServer": "Pylance",
      "python.linting.flake8Enabled": true,
      "python.analysis.typeCheckingMode": "strict"
  }
  ```
- [Python Indent](https://marketplace.visualstudio.com/items?itemName=KevinRose.vsc-python-indent): Helps to indent Python code properly
- [autoDocstring](https://marketplace.visualstudio.com/items?itemName=njpwerner.autodocstring): Helps to generate doc-strings for functions and classes

## NeoVim

If you use NeoVim (great choice) for developing Python, here are some useful plugins and their recommended configurations:

- [nvie/vim-flake8](https://github.com/nvie/vim-flake8): Flake8 Python code linting
- [psf/black](https://github.com/psf/black): Black code formatter integration
- [neoclide/coc.nvim](https://github.com/neoclide/coc.nvim): NeoVim extension loader
  - Use `:CoCInstall coc-pyright` to install the Python language server. Pyright is the OSS equivalent to Pyright in VSCode.
- [dense-analysis/ale](https://github.com/dense-analysis/ale#installation-with-vim-plug): A linting engine for NeoVim
  - Let ALE complete Python files with Black
    ```vim
    let g:ale_fixers = {
    \ '*': ['trim_whitespace', 'prettier'],
    \ 'python': ['black']
    \}
    let g:ale_completion_enabled = 1
    ```
- [kkoomen/vim-doge](https://github.com/kkoomen/vim-doge): Documentation generation tool
  - Bind Alt+D to `:DogeGenerate`: `nnoremap <M-d> :DogeGenerate<CR>`
