package main

import (
	"flag"
	"os"
	"path/filepath"
	"time"

	logger "github.com/sirupsen/logrus"
)

func main() {
	sourceDir := flag.String("source", ".", "root directory of the documentation repository")
	outputDir := flag.String("output", ".", "directory where generated rule files are written")
	externalConfig := flag.String("external-config", "", "path to external-sources.yaml for fetching agents from external repos")
	logLevel := flag.String("log-level", "info", "log level: trace, debug, info, warn, error, fatal")
	flag.Parse()

	level, err := logger.ParseLevel(*logLevel)
	if err != nil {
		logger.WithFields(logger.Fields{
			"provided": *logLevel,
		}).Fatal("invalid log level")
	}
	logger.SetLevel(level)

	groups := ruleGroups()

	logger.WithFields(logger.Fields{
		"source_dir":  *sourceDir,
		"output_dir":  *outputDir,
		"log_level":   *logLevel,
		"group_count": len(groups),
	}).Info("starting rule generation")

	start := time.Now()
	contents := make([]string, len(groups))
	var errorCount int

	for i, group := range groups {
		merged, err := processGroup(*sourceDir, group)
		if err != nil {
			logger.WithFields(logger.Fields{
				"group": group.Name,
				"error": err.Error(),
			}).Error("failed to process group")
			errorCount++
			continue
		}
		contents[i] = merged
	}

	writeErrors := writeAllRules(*outputDir, groups, contents)
	totalErrors := errorCount + writeErrors

	if *externalConfig != "" {
		externalErrors := fetchExternalSources(*externalConfig, *outputDir)
		totalErrors += externalErrors
	}

	logger.WithFields(logger.Fields{
		"groups_processed": len(groups),
		"errors":           totalErrors,
		"duration":         time.Since(start).String(),
	}).Info("rule generation complete")

	if totalErrors > 0 {
		os.Exit(1)
	}
}

// processGroup reads and transforms all source files for a rule group.
func processGroup(sourceDir string, group RuleGroup) (string, error) {
	var parts []string
	for _, src := range group.Sources {
		path := filepath.Join(sourceDir, src)
		data, err := os.ReadFile(path)
		if err != nil {
			logger.WithFields(logger.Fields{
				"group":  group.Name,
				"source": src,
				"error":  err.Error(),
			}).Warn("skipped source file")
			continue
		}
		transformed := transformContent(string(data))
		logger.WithFields(logger.Fields{
			"group":        group.Name,
			"source":       src,
			"raw_bytes":    len(data),
			"output_bytes": len(transformed),
		}).Debug("processed source file")
		parts = append(parts, transformed)
	}
	return mergeContents(parts), nil
}

// writeAllRules writes rule files for all AI assistants (Claude, Cursor, Copilot, and Codex).
// It returns the number of errors encountered during writing.
func writeAllRules(outputDir string, groups []RuleGroup, contents []string) int {
	var errorCount int
	var claudeCount, cursorCount, copilotCount int

	for i, group := range groups {
		if contents[i] == "" {
			logger.WithFields(logger.Fields{
				"group": group.Name,
			}).Debug("skipped group with empty content")
			continue
		}
		if err := writeClaude(outputDir, group, contents[i]); err != nil {
			logger.WithFields(logger.Fields{
				"group": group.Name,
				"error": err.Error(),
			}).Error("failed to write Claude rule")
			errorCount++
		} else {
			claudeCount++
		}
		if err := writeCursor(outputDir, group, contents[i]); err != nil {
			logger.WithFields(logger.Fields{
				"group": group.Name,
				"error": err.Error(),
			}).Error("failed to write Cursor rule")
			errorCount++
		} else {
			cursorCount++
		}
		if err := writeCopilot(outputDir, group, contents[i]); err != nil {
			logger.WithFields(logger.Fields{
				"group": group.Name,
				"error": err.Error(),
			}).Error("failed to write Copilot instruction")
			errorCount++
		} else {
			copilotCount++
		}
	}

	if err := writeCodex(outputDir, groups, contents); err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
		}).Error("failed to write Codex AGENTS.md")
		errorCount++
	}

	if err := writeCodexRules(outputDir); err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
		}).Error("failed to write Codex rules")
		errorCount++
	}

	logger.WithFields(logger.Fields{
		"claude_rules":        claudeCount,
		"cursor_rules":        cursorCount,
		"copilot_instructions": copilotCount,
	}).Info("completed writing rules")

	return errorCount
}
