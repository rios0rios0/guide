package main

import "testing"

func TestStripImages(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "standalone image line removed",
			input:    "before\n![](.assets/flow.png)\nafter",
			expected: "before\n\nafter",
		},
		{
			name:     "image with alt text removed",
			input:    "![Architecture](.assets/arch.png)",
			expected: "",
		},
		{
			name:     "external image preserved",
			input:    "![logo](https://example.com/logo.png)",
			expected: "![logo](https://example.com/logo.png)",
		},
		{
			name:     "table row with only images becomes empty and is removed",
			input:    "| header1 | header2 |\n|:---:|:---:|\n| ![](.assets/a.png) | ![](.assets/b.png) |",
			expected: "| header1 | header2 |\n|:---:|:---:|",
		},
		{
			name:     "text without images unchanged",
			input:    "This is a regular paragraph.",
			expected: "This is a regular paragraph.",
		},
		{
			name:     "mixed line keeps text removes image",
			input:    "See this diagram: ![diagram](.assets/d.png) for details",
			expected: "See this diagram:  for details",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			input := tt.input

			// when
			result := stripImages(input)

			// then
			if result != tt.expected {
				t.Errorf("stripImages()\n  got:  %q\n  want: %q", result, tt.expected)
			}
		})
	}
}

func TestStripReferences(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "removes references section at end",
			input:    "# Title\n\nContent here.\n\n## References\n\n- [Link](https://example.com)\n- [Other](https://other.com)\n",
			expected: "# Title\n\nContent here.\n",
		},
		{
			name:     "no references section unchanged",
			input:    "# Title\n\nContent only.\n",
			expected: "# Title\n\nContent only.\n",
		},
		{
			name:     "references with extra whitespace",
			input:    "Content.\n\n## References\n\n- link\n",
			expected: "Content.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			input := tt.input

			// when
			result := stripReferences(input)

			// then
			if result != tt.expected {
				t.Errorf("stripReferences()\n  got:  %q\n  want: %q", result, tt.expected)
			}
		})
	}
}

func TestStripFootnotes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "inline footnote references removed",
			input:    "Clean Architecture[^1] is important[^2].",
			expected: "Clean Architecture is important.",
		},
		{
			name:     "footnote definition lines removed",
			input:    "Text here.\n\n[^1]: Some reference link\n[^2]: Another reference\n",
			expected: "Text here.\n\n\n\n",
		},
		{
			name:     "no footnotes unchanged",
			input:    "Regular text without footnotes.",
			expected: "Regular text without footnotes.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			input := tt.input

			// when
			result := stripFootnotes(input)

			// then
			if result != tt.expected {
				t.Errorf("stripFootnotes()\n  got:  %q\n  want: %q", result, tt.expected)
			}
		})
	}
}

func TestTransformLinks(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "internal link converted to text",
			input:    "See the [Backend Design](Life-Cycle/Architecture/Backend-Design.md) section.",
			expected: "See the Backend Design section.",
		},
		{
			name:     "relative parent link converted",
			input:    "Refer to [Code Style](../Code-Style.md) guide.",
			expected: "Refer to Code Style guide.",
		},
		{
			name:     "external link unchanged",
			input:    "See [Effective Go](https://go.dev/doc/effective_go) for details.",
			expected: "See [Effective Go](https://go.dev/doc/effective_go) for details.",
		},
		{
			name:     "external .md link unchanged",
			input:    "See [Docs](https://example.com/page.md) here.",
			expected: "See [Docs](https://example.com/page.md) here.",
		},
		{
			name:     "multiple internal links",
			input:    "See [A](a.md) and [B](b.md) here.",
			expected: "See A and B here.",
		},
		{
			name:     "link with parentheses in filename",
			input:    "See [PEP 8](Code-Style/Python/Styling-and-Formatting-(PEP-8).md).",
			expected: "See PEP 8.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			input := tt.input

			// when
			result := transformLinks(input)

			// then
			if result != tt.expected {
				t.Errorf("transformLinks()\n  got:  %q\n  want: %q", result, tt.expected)
			}
		})
	}
}

