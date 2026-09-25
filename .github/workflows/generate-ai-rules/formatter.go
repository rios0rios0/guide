package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	logger "github.com/sirupsen/logrus"
)

// writeClaude writes a rule file in Claude Code format to claude/rules/<name>.md.
func writeClaude(outputDir string, group RuleGroup, content string) error {
	dir := filepath.Join(outputDir, "claude", "rules")
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

// writeCursor writes a rule file in Cursor format to cursor/rules/<name>.mdc.
func writeCursor(outputDir string, group RuleGroup, content string) error {
	dir := filepath.Join(outputDir, "cursor", "rules")
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

// writeCodex writes a small routing file and complete, separately loaded standards.
func writeCodex(outputDir string, groups []RuleGroup, contents []string) error {
	dir := filepath.Join(outputDir, "codex", "instructions")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating instruction directory: %w", err)
	}
	var sb strings.Builder
	sb.WriteString("# Personal engineering standards\n\n")
	sb.WriteString("Before project changes, read documentation and git-flow; before code changes also read architecture, code-style, security, testing, and the relevant language. Read ci-cd before validation or Git operations. Resolve the paths below relative to this AGENTS.md file.\n\n")
	sb.WriteString("New Go projects use manual constructor injection in container.go. Existing Wire migrations belong in dedicated service PRs; encourage Dig migrations. Preserve project-specific loggers. Unit tests are untagged; mocking libraries require the narrow documented external-abstraction exception in testing.md.\n\n")
	sb.WriteString("Invoke Codex skills with $skill-name. Claude slash commands and CLAUDE.md are tool-specific; Codex uses AGENTS.md. Project instructions refine these defaults. Preserve existing authorization safeguards.\n\n")
	sb.WriteString("| Applies to | Full standard |\n| --- | --- |\n")
	for i, group := range groups {
		if contents[i] == "" {
			continue
		}
		name := group.Name + ".md"
		content := strings.ReplaceAll(contents[i], "CLAUDE.md", "AGENTS.md")
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
			return fmt.Errorf("writing instruction %s: %w", name, err)
		}
		scope := group.Globs
		if scope == "" {
			scope = group.Name
		}
		fmt.Fprintf(&sb, "| %s | [%s](instructions/%s) |\n", scope, group.Name, name)
	}
	return os.WriteFile(filepath.Join(outputDir, "codex", "AGENTS.md"), []byte(sb.String()), 0644)
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

// writeCodexRules writes the Codex command execution policy file to codex/rules/default.rules.
func writeCodexRules(outputDir string) error {
	dir := filepath.Join(outputDir, "codex", "rules")
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

// writeCopilot writes a rule file in GitHub Copilot format to copilot/instructions/<name>.instructions.md.
func writeCopilot(outputDir string, group RuleGroup, content string) error {
	dir := filepath.Join(outputDir, "copilot", "instructions")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory %s: %w", dir, err)
	}

	body := formatCopilotFrontmatter(group.Globs) + content

	path := filepath.Join(dir, group.Name+".instructions.md")
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		return err
	}
	logger.WithFields(logger.Fields{
		"path":  path,
		"bytes": len(body),
	}).Debug("wrote Copilot instruction file")
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

// formatCopilotFrontmatter returns the frontmatter string for a GitHub Copilot instruction file.
func formatCopilotFrontmatter(globs string) string {
	if globs == "" {
		return ""
	}
	return fmt.Sprintf("---\napplyTo: \"%s\"\n---\n\n", globs)
}
