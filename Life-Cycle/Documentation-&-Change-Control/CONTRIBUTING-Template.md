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

A project that uses [chlog](https://github.com/luizjhonata/chlog) -- one with a `.chlog.yaml` at its
root -- additionally needs the tool to write changelog fragments. It is a single self-contained
binary, so the project itself never gains a Go dependency -- but `go install github.com/luizjhonata/chlog@latest`
builds it from source and therefore needs a Go toolchain locally. On a Java, Python, or JavaScript
project, point contributors at the prebuilt binary for their platform on the
[releases page](https://github.com/luizjhonata/chlog/releases) (Linux, macOS, and Windows; amd64 and
arm64) and keep the toolchain out of the prerequisites. A project that does not use chlog drops the
prerequisite entirely and edits `CHANGELOG.md` by hand -- see
[Documentation & Change Control](../Documentation-&-Change-Control.md).

## Template

````markdown
# Contributing

Contributions are welcome. By participating, you agree to maintain a respectful and constructive environment.

For coding standards, testing patterns, architecture guidelines, commit conventions, and all
development practices, refer to the **[Development Guide](https://github.com/rios0rios0/guide/wiki)**.

## Prerequisites

- {LANGUAGE}
- [Make](https://www.gnu.org/software/make/)
<!-- Projects that use chlog (a `.chlog.yaml` at the root) -- drop this line if the project edits CHANGELOG.md by hand: -->
- [chlog](https://github.com/luizjhonata/chlog) -- `go install github.com/luizjhonata/chlog@latest` (needs a Go toolchain), or a prebuilt binary from the [releases page](https://github.com/luizjhonata/chlog/releases)
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
6. Record the change in the changelog.
   <!-- Projects that use chlog -- never edit CHANGELOG.md, which is generated from the fragments: -->
   ```bash
   chlog new --kind Added --body "added the thing that was not there before"
   ```
   <!-- Projects without chlog -- add the bullet under `[Unreleased]` in `CHANGELOG.md` instead. -->
7. Commit following the [commit conventions](https://github.com/rios0rios0/guide/wiki/Git-Flow) --
   use the ticket ID as the commit scope
   ([exceptions](https://github.com/rios0rios0/guide/wiki/Git-Flow#ticket-reference))
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
5. Record the change in the changelog: `chlog new --kind Added --body "added ..."`, or a bullet under `[Unreleased] > Added` in `CHANGELOG.md` when the project does not use chlog
-->
````

## References

- [Documentation & Change Control](../Documentation-&-Change-Control.md) -- when and how to update documentation files
- [Git Flow](../Git-Flow.md) -- branch naming and commit message conventions
- [Tests](../Tests.md) -- cross-language testing standards
