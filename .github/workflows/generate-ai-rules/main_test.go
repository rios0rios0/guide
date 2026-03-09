package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	logger "github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	logger.SetLevel(logger.WarnLevel)
	os.Exit(m.Run())
}

func TestProcessGroup(t *testing.T) {
	// given
	tmpDir := t.TempDir()
	indexContent := `# Go

> **TL;DR:** Use Dig for DI.

## Overview

- [Conventions](GoLang/GoLang-Conventions.md)

## Go Proverbs

Clear is better than clever.

## References

- [Effective Go](https://go.dev/doc/effective_go)
`
	subContent := `# Go Conventions

> **TL;DR:** Use snake_case for files.

## File Naming

Use snake_case.

## References

- [Go Wiki](https://go.dev/wiki)
`
	os.MkdirAll(filepath.Join(tmpDir, "Code-Style", "GoLang"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "Code-Style", "GoLang.md"), []byte(indexContent), 0644)
	os.WriteFile(filepath.Join(tmpDir, "Code-Style", "GoLang", "GoLang-Conventions.md"), []byte(subContent), 0644)

	group := RuleGroup{
		Name: "golang",
		Sources: []string{
			"Code-Style/GoLang.md",
			"Code-Style/GoLang/GoLang-Conventions.md",
		},
		Globs: "**/*.go",
	}

	// when
	result, err := processGroup(tmpDir, group)

	// then
	if err != nil {
		t.Fatalf("processGroup() error: %v", err)
	}
	if !strings.Contains(result, "Go Proverbs") {
		t.Error("result should contain Go Proverbs section")
	}
	if !strings.Contains(result, "File Naming") {
		t.Error("result should contain File Naming section from sub-page")
	}
	if strings.Contains(result, "## References") {
		t.Error("result should not contain References section")
	}
	if strings.Contains(result, "Effective Go") {
		t.Error("result should not contain reference links")
	}
	if !strings.Contains(result, "---") {
		t.Error("result should contain separator between merged files")
	}
}

func TestEndToEnd(t *testing.T) {
	// given
	sourceDir := t.TempDir()
	outputDir := t.TempDir()

	// create minimal source files for a few rule groups
	writeTestFile(t, sourceDir, "Code-Style.md", "# Code Style\n\nNaming conventions.\n")
	writeTestFile(t, sourceDir, "Life-Cycle/Git-Flow.md", "# Git Flow\n\nUse feature branches.\n\n## References\n\n- [SemVer](https://semver.org)\n")
	writeTestFile(t, sourceDir, "Life-Cycle/Git-Flow/Merge-Guide.md", "# Merge Guide\n\nRebase before merge.\n")

	// override ruleGroups for this test by directly calling processGroup and write functions
	groups := []RuleGroup{
		{
			Name:        "code-style",
			Description: "General code style",
			Sources:     []string{"Code-Style.md"},
		},
		{
			Name:        "git-flow",
			Description: "Git workflow",
			Sources:     []string{"Life-Cycle/Git-Flow.md", "Life-Cycle/Git-Flow/Merge-Guide.md"},
		},
	}

	contents := make([]string, len(groups))
	for i, group := range groups {
		merged, err := processGroup(sourceDir, group)
		if err != nil {
			t.Fatalf("processGroup(%q) error: %v", group.Name, err)
		}
		contents[i] = merged
	}

	// when
	errCount := writeAllRules(outputDir, groups, contents)
	if errCount != 0 {
		t.Errorf("writeAllRules reported %d errors", errCount)
	}

	// then - verify Claude rules under .ai/claude/rules/
	claudeCodeStyle := filepath.Join(outputDir, ".ai", "claude", "rules", "code-style.md")
	assertFileExists(t, claudeCodeStyle)
	assertFileContains(t, claudeCodeStyle, "Naming conventions")
	assertFileNotContains(t, claudeCodeStyle, "---\npaths:")

	claudeGitFlow := filepath.Join(outputDir, ".ai", "claude", "rules", "git-flow.md")
	assertFileExists(t, claudeGitFlow)
	assertFileContains(t, claudeGitFlow, "feature branches")
	assertFileContains(t, claudeGitFlow, "Rebase before merge")
	assertFileNotContains(t, claudeGitFlow, "## References")

	// then - verify Cursor rules under .ai/cursor/rules/
	cursorCodeStyle := filepath.Join(outputDir, ".ai", "cursor", "rules", "code-style.mdc")
	assertFileExists(t, cursorCodeStyle)
	assertFileContains(t, cursorCodeStyle, "alwaysApply: true")

	cursorGitFlow := filepath.Join(outputDir, ".ai", "cursor", "rules", "git-flow.mdc")
	assertFileExists(t, cursorGitFlow)
	assertFileContains(t, cursorGitFlow, "alwaysApply: true")

	// then - verify Copilot instructions under .ai/copilot/instructions/
	copilotCodeStyle := filepath.Join(outputDir, ".ai", "copilot", "instructions", "code-style.instructions.md")
	assertFileExists(t, copilotCodeStyle)
	assertFileContains(t, copilotCodeStyle, "Naming conventions")
	assertFileNotContains(t, copilotCodeStyle, "applyTo:")

	copilotGitFlow := filepath.Join(outputDir, ".ai", "copilot", "instructions", "git-flow.instructions.md")
	assertFileExists(t, copilotGitFlow)
	assertFileContains(t, copilotGitFlow, "feature branches")
	assertFileContains(t, copilotGitFlow, "Rebase before merge")
	assertFileNotContains(t, copilotGitFlow, "## References")

	// then - verify Codex AGENTS.md under .ai/codex/
	agentsFile := filepath.Join(outputDir, ".ai", "codex", "AGENTS.md")
	assertFileExists(t, agentsFile)
	assertFileContains(t, agentsFile, "Naming conventions")
	assertFileContains(t, agentsFile, "feature branches")

	// then - verify Codex rules under .ai/codex/rules/
	codexRulesFile := filepath.Join(outputDir, ".ai", "codex", "rules", "default.rules")
	assertFileExists(t, codexRulesFile)
	assertFileContains(t, codexRulesFile, "prefix_rule(")
	assertFileContains(t, codexRulesFile, "golangci-lint")
	assertFileContains(t, codexRulesFile, "forbidden")
	assertFileContains(t, codexRulesFile, "prompt")
}

func writeTestFile(t *testing.T, baseDir, relPath, content string) {
	t.Helper()
	fullPath := filepath.Join(baseDir, relPath)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		t.Fatalf("creating directory: %v", err)
	}
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		t.Fatalf("writing test file: %v", err)
	}
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected file to exist: %s", path)
	}
}

func assertFileContains(t *testing.T, path, substr string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	if !strings.Contains(string(data), substr) {
		t.Errorf("%s should contain %q", path, substr)
	}
}

func assertFileNotContains(t *testing.T, path, substr string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	if strings.Contains(string(data), substr) {
		t.Errorf("%s should not contain %q", path, substr)
	}
}
