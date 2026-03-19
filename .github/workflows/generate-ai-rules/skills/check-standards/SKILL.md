---
name: check-standards
description: Review current branch changes against coding standards before submitting a PR. Catches issues before reviewers do.
---

Review all changes on the current branch against the team's coding standards defined in `.claude/rules/` (already loaded into this conversation). This is meant to run BEFORE creating a PR to catch issues before reviewers do.

## Step 1: Get the Changes

- Detect the default branch with `git symbolic-ref refs/remotes/origin/HEAD` (strip `refs/remotes/origin/`). Fall back to `main` if unavailable.
- Get the current branch name: `git branch --show-current`
- List changed files: `git diff <default-branch>...HEAD --name-only`
- Get the full diff: `git diff <default-branch>...HEAD`

## Step 2: Review Each Changed File

For each changed file, read the full file content and check it against the standards in `.claude/rules/`. Apply the rules that are relevant to that file type (e.g., styling rules for style files, architecture rules for services, testing rules for test files, etc.).

## Step 3: Output Structured Review

Format the review as follows:

```
## Standards Check: [branch name]

### Critical Issues (Must Fix)
These WILL get flagged by reviewers:
- [file:line] Description of violation -- which rule it violates

### Warnings (Should Fix)
These will likely get flagged:
- [file:line] Description -- which rule it violates

### Suggestions (Nice to Have)
Optional improvements:
- [file:line] Description

### PR Quality Reminders
- [ ] Screenshots for all new/modified UI screens
- [ ] Evidence of all new features (modals, exports, charts, etc.)
- [ ] Target branch is correct
- [ ] Translations added to all locale files for new user-facing strings
- [ ] Linter passes with no errors

### Verdict: READY / NOT READY
Summary of findings.
```

## Severity Classification

**Critical** (reviewers will reject for this):
- Any violation marked as CRITICAL or IMPORTANT in the rules files
- Hardcoded values that should use tokens or configuration
- Disabled linter rules in production code
- Security vulnerabilities (XSS, injection, etc.)
- Interfaces or types defined in the wrong architectural layer

**Warning** (reviewers will likely comment):
- Magic numbers or hardcoded strings
- Missing error handling or loading state cleanup
- Styled components defined inline with component logic instead of a separate file

**Suggestion** (optional):
- Code organization improvements
- Pattern consistency with other modules
- Performance optimizations
