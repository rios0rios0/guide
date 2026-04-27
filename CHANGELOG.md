# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- refreshed `CLAUDE.md` and `.github/copilot-instructions.md` to document the `claude-code-review.yaml` and `claude.yaml` workflows, and corrected stale static-asset counts in copilot instructions

## [0.3.0] - 2026-04-14

### Added

- added `changelog-guard.sh` hook as a static asset to block commits that add CHANGELOG entries outside the `[Unreleased]` section
- added `fix-ci` slash command for auto-detecting failing CI checks, classifying failures, and pushing fixes
- added checksum comparison to `install-rules.sh` that detects unchanged files, shows byte-level size differences, and warns whether installing adds or removes content
- added hooks installation section to `install-rules.sh` so hooks from the `generated` branch are installed alongside rules, commands, and agents

## [0.2.0] - 2026-03-20

### Added

- added `resolve-pr-comments` slash command for auto-detecting the current branch's PR and resolving review comments

## [0.1.0] - 2026-03-12

### Added

- added 5 new Claude Code agents: changelog-enforcer, code-reviewer, git-workflow, security-auditor, and bulk-operations
- added `pr-review-resolver` agent for automating PR review comment resolution
- added AI-Assisted Workflows cookbook page documenting the RPI (Research, Plan, Implement, Review) methodology
- added CHANGELOG.md following the Keep a Changelog standard
- added dynamic external agent fetch mechanism to pull agents from configured GitHub repositories via `external-sources.yaml`
- added GitHub Copilot as a generation target with `.instructions.md` files using `applyTo` frontmatter

### Changed

- changed `install-rules.sh` to download from the `generated` branch instead of `main`
- moved generated `.ai/` directory from `main` branch to a dedicated `generated` branch to eliminate workflow conflicts
- moved static assets (agents, commands, skills) to `.github/workflows/generate-ai-rules/` as source files on `main`
- removed `.ai/` prefix directory from the `generated` branch; rule files now live directly at `claude/`, `cursor/`, `codex/`, `copilot/`

### Fixed

- fixed external agent paths in `external-sources.yaml` after upstream `wshobson/agents` repo reorganized its plugin structure

