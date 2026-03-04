package main

import (
	"regexp"
	"strings"
)

// Compiled regex patterns for content transformation.
var (
	// imageRegex matches markdown image syntax: ![alt](path)
	imageRegex = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)

	// referenceSectionRegex matches the ## References heading.
	referenceSectionRegex = regexp.MustCompile(`(?m)^## References\s*$`)

	// footnoteRefRegex matches inline footnote references like [^1].
	footnoteRefRegex = regexp.MustCompile(`\[\^\d+\]`)

	// footnoteDefRegex matches footnote definition lines like [^1]: text.
	footnoteDefRegex = regexp.MustCompile(`(?m)^\[\^\d+\]:.*$`)

	// mdLinkRegex matches [text](path.md) markdown links (lazy quantifier for parentheses in paths).
	mdLinkRegex = regexp.MustCompile(`\[([^\]]+)\]\((.+?\.md)\)`)

	// subPageLinkLineRegex matches bullet list items that are .md links.
	subPageLinkLineRegex = regexp.MustCompile(`(?m)^-\s+\[([^\]]+)\]\(([^)]+\.md)\).*$`)

	// multipleBlankLines matches 3+ consecutive newlines (2+ blank lines).
	multipleBlankLines = regexp.MustCompile(`\n{3,}`)
)

// transformContent applies all content transformations to a markdown string.
func transformContent(content string) string {
	content = stripImages(content)
	content = stripReferences(content)
	content = stripFootnotes(content)
	content = stripSubPageLinks(content)
	content = transformLinks(content)
	content = collapseWhitespace(content)
	return content
}

// stripImages removes internal markdown images and cleans up empty table rows.
// External images (http/https URLs) are preserved.
func stripImages(content string) string {
	lines := strings.Split(content, "\n")
	var result []string
	for _, line := range lines {
		cleaned := imageRegex.ReplaceAllStringFunc(line, func(match string) string {
			groups := imageRegex.FindStringSubmatch(match)
			if len(groups) >= 3 && strings.HasPrefix(groups[2], "http") {
				return match
			}
			return ""
		})
		cleaned = strings.TrimRight(cleaned, " ")
		// skip lines that became empty after image removal
		if cleaned == "" {
			result = append(result, cleaned)
			continue
		}
		// skip table rows that only contain pipes and whitespace after image removal
		if isEmptyTableRow(cleaned) {
			continue
		}
		result = append(result, cleaned)
	}
	return strings.Join(result, "\n")
}

// isEmptyTableRow checks if a line is a table row with no meaningful content.
func isEmptyTableRow(line string) bool {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "|") {
		return false
	}
	remaining := strings.ReplaceAll(trimmed, "|", "")
	remaining = strings.TrimSpace(remaining)
	return remaining == ""
}

// stripReferences removes the ## References section and everything after it.
func stripReferences(content string) string {
	loc := referenceSectionRegex.FindStringIndex(content)
	if loc == nil {
		return content
	}
	return strings.TrimRight(content[:loc[0]], "\n") + "\n"
}

// stripFootnotes removes footnote references [^N] and definition lines [^N]: text.
func stripFootnotes(content string) string {
	content = footnoteDefRegex.ReplaceAllString(content, "")
	content = footnoteRefRegex.ReplaceAllString(content, "")
	return content
}

// stripSubPageLinks removes bullet-list lines that link to internal .md files.
func stripSubPageLinks(content string) string {
	return subPageLinkLineRegex.ReplaceAllStringFunc(content, func(match string) string {
		groups := subPageLinkLineRegex.FindStringSubmatch(match)
		if len(groups) >= 3 && strings.HasPrefix(groups[2], "http") {
			return match
		}
		return ""
	})
}

// transformLinks converts [text](path.md) to just the display text for internal links.
func transformLinks(content string) string {
	return mdLinkRegex.ReplaceAllStringFunc(content, func(match string) string {
		groups := mdLinkRegex.FindStringSubmatch(match)
		if len(groups) < 3 {
			return match
		}
		linkPath := groups[2]
		// keep external links unchanged
		if strings.HasPrefix(linkPath, "http") {
			return match
		}
		return groups[1]
	})
}

// collapseWhitespace reduces multiple blank lines and trims trailing whitespace.
func collapseWhitespace(content string) string {
	content = multipleBlankLines.ReplaceAllString(content, "\n\n")
	content = strings.TrimRight(content, "\n ") + "\n"
	return content
}

// mergeContents concatenates multiple transformed contents with a separator.
func mergeContents(contents []string) string {
	var nonEmpty []string
	for _, c := range contents {
		trimmed := strings.TrimSpace(c)
		if trimmed != "" {
			nonEmpty = append(nonEmpty, trimmed)
		}
	}
	if len(nonEmpty) == 0 {
		return ""
	}
	return strings.Join(nonEmpty, "\n\n---\n\n") + "\n"
}
