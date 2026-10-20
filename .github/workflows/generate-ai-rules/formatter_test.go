package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestFormatClaudeFrontmatter(t *testing.T) {
	tests := []struct {
		name     string
		globs    string
		expected string
	}{
		{
			name:     "with globs",
			globs:    "**/*.go",
			expected: "---\npaths:\n  - \"**/*.go\"\n---\n\n",
		},
		{
			name:     "without globs",
			globs:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			globs := tt.globs

			// when
			result := formatClaudeFrontmatter(globs)

			// then
			if result != tt.expected {
				t.Errorf("formatClaudeFrontmatter(%q)\n  got:  %q\n  want: %q", globs, result, tt.expected)
			}
		})
	}
}

func TestFormatCursorFrontmatter(t *testing.T) {
	tests := []struct {
		name        string
		description string
		globs       string
		expected    string
	}{
		{
			name:        "language-specific with globs",
			description: "Go language coding standards",
			globs:       "**/*.go",
			expected:    "---\ndescription: \"Go language coding standards\"\nglobs: \"**/*.go\"\nalwaysApply: false\n---\n\n",
		},
		{
			name:        "cross-cutting without globs",
			description: "General code style conventions",
			globs:       "",
			expected:    "---\ndescription: \"General code style conventions\"\nalwaysApply: true\n---\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			description := tt.description
			globs := tt.globs

			// when
			result := formatCursorFrontmatter(description, globs)

			// then
			if result != tt.expected {
				t.Errorf("formatCursorFrontmatter(%q, %q)\n  got:  %q\n  want: %q", description, globs, result, tt.expected)
			}
		})
	}
}

func TestWriteClaude(t *testing.T) {
	tests := []struct {
		name          string
		group         RuleGroup
		content       string
		expectGlobs   bool
		expectContent string
	}{
		{
			name: "language rule with globs frontmatter",
			group: RuleGroup{
				Name:  "golang",
				Globs: "**/*.go",
			},
			content:       "# Go Standards\n\nUse gofmt.\n",
			expectGlobs:   true,
			expectContent: "---\npaths:\n  - \"**/*.go\"\n---\n\n# Go Standards\n\nUse gofmt.\n",
		},
		{
			name: "cross-cutting rule without frontmatter",
			group: RuleGroup{
				Name: "code-style",
			},
			content:       "# Code Style\n\nNaming conventions.\n",
			expectContent: "# Code Style\n\nNaming conventions.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			tmpDir := t.TempDir()

			// when
			err := writeClaude(tmpDir, tt.group, tt.content)

			// then
			if err != nil {
				t.Fatalf("writeClaude() error: %v", err)
			}
			path := filepath.Join(tmpDir, "claude", "rules", tt.group.Name+".md")
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("reading output file: %v", err)
			}
			if string(data) != tt.expectContent {
				t.Errorf("file content\n  got:  %q\n  want: %q", string(data), tt.expectContent)
			}
		})
	}
}

func TestWriteCursor(t *testing.T) {
	tests := []struct {
		name         string
		group        RuleGroup
		content      string
		expectPrefix string
	}{
		{
			name: "language rule with globs",
			group: RuleGroup{
				Name:        "golang",
				Description: "Go coding standards",
				Globs:       "**/*.go",
			},
			content:      "# Go\n",
			expectPrefix: "---\ndescription: \"Go coding standards\"\nglobs: \"**/*.go\"\nalwaysApply: false\n---\n\n",
		},
		{
			name: "cross-cutting rule always apply",
			group: RuleGroup{
				Name:        "testing",
				Description: "Testing standards",
			},
			content:      "# Testing\n",
			expectPrefix: "---\ndescription: \"Testing standards\"\nalwaysApply: true\n---\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			tmpDir := t.TempDir()

			// when
			err := writeCursor(tmpDir, tt.group, tt.content)

			// then
			if err != nil {
				t.Fatalf("writeCursor() error: %v", err)
			}
			path := filepath.Join(tmpDir, "cursor", "rules", tt.group.Name+".mdc")
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("reading output file: %v", err)
			}
			if !strings.HasPrefix(string(data), tt.expectPrefix) {
				t.Errorf("file should start with frontmatter\n  got:  %q\n  want prefix: %q", string(data), tt.expectPrefix)
			}
		})
	}
}

