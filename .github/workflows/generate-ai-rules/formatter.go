package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	logger "github.com/sirupsen/logrus"
)

const codexMaxSize = 32 * 1024 // 32 KiB

// writeClaude writes a rule file in Claude Code format to .claude/rules/<name>.md.
func writeClaude(outputDir string, group RuleGroup, content string) error {
	dir := filepath.Join(outputDir, ".claude", "rules")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory %s: %w", dir, err)
	}

	var body string
	if group.Globs != "" {
		body = fmt.Sprintf("---\npaths:\n  - \"%s\"\n---\n\n%s", group.Globs, content)
	} else {
		body = content
	}

	path := filepath.Join(dir, group.Name+".md")
	return os.WriteFile(path, []byte(body), 0644)
}

// writeCursor writes a rule file in Cursor format to .cursor/rules/<name>.mdc.
func writeCursor(outputDir string, group RuleGroup, content string) error {
	dir := filepath.Join(outputDir, ".cursor", "rules")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory %s: %w", dir, err)
	}

	var frontmatter string
	if group.Globs != "" {
		frontmatter = fmt.Sprintf("---\ndescription: \"%s\"\nglobs: \"%s\"\nalwaysApply: false\n---\n\n",
			group.Description, group.Globs)
	} else {
		frontmatter = fmt.Sprintf("---\ndescription: \"%s\"\nalwaysApply: true\n---\n\n",
			group.Description)
	}

	path := filepath.Join(dir, group.Name+".mdc")
	return os.WriteFile(path, []byte(frontmatter+content), 0644)
}

// writeCodex writes a single AGENTS.md file for Codex by concatenating all rule groups.
func writeCodex(outputDir string, groups []RuleGroup, contents []string) error {
	var sb strings.Builder
	for i := range groups {
		if i > 0 {
			sb.WriteString("\n---\n\n")
		}
		sb.WriteString(contents[i])
	}

	body := sb.String()
	if len(body) > codexMaxSize {
		logger.Warnf("AGENTS.md size (%d bytes) exceeds Codex limit of %d bytes", len(body), codexMaxSize)
	}

	path := filepath.Join(outputDir, "AGENTS.md")
	return os.WriteFile(path, []byte(body), 0644)
}

// formatClaudeFrontmatter returns the frontmatter string for a Claude rule file.
func formatClaudeFrontmatter(globs string) string {
	if globs == "" {
		return ""
	}
	return fmt.Sprintf("---\npaths:\n  - \"%s\"\n---\n\n", globs)
}

// formatCursorFrontmatter returns the frontmatter string for a Cursor rule file.
func formatCursorFrontmatter(description string, globs string) string {
	if globs != "" {
		return fmt.Sprintf("---\ndescription: \"%s\"\nglobs: \"%s\"\nalwaysApply: false\n---\n\n",
			description, globs)
	}
	return fmt.Sprintf("---\ndescription: \"%s\"\nalwaysApply: true\n---\n\n", description)
}
