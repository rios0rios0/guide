# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A documentation repository that serves as the single source of truth for software engineering standards. It contains 80+ Markdown files organized hierarchically and includes two Go-based automation tools that transform this documentation into GitHub Wiki pages and AI assistant rule files (for Claude Code, Cursor, Codex, and GitHub Copilot).

## Build, Test, and Validate

There is no Makefile. The two Go tools are built and tested independently:

```bash
# update-wiki tool (Go 1.26.0) — syncs docs to GitHub Wiki
cd .github/workflows/update-wiki
go build -o update-wiki ./...
go test ./...

# generate-ai-rules tool (Go 1.24.7) — generates AI rule files from docs
cd .github/workflows/generate-ai-rules
go build -o generate-ai-rules ./...
go test ./...

# TOC sync validation (run from repo root)
bash .github/workflows/sync-docs/check-toc-sync.sh
```

Build times are ~1 second each. Set timeouts to 60+ seconds for safety.

## Architecture

### Documentation Flow

```
Markdown files (source of truth)
  ├─→ update-wiki tool ─→ GitHub Wiki (flattened links, absolute image URLs)
  ├─→ generate-ai-rules tool ─→ 'generated' branch (claude/, cursor/, codex/, copilot/)
  │   ├── Claude Code, Cursor, Codex, GitHub Copilot rule files
  │   ├── Static assets (agents, commands, skills) ─→ also copied to 'generated' branch
  │   └── External agents (fetched from external-sources.yaml) ─→ merged into agents/
  └─→ install-rules.sh ─→ downloads from 'generated' branch to ~/.claude/, ~/.cursor/, etc.
```

The generated rule directories (`claude/`, `cursor/`, `codex/`, `copilot/`) do **not** exist on `main`. They live only on the `generated` branch, which is updated automatically by the `generate-ai-rules.yaml` workflow.

### Go Tools

- **`update-wiki`** (`.github/workflows/update-wiki/main.go`): Converts relative markdown links to Wiki-flat format and resolves image paths to GitHub raw content URLs.
- **`generate-ai-rules`** (`.github/workflows/generate-ai-rules/`): Parses documentation, extracts sections, and formats them into AI-assistant-specific rule files. Has separate `config.go`, `parser.go`, `formatter.go`, and `external.go` modules with corresponding tests. Supports fetching agents from external GitHub repos via `external-sources.yaml`.

### Static Assets (on `main`)

Hand-written files that the workflow copies to the `generated` branch alongside auto-generated rules:

- `.github/workflows/generate-ai-rules/agents/` — 6 Claude Code agent files
- `.github/workflows/generate-ai-rules/commands/` — 5 Claude Code slash commands
- `.github/workflows/generate-ai-rules/skills/` — 4 Cursor skills

### Generated Output (on `generated` branch)

The `generated` branch contains the distributable rule files:
- `claude/rules/` — 14 rule files (`.md`) — auto-generated from docs
- `claude/commands/` — 5 slash commands — copied from static assets
- `claude/agents/` — 6 static agents + external agents fetched from configured repos
- `cursor/rules/` — 14 rule files (`.mdc`) — auto-generated from docs
- `cursor/skills/` — 4 skills — copied from static assets
- `copilot/instructions/` — 14 instruction files (`.instructions.md`) — auto-generated with `applyTo` frontmatter
- `codex/` — `AGENTS.md` and `rules/default.rules` — auto-generated

## Critical Constraints

### Navigation File Sync

`README.md`, `Home.md`, and `_Sidebar.md` all contain a table of contents that **must stay in sync** (same labels and hierarchy; link targets may differ). The `sync-docs.yaml` workflow enforces this on every PR. Always run `bash .github/workflows/sync-docs/check-toc-sync.sh` after modifying any of these three files.

### EditorConfig

All files: 2-space indentation, LF line endings, UTF-8, final newline, trim trailing whitespace.

### File Naming

Documentation files use hyphens (e.g., `Backend-Design.md`). Directories use `&` for compound names (e.g., `Agile-&-Culture/`).

## GitHub Actions Workflows

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| `update-wiki.yml` | Push to `main`, manual | Syncs docs to GitHub Wiki |
| `generate-ai-rules.yaml` | Push to `main` (doc paths), manual | Regenerates AI rules on `generated` branch |
| `sync-docs.yaml` | PR with `.md` changes | Validates TOC sync |

## Install Script

`install-rules.sh` downloads generated rule files from the `generated` branch on GitHub and installs them globally (`~/.claude/`, `~/.cursor/`, etc.) or into a specific project directory. Supports `--force` to skip overwrite prompts.
