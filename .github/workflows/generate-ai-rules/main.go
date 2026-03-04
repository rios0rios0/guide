package main

import (
	"flag"
	"os"
	"path/filepath"

	logger "github.com/sirupsen/logrus"
)

func main() {
	sourceDir := flag.String("source", ".", "root directory of the documentation repository")
	outputDir := flag.String("output", ".", "directory where generated rule files are written")
	flag.Parse()

	groups := ruleGroups()
	contents := make([]string, len(groups))

	for i, group := range groups {
		merged, err := processGroup(*sourceDir, group)
		if err != nil {
			logger.Errorf("processing group %q: %v", group.Name, err)
			continue
		}
		contents[i] = merged
	}

	writeAllRules(*outputDir, groups, contents)
}

// processGroup reads and transforms all source files for a rule group.
func processGroup(sourceDir string, group RuleGroup) (string, error) {
	var parts []string
	for _, src := range group.Sources {
		path := filepath.Join(sourceDir, src)
		data, err := os.ReadFile(path)
		if err != nil {
			logger.Warnf("skipping %s: %v", src, err)
			continue
		}
		transformed := transformContent(string(data))
		parts = append(parts, transformed)
	}
	return mergeContents(parts), nil
}

// writeAllRules writes rule files for all three AI assistants.
func writeAllRules(outputDir string, groups []RuleGroup, contents []string) {
	for i, group := range groups {
		if contents[i] == "" {
			continue
		}
		if err := writeClaude(outputDir, group, contents[i]); err != nil {
			logger.Errorf("writing Claude rule %q: %v", group.Name, err)
		}
		if err := writeCursor(outputDir, group, contents[i]); err != nil {
			logger.Errorf("writing Cursor rule %q: %v", group.Name, err)
		}
	}

	if err := writeCodex(outputDir, groups, contents); err != nil {
		logger.Errorf("writing Codex AGENTS.md: %v", err)
	}

	logger.Infof("generated rules for %d groups", len(groups))
}
