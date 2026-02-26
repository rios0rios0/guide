# Python

> **TL;DR:** Follow the Zen of Python -- prioritize readability, simplicity, and explicitness. Use [Black](https://black.readthedocs.io/) for formatting, [isort](https://pycqa.github.io/isort/) for imports, [Flake8](https://flake8.pycqa.org/) for linting, type hints on all functions, [Loguru](https://loguru.readthedocs.io/) for logging, [pytest](https://docs.pytest.org/) for testing, and [PDM](https://pdm-project.org/) for packaging.

## Overview

This series of pages outlines best practices for Python development, with detailed explanations of recommended approaches and the reasoning behind them. For the general baseline, refer to the [Code Style](../Code-Style.md) guide. See the sub-pages for specific topics:

- [Conventions](Python/Python-Conventions.md)
- [Formatting and Linting](Python/Python-Formatting-and-Linting.md)
- [Type System](Python/Python-Type-System.md)
- [Logging](Python/Python-Logging.md)
- [Testing](Python/Python-Testing.md)
- [Project Structure](Python/Python-Project-Structure.md)

## The Zen of Python

The Zen of Python defines the core philosophy for writing Pythonic code. Access it anytime by running `import this` in a Python REPL:

```
Beautiful is better than ugly.
Explicit is better than implicit.
Simple is better than complex.
Complex is better than complicated.
Flat is better than nested.
Sparse is better than dense.
Readability counts.
Special cases aren't special enough to break the rules.
Although practicality beats purity.
Errors should never pass silently.
Unless explicitly silenced.
In the face of ambiguity, refuse the temptation to guess.
There should be one-- and preferably only one --obvious way to do it.
Although that way may not be obvious at first unless you're Dutch.
Now is better than never.
Although never is often better than *right* now.
If the implementation is hard to explain, it's a bad idea.
If the implementation is easy to explain, it may be a good idea.
Namespaces are one honking great idea -- let's do more of those!
```

## References

- [PEP 20 -- The Zen of Python](https://peps.python.org/pep-0020/)
- [Python Documentation](https://docs.python.org/3/)
