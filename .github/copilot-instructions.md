# Software Development Guide Repository

This is a comprehensive documentation repository containing software development standards, best practices, and guidelines organized across multiple categories including Code Style, Life Cycle processes, and practical Cookbooks.

Always reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.

## Working Effectively

### Quick Setup and Validation
- Validate repository structure: `ls -la README.md Home.md .editorconfig .github/workflows/update-wiki.yml`
- Build the wiki sync tool:
  ```bash
  cd .github/workflows/update-wiki
  go build -o update-wiki ./...
  ```
  - Build time: ~1 second. NEVER CANCEL. Set timeout to 30+ seconds for safety.
- Build the AI rules generator:
  ```bash
  cd .github/workflows/generate-ai-rules
  go build -o generate-ai-rules ./...
  ```
- Run basic validation: `find . -name "*.md" | wc -l` (should show 79+ Markdown files)

### Core Build Process
The repository has two Go-based tools under `.github/workflows/`:

#### update-wiki (Go 1.26.0)
Syncs documentation to GitHub Wiki:
- Build location: `.github/workflows/update-wiki/`
- Build command: `go build -o update-wiki ./...`
- Expected build time: ~1 second

#### generate-ai-rules (Go 1.24.7)
Generates AI assistant rule files for Claude Code, Cursor, Codex, and GitHub Copilot from the guide's own documentation.
- Build location: `.github/workflows/generate-ai-rules/`
- Build command: `go build -o generate-ai-rules ./...`
- Run: `./generate-ai-rules`
- Expected build time: ~1 second

**NEVER CANCEL BUILD COMMANDS** - Even though builds are fast (~1s), set timeouts of 30+ seconds.

### GitHub Actions Workflows
Three workflows are defined under `.github/workflows/`:

| File | Trigger | Purpose |
|------|---------|---------|
| `update-wiki.yml` | Push to `main`, manual dispatch | Syncs docs to GitHub Wiki |
| `generate-ai-rules.yaml` | Push to `main` (docs paths), manual dispatch | Regenerates AI rules on `generated` branch |
| `sync-docs.yaml` | Pull request (any `.md` change) | Validates TOC sync across README.md, Home.md, _Sidebar.md |

**NEVER CANCEL**: Full workflow takes ~2-3 minutes including setup. Set timeout to 10+ minutes.

## Validation

### Always Validate Changes
Run this complete validation sequence after making any changes:

1. **Repository Structure Check**:
   ```bash
   ls -la README.md Home.md .editorconfig .github/workflows/update-wiki.yml
   ```
   All files must exist.

2. **Build Validation** (update-wiki):
   ```bash
   cd .github/workflows/update-wiki
   go mod verify
   go build -o update-wiki ./...
   ```
   NEVER CANCEL: Set timeout to 60+ seconds. Actual time: ~1 second.

3. **Build Validation** (generate-ai-rules):
   ```bash
   cd .github/workflows/generate-ai-rules
   go mod verify
   go build -o generate-ai-rules ./...
   ```

4. **TOC Sync Validation**:
   ```bash
   bash .github/workflows/sync-docs/check-toc-sync.sh
   ```
   The Summary sections in `README.md`, `Home.md`, and `_Sidebar.md` must stay in sync.
   Run this whenever you modify any of those three navigation files.

5. **Content Validation**:
   ```bash
   # Validate markdown files exist and are readable
   find . -name "*.md" -exec head -1 {} \; > /dev/null
   # Expected: 79+ files, no errors
   ```

6. **Manual Validation Scenarios**:
   - **Navigation test**: Open `README.md` and verify all links in the Summary section point to existing files
   - **Structure test**: Verify key directories exist: `Code-Style/`, `Life-Cycle/`, `Cookbooks/`, `Agile-&-Culture/`
   - **Build test**: Ensure `go build` succeeds in both workflow tool directories

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

