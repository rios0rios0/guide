# Language Guide Template

> This document defines the standard structure that every language-specific guide must follow. Use it as a checklist when adding a new language or restructuring an existing one.

## File Structure

Every language guide under `Code-Style/` must consist of exactly **7 files**: one index page and six sub-pages. The naming convention uses the language name as a prefix for all sub-pages.

```
Code-Style/<Language>.md
Code-Style/<Language>/
    <Language>-Conventions.md
    <Language>-Formatting-and-Linting.md
    <Language>-Type-System.md
    <Language>-Logging.md
    <Language>-Testing.md
    <Language>-Project-Structure.md
```

## Page Specifications

### 1. `<Language>.md` -- Index + Philosophy

The entry point for the language guide. Must contain:

- **TL;DR** -- A single sentence summarizing the key tools and decisions for this language.
- **Overview** -- Brief introduction and links to all sub-pages.
- **Philosophy** -- The language's guiding principles (e.g., Zen of Python, Go Proverbs, Rust's zero-cost abstractions).
- **References** -- Links to official documentation and community resources.

### 2. `<Language>-Conventions.md` -- Naming & Component Patterns

Defines how to name things and implement the architecture patterns in this language. Must contain:

- **File Naming** -- The casing convention for file names (e.g., `snake_case`, `kebab-case`).
- **General Conventions** -- Language-specific idioms, syntax preferences, and code style rules that go beyond what the formatter enforces.
- **Entities** -- Rules for defining domain entities in this language.
- **Commands** -- Naming pattern for command files, classes/structs, and methods. Must follow `<Operation><Entity>Command` with an `Execute` method.
- **Controllers** -- Naming pattern for controller files, classes/structs, and methods. Must follow `<Operation><Entity>Controller` with an `Execute` method.
- **Services** -- Whether this layer is used, and if so, its naming pattern.
- **Repositories** -- Contract (domain layer) and implementation (infrastructure layer) naming patterns.
- **Mappers** -- Repository mappers and controller mappers (request/response) naming patterns.
- **Models** -- Infrastructure-layer DTO naming conventions.
- **Dependency Injection** -- The recommended DI tool or approach for this language.
- **References**

### 3. `<Language>-Formatting-and-Linting.md` -- Formatting & Linting Tools

Defines the mandatory toolchain for code formatting and static analysis. Must contain:

- **Formatter** -- The mandatory code formatter (e.g., Black, gofmt, Prettier, rustfmt) and why it was chosen.
- **Linter(s)** -- The mandatory linter(s) (e.g., Flake8, golangci-lint, ESLint, clippy) and their configuration.
- **Import Ordering** -- Rules for organizing imports (e.g., isort, goimports).
- **Editor Configuration** -- Recommended settings for common editors (VS Code, NeoVim).
- **References**

### 4. `<Language>-Type-System.md` -- Type Safety

Defines how to leverage the language's type system for safety and clarity. Must contain:

- **Overview** -- Whether the language is statically or dynamically typed, and the team's approach to type safety.
- **Type Annotation Requirements** -- What must be annotated (function signatures, variables, return types).
- **Key Type Patterns** -- Generics, union types, interfaces, traits, or other type constructs relevant to the language.
- **Prohibited Patterns** -- Type-unsafe patterns to avoid (e.g., `any` in TypeScript, missing type hints in Python).
- **References**

### 5. `<Language>-Logging.md` -- Logging

Defines the mandatory logging library and patterns. Must contain:

- **Mandatory Library** -- The required logging library (e.g., Loguru, Logrus, tracing).
- **Installation** -- How to install the library.
- **Import Convention** -- The standard import alias or pattern.
- **Log Levels** -- A table of available log levels with descriptions and usage guidance.
- **Structured Logging** -- How to attach contextual fields to log entries.
- **Prohibited Patterns** -- What must not be used (standard library logger, print statements, string interpolation in logs).
- **References**

### 6. `<Language>-Testing.md` -- Testing

Defines the testing framework, structure, and conventions. Must contain:

- **Testing Framework** -- The required framework and assertion library.
- **BDD Structure** -- The mandatory `given` / `when` / `then` comment block pattern with examples.
- **File Placement** -- Where test files live relative to production code.
- **File Naming** -- The naming convention for test files.
- **Unit Tests** -- Structure, parallelism requirements, and examples for command and controller tests.
- **Integration Tests** -- Suite structure, setup/teardown lifecycle, and examples for repository tests.
- **Test Doubles** -- Which doubles to use (stubs, dummies, fakers, in-memory) and when. Mocks should be avoided when possible.
- **Builders** -- The builder pattern for constructing test data.
- **References**

### 7. `<Language>-Project-Structure.md` -- Directory Layout & Packaging

Defines the standard project directory tree and dependency management. Must contain:

- **Directory Structure** -- The complete directory tree with explanations for each directory.
- **Package Manager** -- The recommended package manager and its configuration file.
- **Dependency Management** -- How to add, update, and lock dependencies.
- **Build & Distribution** -- How to build and distribute the package.
- **Key Configuration Files** -- Description of each configuration file in the project root.
- **References**

## Checklist for Adding a New Language

1. Create `Code-Style/<Language>.md` with TL;DR, overview, philosophy, and links to sub-pages.
2. Create `Code-Style/<Language>/<Language>-Conventions.md` with naming patterns for all architectural components.
3. Create `Code-Style/<Language>/<Language>-Formatting-and-Linting.md` with the mandatory formatter and linter.
4. Create `Code-Style/<Language>/<Language>-Type-System.md` with type safety guidelines.
5. Create `Code-Style/<Language>/<Language>-Logging.md` with the mandatory logging library.
6. Create `Code-Style/<Language>/<Language>-Testing.md` with framework, BDD structure, and examples.
7. Create `Code-Style/<Language>/<Language>-Project-Structure.md` with directory layout and packaging.
8. Add the language and all sub-pages to `Home.md` under the Code Style section.
9. Verify all cross-references link correctly.

## References

- [Code Style](../Code-Style.md) -- Baseline conventions shared across all languages.
- [Backend Design](../Life-Cycle/Architecture/Backend-Design.md) -- Architectural layers referenced by the conventions page.
- [Tests](../Life-Cycle/Tests.md) -- Cross-language testing standards (BDD, test doubles, builders).
