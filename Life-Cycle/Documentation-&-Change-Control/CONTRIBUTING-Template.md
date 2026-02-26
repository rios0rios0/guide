# CONTRIBUTING Template

> **TL;DR:** Copy this template when creating a new project. Replace all `{PLACEHOLDERS}` with project-specific values. Remove optional sections that do not apply. Every project must have a CONTRIBUTING.md that follows this structure.

## Overview

This template defines the standard structure for all project CONTRIBUTING files. It ensures every repository provides clear, consistent onboarding for contributors. Copy the raw template below, replace placeholders, and remove optional sections wrapped in HTML comments.

## Placeholder Reference

| Placeholder         | Description                          | Example                                                                   |
|---------------------|--------------------------------------|---------------------------------------------------------------------------|
| `{LANGUAGE}`        | Primary programming language         | `Go 1.23`, `Java 21`, `Python 3.12`, `Node.js 20`                         |
| `{INSTALL_COMMAND}` | Dependency install command           | `go mod download`, `./gradlew dependencies`, `pdm install`, `npm install` |
| `{EXTENSION_TYPE}`  | Plugin/provider type (if applicable) | `Provider`, `Updater`, `Plugin`                                           |

## Language-Specific Prerequisites

Use the prerequisite block that matches your project:

| Language                  | Prerequisites                                                                          |
|---------------------------|----------------------------------------------------------------------------------------|
| **Go**                    | Go 1.23+, Make                                                                         |
| **Java**                  | Java 21+ (Eclipse Temurin), Gradle (via wrapper), Make, Docker (for integration tests) |
| **Python**                | Python 3.12+, PDM, Make                                                                |
| **JavaScript/TypeScript** | Node.js 20+, npm, Make                                                                 |

## Template

````markdown
# Contributing

Contributions are welcome. By participating, you agree to maintain a respectful and constructive environment.

For coding standards, testing patterns, architecture guidelines, commit conventions, and all
development practices, refer to the **[Development Guide](https://github.com/rios0rios0/guide/wiki)**.

## Prerequisites

- {LANGUAGE}
- [Make](https://www.gnu.org/software/make/)
<!-- Add any other tools required by your project -->
<!-- Java projects: -->
<!-- - Docker (for integration tests with TestContainers) -->
<!-- Python projects: -->
<!-- - [PDM](https://pdm-project.org/) -->

## Development Workflow

1. Fork and clone the repository
2. Create a branch: `git checkout -b feat/my-change`
3. Install dependencies:
   ```bash
   {INSTALL_COMMAND}
   ```
4. Make your changes
5. Validate:
   ```bash
   make lint
   make test
   make sast
   ```
6. Update `CHANGELOG.md` under `[Unreleased]`
7. Commit following the [commit conventions](https://github.com/rios0rios0/guide/wiki/Git-Flow)
8. Open a pull request against `main`

<!-- OPTIONAL: only when the project requires environment variables or local services -->
<!--
## Local Environment

Copy `.env.example` to `.env` and fill in the required values:

```bash
cp .env.example .env
```

Start local services:

```bash
docker compose -f compose.dev.yaml up -d
```

| Variable | Description | Required |
|----------|-------------|----------|
| `DB_HOST` | Database hostname | Yes |
| `DB_PORT` | Database port | Yes |
-->

<!-- OPTIONAL: only for projects where contributors extend functionality -->
<!--
## Adding a New {EXTENSION_TYPE}

1. Create the implementation file following the [naming conventions](https://github.com/rios0rios0/guide/wiki/Code-Style)
2. Implement the required interface/contract
3. Register it in the dependency injection wiring
4. Add tests following the [testing guide](https://github.com/rios0rios0/guide/wiki/Tests)
5. Update `CHANGELOG.md` with an entry under `[Unreleased] > Added`
-->
````

## References

- [Documentation & Change Control](../Documentation-&-Change-Control.md) -- when and how to update documentation files
- [Git Flow](../Git-Flow.md) -- branch naming and commit message conventions
- [Tests](../Tests.md) -- cross-language testing standards
