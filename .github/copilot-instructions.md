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
- Run basic validation: `find . -name "*.md" | wc -l` (should show 41+ markdown files)

### Core Build Process
The repository uses a Go-based build system that syncs documentation to GitHub Wiki:
- **NEVER CANCEL BUILD COMMANDS** - Even though builds are fast (~1s), set timeouts of 30+ seconds
- Build location: `.github/workflows/update-wiki/`
- Build command: `go build -o update-wiki ./...`
- Expected build time: 1 second maximum
- Go module download: `go mod download` (takes ~0.3 seconds)
- Test command: `go test ./...` (reports no test files, completes in ~1 second)

### GitHub Actions Workflow
- File: `.github/workflows/update-wiki.yml`
- Triggers: Push to main branch or manual workflow dispatch
- Runtime: Uses Go 1.23.4+
- Process: Builds Go program, syncs content to GitHub Wiki
- **NEVER CANCEL**: Full workflow takes ~2-3 minutes including setup. Set timeout to 10+ minutes.

## Validation

### Always Validate Changes
Run this complete validation sequence after making any changes:

1. **Repository Structure Check**:
   ```bash
   ls -la README.md Home.md .editorconfig .github/workflows/update-wiki.yml
   ```
   All files must exist.

2. **Build Validation**:
   ```bash
   cd .github/workflows/update-wiki
   go mod verify
   go build -o update-wiki ./...
   ```
   NEVER CANCEL: Set timeout to 60+ seconds. Actual time: ~1 second.

3. **Content Validation**:
   ```bash
   # Validate markdown files exist and are readable
   find . -name "*.md" -exec head -1 {} \; > /dev/null
   # Expected: 41+ files, no errors
   ```

4. **Manual Validation Scenarios**:
   - **Navigation test**: Open `README.md` and verify all links in the Summary section point to existing files
   - **Structure test**: Verify key directories exist: `Code-Style/`, `Life-Cycle/`, `Cookbooks/`, `Agile-&-Culture/`
   - **Build test**: Ensure `go build` in `.github/workflows/update-wiki/` succeeds without errors

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
4. Always validate build process after changes

### Repository Navigation
```
/home/runner/work/guide/guide/
├── README.md              # Main repository overview (identical to Home.md)
├── Home.md                # Navigation hub with full summary
├── .editorconfig          # Formatting standards
├── .github/workflows/     # Build automation
├── Code-Style/            # 13 files: Language-specific guidelines
│   ├── GoLang/           # Go coding standards and testing
│   ├── JavaScript/       # JS standards and testing frameworks
│   ├── Python/           # Python PEP compliance and best practices
│   └── Java/             # Java conventions and patterns
├── Life-Cycle/           # 8 files: Development process guides
│   ├── Architecture/     # Backend and frontend design principles
│   ├── Git-Flow/         # Version control workflows
│   └── Security/         # Security best practices
├── Cookbooks/            # 10 files: Practical setup guides
│   ├── Tools-&-Setup/    # Development environment setup
│   └── OpenSearch-Dashboards/ # Specialized tooling guides
└── Agile-&-Culture/      # 1 file: PDCA methodology
```

### Frequently Accessed Locations
- **Main navigation**: `Home.md` or `README.md` (identical content)
- **Testing guidelines**: `Life-Cycle/Tests.md`
- **Code style references**: `Code-Style/<language>/` directories
- **Setup guides**: `Cookbooks/Tools-&-Setup/`
- **CI/CD information**: `Life-Cycle/CI-&-CD.md`

### Build System Details
```bash
# Repository root commands
ls -la                     # View all files and directories
find . -name "*.md"        # List all documentation files (41+)

# Build system commands (run from .github/workflows/update-wiki/)
go mod verify              # Verify dependencies (0.3s)
go mod download            # Download dependencies (0.3s) 
go build -o update-wiki ./... # Build binary (~1s)
go test ./...              # Run tests (1s, no test files)
go clean                   # Clean build artifacts (0.02s)
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
- **Primary content**: 41+ markdown files across 22+ directories
- **Build output**: GitHub Wiki synchronization
- **No traditional tests**: Content validation via structure checks
- **Dependencies**: Go 1.23.4+ for build system only

### Change Validation Workflow
Always complete this checklist when making changes:
1. ✅ Verify file structure with `ls` commands
2. ✅ Run `go build` in update-wiki directory  
3. ✅ Check markdown file count matches expected (~41 files)
4. ✅ Manually verify key navigation links work
5. ✅ Confirm .editorconfig compliance (2-space indents, LF endings)
6. ✅ Test that modified files render correctly as markdown

The repository is optimized for documentation maintenance and automated wiki synchronization rather than traditional software development workflows.