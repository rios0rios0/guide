# Documentation & Change Control

> **TL;DR:** Every project must maintain a `CHANGELOG.md` (following [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)) and a `README.md`. Every code change must be accompanied by the corresponding documentation update -- changelog **always**, README and other docs (e.g., `.github/copilot-instructions.md`) **whenever behavior, configuration, or architecture changes**.

## Overview

Documentation and change control are integral parts of the engineering workflow, not afterthoughts. A well-maintained changelog and README provide traceability for stakeholders, reduce onboarding friction for new team members, and ensure that the state of the project is always understandable from its documentation alone.

**Every change introduced to a project must include updates to the relevant documentation files.** This is enforced as part of the development workflow, not as a separate task.

## Required Documentation Files

Every project must contain at minimum:

| File                                  | Purpose                                                      | Update Frequency                                             |
|---------------------------------------|--------------------------------------------------------------|--------------------------------------------------------------|
| **`CHANGELOG.md`**                    | Records all notable changes, organized by version            | **Every change** (always)                                    |
| **`README.md`**                       | Describes the project, its usage, setup, and architecture    | When behavior, configuration, CLI, or setup changes          |
| **`.github/copilot-instructions.md`** | AI assistant context for the project structure and workflows | When architecture, commands, or development workflow changes |

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

## References

- [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
- [Semantic Versioning](https://semver.org/)
- [Common Changelog](https://common-changelog.org/)
- [The Art of Writing a Great Changelog](https://softwareforprogress.org/learn/the-art-of-writing-a-great-changelog/)
- [GitHub Copilot Instructions](https://docs.github.com/en/copilot/customizing-copilot/adding-repository-custom-instructions-for-github-copilot)