func TestStripSubPageLinks(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "bullet list of internal links removed",
			input:    "Overview:\n\n- [Conventions](GoLang/GoLang-Conventions.md)\n- [Testing](GoLang/GoLang-Testing.md)\n\nNext section.",
			expected: "Overview:\n\n\n\n\nNext section.",
		},
		{
			name:     "bullet with trailing description removed",
			input:    "- [Backend Design](Architecture/Backend-Design.md) -- Layers and actors\n",
			expected: "\n",
		},
		{
			name:     "external link bullet preserved",
			input:    "- [Google](https://google.com)\n",
			expected: "- [Google](https://google.com)\n",
		},
		{
			name:     "inline link not affected",
			input:    "See the [Backend Design](Architecture/Backend-Design.md) section.\n",
			expected: "See the [Backend Design](Architecture/Backend-Design.md) section.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			input := tt.input

			// when
			result := stripSubPageLinks(input)

			// then
			if result != tt.expected {
				t.Errorf("stripSubPageLinks()\n  got:  %q\n  want: %q", result, tt.expected)
			}
		})
	}
}

func TestCollapseWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "multiple blank lines collapsed to single blank line",
			input:    "line1\n\n\n\n\nline2",
			expected: "line1\n\nline2\n",
		},
		{
			name:     "three newlines collapsed to two",
			input:    "line1\n\n\nline2",
			expected: "line1\n\nline2\n",
		},
		{
			name:     "trailing whitespace trimmed",
			input:    "content\n\n\n",
			expected: "content\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			input := tt.input

			// when
			result := collapseWhitespace(input)

			// then
			if result != tt.expected {
				t.Errorf("collapseWhitespace()\n  got:  %q\n  want: %q", result, tt.expected)
			}
		})
	}
}

func TestMergeContents(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "two contents merged with separator",
			input:    []string{"# Part 1\n\nContent A.\n", "# Part 2\n\nContent B.\n"},
			expected: "# Part 1\n\nContent A.\n\n---\n\n# Part 2\n\nContent B.\n",
		},
		{
			name:     "empty contents skipped",
			input:    []string{"Content A.\n", "", "  \n", "Content B.\n"},
			expected: "Content A.\n\n---\n\nContent B.\n",
		},
		{
			name:     "single content no separator",
			input:    []string{"# Only Part\n\nContent.\n"},
			expected: "# Only Part\n\nContent.\n",
		},
		{
			name:     "all empty returns empty",
			input:    []string{"", "  "},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			input := tt.input

			// when
			result := mergeContents(input)

			// then
			if result != tt.expected {
				t.Errorf("mergeContents()\n  got:  %q\n  want: %q", result, tt.expected)
			}
		})
	}
}

func TestTransformContent(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "full pipeline on realistic content",
			input: `# Go

> **TL;DR:** Use Dig for DI, golangci-lint for linting.

## Overview

See sub-pages for specific topics:

- [Conventions](GoLang/GoLang-Conventions.md)
- [Testing](GoLang/GoLang-Testing.md)

## Go Proverbs

The [Go Proverbs](https://go-proverbs.github.io/) capture the design philosophy.

![](.assets/proverbs.png)

Clean Architecture[^1] is key.

[^1]: Reference link

## References

- [Effective Go](https://go.dev/doc/effective_go)
`,
			expected: `# Go

> **TL;DR:** Use Dig for DI, golangci-lint for linting.

## Overview

See sub-pages for specific topics:

## Go Proverbs

The [Go Proverbs](https://go-proverbs.github.io/) capture the design philosophy.

Clean Architecture is key.
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			input := tt.input

			// when
			result := transformContent(input)

			// then
			if result != tt.expected {
				t.Errorf("transformContent()\n  got:  %q\n  want: %q", result, tt.expected)
			}
		})
	}
}
