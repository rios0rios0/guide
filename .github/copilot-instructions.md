# Software Development Guide Repository

This is a comprehensive documentation repository containing software development standards, best practices, and guidelines organized across multiple categories including Code Style, Life Cycle processes, and practical Cookbooks.

## Working Effectively

### Quick Setup and Validation
- Build the wiki sync tool:
  ```bash
  cd .github/workflows/update-wiki
  go build -o update-wiki ./...
  ```
- Build the AI rules generator:
  ```bash
  cd .github/workflows/generate-ai-rules
  go build -o generate-ai-rules ./...
  ```

### Core Build Process
The repository has two Go-based tools under `.github/workflows/`:

#### update-wiki (Go 1.26.2)
Syncs documentation to GitHub Wiki:
- Build location: `.github/workflows/update-wiki/`
- Build command: `go build -o update-wiki ./...`
- Expected build time: ~1 second

#### generate-ai-rules (Go 1.26.2)
Generates AI assistant rule files for Claude Code, Cursor, Codex, and GitHub Copilot from the guide's own documentation.
- Build location: `.github/workflows/generate-ai-rules/`
- Build command: `go build -o generate-ai-rules ./...`
- Run from the repository root: `./.github/workflows/generate-ai-rules/generate-ai-rules` (it reads docs from `--source`, default `.`, and writes the git-ignored `claude/`, `cursor/`, `codex/`, `copilot/` directories)
- Expected build time: ~1 second
- `config.go` lists each rule file's source docs by path; renaming or moving one of those docs means updating it.
- `scaffold_test.go` compiles and runs the Go examples in `commands/scaffold-go-project.md`, so `go test ./...` fails when those examples break.

### GitHub Actions Workflows
Seven workflows are defined under `.github/workflows/`:

| File | Trigger | Purpose |
|------|---------|---------|
| `update-wiki.yml` | Push to `main`, manual dispatch | Syncs docs to GitHub Wiki |
| `generate-ai-rules.yaml` | Push to `main` (docs paths), manual dispatch | Regenerates AI rules on `generated` branch |
| `sync-docs.yaml` | Pull request (any `.md` change) | Validates TOC sync across README.md, Home.md, _Sidebar.md |
| `checks.yaml` | Pull request to `main` | Runs shared checks via a reusable workflow in `rios0rios0/pipelines` |
| `claude-review.yaml` | PR opened/synchronize/ready_for_review/reopened | AI-assisted PR review; calls a reusable workflow in `rios0rios0/pipelines` |
| `claude-mention.yaml` | Issue/PR comments, issue events, PR reviews | Claude Code automation for issues and PRs; calls a reusable workflow in `rios0rios0/pipelines` |
| `release.yaml` | Push to `main` | Automated release via reusable workflow |

## Validation

### Always Validate Changes
Run this validation sequence after making any changes:

1. **Build and Test** (update-wiki):
   ```bash
   cd .github/workflows/update-wiki
   go mod verify
   go build -o update-wiki ./...
   go test ./...
   ```

2. **Build and Test** (generate-ai-rules):
   ```bash
   cd .github/workflows/generate-ai-rules
   go mod verify
   go build -o generate-ai-rules ./...
   go test ./...
   ```

3. **TOC Sync Validation**:
   ```bash
   bash .github/workflows/sync-docs/check-toc-sync.sh
   ```
   The Summary sections in `README.md`, `Home.md`, and `_Sidebar.md` must stay in sync.
   Run this whenever you modify any of those three navigation files.

4. **Navigation links**: verify that every link in the `README.md` Summary section points to an existing file.

### EditorConfig Compliance
Always follow the `.editorconfig` settings:
- Indent style: spaces (not tabs)
- Indent size: 2 spaces
- End of line: LF
- Insert final newline: true
- Trim trailing whitespace: true

## Common Tasks

### Adding New Documentation
1. Create markdown file in appropriate directory structure
2. Follow existing naming conventions (use hyphens, not spaces)
3. Add navigation links to parent directory's index file
4. **Update all three navigation files**: `README.md` (Summary section), `Home.md` (Summary section), and `_Sidebar.md` must remain in sync
5. Run `bash .github/workflows/sync-docs/check-toc-sync.sh` to verify sync
6. Always validate build process after changes

### Installing AI Rules
Use [aisync](https://github.com/rios0rios0/aisync) to install and sync AI rule files across devices:
```bash
aisync init
aisync source add guide --source-repo rios0rios0/guide --branch generated
aisync pull
```

Or as a Claude Code plugin: `/plugin marketplace add rios0rios0/guide`

### Frequently Accessed Locations
- **GitHub repo page**: `README.md` (includes badges)
- **Wiki navigation**: `Home.md` (full navigation tree for GitHub Wiki)
- **Testing guidelines**: `Life-Cycle/Tests.md`
- **Code style references**: `Code-Style/<language>/` directories
- **Setup guides**: `Cookbooks/Tools-&-Setup/`
- **CI/CD information**: `Life-Cycle/CI-&-CD.md`
- **AI rules**: `generated` branch (auto-generated from docs)
- **AI rule sources**: `.github/workflows/generate-ai-rules/agents/`, `commands/`, `skills/`, `hooks/`

## Important Notes

### Repository Characteristics
- **Type**: Documentation repository (not traditional software)
- **Primary content**: 80+ Markdown files across 25+ directories
- **Build output**: GitHub Wiki synchronization + AI rule files for Claude Code, Cursor, Codex, and GitHub Copilot (on `generated` branch)
- **Dependencies**: Go 1.26.2 for update-wiki and generate-ai-rules
- **Tests**: Both Go modules include test files (`*_test.go`)

### Navigation File Sync Requirement
`README.md`, `Home.md`, and `_Sidebar.md` all contain a table of contents. These **must always be kept in sync** (same labels and hierarchy). Differences in link targets are allowed. The `sync-docs.yaml` workflow enforces this on every pull request.

<!-- chlog:start -->
## Changelog (chlog) — MANDATORY

If the repository you are working in uses chlog (a `.chlog.yaml` or `.chlog.yml`
config file, or a `.changes/` directory, exists at the project root), the
following is binding and ALWAYS applies: whenever you make ANY change, you MUST
create a changelog fragment as part of the same change — automatically, without
being asked, before committing.

- Do NOT edit CHANGELOG.md directly; it is generated from fragments.
- Create the fragment with:
  `chlog new --kind <Kind> --body '<past-tense description>'`
- Write an apostrophe inside the single-quoted body as `'\''`.
- Valid kinds: Added, Changed, Deprecated, Removed, Fixed, Security
- Choose the kind that best matches the change (e.g., new feature → Added,
  bug fix → Fixed, behavior change → Changed, removal → Removed, security fix → Security).
- If the change is backward-INCOMPATIBLE with the public API (a breaking
  change), you MUST add the `--breaking` flag:
  `chlog new --kind <Kind> --breaking --body '<past-tense description>'`.
  This is the ONLY thing that triggers a major version bump — the kind alone
  never does (per SemVer, major = incompatible change). When unsure whether a
  change breaks compatibility, ask the user instead of guessing.
- Fragments are YAML files in `.changes/unreleased/`; stage them with your commit.
- `chlog check` fails the build when a fragment is missing — never skip it.
<!-- chlog:end -->