func TestFormatCopilotFrontmatter(t *testing.T) {
	tests := []struct {
		name     string
		globs    string
		expected string
	}{
		{
			name:     "with globs",
			globs:    "**/*.go",
			expected: "---\napplyTo: \"**/*.go\"\n---\n\n",
		},
		{
			name:     "without globs",
			globs:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			globs := tt.globs

			// when
			result := formatCopilotFrontmatter(globs)

			// then
			if result != tt.expected {
				t.Errorf("formatCopilotFrontmatter(%q)\n  got:  %q\n  want: %q", globs, result, tt.expected)
			}
		})
	}
}

func TestWriteCopilot(t *testing.T) {
	tests := []struct {
		name          string
		group         RuleGroup
		content       string
		expectContent string
	}{
		{
			name: "language rule with applyTo frontmatter",
			group: RuleGroup{
				Name:  "golang",
				Globs: "**/*.go",
			},
			content:       "# Go Standards\n\nUse gofmt.\n",
			expectContent: "---\napplyTo: \"**/*.go\"\n---\n\n# Go Standards\n\nUse gofmt.\n",
		},
		{
			name: "cross-cutting rule without frontmatter",
			group: RuleGroup{
				Name: "code-style",
			},
			content:       "# Code Style\n\nNaming conventions.\n",
			expectContent: "# Code Style\n\nNaming conventions.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			tmpDir := t.TempDir()

			// when
			err := writeCopilot(tmpDir, tt.group, tt.content)

			// then
			if err != nil {
				t.Fatalf("writeCopilot() error: %v", err)
			}
			path := filepath.Join(tmpDir, "copilot", "instructions", tt.group.Name+".instructions.md")
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("reading output file: %v", err)
			}
			if string(data) != tt.expectContent {
				t.Errorf("file content\n  got:  %q\n  want: %q", string(data), tt.expectContent)
			}
		})
	}
}

func TestWriteCodex(t *testing.T) {
	// given
	tmpDir := t.TempDir()
	groups := []RuleGroup{
		{Name: "code-style", Description: "Code style"},
		{Name: "golang", Description: "Go standards"},
	}
	contents := []string{
		"# Code Style\n\nNaming conventions.\n",
		"# Go\n\nUse gofmt.\n",
	}

	// when
	err := writeCodex(tmpDir, groups, contents)

	// then
	if err != nil {
		t.Fatalf("writeCodex() error: %v", err)
	}
	path := filepath.Join(tmpDir, "codex", "AGENTS.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading AGENTS.md: %v", err)
	}
	content := string(data)
	for i, group := range groups {
		if !strings.Contains(content, "instructions/"+group.Name+".md") {
			t.Errorf("missing instruction route for %s", group.Name)
		}
		full, err := os.ReadFile(filepath.Join(tmpDir, "codex", "instructions", group.Name+".md"))
		if err != nil || string(full) != contents[i] {
			t.Errorf("incomplete instruction %s: %v", group.Name, err)
		}
	}
	if len(data) >= 32*1024 {
		t.Fatal("routing file exceeds default instruction budget")
	}
}

func TestFormatCodexRules(t *testing.T) {
	// given
	rules := []CodexRule{
		{Pattern: []string{"make", "lint"}, Decision: "allow", Justification: "Use Makefile"},
		{Pattern: []string{"golangci-lint"}, Decision: "forbidden", Justification: "Use make lint"},
	}

	// when
	result := formatCodexRules(rules)

	// then
	if !strings.Contains(result, "prefix_rule(") {
		t.Error("result should contain prefix_rule(")
	}
	if !strings.Contains(result, `pattern = ["make", "lint"]`) {
		t.Error("result should contain make lint pattern")
	}
	if !strings.Contains(result, `decision = "allow"`) {
		t.Error("result should contain allow decision")
	}
	if !strings.Contains(result, `decision = "forbidden"`) {
		t.Error("result should contain forbidden decision")
	}
	if !strings.Contains(result, `pattern = ["golangci-lint"]`) {
		t.Error("result should contain golangci-lint pattern")
	}
}

