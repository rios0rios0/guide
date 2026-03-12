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
API_URL="https://api.github.com/repos/${REPO}/contents"

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

# list_remote_entries fetches entry names from a directory on the generated branch via the GitHub API.
list_remote_entries() {
  path="$1"
  response=$(curl -fsSL "${API_URL}/${path}?ref=${BRANCH}" 2>/dev/null) || {
    echo "Error: failed to list files at ${path}" >&2
    return 1
  }
  echo "$response" \
    | grep -o '"name": *"[^"]*"' \
    | sed 's/"name": *"//;s/"//'
}

echo "Installing AI assistant rules into: $TARGET_DIR"
echo ""

# Claude Code rules (.claude/rules/*.md)
echo "Claude Code rules:"
for file in $(list_remote_entries "claude/rules"); do
  download_file "${BASE_URL}/claude/rules/${file}" "${TARGET_DIR}/.claude/rules/${file}"
done
echo ""

# Claude Code commands (.claude/commands/*.md)
echo "Claude Code commands:"
for file in $(list_remote_entries "claude/commands"); do
  download_file "${BASE_URL}/claude/commands/${file}" "${TARGET_DIR}/.claude/commands/${file}"
done
echo ""

# Claude Code agents (.claude/agents/*.md)
echo "Claude Code agents:"
for file in $(list_remote_entries "claude/agents"); do
  download_file "${BASE_URL}/claude/agents/${file}" "${TARGET_DIR}/.claude/agents/${file}"
done
echo ""

# Cursor rules (.cursor/rules/*.mdc)
echo "Cursor rules:"
for file in $(list_remote_entries "cursor/rules"); do
  download_file "${BASE_URL}/cursor/rules/${file}" "${TARGET_DIR}/.cursor/rules/${file}"
done
echo ""

# Cursor skills (.cursor/skills/<name>/SKILL.md)
echo "Cursor skills:"
for skill_dir in $(list_remote_entries "cursor/skills"); do
  download_file "${BASE_URL}/cursor/skills/${skill_dir}/SKILL.md" "${TARGET_DIR}/.cursor/skills/${skill_dir}/SKILL.md"
done
echo ""

# GitHub Copilot instructions (.github/instructions/*.instructions.md)
echo "GitHub Copilot instructions:"
for file in $(list_remote_entries "copilot/instructions"); do
  download_file "${BASE_URL}/copilot/instructions/${file}" "${TARGET_DIR}/.github/instructions/${file}"
done
echo ""

# Codex AGENTS.md
echo "Codex instructions:"
download_file "${BASE_URL}/codex/AGENTS.md" "${TARGET_DIR}/AGENTS.md"
echo ""

# Codex rules (.codex/rules/default.rules)
echo "Codex rules:"
download_file "${BASE_URL}/codex/rules/default.rules" "${TARGET_DIR}/.codex/rules/default.rules"
echo ""

echo "Done."
