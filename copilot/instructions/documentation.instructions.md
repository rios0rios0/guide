# Documentation & Change Control

> **TL;DR:** Every project must maintain a `CHANGELOG.md` (following [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)) and a `README.md`. Every code change must be accompanied by the corresponding documentation update -- changelog **always**, README and other docs (e.g., `.github/copilot-instructions.md`) **whenever behavior, configuration, or architecture changes**.

## Overview

Documentation and change control are integral parts of the engineering workflow, not afterthoughts. A well-maintained changelog and README provide traceability for stakeholders, reduce onboarding friction for new team members, and ensure that the state of the project is always understandable from its documentation alone.

**Every change introduced to a project must include updates to the relevant documentation files.** This is enforced as part of the development workflow, not as a separate task.

## Required Documentation Files

Every project must contain at minimum:

| File                                  | Purpose                                                      | Update Frequency                                             |
|---------------------------------------|--------------------------------------------------------------|--------------------------------------------------------------|
| **`README.md`**                       | Describes the project, its usage, setup, and architecture    | When behavior, configuration, CLI, or setup changes          |
| **`CONTRIBUTING.md`**                 | Guides contributors on prerequisites, workflow, and standards | When prerequisites, workflow, or project structure changes    |
| **`CHANGELOG.md`**                    | Records all notable changes, organized by version            | **Every change** (always)                                    |
| **`.github/copilot-instructions.md`** | AI assistant context for the project structure and workflows | When architecture, commands, or development workflow changes |

**Templates are available for standardized project setup:**

## Changelog Standard