func TestFormatStarlarkList(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "single element",
			input:    []string{"golangci-lint"},
			expected: `["golangci-lint"]`,
		},
		{
			name:     "multiple elements",
			input:    []string{"git", "push", "--force"},
			expected: `["git", "push", "--force"]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			input := tt.input

			// when
			result := formatStarlarkList(input)

			// then
			if result != tt.expected {
				t.Errorf("formatStarlarkList()\n  got:  %q\n  want: %q", result, tt.expected)
			}
		})
	}
}

func TestWriteCodexRules(t *testing.T) {
	// given
	tmpDir := t.TempDir()

	// when
	err := writeCodexRules(tmpDir)

	// then
	if err != nil {
		t.Fatalf("writeCodexRules() error: %v", err)
	}
	path := filepath.Join(tmpDir, "codex", "rules", "default.rules")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading default.rules: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "prefix_rule(") {
		t.Error("default.rules should contain prefix_rule(")
	}
	if !strings.Contains(content, "make") {
		t.Error("default.rules should contain make rules")
	}
	if !strings.Contains(content, "golangci-lint") {
		t.Error("default.rules should contain golangci-lint rule")
	}
	if !strings.Contains(content, "forbidden") {
		t.Error("default.rules should contain forbidden decisions")
	}
	if !strings.Contains(content, "prompt") {
		t.Error("default.rules should contain prompt decisions for git force-push")
	}
}

// The per-tool SAST rules are copied by hand from the CI/CD page, so a tool
// that leaves the pipelines repository must leave both places together --
// including the ban that would otherwise keep sending agents to its target.
func TestCodexRulesMatchDocumentedSASTTargets(t *testing.T) {
	// given
	documented := documentedSASTTargets(t)

	// when
	allowed, redirects := codexMakeTargets(codexRules())

	// then
	for tool := range documented {
		if !allowed[tool] {
			t.Errorf("codexRules() should allow the documented `make %s`", tool)
		}
		if redirects[tool] != tool {
			t.Errorf("codexRules() should forbid calling %s directly in favor of `make %s`", tool, tool)
		}
	}
	for target := range allowed {
		if !isGeneralMakeTarget(target) && !documented[target] {
			t.Errorf("codexRules() should not allow `make %s`, which Life-Cycle/CI-&-CD.md does not document", target)
		}
	}
	for tool, target := range redirects {
		if !isGeneralMakeTarget(target) && !documented[target] {
			t.Errorf("codexRules() should not send %s to `make %s`, which Life-Cycle/CI-&-CD.md does not document", tool, target)
		}
	}
}

// documentedSASTTargets returns the per-tool SAST targets Life-Cycle/CI-&-CD.md lists.
func documentedSASTTargets(t *testing.T) map[string]bool {
	t.Helper()
	source, err := os.ReadFile(filepath.Join("..", "..", "..", "Life-Cycle", "CI-&-CD.md"))
	if err != nil {
		t.Fatal(err)
	}
	fence := strings.Repeat(string(rune(96)), 3)
	block := regexp.MustCompile("(?s)Individual SAST tools[^\n]*\n\n" + fence + "bash\n(.*?)" + fence).FindStringSubmatch(string(source))
	if len(block) != 2 {
		t.Fatal("Life-Cycle/CI-&-CD.md should list the per-tool SAST targets in a bash block")
	}
	targets := map[string]bool{}
	for _, match := range regexp.MustCompile(`(?m)^make ([a-z-]+)`).FindAllStringSubmatch(block[1], -1) {
		targets[match[1]] = true
	}
	if len(targets) == 0 {
		t.Fatal("Life-Cycle/CI-&-CD.md should document at least one per-tool SAST target")
	}
	return targets
}

// codexMakeTargets returns the Makefile targets the rules allow and, for every
// forbidden command, the Makefile target its justification sends the agent to.
func codexMakeTargets(rules []CodexRule) (map[string]bool, map[string]string) {
	allowed := map[string]bool{}
	redirects := map[string]string{}
	target := regexp.MustCompile("`make ([a-z-]+)`")
	for _, rule := range rules {
		if rule.Decision == "allow" && len(rule.Pattern) == 2 && rule.Pattern[0] == "make" {
			allowed[rule.Pattern[1]] = true
		}
		if match := target.FindStringSubmatch(rule.Justification); rule.Decision == "forbidden" && match != nil {
			redirects[strings.Join(rule.Pattern, " ")] = match[1]
		}
	}
	return allowed, redirects
}

// isGeneralMakeTarget reports whether target is one every project has rather
// than a per-tool SAST target.
func isGeneralMakeTarget(target string) bool {
	return target == "lint" || target == "test" || target == "sast"
}
