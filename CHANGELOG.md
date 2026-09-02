# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

This file is not edited by hand. Every change writes its own fragment under
`.changes/unreleased/` with [chlog](https://github.com/luizjhonata/chlog), and a release compiles
the pending fragments into a version section here — so two branches each adding an entry no
longer touch the same lines, and a rebase that used to conflict on this file now conflicts on
nothing.

## [Unreleased]

## [0.6.0] - 2026-09-02

### Added

- added the `checks` workflow, so pull requests here run the shared `code-check > quality:basic-checks` gate (rebase status and the changelog rule) that every repository with a language pipeline already gets as that pipeline's first job. This repository has no build to attach it to, so it had no changelog enforcement at all — which is how the weekly configuration and documentation refresh hand-edited a generated `CHANGELOG.md` across the fleet before anything objected

### Removed

- removed the `.chlog.yaml` added alongside the `checks` workflow. It carried nothing but chlog's own defaults, and the changelog gate this repository now runs already detects chlog from the `.changes/unreleased/` directory — so the file was a second copy of the defaults with nothing to say and another place to drift

## [0.5.2] - 2026-09-01

### Changed

- refreshed `CLAUDE.md` and `.github/copilot-instructions.md` to rename the Claude workflows from `claude-code-review.yaml`/`claude.yaml` to `claude-review.yaml`/`claude-mention.yaml` and note they call the reusable workflows in `rios0rios0/pipelines`, and refreshed `.github/skills/code-review/SKILL.md` to replace the "no build step" claim with the `go build`/`go test` commands for the `update-wiki` and `generate-ai-rules` tools

## [0.5.1] - 2026-08-28

### Changed

- changed the Claude workflows to call the reusable workflows in `rios0rios0/pipelines` instead of `rios0rios0/.github`, renamed them to `claude-review.yaml` and `claude-mention.yaml`, matching the `reusable-claude-review.yaml` / `reusable-claude-mention.yaml` definitions they call, and changed them to pass `CLAUDE_CODE_OAUTH_TOKEN` explicitly rather than with `secrets: inherit`, which fails Semgrep's `yaml.github-actions.security.secrets-inherit` rule

### Fixed

- restored the `.changes/unreleased/` directory with a `.gitkeep`, so the release tooling keeps recognizing this project as [chlog](https://github.com/luizjhonata/chlog)-based after a release consumes the last fragment. Git tracks files rather than directories, so the bump commit that removed the final fragment removed the directory too, and the next run read the empty `[Unreleased]` section as "nothing to release"
- restored the `id-token: write` permission on both Claude workflow callers. Without it the caller grants less than the reusable workflow declares, which GitHub rejects before the job starts -- runs ended in `startup_failure`. The action needs the scope because `setupGitHubToken()` exchanges a GitHub OIDC token for the GitHub App token it posts with, unless a `github_token` is passed explicitly.

### Removed

- removed the unused `id-token: write` permission from the Claude workflow callers, and changed `claude-review.yaml`'s display name to `Claude Review` so it matches its file name and its `Claude Mention` sibling. `anthropics/claude-code-action` needs `id-token: write` only for workload identity federation or the Bedrock / Vertex / Foundry OIDC paths; these authenticate with `claude_code_oauth_token`, so the scope allowed minting OIDC tokens for any audience without ever being used.

## [0.5.0] - 2026-08-26

### Added

- added a `Ticket Reference` section to `Git Flow` making the ticket ID mandatory in branch names and commit scopes, and enumerating the three changes exempt from it because no human opens a ticket for them: dependency upgrades from [autoupdate](https://github.com/rios0rios0/autoupdate) (`chore/autoupdate-YYYY-MM-DD`, `chore(deps)`), release bumps from [autobump](https://github.com/rios0rios0/autobump) (`chore/bump-x.y.z`, `chore(bump)`), and configuration and documentation refreshes from [config-automation](https://github.com/rios0rios0/config-automation) (`chore/config-and-docs-refresh`, `chore(refresh)`). The exemption follows the change rather than the author, waives only the ticket, and is mirrored in the `git-workflow` agent and the `CONTRIBUTING` template.
- added a tailored `code-review` skill under `.github/skills/` so GitHub Copilot reviews changes against the [rios0rios0/guide](https://github.com/rios0rios0/guide/wiki) standards and this repository's own load-bearing invariants
- added the chlog prerequisite and the fragment step to this repository's own `CONTRIBUTING.md`, which had adopted chlog without telling its contributors -- the standards repository now follows the standard it publishes.

### Changed

- changed the `.chlog.yaml` example to follow the [YAML conventions](https://github.com/rios0rios0/guide/wiki/YAML) it publishes -- every string is single-quoted, with a note that the double quotes inside `versionFormat` belong to the Go template and survive a single-quoted YAML scalar.
- changed the changelog standard itself: `Documentation & Change Control` now documents the chlog fragment flow (one fragment per change, `CHANGELOG.md` generated at release time, `--breaking` as the only major bump), `CHANGELOG Formatting` now governs the body of a fragment, and the `CONTRIBUTING` template, `Git Flow`, the `changelog-enforcer` agent, the `git-workflow` agent, `fix-ci`, `fix-guardrails`, and the `changelog-guard` hook were retargeted to match — so every repository consuming these standards is told to write a fragment rather than edit `CHANGELOG.md` by hand
- changed the changelog standard to apply chlog **when the project uses it** rather than unconditionally: `Documentation & Change Control` now documents both modes side by side -- a chlog project writes a fragment under `.changes/unreleased/`, a project without it edits `[Unreleased]` in `CHANGELOG.md` by hand -- with a detection table (`.chlog.yaml`/`.chlog.yml`/`.changes/`) matching the one the shared pipelines basic-checks gate uses, shared writing rules, and an adoption checklist. `CHANGELOG Formatting`, the `CONTRIBUTING` template, `Git Flow`, the `changelog-enforcer` and `git-workflow` agents, `fix-ci`, `fix-guardrails`, and the `changelog-guard` hook all branch the same way, so a repository that has not adopted chlog is no longer told to run a tool it does not have.
- changed the changelog to [chlog](https://github.com/luizjhonata/chlog) fragments: a change now writes its own YAML file under `.changes/unreleased/` through `chlog new --kind <Kind> --body "..."`, and `CHANGELOG.md` is GENERATED from them at release time by `chlog batch auto && chlog merge`. That is the one thing a single shared file cannot do — two branches each adding an entry no longer touch the same lines, so a rebase that used to conflict on `CHANGELOG.md` now conflicts on nothing. The `[Unreleased]` section was empty, so nothing had to be carried across. AutoBump already reads the fragments directly, so the release flow is unchanged.
- changed the chlog install guidance to say what `go install` actually costs: it builds from source and so needs a Go toolchain on the machine, which a Java, Python, or TypeScript project has no reason to require. Those projects are now pointed at the prebuilt binaries chlog publishes for Linux, macOS, and Windows on amd64 and arm64. Both the Documentation & Change Control standard and the CONTRIBUTING template carry the alternative.

### Fixed

- fixed the wiki link transformer dropping `#anchor` fragments: `](Life-Cycle/Git-Flow.md#ticket-reference)` matched no rewrite rule and reached the Wiki as a broken relative path. It now flattens to `](Git-Flow#ticket-reference)`, repairing the existing cross-page anchor in the Python conventions guide as well.

## [0.4.3] - 2026-07-16

### Fixed

- fixed the documented release-bump branch name from `bump/x.y.z` to `chore/bump-x.y.z` in the Documentation & Change Control guide

## [0.4.2] - 2026-06-09

### Changed

- changed Go testing convention to remove mandatory `//go:build unit` tag from unit tests — build tags are now required only for integration, e2e, and other non-unit test types

## [0.4.1] - 2026-06-03

### Changed

- refreshed `CLAUDE.md` to correct the generated rule-file count from 14 to 15 per assistant, matching the 15 rule groups defined in `generate-ai-rules` `config.go`

## [0.4.0] - 2026-05-19

### Added

- added `.claude-plugin/marketplace.json` for Claude Code plugin marketplace (`/plugin marketplace add rios0rios0/guide`)
- added `aisync-source.yaml` generation on the `generated` branch for [aisync](https://github.com/rios0rios0/aisync) users
- added `codex/AGENTS.md` mapping to generated `aisync-source.yaml` so aisync installs the full Codex output, not just `codex/rules`

### Changed

- refreshed `CLAUDE.md` and `.github/copilot-instructions.md` to document the `release.yaml` workflow and correct the `claude-code-review.yaml` trigger list

### Removed

- removed `claude/skills` mapping from generated `aisync-source.yaml` — the `generated` branch has no `claude/skills/` directory; skills live under `cursor/skills/`
- removed `external-sources.yaml` and external source fetching from `generate-ai-rules` tool — each external repo is now an independent source that users add directly to their aisync config
- removed `external.go` and `external_test.go` from `generate-ai-rules` tool
- removed `install-rules.sh` — distribution now handled by [aisync](https://github.com/rios0rios0/aisync)
- removed unused `gopkg.in/yaml.v3` direct dependency from `generate-ai-rules` `go.mod` after deleting `external.go`

## [0.3.2] - 2026-04-29

### Changed

- bumped `update-wiki` and `generate-ai-rules` Go modules to `1.26.2`
- switched `update-wiki.yml` and `generate-ai-rules.yaml` workflows to `go-version-file` so the workflow Go version always tracks each subproject's `go.mod`, preventing future drift

### Fixed

- aligned the workflow Go version with each subproject's `go.mod`, fixing build failures in `Generate AI Rules` and `Update Wiki` caused by workflow/toolchain version drift

## [0.3.1] - 2026-04-28

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

