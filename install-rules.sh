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

# checksum computes a SHA-256 hash for a file (portable across Linux and macOS).
checksum() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | cut -d ' ' -f1
  else
    shasum -a 256 "$1" | cut -d ' ' -f1
  fi
}

# file_size returns the size in bytes of a file (portable across Linux and macOS).
file_size() {
  wc -c < "$1" | tr -d ' '
}

# download_file fetches a URL to a local path, comparing checksums when the file already exists.
download_file() {
  url="$1"
  dest="$2"
  dir="$(dirname "$dest")"
  mkdir -p "$dir"

  # Download to a temporary file first
  tmp="$(mktemp "${TMPDIR:-/tmp}/install-rules.XXXXXX")"
  if ! curl -fsSL -o "$tmp" "$url"; then
    echo "  Warning: failed to download $url" >&2
    rm -f "$tmp"
    return
  fi

  # If the local file does not exist, install directly
  if [ ! -f "$dest" ]; then
    mv "$tmp" "$dest"
    chmod 644 "$dest"
    # Ensure hook scripts are executable
    case "$dest" in *.sh) chmod 755 "$dest" ;; esac
    echo "  Installed (new): $dest"
    return
  fi

  # Compare checksums
  local_hash="$(checksum "$dest")"
  remote_hash="$(checksum "$tmp")"

  if [ "$local_hash" = "$remote_hash" ]; then
    rm -f "$tmp"
    echo "  Up to date: $dest"
    return
  fi

  # Files differ -- show size comparison to help the user decide
  local_bytes="$(file_size "$dest")"
  remote_bytes="$(file_size "$tmp")"
  diff_bytes=$((remote_bytes - local_bytes))

  if [ "$diff_bytes" -gt 0 ]; then
    size_note="remote is ${diff_bytes} bytes larger (local: ${local_bytes}B, remote: ${remote_bytes}B) -- installing adds content"
  elif [ "$diff_bytes" -lt 0 ]; then
    abs_diff=$(( -diff_bytes ))
    size_note="local is ${abs_diff} bytes larger (local: ${local_bytes}B, remote: ${remote_bytes}B) -- installing removes content"
  else
    size_note="same size but different content (${local_bytes}B)"
  fi

  if [ "$OVERWRITE_ALL" = "yes" ]; then
    mv "$tmp" "$dest"
    chmod 644 "$dest"
    case "$dest" in *.sh) chmod 755 "$dest" ;; esac
    echo "  Updated: $dest ($size_note)"
    return
  fi

  printf "  Changed: %s\n" "$dest"
  printf "    %s\n" "$size_note"
  printf "    Overwrite? [y/n/a] "
  read -r answer </dev/tty
  case "$answer" in
    y|Y) mv "$tmp" "$dest"; chmod 644 "$dest"; case "$dest" in *.sh) chmod 755 "$dest" ;; esac; echo "  Updated: $dest" ;;
    a|A) OVERWRITE_ALL="yes"; mv "$tmp" "$dest"; chmod 644 "$dest"; case "$dest" in *.sh) chmod 755 "$dest" ;; esac; echo "  Updated: $dest" ;;
    *)   rm -f "$tmp"; echo "  Skipped: $dest" ;;
  esac
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

# Claude Code hooks (.claude/hooks/*)
echo "Claude Code hooks:"
for file in $(list_remote_entries "claude/hooks"); do
  download_file "${BASE_URL}/claude/hooks/${file}" "${TARGET_DIR}/.claude/hooks/${file}"
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