All changelogs must follow the [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) format, combined with [Semantic Versioning](https://semver.org/).

### Format

```markdown
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- added new feature X that does Y

### Changed
- changed behavior of Z to handle edge case W

### Fixed
- fixed a bug where A caused B

## [1.0.0] - 2026-01-15

### Added
- added initial project setup with Clean Architecture
```

### Change Categories

| Category       | When to Use                                       |
|----------------|---------------------------------------------------|
| **Added**      | New features, new files, new capabilities         |
| **Changed**    | Modifications to existing functionality           |
| **Deprecated** | Features that will be removed in a future version |
| **Removed**    | Features that were removed                        |
| **Fixed**      | Bug fixes                                         |
| **Security**   | Vulnerability fixes                               |

### Writing Rules

- **Write for humans, not machines.** Describe *what changed and why*, not implementation details.
- **Use simple past tense.** "added", "changed", "fixed", "removed" -- consistent with the commit message standard.
- **Start each entry with a lowercase verb.** Example: `- added automatic Dockerfile image tag update`.
- **Be specific.** Bad: `- updated dependencies`. Good: `- added JavaScript updater supporting npm, yarn, and pnpm projects`.
- **Link to issues or PRs** when the change is non-trivial.
- **Group related changes** in a single entry rather than listing every file touched.

### The `[Unreleased]` Section

All in-progress changes go under `[Unreleased]`. When a release is cut:

1. Create a branch `bump/x.x.x`.
2. Move entries from `[Unreleased]` to a new version heading with the release date.
3. Open a Pull Request targeting `main`.
4. After merge, create a Git tag for the version.

## README Standard

The `README.md` must accurately describe the current state of the project. Update it whenever:

- A new feature changes how users interact with the project.
- CLI commands, flags, or configuration options are added, changed, or removed.
- Setup instructions, prerequisites, or environment requirements change.
- The project structure or architecture changes significantly.
- New dependencies or integrations are introduced.

### Recommended Sections

| Section                              | Purpose                                       |
|--------------------------------------|-----------------------------------------------|
| **Title and description**            | One-line summary of what the project does     |
| **Quick start / Installation**       | How to get running in under 5 minutes         |
| **Usage**                            | Commands, configuration, and examples         |
| **Architecture / Project structure** | High-level overview of directories and layers |
| **Development**                      | How to build, test, and contribute            |
| **References**                       | Links to external documentation               |

## AI Assistant Instructions

Projects that use AI-assisted development (GitHub Copilot, Cursor, etc.) should maintain a `.github/copilot-instructions.md` file. This file provides the AI with project-specific context about:

- Project purpose and architecture
- Build, test, and lint commands with expected timings
- Repository structure and key files
- Development workflow and validation steps
- Testing infrastructure and conventions

Update this file whenever the development workflow, architecture, or key commands change.

## Workflow Integration

Documentation updates must be part of the same commit or PR that introduces the change:

1. **Write the code change.**
2. **Update `CHANGELOG.md`** -- add an entry under `[Unreleased]` describing what changed.
3. **Update `README.md`** -- if the change affects usage, setup, or architecture.
4. **Update `.github/copilot-instructions.md`** -- if the change affects build commands, project structure, or development workflow.
5. **Commit everything together.** Documentation and code ship as one unit.

**Never merge a PR that introduces user-facing or architectural changes without the corresponding documentation update.**

## Automation with AutoBump

[AutoBump](https://github.com/rios0rios0/autobump) is a CLI tool that automates the release step of the changelog workflow. When the `[Unreleased]` section is ready to ship, AutoBump detects the project language, moves unreleased entries into a new versioned section with the current date, updates language-specific version files (e.g., `go.mod`, `package.json`, `pyproject.toml`, `build.gradle`), commits, pushes, and opens a merge/pull request -- all in a single command.

It supports Go, Java, Python, TypeScript, and C# projects, with automatic language detection, and works across GitHub, GitLab, and Azure DevOps.

**AutoBump does not replace the discipline of writing changelog entries.** The `CHANGELOG.md`, `README.md`, and other documentation files must already exist and be maintained by the team as part of every change. AutoBump only automates the versioning and release ceremony -- not the content creation.

---

# README Template

> **TL;DR:** Copy this template when creating a new project. Replace all `{PLACEHOLDERS}` with project-specific values. Remove optional sections that do not apply. Every project must have a README that follows this structure.

## Overview

This template defines the standard structure for all project README files. It ensures consistency across repositories and covers the sections described in the Documentation & Change Control guide. Copy the raw template below, replace placeholders, and remove optional sections wrapped in HTML comments.

## Placeholder Reference

| Placeholder      | Description                        | Example                |
|------------------|------------------------------------|------------------------|
| `{project-name}` | Human-readable project name        | `AutoBump`             |
| `{ORG}`          | GitHub organization or user        | `rios0rios0`           |
| `{REPO}`         | Repository name                    | `autobump`             |
| `{PACKAGE}`      | Published package name (npm, PyPI) | `@rios0rios0/autobump` |

## Template

````markdown
<h1 align="center">{project-name}</h1>
<p align="center">
    <a href="https://github.com/{ORG}/{REPO}/releases/latest">
        <img src="https://img.shields.io/github/release/{ORG}/{REPO}.svg?style=for-the-badge&logo=github" alt="Latest Release"/></a>
    <a href="https://github.com/{ORG}/{REPO}/blob/main/LICENSE">
        <img src="https://img.shields.io/github/license/{ORG}/{REPO}.svg?style=for-the-badge&logo=github" alt="License"/></a>
    <a href="https://github.com/{ORG}/{REPO}/actions/workflows/default.yaml">
        <img src="https://img.shields.io/github/actions/workflow/status/{ORG}/{REPO}/default.yaml?branch=main&style=for-the-badge&logo=github" alt="Build Status"/></a>
    <!-- Add SonarCloud badges when configured -->
    <!--
    <a href="https://sonarcloud.io/summary/overall?id={ORG}_{REPO}">
        <img src="https://img.shields.io/sonar/coverage/{ORG}_{REPO}?server=https%3A%2F%2Fsonarcloud.io&style=for-the-badge&logo=sonarqubecloud" alt="Coverage"/></a>
    <a href="https://sonarcloud.io/summary/overall?id={ORG}_{REPO}">
        <img src="https://img.shields.io/sonar/quality_gate/{ORG}_{REPO}?server=https%3A%2F%2Fsonarcloud.io&style=for-the-badge&logo=sonarqubecloud" alt="Quality Gate"/></a>
    -->
    <!-- Language-specific badges (uncomment the one that applies) -->
    <!-- Go: -->
    <!-- <a href="https://pkg.go.dev/github.com/{ORG}/{REPO}"><img src="https://img.shields.io/badge/go-reference-007d9c?style=for-the-badge&logo=go" alt="Go Reference"/></a> -->
    <!-- Java: -->
    <!-- <a href="https://central.sonatype.com/artifact/{GROUP_ID}/{ARTIFACT_ID}"><img src="https://img.shields.io/maven-central/v/{GROUP_ID}/{ARTIFACT_ID}?style=for-the-badge&logo=apachemaven" alt="Maven Central"/></a> -->
    <!-- Python: -->
    <!-- <a href="https://pypi.org/project/{PACKAGE}"><img src="https://img.shields.io/pypi/v/{PACKAGE}?style=for-the-badge&logo=pypi" alt="PyPI"/></a> -->
    <!-- JavaScript/TypeScript: -->
    <!-- <a href="https://www.npmjs.com/package/{PACKAGE}"><img src="https://img.shields.io/npm/v/{PACKAGE}?style=for-the-badge&logo=npm" alt="npm"/></a> -->
    <!-- OpenSSF: -->
    <!-- <a href="https://www.bestpractices.dev/projects/{ID}"><img src="https://img.shields.io/cii/level/{ID}?style=for-the-badge&logo=opensourceinitiative" alt="OpenSSF Best Practices"/></a> -->
</p>

One to two sentence description of what this project does and why it exists.

## Features

- **Feature One**: brief explanation
- **Feature Two**: brief explanation
- **Feature Three**: brief explanation

<!-- OPTIONAL: only for multi-ecosystem tools -->
<!--
## Supported Ecosystems

| Ecosystem | Detection | Version File |
|-----------|-----------|--------------|
| Go        | `go.mod`  | `go.mod`     |
-->

## Installation

<!-- Keep only the section that applies to your project's language -->

<!-- Go CLI or library -->
```bash
go install github.com/{ORG}/{REPO}@latest
# or for libraries:
go get github.com/{ORG}/{REPO}
```

<!-- Java (Gradle) -->
```groovy
implementation '{GROUP_ID}:{ARTIFACT_ID}:{VERSION}'
```

<!-- Python -->
```bash
pdm add {PACKAGE}
# or:
pip install {PACKAGE}
```

<!-- JavaScript/TypeScript -->
```bash
npm install {PACKAGE}
```

Download pre-built binaries from the [releases page](https://github.com/{ORG}/{REPO}/releases).

<!-- OPTIONAL: only for CLIs with config files -->
<!--
## Configuration

Create `~/.config/{REPO}.yaml`:

```yaml
key: 'value'
```
-->

## Usage

```bash
{REPO} [flags] [arguments]
```

Brief explanation of the primary workflow.

<!-- OPTIONAL: only for libraries and complex CLIs -->
<!--
## Architecture

```
{REPO}/
├── domain/           # core business objects and contracts
├── infrastructure/   # implementations
└── ...
```
-->

<!-- OPTIONAL: only for libraries exposing interfaces -->
<!--
## API Reference

- **`InterfaceName`**: what it does
-->

## Contributing

Contributions are welcome. See CONTRIBUTING.md for guidelines.

## License

See [LICENSE](LICENSE) file for details.
````

---

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
