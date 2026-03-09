#!/bin/sh
# install-rules.sh — Download and install AI assistant rule files from the guide repository.
#
# Usage:
#   ./install-rules.sh [target-dir]
#
# If target-dir is provided, rules are installed at project level:
#   <target>/.claude/rules/  <target>/.cursor/rules/  <target>/AGENTS.md  <target>/.codex/rules/
#
# If omitted, rules are installed globally:
#   ~/.claude/rules/  ~/.cursor/rules/  ~/AGENTS.md  ~/.codex/rules/

set -e

REPO="rios0rios0/guide"
BRANCH="generated"
BASE_URL="https://raw.githubusercontent.com/${REPO}/${BRANCH}"

RULE_NAMES="architecture bulk-operations ci-cd code-style design-patterns documentation git-flow golang java javascript python security testing yaml"

OVERWRITE_ALL=""
if [ "${1:-}" = "--force" ]; then
  OVERWRITE_ALL="yes"
  shift
fi

# Determine target directory
if [ -n "${1:-}" ]; then
  TARGET_DIR="$(cd "$1" 2>/dev/null && pwd)" || { echo "Error: directory '$1' does not exist."; exit 1; }
else
  TARGET_DIR="$HOME"
fi

# prompt_overwrite checks if file exists and prompts user for confirmation.
# Returns 0 if write should proceed, 1 if skipped.
prompt_overwrite() {
  file="$1"
  if [ ! -f "$file" ]; then
    return 0
  fi
  if [ "$OVERWRITE_ALL" = "yes" ]; then
    return 0
  fi
  printf "File %s already exists. Overwrite? [y/n/a] " "$file"
  read -r answer </dev/tty
  case "$answer" in
    y|Y) return 0 ;;
    a|A) OVERWRITE_ALL="yes"; return 0 ;;
    *)   return 1 ;;
  esac
}

# download_file fetches a URL to a local path.
download_file() {
  url="$1"
  dest="$2"
  dir="$(dirname "$dest")"
  mkdir -p "$dir"
  if prompt_overwrite "$dest"; then
    if curl -fsSL -o "$dest" "$url"; then
      echo "  Installed: $dest"
    else
      echo "  Warning: failed to download $url" >&2
    fi
  else
    echo "  Skipped: $dest"
  fi
}

echo "Installing AI assistant rules into: $TARGET_DIR"
echo ""

# Claude Code rules (.claude/rules/*.md)
echo "Claude Code rules:"
for name in $RULE_NAMES; do
  download_file "${BASE_URL}/.ai/claude/rules/${name}.md" "${TARGET_DIR}/.claude/rules/${name}.md"
done
echo ""

# Claude Code commands (.claude/commands/*.md)
COMMAND_NAMES="scaffold-go-project scaffold-frontend-project scaffold-python-package fix-guardrails sync-repos"
echo "Claude Code commands:"
for name in $COMMAND_NAMES; do
  download_file "${BASE_URL}/.ai/claude/commands/${name}.md" "${TARGET_DIR}/.claude/commands/${name}.md"
done
echo ""

# Claude Code agents (.claude/agents/*.md)
AGENT_NAMES="chezmoi changelog-enforcer code-reviewer git-workflow security-auditor bulk-operations"
echo "Claude Code agents:"
for name in $AGENT_NAMES; do
  download_file "${BASE_URL}/.ai/claude/agents/${name}.md" "${TARGET_DIR}/.claude/agents/${name}.md"
done
echo ""

# Cursor rules (.cursor/rules/*.mdc)
echo "Cursor rules:"
for name in $RULE_NAMES; do
  download_file "${BASE_URL}/.ai/cursor/rules/${name}.mdc" "${TARGET_DIR}/.cursor/rules/${name}.mdc"
done
echo ""

# Cursor skills (.cursor/skills/<name>/SKILL.md)
SKILL_NAMES="scaffold-go-project scaffold-frontend-project scaffold-python-package fix-guardrails"
echo "Cursor skills:"
for name in $SKILL_NAMES; do
  download_file "${BASE_URL}/.ai/cursor/skills/${name}/SKILL.md" "${TARGET_DIR}/.cursor/skills/${name}/SKILL.md"
done
echo ""

# Codex AGENTS.md
echo "Codex instructions:"
download_file "${BASE_URL}/.ai/codex/AGENTS.md" "${TARGET_DIR}/AGENTS.md"
echo ""

# Codex rules (.codex/rules/default.rules)
echo "Codex rules:"
download_file "${BASE_URL}/.ai/codex/rules/default.rules" "${TARGET_DIR}/.codex/rules/default.rules"
echo ""

echo "Done."
