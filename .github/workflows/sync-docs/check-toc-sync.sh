#!/usr/bin/env bash
set -euo pipefail

# Normalizes TOC lines from stdin: keeps only markdown list items with links,
# strips link targets, leaving only indentation + "- [Label]".
normalize_toc() {
  grep -E '^\s*- \[' | sed 's/\](.*/]/'
}

# --- Extract from README.md ---
readme_toc=$(sed -n '/^## Summary$/,/^## References$/p' README.md \
  | head -n -1 \
  | tail -n +2 \
  | normalize_toc)

# --- Extract from Home.md ---
home_toc=$(sed -n '/^## Summary$/,/^## References$/p' Home.md \
  | head -n -1 \
  | tail -n +2 \
  | normalize_toc)

# --- Extract from _Sidebar.md ---
sidebar_toc=$(normalize_toc < _Sidebar.md)

# --- Compare all three pairwise ---
errors=0

if [ "$readme_toc" != "$home_toc" ]; then
  echo "ERROR: README.md and Home.md have divergent TOC structures."
  echo ""
  echo "--- README.md"
  echo "+++ Home.md"
  diff --unified <(echo "$readme_toc") <(echo "$home_toc") || true
  echo ""
  errors=$((errors + 1))
fi

if [ "$readme_toc" != "$sidebar_toc" ]; then
  echo "ERROR: README.md and _Sidebar.md have divergent TOC structures."
  echo ""
  echo "--- README.md"
  echo "+++ _Sidebar.md"
  diff --unified <(echo "$readme_toc") <(echo "$sidebar_toc") || true
  echo ""
  errors=$((errors + 1))
fi

if [ "$home_toc" != "$sidebar_toc" ]; then
  echo "ERROR: Home.md and _Sidebar.md have divergent TOC structures."
  echo ""
  echo "--- Home.md"
  echo "+++ _Sidebar.md"
  diff --unified <(echo "$home_toc") <(echo "$sidebar_toc") || true
  echo ""
  errors=$((errors + 1))
fi

if [ "$errors" -gt 0 ]; then
  echo "FAILED: TOC structures are out of sync across $errors file pair(s)."
  echo ""
  echo "The menu labels and hierarchy (indentation) must be identical in:"
  echo "  - README.md (## Summary section)"
  echo "  - Home.md (## Summary section)"
  echo "  - _Sidebar.md (entire file)"
  echo ""
  echo "Note: Link targets may differ (README/Home use path/File.md; Sidebar uses PageName)."
  echo "Only the labels in [brackets] and their indentation levels must match."
  exit 1
fi

echo "SUCCESS: All three TOC structures are in sync."