### Repository Navigation
```
/
├── README.md                         # GitHub repo landing page (badges + Summary)
├── Home.md                           # GitHub Wiki navigation hub (Summary + Context)
├── Onboarding.md                     # Developer onboarding guide
├── _Sidebar.md                       # GitHub Wiki sidebar navigation
├── _Footer.md                        # GitHub Wiki footer
├── .claude-plugin/marketplace.json   # Claude Code plugin marketplace definition
├── .editorconfig                     # Formatting standards
├── .github/workflows/                # CI/CD automation
│   ├── update-wiki.yml               # Syncs docs to GitHub Wiki (Go 1.26.0)
│   ├── generate-ai-rules.yaml        # Generates AI rules on 'generated' branch (Go 1.24.7)
│   ├── generate-ai-rules/            # Go tool + static assets
│   │   ├── agents/                   # Claude Code agent source files (6 static agents)
│   │   ├── commands/                 # Claude Code command source files (5 commands)
│   │   ├── skills/                   # Cursor skill source files (4 skills)
│   │   ├── └── *.go                      # Go tool (config, parser, formatter)
│   └── sync-docs.yaml                # Validates TOC sync on PRs
├── Agile-&-Culture/                  # Agile methodology guides
│   └── PDCA.md                       # PDCA cycle methodology
├── Life-Cycle/                       # Development process guides
│   ├── Architecture/                 # Backend and frontend design principles
│   ├── Git-Flow/                     # Version control workflows
│   ├── Documentation-&-Change-Control/ # Templates and change control
│   └── *.md                          # Tests, CI-&-CD, Security, Git-Flow, Architecture
├── Code-Style/                       # Language-specific coding guidelines
│   ├── GoLang/                       # Go: conventions, formatting, types, logging, testing, structure
│   ├── Java/                         # Java: conventions, formatting, types, logging, testing, structure
│   ├── JavaScript/                   # JavaScript: testing
│   ├── Python/                       # Python: conventions, formatting, types, logging, testing, structure
│   └── *.md                          # YAML, Language-Guide-Template, language index files
└── Cookbooks/                        # Practical setup and technique guides
    └── Tools-&-Setup/                # WSL, Azure CLI, Azure Functions, Database Sync
```

### Frequently Accessed Locations
- **GitHub repo page**: `README.md` (includes badges)
- **Wiki navigation**: `Home.md` (full navigation tree for GitHub Wiki)
- **Testing guidelines**: `Life-Cycle/Tests.md`
- **Code style references**: `Code-Style/<language>/` directories
- **Setup guides**: `Cookbooks/Tools-&-Setup/`
- **CI/CD information**: `Life-Cycle/CI-&-CD.md`
- **AI rules**: `generated` branch (auto-generated from docs + external sources)
- **AI rule sources**: `.github/workflows/generate-ai-rules/agents/`, `commands/`, `skills/`, `hooks/`

### Build System Details
```bash
# Repository root commands
ls -la                     # View all files and directories
find . -name "*.md"        # List all documentation files (79+)

# update-wiki build (run from .github/workflows/update-wiki/)
go mod verify              # Verify dependencies (0.3s)
go mod download            # Download dependencies (0.3s)
go build -o update-wiki ./... # Build binary (~1s)
go test ./...              # Run tests (~1s)
go clean                   # Clean build artifacts (0.02s)

# generate-ai-rules build (run from .github/workflows/generate-ai-rules/)
go mod verify              # Verify dependencies (0.3s)
go mod download            # Download dependencies (0.3s)
go build -o generate-ai-rules ./... # Build binary (~1s)
go test ./...              # Run tests (~1s)
go clean                   # Clean build artifacts (0.02s)

# TOC sync check (run from repository root)
bash .github/workflows/sync-docs/check-toc-sync.sh
```

## Important Notes

### Critical Timing Expectations
- **Go module operations**: 0.3 seconds typical, set 30+ second timeout
- **Go build process**: 1 second typical, set 60+ second timeout
- **Full validation**: 1-2 seconds typical, set 30+ second timeout
- **GitHub Actions workflow**: 2-3 minutes total, set 10+ minute timeout
- **NEVER CANCEL ANY BUILD OPERATION** - Always wait for completion

### Repository Characteristics
- **Type**: Documentation repository (not traditional software)
- **Primary content**: 79+ Markdown files across 25+ directories
- **Build output**: GitHub Wiki synchronization + AI rule files for Claude Code, Cursor, Codex, and GitHub Copilot (on `generated` branch)
- **Dependencies**: Go 1.26.0 for update-wiki; Go 1.24.7 for generate-ai-rules (+ gopkg.in/yaml.v3 for external sources)
- **Tests**: Both Go modules include test files (`*_test.go`)

### Navigation File Sync Requirement
`README.md`, `Home.md`, and `_Sidebar.md` all contain a table of contents. These **must always be kept in sync** (same labels and hierarchy). Differences in link targets are allowed. The `sync-docs.yaml` workflow enforces this on every pull request.

### Change Validation Workflow
Always complete this checklist when making changes:
1. ✅ Verify file structure with `ls` commands
2. ✅ Run `go build` in both workflow tool directories
3. ✅ Run `bash .github/workflows/sync-docs/check-toc-sync.sh` after any TOC change
4. ✅ Check Markdown file count matches expected (~79 files)
5. ✅ Manually verify key navigation links work
6. ✅ Confirm .editorconfig compliance (2-space indents, LF endings)
7. ✅ Test that modified files render correctly as Markdown

The repository is optimized for documentation maintenance and automated wiki synchronization rather than traditional software development workflows.
