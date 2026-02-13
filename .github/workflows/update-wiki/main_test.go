package main

import "testing"

const testBaseURL = "https://raw.githubusercontent.com/wiki/rios0rios0/guide"

func TestReplaceLinks(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "link without directory path",
			input:    "- [Home](Home.md)",
			expected: "- [Home](Home)",
		},
		{
			name:     "link without directory path (Onboarding)",
			input:    "- [Onboarding](Onboarding.md)",
			expected: "- [Onboarding](Onboarding)",
		},
		{
			name:     "link with single directory",
			input:    "- [PDCA](Agile-&-Culture/PDCA.md)",
			expected: "- [PDCA](PDCA)",
		},
		{
			name:     "link with nested directory",
			input:    "- [Backend](Life-Cycle/Architecture/Backend-Design.md)",
			expected: "- [Backend](Backend-Design)",
		},
		{
			name:     "external link unchanged",
			input:    "- [Google](https://google.com)",
			expected: "- [Google](https://google.com)",
		},
		{
			name:     "external link with .md unchanged",
			input:    "- [Docs](https://example.com/page.md)",
			expected: "- [Docs](https://example.com/page.md)",
		},
		{
			name:     "multiple links on one line",
			input:    "See [Home](Home.md) and [Tests](Life-Cycle/Tests.md) for details.",
			expected: "See [Home](Home) and [Tests](Tests) for details.",
		},
		{
			name:     "link without .md unchanged",
			input:    "- [Home](Home)",
			expected: "- [Home](Home)",
		},
		{
			name:     "link with special characters in path",
			input:    "- [CI & CD](Life-Cycle/CI-&-CD.md)",
			expected: "- [CI & CD](CI-&-CD)",
		},
		{
			name:     "link with parentheses in filename",
			input:    "- [PEP 8](Code-Style/Python/Styling-and-Formatting-(PEP-8).md)",
			expected: "- [PEP 8](Styling-and-Formatting-(PEP-8))",
		},
		{
			name:     "inline link within paragraph text",
			input:    "Refer to the [Backend Design](Life-Cycle/Architecture/Backend-Design.md) section.",
			expected: "Refer to the [Backend Design](Backend-Design) section.",
		},
		{
			name:     "wiki image syntax not affected",
			input:    "[[https://example.com/image.png]]",
			expected: "[[https://example.com/image.png]]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			input := tt.input

			// when
			result := replaceLinks(input)

			// then
			if result != tt.expected {
				t.Errorf("replaceLinks(%q)\n  got:  %q\n  want: %q", input, result, tt.expected)
			}
		})
	}
}

func TestReplaceImages(t *testing.T) {
	tests := []struct {
		name     string
		fileDir  string
		input    string
		expected string
	}{
		{
			name:     "image in root directory",
			fileDir:  ".",
			input:    "![](.assets/flow-view.png)",
			expected: "[[" + testBaseURL + "/.assets/flow-view.png]]",
		},
		{
			name:     "image in nested directory",
			fileDir:  "Life-Cycle/Architecture",
			input:    "![](.assets/requests_flow.png)",
			expected: "[[" + testBaseURL + "/Life-Cycle/Architecture/.assets/requests_flow.png]]",
		},
		{
			name:     "image with parent directory traversal",
			fileDir:  "Life-Cycle/Git-Flow",
			input:    "![](../.assets/feature-branches.svg)",
			expected: "[[" + testBaseURL + "/Life-Cycle/.assets/feature-branches.svg]]",
		},
		{
			name:     "image with alt text",
			fileDir:  "Life-Cycle",
			input:    "![Architecture](.assets/clean-architecture.png)",
			expected: "[[" + testBaseURL + "/Life-Cycle/.assets/clean-architecture.png|alt=Architecture]]",
		},
		{
			name:    "multiple images on one line (table row)",
			fileDir: "Life-Cycle",
			input:   "| ![](.assets/not-clean.png) | ![](.assets/clean.png) |",
			expected: "| [[" + testBaseURL + "/Life-Cycle/.assets/not-clean.png]]" +
				" | [[" + testBaseURL + "/Life-Cycle/.assets/clean.png]] |",
		},
		{
			name:     "external image URL unchanged",
			fileDir:  ".",
			input:    "![logo](https://example.com/logo.png)",
			expected: "![logo](https://example.com/logo.png)",
		},
		{
			name:     "no image in line unchanged",
			fileDir:  ".",
			input:    "This is a regular paragraph.",
			expected: "This is a regular paragraph.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			input := tt.input

			// when
			result := replaceImages(input, tt.fileDir, testBaseURL)

			// then
			if result != tt.expected {
				t.Errorf("replaceImages(%q, %q)\n  got:  %q\n  want: %q", input, tt.fileDir, result, tt.expected)
			}
		})
	}
}
