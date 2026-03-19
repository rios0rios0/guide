Detect the current branch's pull request and resolve its unresolved review comments. This command auto-detects the PR from the current Git branch, fetches unresolved review threads via GraphQL, classifies each comment, implements valid changes, replies with explanations, pushes, and resolves addressed threads.

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

### Step 4 -- Fetch unresolved review threads

Use the GitHub GraphQL API:

```bash
gh api graphql -f query='
  query($owner: String!, $repo: String!, $pr: Int!) {
    repository(owner: $owner, name: $repo) {
      pullRequest(number: $pr) {
        reviewThreads(first: 100) {
          nodes {
            id
            isResolved
            comments(first: 10) {
              nodes {
                databaseId
                body
                path
                line
                startLine
                author { login }
              }
            }
          }
        }
      }
    }
  }
' -f owner="{owner}" -f repo="{repo}" -F pr={PR_NUMBER}
```

Filter to threads where `isResolved` is `false`. If no unresolved threads exist, report "No unresolved review threads found on PR #{PR_NUMBER}." and stop.

### Step 5 -- Process each unresolved thread

Sort threads by file path (alphabetical) so changes to the same file are grouped. For each thread:

#### 5a. Understand the comment

Read the full conversation in the thread (all comments, not just the first). The last comment is usually the most relevant. Extract:
- The file path and line number(s) referenced
- The suggested change or concern
- Whether it contains a GitHub suggestion block (` ```suggestion `)

#### 5b. Read the local file

Read the referenced file at the specified line range. If the file no longer exists, mark the thread for resolution with the reply: "The file `<path>` referenced in this comment no longer exists in the current branch. This comment may be outdated."

#### 5c. Classify and act

Follow this decision tree **in order**:

| Priority | Condition | Action |
|----------|-----------|--------|
| 1 | Code already matches the suggestion | Reply "Already addressed in the current code." then resolve |
| 2 | Comment contains a ` ```suggestion ` block | Extract the suggested code and apply it directly with Edit then resolve |
| 3 | Comment clearly describes an actionable code change | Implement the change with Edit then resolve |
| 4 | Comment is a style preference with no standard backing it | Decline with technical reasoning |
| 5 | Comment is factually incorrect | Decline with explanation |
| 6 | Comment is ambiguous or could be interpreted multiple ways | Ask for clarification |

**Be conservative.** When in doubt, ask for clarification rather than making a wrong change.

#### 5d. Handle line drift

If the line numbers in the comment do not match the current file content, search the file for the code context mentioned in the comment. If found, apply the fix at the correct location. If not found, reply asking the reviewer to verify the line reference.

### Step 6 -- Stage changed files

```bash
git add <file1> <file2> ...
```

Stage only the specific files that were modified. Never use `git add -A` or `git add .`.

### Step 7 -- Verify changes

If a `Makefile` exists in the repo root, run:

```bash
make lint
```

If lint fails due to a change you made, attempt to fix the lint issue. If the lint failure cannot be resolved, revert the problematic change, update its thread classification to "declined" with the reason: "The suggested change introduces lint failures that cannot be resolved without broader refactoring."

### Step 8 -- Commit

Follow the Git Flow commit convention:

```bash
git commit -m "$(cat <<'EOF'
fix(pr-review): addressed PR review comments

- <summary of change 1>
- <summary of change 2>

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
EOF
)"
```

If only one change was made, the bullet list in the body can be omitted. Group all review-comment changes into a single commit.

### Step 9 -- Push

```bash
git push
```

The branch already tracks the remote. If push fails, report the error and **do not** reply to or resolve any threads -- the changes were not delivered.

### Step 10 -- Reply to each thread

For each processed thread, reply using the REST API:

```bash
gh api "repos/{owner}/{repo}/pulls/{PR_NUMBER}/comments" \
  --method POST \
  -f body="<reply>" \
  -F in_reply_to=<databaseId>
```

Where `<databaseId>` is the `databaseId` of the first comment in the thread.

### Step 11 -- Resolve addressed threads

For each thread that was addressed (change implemented or already resolved), resolve it:

```bash
gh api graphql -f query='
  mutation($threadId: ID!) {
    resolveReviewThread(input: {threadId: $threadId}) {
      thread { isResolved }
    }
  }
' -f threadId="<THREAD_NODE_ID>"
```

Do **not** resolve threads where you declined the suggestion or asked for clarification -- those need the reviewer's response.

## Reply Templates

| Outcome | Template |
|---------|----------|
| **Addressed** | "Addressed in `<short-sha>`. <brief description of what changed>." |
| **Already done** | "Already addressed in the current code. Resolving this thread." |
| **Declined** | "Considered this suggestion but keeping the current approach because <specific technical reason>. Happy to discuss further." |
| **Clarification** | "Could you clarify what you mean by <quoted part>? Specifically, <what is ambiguous and what options were considered>." |
| **File missing** | "The file `<path>` referenced in this comment no longer exists in the current branch. This comment may be outdated." |
| **Lint failure** | "The suggested change introduces lint failures. Keeping the current implementation. Details: <lint error summary>." |

Rules for replies:
- Keep replies concise (1-3 sentences)
- Reference the specific commit SHA when a change is made
- Never be dismissive -- acknowledge the reviewer's perspective
- When declining, always provide a concrete technical reason

## Error Handling

| Situation | Action |
|-----------|--------|
| `gh` CLI not available or not authenticated | Report error and stop. Requires `gh auth status` to succeed. |
| No PR found for the current branch | Report "No pull request found for the current branch." and stop. |
| File referenced in comment no longer exists | Reply with the "File missing" template and resolve the thread. |
| Line numbers do not match current file | Search for the code context in the file. If found, apply fix there. If not, reply asking reviewer to verify. |
| `git push` fails | Report the error. Do not reply to or resolve any threads. |
| `make lint` fails after changes | Fix lint issues if possible. If not, revert the change and decline with the "Lint failure" template. |
| No unresolved threads found | Report "No unresolved review threads found on PR #{PR_NUMBER}." and exit. |
| Thread has a multi-comment conversation | Read the entire thread. The last comment is usually the most relevant action item. |
| Comment references code outside the PR diff | Read the file anyway and apply the change if it makes sense. Note in the reply that the change is outside the original diff. |

## Progress Reporting

After processing all threads, print a summary:

```
## PR #<N> Review Resolution Summary

| Thread | File | Action | Status |
|--------|------|--------|--------|
| 1 | src/main.go:42 | Applied suggestion | Resolved |
| 2 | src/util.go:18 | Declined (style preference) | Open |
| 3 | README.md:5 | Asked for clarification | Open |

Commit: <short-sha>
Threads resolved: X / Y
```
