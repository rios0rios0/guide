# Documentation & Change Control

> **TL;DR:** Every project must maintain a `README.md` and a `CHANGELOG.md`. *How* the changelog is written depends on the project: one that uses [chlog](https://github.com/luizjhonata/chlog) writes a YAML fragment per change under `.changes/unreleased/` and never edits `CHANGELOG.md` by hand; one that has not adopted chlog edits the `[Unreleased]` section of `CHANGELOG.md` directly. Either way the changelog entry ships with the change -- **always** -- and README and other docs (e.g., `.github/copilot-instructions.md`) are updated **whenever behavior, configuration, or architecture changes**.

## Overview

Documentation and change control are integral parts of the engineering workflow, not afterthoughts. A well-maintained changelog and README provide traceability for stakeholders, reduce onboarding friction for new team members, and ensure that the state of the project is always understandable from its documentation alone.

**Every change introduced to a project must include updates to the relevant documentation files.** This is enforced as part of the development workflow, not as a separate task.

## Required Documentation Files

Every project must contain at minimum:

| File                                  | Purpose                                                                                                     | Required when                     | Update Frequency                                      |
|---------------------------------------|-------------------------------------------------------------------------------------------------------------|-----------------------------------|-------------------------------------------------------|
| **`README.md`**                       | Describes the project, its usage, setup, and architecture                                                   | Always                            | When behavior, configuration, CLI, or setup changes   |
| **`CONTRIBUTING.md`**                 | Guides contributors on prerequisites, workflow, and standards                                               | Always                            | When prerequisites, workflow, or structure changes    |
| **`CHANGELOG.md`**                    | Records all notable changes, organized by version                                                           | Always                            | See [Changelog Standard](#changelog-standard)         |
| **`.chlog.yaml`**                     | chlog's configuration -- the kinds and the bump each infers, and the marker CI detects                      | When the project uses chlog       | When the kinds change (rarely)                        |
| **`.changes/unreleased/`**            | One YAML changelog fragment per change, written by `chlog new`                                              | When the project uses chlog       | **Every change** (always)                             |
| **`.github/copilot-instructions.md`** | AI assistant context for the project structure and workflows                                                | Always                            | When architecture, commands, or workflow changes      |

**Templates are available for standardized project setup:**

## Changelog Standard

Every project keeps a `CHANGELOG.md` in the [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
format, versioned with [Semantic Versioning](https://semver.org/). What differs between projects is
**who writes that file**:

| Mode                | Where a change is recorded                      | Who writes `CHANGELOG.md`            |
|---------------------|-------------------------------------------------|--------------------------------------|
| **With chlog**      | A YAML fragment under `.changes/unreleased/`    | `chlog merge`, at release time       |
| **Without chlog**   | A bullet under `[Unreleased]` in `CHANGELOG.md` | The author of the change, by hand    |

Both modes produce the same released file and obey the same [writing rules](#writing-rules). Only
the mechanics differ, so a contributor moving between repositories has one thing to check first:
which mode this repository is in.

### Deciding Which Mode Applies

Look at the project root:

| Signal at the project root                       | Mode                                                       |
|--------------------------------------------------|-------------------------------------------------------------|
| `.chlog.yaml` (or `.chlog.yml`) exists           | **chlog** -- write a fragment                               |
| No config file, but `.changes/` exists           | **chlog** -- write a fragment, and add the missing config   |
| Neither exists                                   | **No chlog** -- edit `[Unreleased]` in `CHANGELOG.md`       |

The shared [`rios0rios0/pipelines`](https://github.com/rios0rios0/pipelines) basic-checks gate
decides the same way, and it keys on **`.chlog.yaml` specifically**. A project that adopts chlog
must therefore commit that file: with the fragments in place but no config, CI still asks for a
`CHANGELOG.md` diff and the fragment does not satisfy it. Its values repeating chlog's built-in
defaults is not a reason to delete it -- the file is the marker that flips the gate, whatever it
contains.

## Changelog with chlog

[chlog](https://github.com/luizjhonata/chlog) is a single Go binary that compiles per-change
fragments into `CHANGELOG.md`:

```bash
go install github.com/luizjhonata/chlog@latest
```

`go install` builds from source, so it needs a Go toolchain on the contributor's machine -- which a
Java, Python, or TypeScript project has no reason to require. Those projects should take the
prebuilt binary for their platform from the
[releases page](https://github.com/luizjhonata/chlog/releases) instead; chlog publishes one for
Linux, macOS, and Windows on both amd64 and arm64. Either way the tool is a single self-contained
binary, and the project itself never gains a Go dependency.

### One Fragment per Change

Every change writes its **own** YAML file under `.changes/unreleased/`:

```bash
chlog new --kind Added --body "added JavaScript updater supporting npm, yarn, and pnpm projects"
chlog new --kind Changed --breaking --body "**BREAKING CHANGE:** changed `Input` to take its value from props"
```

which produces a file the tool names for you:

```yaml
kind: 'Added'
body: 'added JavaScript updater supporting npm, yarn, and pnpm projects'
time: '2026-01-15T09:41:02.117823941Z'
```

This buys the one thing a single shared file cannot give: two branches each adding an entry no
longer touch the same lines, so a rebase that used to conflict on `CHANGELOG.md` now conflicts on
nothing.

**Never edit `CHANGELOG.md` by hand in this mode.** The only writer is `chlog merge`, at release time.

### Kinds

The kind is the Keep a Changelog category, and it carries the version bump a release infers from it:

| Kind           | When to Use                                       | Infers  |
|----------------|---------------------------------------------------|---------|
| **Added**      | New features, new files, new capabilities         | `minor` |
| **Changed**    | Modifications to existing functionality           | `minor` |
| **Deprecated** | Features that will be removed in a future version | `minor` |
| **Removed**    | Features that were removed                        | `minor` |
| **Fixed**      | Bug fixes                                         | `patch` |
| **Security**   | Vulnerability fixes                               | `patch` |

**No kind infers a major.** Under SemVer a major bump means a backward-incompatible change, which is
a property of the change and not of its category -- so it is signalled per fragment with
`chlog new --breaking`, and never inferred from a label. Keep the `**BREAKING CHANGE:**` prefix in
the body as well: the flag drives the version, the prefix tells the reader.

### Configuration

Every project that uses chlog carries a `.chlog.yaml` at its root. Spell chlog's defaults out rather
than relying on them: the file is what CI detects, and it makes the kinds and their bump levels
readable without going to the tool. It follows the YAML conventions like
every other YAML file -- the `.yaml` extension, and single quotes around every string:

```yaml
changesDir: '.changes'
unreleasedDir: 'unreleased'
changelogPath: 'CHANGELOG.md'
versionFormat: '## [{{.Version}}] - {{.Time.Format "2006-01-02"}}'
kindFormat: '### {{.Kind}}'
changeFormat: '- {{.Body}}'
kinds:
  - label: 'Added'
    auto: 'minor'
  - label: 'Changed'
    auto: 'minor'
  - label: 'Deprecated'
    auto: 'minor'
  - label: 'Removed'
    auto: 'minor'
  - label: 'Fixed'
    auto: 'patch'
  - label: 'Security'
    auto: 'patch'
```

The double quotes inside `versionFormat` belong to the Go template, not to YAML: a single-quoted
scalar takes them literally, which is exactly what the template needs.

### The Generated File

`CHANGELOG.md` keeps its familiar shape -- an empty `[Unreleased]` heading and one section per
released version:

```markdown
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

This file is not edited by hand. Every change writes its own fragment under
`.changes/unreleased/` with [chlog](https://github.com/luizjhonata/chlog), and a release
compiles the pending fragments into a version section here.

## [Unreleased]

## [1.0.0] - 2026-01-15

### Added
- added initial project setup with Clean Architecture
```

### Releasing

The pending fragments become a version section at release time:

1. Create a branch `chore/bump-x.y.z`.
2. Run `chlog batch auto && chlog merge`. `batch auto` derives the version from the pending kinds and
   their `breaking` flags; `merge` folds the compiled batch into `CHANGELOG.md` and empties
   `.changes/unreleased/`.
3. Open a Pull Request targeting `main`.
4. After merge, create a Git tag for the version.

[AutoBump](https://github.com/rios0rios0/autobump) performs steps 2 and 3 -- see
[Automation with AutoBump](#automation-with-autobump).

### CI Enforcement

The shared [`rios0rios0/pipelines`](https://github.com/rios0rios0/pipelines) basic-checks gate is
chlog-aware: when `.chlog.yaml` is present it requires a **new fragment** on an ordinary branch, and
flips to requiring an updated `CHANGELOG.md` on a `chore/bump-*` (or `bump/*`) branch, where
`chlog merge` has already consumed the fragments. No per-repository CI job is needed.
`chlog hook install --local` runs the same check on every commit in a local clone, and `chlog check`
runs it on demand.

### AI Assistants

`chlog ai setup` injects a marked block (`<!-- chlog:start -->` ... `<!-- chlog:end -->`) into
`CLAUDE.md` and `.github/copilot-instructions.md` telling the assistant to write a fragment on every
change and never to touch `CHANGELOG.md`. Re-running the command updates the block in place, so it
is safe to run again after an upgrade. The injected block is conditional on the repository actually
using chlog, so it is harmless in a repository that has not adopted it yet.

## Changelog without chlog

A project that has not adopted chlog keeps the same file, written by hand. Everything a change adds
goes under `[Unreleased]`, grouped by category:

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

The categories are the same ones chlog calls kinds:

| Category       | When to Use                                       |
|----------------|---------------------------------------------------|
| **Added**      | New features, new files, new capabilities         |
| **Changed**    | Modifications to existing functionality           |
| **Deprecated** | Features that will be removed in a future version |
| **Removed**    | Features that were removed                        |
| **Fixed**      | Bug fixes                                         |
| **Security**   | Vulnerability fixes                               |

A backward-incompatible change is marked in the entry itself, since there is no `--breaking` flag to
carry it:

```markdown
- **BREAKING CHANGE:** changed `Input` to take its value from props
```

### The `[Unreleased]` Section

All in-progress changes go under `[Unreleased]`. When a release is cut:

1. Create a branch `chore/bump-x.y.z`.
2. Move the entries from `[Unreleased]` to a new version heading with the release date.
3. Open a Pull Request targeting `main`.
4. After merge, create a Git tag for the version.

### CI Enforcement

With no `.chlog.yaml` at the root, the basic-checks gate falls back to requiring that the PR touch
`CHANGELOG.md`. A "CHANGELOG.md was NOT modified" failure therefore means the project is in this
mode: add the entry under `[Unreleased]` rather than reaching for `chlog new`.

## Writing Rules

These govern the text of the entry -- the `--body` of a fragment, or the bullet typed under
`[Unreleased]`. They are identical in both modes:

- **Write for humans, not machines.** Describe *what changed and why*, not implementation details.
- **Use simple past tense.** "added", "changed", "fixed", "removed" -- consistent with the commit message standard.
- **Start each entry with a lowercase verb.** Example: `added automatic Dockerfile image tag update`.
- **Be specific.** Bad: `updated dependencies`. Good: `added JavaScript updater supporting npm, yarn, and pnpm projects`.
- **Link to issues or PRs** when the change is non-trivial.
- **Group related changes** in a single entry rather than listing every file touched.
- **One entry per change, not per commit.** A branch that does one thing carries one entry.

See CHANGELOG Formatting for capitalization
and backtick rules.

## Adopting chlog

Moving a project from the hand-written file to fragments:

1. Release or carry over whatever sits under `[Unreleased]`. Anything left there is invisible to
   `chlog batch`, which only reads fragments -- so either cut a release first, or re-create each
   pending entry with `chlog new`.
2. Add `.chlog.yaml` at the root, as spelled out in [Configuration](#configuration). This is the step
   that flips CI; without it the gate keeps asking for a `CHANGELOG.md` diff.
3. Leave the released sections of `CHANGELOG.md` untouched -- chlog appends above them.
4. Add the header note to `CHANGELOG.md` saying the file is generated, so the next contributor does
   not hand-edit it.
5. Run `chlog ai setup` to update `CLAUDE.md` and `.github/copilot-instructions.md`.
6. Update `CONTRIBUTING.md`: add chlog to the prerequisites and replace the "update CHANGELOG.md"
   step with `chlog new`.

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
2. **Record the change in the changelog** -- `chlog new --kind <Kind> --body "..."` when the project
   uses chlog, or a bullet under `[Unreleased]` in `CHANGELOG.md` when it does not. See
   [Deciding Which Mode Applies](#deciding-which-mode-applies).
3. **Update `README.md`** -- if the change affects usage, setup, or architecture.
4. **Update `.github/copilot-instructions.md`** -- if the change affects build commands, project structure, or development workflow.
5. **Commit everything together.** Documentation and code ship as one unit.

**Never merge a PR that introduces user-facing or architectural changes without the corresponding documentation update.**

## Automation with AutoBump

[AutoBump](https://github.com/rios0rios0/autobump) is a CLI tool that automates the release step of the changelog workflow. When the pending changes are ready to ship, AutoBump detects the project language, compiles the pending entries into a new versioned section with the current date, updates language-specific version files (e.g., `go.mod`, `package.json`, `pyproject.toml`, `build.gradle`), commits, pushes, and opens a merge/pull request -- all in a single command.

It handles both modes: it detects the chlog layout (a `.chlog.yaml` or the fragment directory), reads the fragments under `.changes/unreleased/` directly, and consumes them -- otherwise it reads the `[Unreleased]` section. Adopting chlog therefore changes nothing about the release flow.

It supports Go, Java, Python, TypeScript, and C# projects, with automatic language detection, and works across GitHub, GitLab, and Azure DevOps.

**AutoBump does not replace the discipline of writing changelog entries.** The fragments (or the `[Unreleased]` bullets), `README.md`, and other documentation files must already exist and be maintained by the team as part of every change. AutoBump only automates the versioning and release ceremony -- not the content creation.

Because AutoBump opens its own release pull requests, those commits are exempt from the ticket
reference every other commit carries -- see
[Ticket Reference](Git-Flow.md#ticket-reference) in the Git Flow guide.

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

A project that uses [chlog](https://github.com/luizjhonata/chlog) -- one with a `.chlog.yaml` at its
root -- additionally needs the tool to write changelog fragments. It is a single self-contained
binary, so the project itself never gains a Go dependency -- but `go install github.com/luizjhonata/chlog@latest`
builds it from source and therefore needs a Go toolchain locally. On a Java, Python, or JavaScript
project, point contributors at the prebuilt binary for their platform on the
[releases page](https://github.com/luizjhonata/chlog/releases) (Linux, macOS, and Windows; amd64 and
arm64) and keep the toolchain out of the prerequisites. A project that does not use chlog drops the
prerequisite entirely and edits `CHANGELOG.md` by hand -- see
Documentation & Change Control.

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
