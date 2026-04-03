#!/usr/bin/env bash
# Claude Code PreToolUse hook: blocks git commits that add CHANGELOG entries
# outside the [Unreleased] section.
#
# Input: JSON on stdin with { tool, input: { command } }
# Output: exit 2 + JSON on stderr to block with a message.

set -euo pipefail

if ! command -v jq >/dev/null 2>&1; then
  echo '{"decision":"block","reason":"Missing required dependency: jq. Install jq to enable the changelog-guard hook."}' >&2
  exit 2
fi

INPUT="$(cat)"
COMMAND="$(echo "$INPUT" | jq -r '.input.command // empty')"

# Only check Bash calls that contain "git commit"
if ! echo "$COMMAND" | grep -q "git commit"; then
  exit 0
fi

# Only run repository checks when we're inside a Git work tree.
# Fail open if Git context is unavailable.
if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  exit 0
fi

# Check if CHANGELOG.md is staged
if ! git diff --cached --name-only 2>/dev/null | grep -q "^CHANGELOG.md$"; then
  exit 0
fi

# Get full-context diff so section headers are always visible.
# Avoid aborting under set -e if git diff fails unexpectedly; fail open instead.
if ! DIFF="$(git diff --cached -U9999 -- CHANGELOG.md 2>/dev/null)"; then
  exit 0
fi

# Parse the diff to check if additions land in a released section.
# Strategy: walk the diff hunks. Track which section we're in by looking at
# context and added lines that match version headers. If we see an added line
# (starting with +) that isn't a header and we're inside a released section,
# that's a violation.
IN_RELEASED=false
VIOLATION_LINES=""

while IFS= read -r line; do
  # Skip diff metadata
  case "$line" in
    "---"*|"+++"*|"diff "*|"index "*) continue ;;
  esac

  # Detect section headers from both context and added lines
  clean="${line#[+ -]}"

  if echo "$clean" | grep -qE '^## \[Unreleased\]'; then
    IN_RELEASED=false
    continue
  fi

  if echo "$clean" | grep -qE '^## \[[0-9]+\.[0-9]+'; then
    IN_RELEASED=true
    continue
  fi

  # Check added lines (not headers) while in a released section
  if [ "$IN_RELEASED" = true ] && echo "$line" | grep -qE '^\+[^+]'; then
    content="${line#+}"
    # Skip blank lines and category headers (### Added, etc.)
    trimmed="$(echo "$content" | sed 's/^[[:space:]]*//')"
    if [ -n "$trimmed" ] && ! echo "$trimmed" | grep -qE '^### '; then
      VIOLATION_LINES="${VIOLATION_LINES}\n  ${trimmed}"
    fi
  fi
done <<< "$DIFF"

if [ -n "$VIOLATION_LINES" ]; then
  MSG="CHANGELOG entries are being added to an already-released version section. Released sections are immutable. Move these entries under [Unreleased]:${VIOLATION_LINES}"
  echo '{"decision": "block", "reason": "'"$(echo "$MSG" | sed 's/"/\\"/g' | tr '\n' ' ')"'"}' >&2
  exit 2
fi

exit 0
