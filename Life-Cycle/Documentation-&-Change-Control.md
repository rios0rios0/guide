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
| **`CLAUDE.md`**                       | Claude Code project rules and conventions                    | When architecture, commands, or development workflow changes |
| **`.cursorrules`**                    | Cursor project rules and conventions                         | When architecture, commands, or development workflow changes |
| **`AGENTS.md`**                       | OpenAI Codex project rules and conventions                   | When architecture, commands, or development workflow changes |

**Templates are available for standardized project setup:**

- [README Template](Documentation-&-Change-Control/README-Template.md) -- copy and customize for new projects
- [CONTRIBUTING Template](Documentation-&-Change-Control/CONTRIBUTING-Template.md) -- copy and customize for new projects
- [CLAUDE.md Template](Documentation-&-Change-Control/CLAUDE-Template.md) -- Claude Code rules for new projects
- [Cursorrules Template](Documentation-&-Change-Control/Cursorrules-Template.md) -- Cursor rules for new projects
- [AGENTS.md Template](Documentation-&-Change-Control/AGENTS-Template.md) -- OpenAI Codex rules for new projects

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

Projects that use AI-assisted development must maintain instruction files for each AI tool in use. These files provide AI assistants with project-specific context and enforce the team's development standards automatically.

### Supported AI Tools

| Tool           | File                              | Purpose                                                           |
|----------------|-----------------------------------|-------------------------------------------------------------------|
| GitHub Copilot | `.github/copilot-instructions.md` | Copilot context for project structure and workflows               |
| Claude Code    | `CLAUDE.md`                       | Claude Code rules for architecture, testing, and code conventions |
| Cursor         | `.cursorrules`                    | Cursor rules for architecture, testing, and code conventions      |
| OpenAI Codex   | `AGENTS.md`                       | Codex rules for architecture, testing, and code conventions       |

### What to Include

All AI instruction files should provide context about:

- Project purpose and architecture
- Build, test, and lint commands
- Repository structure and key files
- Development workflow and validation steps
- Testing infrastructure and conventions
- Commit message format and branch naming
- YAML conventions and security requirements

### Installation

Use the `install-ai-rules.sh` script from the [guide repository](https://github.com/rios0rios0/guide) to install all AI rule files into a project:

```bash
curl -fsSL https://raw.githubusercontent.com/rios0rios0/guide/main/install-ai-rules.sh | bash -s -- .
```

Or clone the guide and run locally:

```bash
./install-ai-rules.sh /path/to/your/project
```

The script generates `CLAUDE.md`, `.cursorrules`, and `AGENTS.md` with the team's full conventions embedded. It overwrites existing files with the same name. No configuration needed -- the AI agents will immediately know all architecture, testing, code style, security, and documentation conventions.

### Templates

- [CLAUDE.md Template](Documentation-&-Change-Control/CLAUDE-Template.md)
- [Cursorrules Template](Documentation-&-Change-Control/Cursorrules-Template.md)
- [AGENTS.md Template](Documentation-&-Change-Control/AGENTS-Template.md)

Update these files whenever the development workflow, architecture, or key commands change.

## Workflow Integration

Documentation updates must be part of the same commit or PR that introduces the change:

1. **Write the code change.**
2. **Update `CHANGELOG.md`** -- add an entry under `[Unreleased]` describing what changed.
3. **Update `README.md`** -- if the change affects usage, setup, or architecture.
4. **Update AI instruction files** (`CLAUDE.md`, `.cursorrules`, `AGENTS.md`, `.github/copilot-instructions.md`) -- if the change affects build commands, project structure, or development workflow.
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
