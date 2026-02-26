# README Template

> **TL;DR:** Copy this template when creating a new project. Replace all `{PLACEHOLDERS}` with project-specific values. Remove optional sections that do not apply. Every project must have a README that follows this structure.

## Overview

This template defines the standard structure for all project README files. It ensures consistency across repositories and covers the sections described in the [Documentation & Change Control](../Documentation-&-Change-Control.md) guide. Copy the raw template below, replace placeholders, and remove optional sections wrapped in HTML comments.

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

Contributions are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

See [LICENSE](LICENSE) file for details.
````

## References

- [Documentation & Change Control](../Documentation-&-Change-Control.md) -- when and how to update documentation files
- [Make a README](https://www.makeareadme.com/)
- [Awesome README](https://github.com/matiassingers/awesome-readme)
