package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	logger "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// ExternalConfig holds the list of external sources to fetch artifacts from.
type ExternalConfig struct {
	Sources []ExternalSource `yaml:"sources"`
}

// ExternalSource defines a single external GitHub repository and the artifacts to fetch.
type ExternalSource struct {
	Repo     string             `yaml:"repo"`
	Branch   string             `yaml:"branch"`
	Agents   []ExternalArtifact `yaml:"agents"`
	Commands []ExternalArtifact `yaml:"commands"`
	Skills   []ExternalArtifact `yaml:"skills"`
}

// ExternalArtifact maps a source file path in an external repo to a target filename.
type ExternalArtifact struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

// httpClient is the shared HTTP client for external fetches.
var httpClient = &http.Client{Timeout: 15 * time.Second}

// loadExternalConfig reads and parses the external sources YAML configuration file.
func loadExternalConfig(path string) (*ExternalConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading external config %s: %w", path, err)
	}

	var config ExternalConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing external config %s: %w", path, err)
	}

	return &config, nil
}

// fetchExternalSources loads the config and fetches all external artifacts into the output directory.
func fetchExternalSources(configPath, outputDir string) int {
	config, err := loadExternalConfig(configPath)
	if err != nil {
		logger.WithFields(logger.Fields{
			"config": configPath,
			"error":  err.Error(),
		}).Error("failed to load external sources config")
		return 1
	}

	var errorCount int
	var fetchedCount int

	for _, source := range config.Sources {
		branch := source.Branch
		if branch == "" {
			branch = "main"
		}

		for _, agent := range source.Agents {
			if err := fetchArtifact(source.Repo, branch, agent, outputDir, "agents"); err != nil {
				logger.WithFields(logger.Fields{
					"repo":   source.Repo,
					"source": agent.Source,
					"error":  err.Error(),
				}).Error("failed to fetch external agent")
				errorCount++
			} else {
				fetchedCount++
			}
		}

		for _, cmd := range source.Commands {
			if err := fetchArtifact(source.Repo, branch, cmd, outputDir, "commands"); err != nil {
				logger.WithFields(logger.Fields{
					"repo":   source.Repo,
					"source": cmd.Source,
					"error":  err.Error(),
				}).Error("failed to fetch external command")
				errorCount++
			} else {
				fetchedCount++
			}
		}

		for _, skill := range source.Skills {
			if err := fetchArtifact(source.Repo, branch, skill, outputDir, "skills"); err != nil {
				logger.WithFields(logger.Fields{
					"repo":   source.Repo,
					"source": skill.Source,
					"error":  err.Error(),
				}).Error("failed to fetch external skill")
				errorCount++
			} else {
				fetchedCount++
			}
		}
	}

	logger.WithFields(logger.Fields{
		"fetched": fetchedCount,
		"errors":  errorCount,
	}).Info("completed fetching external sources")

	return errorCount
}

// validateTarget checks that the artifact target path is safe for the given artifact type.
// Agents and commands must be plain filenames. Skills must follow the <name>/SKILL.md pattern.
func validateTarget(target, artifactType string) error {
	cleaned := filepath.Clean(target)
	if filepath.IsAbs(cleaned) {
		return fmt.Errorf("invalid artifact target %q: must be a relative path", target)
	}
	if strings.Contains(cleaned, "..") {
		return fmt.Errorf("invalid artifact target %q: path traversal not allowed", target)
	}

	switch artifactType {
	case "agents", "commands":
		if filepath.Base(cleaned) != cleaned {
			return fmt.Errorf("invalid artifact target %q: %s targets must be plain filenames", target, artifactType)
		}
	case "skills":
		parts := strings.Split(cleaned, string(filepath.Separator))
		if len(parts) != 2 || parts[1] != "SKILL.md" {
			return fmt.Errorf("invalid artifact target %q: skills targets must follow the <name>/SKILL.md pattern", target)
		}
	default:
		return fmt.Errorf("unknown artifact type %q", artifactType)
	}

	return nil
}

// artifactOutputDir returns the directory where the fetched artifact file should be written.
// Agents and commands go to claude/<type>/, skills go to cursor/skills/<name>/.
func artifactOutputDir(outputDir, artifactType, target string) string {
	if artifactType == "skills" {
		skillName := strings.Split(filepath.Clean(target), string(filepath.Separator))[0]
		return filepath.Join(outputDir, "cursor", "skills", skillName)
	}
	return filepath.Join(outputDir, "claude", artifactType)
}

// fetchArtifact downloads a single artifact from a GitHub repo and writes it to the output directory.
func fetchArtifact(repo, branch string, artifact ExternalArtifact, outputDir, artifactType string) error {
	if err := validateTarget(artifact.Target, artifactType); err != nil {
		return err
	}

	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s", repo, branch, artifact.Source)

	body, err := httpGet(url)
	if err != nil {
		return fmt.Errorf("fetching %s: %w", url, err)
	}

	dir := artifactOutputDir(outputDir, artifactType, artifact.Target)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory %s: %w", dir, err)
	}

	path := filepath.Join(dir, filepath.Base(artifact.Target))
	if err := os.WriteFile(path, body, 0644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}

	logger.WithFields(logger.Fields{
		"repo":   repo,
		"source": artifact.Source,
		"target": path,
		"bytes":  len(body),
	}).Debug("fetched external artifact")

	return nil
}

// httpGet performs an HTTP GET with one retry on failure.
func httpGet(url string) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt < 2; attempt++ {
		resp, err := httpClient.Get(url)
		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode != http.StatusOK {
			_, _ = io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
			lastErr = fmt.Errorf("HTTP %d for %s", resp.StatusCode, url)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("reading response body: %w", err)
			continue
		}
		return body, nil
	}
	return nil, lastErr
}
