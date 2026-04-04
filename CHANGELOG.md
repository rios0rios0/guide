# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- added `.claude-plugin/marketplace.json` for Claude Code plugin marketplace (`/plugin marketplace add rios0rios0/guide`)
- added `aisync-source.yaml` generation on the `generated` branch for [aisync](https://github.com/rios0rios0/aisync) users
- added `fix-ci` slash command for auto-detecting failing CI checks, classifying failures, and pushing fixes
- added `changelog-guard.sh` hook as a static asset to block commits that add CHANGELOG entries outside the `[Unreleased]` section

### Removed

- removed `install-rules.sh` — distribution now handled by [aisync](https://github.com/rios0rios0/aisync)
- removed `external-sources.yaml` and external source fetching from `generate-ai-rules` tool — each external repo is now an independent source that users add directly to their aisync config
- removed `external.go` and `external_test.go` from `generate-ai-rules` tool

## [0.2.0] - 2026-03-20

### Added

- added `resolve-pr-comments` slash command for auto-detecting the current branch's PR and resolving review comments

## [0.1.0] - 2026-03-12

### Added

- added `pr-review-resolver` agent for automating PR review comment resolution
- added GitHub Copilot as a generation target with `.instructions.md` files using `applyTo` frontmatter
- added AI-Assisted Workflows cookbook page documenting the RPI (Research, Plan, Implement, Review) methodology
- added dynamic external agent fetch mechanism to pull agents from configured GitHub repositories via `external-sources.yaml`
- added 5 new Claude Code agents: changelog-enforcer, code-reviewer, git-workflow, security-auditor, and bulk-operations
- added CHANGELOG.md following the Keep a Changelog standard

### Changed

- removed `.ai/` prefix directory from the `generated` branch; rule files now live directly at `claude/`, `cursor/`, `codex/`, `copilot/`
- moved generated `.ai/` directory from `main` branch to a dedicated `generated` branch to eliminate workflow conflicts
- moved static assets (agents, commands, skills) to `.github/workflows/generate-ai-rules/` as source files on `main`
- changed `install-rules.sh` to download from the `generated` branch instead of `main`

### Fixed

- fixed external agent paths in `external-sources.yaml` after upstream `wshobson/agents` repo reorganized its plugin structure

