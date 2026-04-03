Detect the current branch's pull request, retrieve its failing CI workflow logs, fix the root causes, commit, and push. This command auto-detects the PR from the current Git branch, fetches check run statuses, downloads failure logs, classifies each failure, implements fixes, and pushes a single commit.

For commit conventions, refer to the Git Flow rule. For lint verification, refer to the CI/CD rule.

## Pre-flight

### Step 1 -- Verify GitHub CLI authentication

```bash
gh auth status
```

If this fails, report the error and stop. The `gh` CLI must be installed and authenticated.

### Step 2 -- Detect repository context

```bash
gh repo view --json nameWithOwner -q '.nameWithOwner'
```

Split the result into `{owner}` and `{repo}`.

### Step 3 -- Auto-detect PR from current branch

```bash
gh pr view --json number,headRefName -q '.number'
```

This returns the PR number associated with the current branch. If it fails (no PR exists for the current branch), report: "No pull request found for the current branch. Switch to a branch with an open PR or create one first." and stop.

Store the returned number as `{PR_NUMBER}`. Then ensure the local branch is up to date:

```bash
git pull --rebase
```

## Procedure

### Step 4 -- Fetch check run statuses

List all check runs for the PR and identify failures:

```bash
gh pr checks {PR_NUMBER}
```

Filter to checks with `fail` status. If no failing checks exist, report "All CI checks are passing on PR #{PR_NUMBER}." and stop.

### Step 5 -- Retrieve failure logs

For each failing check, retrieve the workflow run ID and job ID from the check URL, then download the failed job logs:

```bash
gh run view {RUN_ID} --job {JOB_ID} --log-failed
```

If the log output is very large, use `tail -100` to focus on the error section. Extract:
- The specific error messages and file/line references
- The tool or stage that failed (e.g., `golangci-lint`, `make test`, `basic-checks`)
- Whether the failure is in code you can fix or an infrastructure/flaky issue

### Step 6 -- Classify each failure

For each failing check, classify and plan the fix:

| Category | Examples | Action |
|----------|----------|--------|
| **Lint errors** | `golangci-lint`, `eslint`, `flake8`, `revive` | Read the referenced files, apply fixes |
| **Test failures** | `make test`, `go test`, `pytest` | Read failing test and source code, fix the bug or update the test |
| **Build errors** | `go build`, `npm build`, type errors | Read the referenced files, fix compilation/type errors |
| **CHANGELOG missing** | `basic-checks`, "CHANGELOG.md was NOT modified" | Add an entry under `[Unreleased]` in `CHANGELOG.md` |
| **Formatting** | `gofmt`, `prettier`, `black` | Run the formatter on the affected files |
| **SAST findings** | `semgrep`, `trivy`, `hadolint`, `gitleaks` | Read the finding, fix the security issue or add a justified suppression |
| **Infrastructure/flaky** | Network timeouts, runner issues, rate limits | Report to user -- cannot fix from code |

### Step 7 -- Implement fixes

For each fixable failure:

1. **Read the file(s)** referenced in the error log
2. **Apply the fix** using the Edit tool
3. **Track what was changed** for the commit message

Group related fixes together. Do not make unrelated changes.

### Step 8 -- Verify fixes locally

Run the relevant verification commands based on what failed:

- If lint failed: run `make lint` (or the project's lint command)
- If tests failed: run `make test` (or the project's test command)
- If build failed: run `make build` (or the project's build command)
- If SAST failed: run `make sast` (or the specific SAST tool)

If local verification still fails, iterate on the fix. If the failure cannot be resolved, report it to the user with the error details.

### Step 9 -- Stage changed files

```bash
git add <file1> <file2> ...
```

Stage only the specific files that were modified. Never use `git add -A` or `git add .`.

### Step 10 -- Commit

Follow the Git Flow commit convention:

```bash
git commit -m "$(cat <<'EOF'
fix(ci): resolved CI pipeline failures

- <summary of fix 1>
- <summary of fix 2>

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
EOF
)"
```

If only one fix was made, the bullet list in the body can be omitted. Group all CI fixes into a single commit.

### Step 11 -- Push

```bash
git push
```

The branch already tracks the remote. If push fails, report the error and stop.

## Error Handling

| Situation | Action |
|-----------|--------|
| `gh` CLI not available or not authenticated | Report error and stop. |
| No PR found for the current branch | Report "No pull request found for the current branch." and stop. |
| All checks passing | Report "All CI checks are passing on PR #{PR_NUMBER}." and stop. |
| Log download fails | Try alternative: `gh run view {RUN_ID} --log-failed`. If that also fails, report the error. |
| Infrastructure/flaky failure (not code-related) | Report to user: "The failure in `{check_name}` appears to be an infrastructure issue (network timeout, runner problem, etc.) and cannot be fixed from code. Consider re-running the workflow." |
| Fix introduces new failures | Iterate on the fix. If unresolvable after 2 attempts, report to user with details. |
| Multiple workflow runs exist | Use the most recent run. Identify it from `gh pr checks` output. |

## Progress Reporting

After processing all failures, print a summary:

```
## PR #{PR_NUMBER} CI Fix Summary

| Check | Failure | Action | Status |
|-------|---------|--------|--------|
| golangci-lint | gochecknoglobals in clear_history.go | Moved globals into function scope | Fixed |
| basic-checks | CHANGELOG.md not modified | Added [Unreleased] entry | Fixed |
| test:all | TestFoo assertion failure | Fixed expected value | Fixed |
| sast:semgrep | SQL injection finding | Cannot fix -- infrastructure issue | Reported |

Commit: <short-sha>
Checks fixed: X / Y
```
