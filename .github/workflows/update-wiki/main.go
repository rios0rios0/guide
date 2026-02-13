package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	logger "github.com/sirupsen/logrus"
)

const wikiDir = "wiki"

// imageRegex matches markdown images: ![alt](path)
var imageRegex = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)

// mdLinkTargetRegex matches markdown link targets ending in .md: ](path.md)
// Uses a lazy quantifier (.+?) to correctly handle filenames containing parentheses
// (e.g., "Styling-and-Formatting-(PEP-8).md"). This works for all occurrences on a line
// and for both [text](path.md) and ![alt](path.md). Since images are processed first
// (converting to HTML <img> tags), image paths will no longer match by the time this
// regex runs.
var mdLinkTargetRegex = regexp.MustCompile(`\]\((.+?\.md)\)`)

func main() {
	// GITHUB_REPOSITORY is automatically set by GitHub Actions (e.g., "rios0rios0/guide")
	githubRepo := os.Getenv("GITHUB_REPOSITORY")
	if githubRepo == "" {
		logger.Warn("GITHUB_REPOSITORY not set; falling back to 'rios0rios0/guide'")
		githubRepo = "rios0rios0/guide"
	}

	// base URL for raw image content in the GitHub Wiki repository
	rawBaseURL := fmt.Sprintf("https://raw.githubusercontent.com/wiki/%s", githubRepo)

	// clear the "wikiDir" directory excluding the .git folder
	err := exec.Command("sh", "-c",
		fmt.Sprintf("find %s -mindepth 1 ! -regex '^%s/\\.git\\(/.*\\)?' -delete", wikiDir, wikiDir)).Run()
	if err != nil {
		logger.Errorf("Error clearing '%s' directory: %v\n", wikiDir, err)
		return
	}

	// copy files and folders from root to "wikiDir" directory, excluding .git and .github folders
	err = exec.Command("sh", "-c",
		fmt.Sprintf("rsync -av --exclude='.git' --exclude='.github' --exclude='%s' "+
			"--exclude='.editorconfig' --exclude='README.md' ./ %s", wikiDir, wikiDir)).Run()
	if err != nil {
		logger.Errorf("Error copying files: %v\n", err)
		return
	}

	// walking through the directory and processing each file
	err = filepath.Walk(wikiDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			processFile(path, rawBaseURL)
		}
		return nil
	})

	// log an error if the directory walk fails
	if err != nil {
		logger.Errorf("Error walking through '%s' directory: %v\n", wikiDir, err)
	}
}

// processFile reads, processes, and writes back the file content.
// The rawBaseURL is used to construct absolute image URLs for the wiki.
func processFile(path string, rawBaseURL string) {
	content, err := os.ReadFile(path)
	if err != nil {
		logger.Errorf("Error reading file %s: %v\n", path, err)
		return
	}

	// compute the directory of this file relative to the wiki root
	// e.g., "wiki/Life-Cycle/Architecture/Backend-Design.md" -> "Life-Cycle/Architecture"
	relPath, _ := filepath.Rel(wikiDir, path)
	fileDir := filepath.Dir(relPath)

	text := string(content)
	text = replaceImages(text, fileDir, rawBaseURL)
	text = replaceLinks(text)

	// write the modified content back to the file
	err = os.WriteFile(path, []byte(text), 0644)
	if err != nil {
		logger.Errorf("Error writing file %s: %v\n", path, err)
	}
}

// replaceImages converts markdown image syntax to GitHub Wiki image syntax with
// absolute URLs. This is necessary because GitHub Wiki renders pages as flat URLs,
// so relative image paths do not resolve correctly.
//
// GitHub Wiki image format: [[URL|alt=description]]
//
// Examples (with rawBaseURL = "https://raw.githubusercontent.com/wiki/rios0rios0/guide"):
//
//	File in "Life-Cycle/Architecture", image ref ".assets/flow.png"
//	  -> [[https://raw.githubusercontent.com/wiki/rios0rios0/guide/Life-Cycle/Architecture/.assets/flow.png]]
//
//	File in ".", image ref ".assets/flow.png", alt text "diagram"
//	  -> [[https://raw.githubusercontent.com/wiki/rios0rios0/guide/.assets/flow.png|alt=diagram]]
func replaceImages(text string, fileDir string, rawBaseURL string) string {
	return imageRegex.ReplaceAllStringFunc(text, func(match string) string {
		groups := imageRegex.FindStringSubmatch(match)
		if len(groups) < 3 {
			return match
		}
		altText := groups[1]
		imagePath := groups[2]

		// skip external images (already absolute URLs)
		if strings.HasPrefix(imagePath, "http") {
			return match
		}

		// resolve the image path relative to the file's directory
		// e.g., fileDir="Life-Cycle/Architecture", imagePath=".assets/flow.png"
		//   -> "Life-Cycle/Architecture/.assets/flow.png"
		// e.g., fileDir="Life-Cycle", imagePath="../.assets/branches.svg"
		//   -> ".assets/branches.svg" (after Clean)
		resolvedPath := filepath.Join(fileDir, imagePath)
		resolvedPath = filepath.Clean(resolvedPath)

		// construct the full raw URL
		fullURL := fmt.Sprintf("%s/%s", rawBaseURL, resolvedPath)

		// return GitHub Wiki image syntax: [[URL|alt=text]] or [[URL]]
		if altText != "" {
			return fmt.Sprintf("[[%s|alt=%s]]", fullURL, altText)
		}
		return fmt.Sprintf("[[%s]]", fullURL)
	})
}

// replaceLinks removes the .md extension from all internal markdown links and
// flattens directory paths to just the page name (GitHub Wiki pages are flat).
// Handles multiple links per line and links with or without directory paths.
//
// Examples:
//
//	[Home](Home.md)                              -> [Home](Home)
//	[Onboarding](Onboarding.md)                  -> [Onboarding](Onboarding)
//	[Git Flow](Life-Cycle/Git-Flow.md)           -> [Git Flow](Git-Flow)
//	[Backend](Life-Cycle/Architecture/Backend.md) -> [Backend](Backend)
//	[Google](https://google.com)                  -> [Google](https://google.com)  (unchanged)
func replaceLinks(text string) string {
	return mdLinkTargetRegex.ReplaceAllStringFunc(text, func(match string) string {
		groups := mdLinkTargetRegex.FindStringSubmatch(match)
		if len(groups) < 2 {
			return match
		}
		linkPath := groups[1]

		// skip external links
		if strings.HasPrefix(linkPath, "http") {
			return match
		}

		// extract the base name and remove the .md extension
		// filepath.Base flattens "Life-Cycle/Git-Flow.md" -> "Git-Flow.md"
		linkName := filepath.Base(linkPath)
		linkName = strings.TrimSuffix(linkName, ".md")
		return "](" + linkName + ")"
	})
}
