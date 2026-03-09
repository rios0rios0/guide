package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	logger "github.com/sirupsen/logrus"
)

const codexMaxSize = 32 * 1024 // 32 KiB

// writeClaude writes a rule file in Claude Code format to .ai/claude/rules/<name>.md.
func writeClaude(outputDir string, group RuleGroup, content string) error {
	dir := filepath.Join(outputDir, ".ai", "claude", "rules")
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
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		return err
	}
	logger.WithFields(logger.Fields{
		"path":  path,
		"bytes": len(body),
	}).Debug("wrote Claude rule file")
	return nil
}

// writeCursor writes a rule file in Cursor format to .ai/cursor/rules/<name>.mdc.
func writeCursor(outputDir string, group RuleGroup, content string) error {
	dir := filepath.Join(outputDir, ".ai", "cursor", "rules")
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

	body := frontmatter + content
	path := filepath.Join(dir, group.Name+".mdc")
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		return err
	}
	logger.WithFields(logger.Fields{
		"path":  path,
		"bytes": len(body),
	}).Debug("wrote Cursor rule file")
	return nil
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
		logger.WithFields(logger.Fields{
			"size_bytes":  len(body),
			"limit_bytes": codexMaxSize,
		}).Warn("AGENTS.md size exceeds Codex limit")
	}

	dir := filepath.Join(outputDir, ".ai", "codex")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory %s: %w", dir, err)
	}

	path := filepath.Join(dir, "AGENTS.md")
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		return err
	}
	logger.WithFields(logger.Fields{
		"path":  path,
		"bytes": len(body),
	}).Debug("wrote Codex AGENTS.md")
	return nil
}

// CodexRule represents a single prefix_rule entry for Codex command execution policies.
type CodexRule struct {
	Pattern       []string // command prefix to match
	Decision      string   // "allow", "prompt", or "forbidden"
	Justification string   // human-readable reason
}

// codexRules returns the list of Codex command execution policy rules derived from
// CI/CD, security, and Git Flow guidelines.
func codexRules() []CodexRule {
	return []CodexRule{
		// Enforce Makefile targets (CI/CD: "Never call tool binaries directly")
		{Pattern: []string{"make", "lint"}, Decision: "allow", Justification: "Linting through Makefile is the approved approach"},
		{Pattern: []string{"make", "test"}, Decision: "allow", Justification: "Testing through Makefile is the approved approach"},
		{Pattern: []string{"make", "sast"}, Decision: "allow", Justification: "SAST through Makefile is the approved approach"},
		{Pattern: []string{"make", "semgrep"}, Decision: "allow", Justification: "SAST through Makefile is the approved approach"},
		{Pattern: []string{"make", "trivy"}, Decision: "allow", Justification: "SAST through Makefile is the approved approach"},
		{Pattern: []string{"make", "hadolint"}, Decision: "allow", Justification: "SAST through Makefile is the approved approach"},
		{Pattern: []string{"make", "gitleaks"}, Decision: "allow", Justification: "SAST through Makefile is the approved approach"},

		// Forbid direct linter/test runner invocation — Go
		{Pattern: []string{"golangci-lint"}, Decision: "forbidden", Justification: "Do not call golangci-lint directly; use `make lint` which loads the correct configuration through the pipelines repository scripts"},

		// Forbid direct linter/test runner invocation — Python
		{Pattern: []string{"pytest"}, Decision: "forbidden", Justification: "Do not call pytest directly; use `make test` which loads the correct configuration through the pipelines repository scripts"},
		{Pattern: []string{"black"}, Decision: "forbidden", Justification: "Do not call black directly; use `make lint` which loads the correct configuration through the pipelines repository scripts"},
		{Pattern: []string{"ruff"}, Decision: "forbidden", Justification: "Do not call ruff directly; use `make lint` which loads the correct configuration through the pipelines repository scripts"},

		// Forbid direct linter/test runner invocation — JavaScript/TypeScript
		{Pattern: []string{"eslint"}, Decision: "forbidden", Justification: "Do not call eslint directly; use `make lint` which loads the correct configuration through the pipelines repository scripts"},
		{Pattern: []string{"prettier"}, Decision: "forbidden", Justification: "Do not call prettier directly; use `make lint` which loads the correct configuration through the pipelines repository scripts"},
		{Pattern: []string{"jest"}, Decision: "forbidden", Justification: "Do not call jest directly; use `make test` which loads the correct configuration through the pipelines repository scripts"},

		// Forbid direct linter/test runner invocation — Java
		{Pattern: []string{"checkstyle"}, Decision: "forbidden", Justification: "Do not call checkstyle directly; use `make lint` which loads the correct configuration through the pipelines repository scripts"},

		// SAST tools — must go through their respective Makefile targets
		{Pattern: []string{"semgrep"}, Decision: "forbidden", Justification: "Do not call semgrep directly; use `make semgrep` which loads the correct configuration through the pipelines repository scripts"},
		{Pattern: []string{"trivy"}, Decision: "forbidden", Justification: "Do not call trivy directly; use `make trivy` which loads the correct configuration through the pipelines repository scripts"},
		{Pattern: []string{"gitleaks"}, Decision: "forbidden", Justification: "Do not call gitleaks directly; use `make gitleaks` which loads the correct configuration through the pipelines repository scripts"},
		{Pattern: []string{"hadolint"}, Decision: "forbidden", Justification: "Do not call hadolint directly; use `make hadolint` which loads the correct configuration through the pipelines repository scripts"},

		// Git safety (Git Flow: force-push requires caution)
		{Pattern: []string{"git", "push", "--force"}, Decision: "prompt", Justification: "Force pushing rewrites remote history. Confirm this is intentional."},
		{Pattern: []string{"git", "push", "-f"}, Decision: "prompt", Justification: "Force pushing rewrites remote history. Confirm this is intentional."},
	}
}

// formatCodexRules generates the Starlark content for a Codex .rules file.
func formatCodexRules(rules []CodexRule) string {
	var sb strings.Builder
	sb.WriteString("# Codex command execution policies\n")
	sb.WriteString("# Generated from the development guide — do not edit manually.\n")
	sb.WriteString("# See: https://developers.openai.com/codex/rules/\n\n")

	for i, rule := range rules {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString("prefix_rule(\n")
		sb.WriteString(fmt.Sprintf("    pattern = %s,\n", formatStarlarkList(rule.Pattern)))
		sb.WriteString(fmt.Sprintf("    decision = %q,\n", rule.Decision))
		sb.WriteString(fmt.Sprintf("    justification = %q,\n", rule.Justification))
		sb.WriteString(")\n")
	}
	return sb.String()
}

// formatStarlarkList formats a string slice as a Starlark list literal.
func formatStarlarkList(items []string) string {
	quoted := make([]string, len(items))
	for i, item := range items {
		quoted[i] = fmt.Sprintf("%q", item)
	}
	return "[" + strings.Join(quoted, ", ") + "]"
}

// writeCodexRules writes the Codex command execution policy file to .ai/codex/rules/default.rules.
func writeCodexRules(outputDir string) error {
	dir := filepath.Join(outputDir, ".ai", "codex", "rules")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory %s: %w", dir, err)
	}

	body := formatCodexRules(codexRules())
	path := filepath.Join(dir, "default.rules")
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		return err
	}
	logger.WithFields(logger.Fields{
		"path":  path,
		"bytes": len(body),
	}).Debug("wrote Codex rules file")
	return nil
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
