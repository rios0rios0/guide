# Go

> **TL;DR:** Use `snake_case` for file names, a short abbreviation of the type as the method receiver (e.g., `c` for `Client`), [Dig](https://github.com/uber-go/dig) for dependency injection, [golangci-lint](https://golangci-lint.run/) for linting, [Logrus](https://github.com/sirupsen/logrus) for logging, and [testify](https://github.com/stretchr/testify) for testing. Entities must be framework-agnostic.

## Overview

This series of pages outlines Go-specific conventions, with detailed explanations of recommended approaches and the reasoning behind them. For the general baseline, refer to the [Code Style](../Code-Style.md) guide. See the sub-pages for specific topics:

- [Conventions](GoLang/GoLang-Conventions.md)
- [Formatting and Linting](GoLang/GoLang-Formatting-and-Linting.md)
- [Type System](GoLang/GoLang-Type-System.md)
- [Logging](GoLang/GoLang-Logging.md)
- [Testing](GoLang/GoLang-Testing.md)
- [Project Structure](GoLang/GoLang-Project-Structure.md)

## Go Proverbs

The [Go Proverbs](https://go-proverbs.github.io/) capture the language's design philosophy:

- Don't communicate by sharing memory, share memory by communicating.
- Concurrency is not parallelism.
- The bigger the interface, the weaker the abstraction.
- Make the zero value useful.
- A little copying is better than a little dependency.
- Clear is better than clever.
- Errors are values.
- Don't just check errors, handle them gracefully.

## References

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Proverbs](https://go-proverbs.github.io/)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
