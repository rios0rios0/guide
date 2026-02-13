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

// imageRegex matches markdown images: ![alt](path)
var imageRegex = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)

// mdLinkTargetRegex matches markdown link targets ending in .md: ](path.md)
// Uses a lazy quantifier (.+?) to correctly handle filenames containing parentheses
// (e.g., "Styling-and-Formatting-(PEP-8).md"). This works for all occurrences on a line
// and for both [text](path.md) and ![alt](path.md). Since images are processed first
// (replacing paths to .assets/name.ext), image paths will no longer end in .md by the
// time this regex runs.
var mdLinkTargetRegex = regexp.MustCompile(`\]\((.+?\.md)\)`)

func main() {
	wikiDir := "wiki"

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
			processFile(path)
		}
		return nil
	})

	// log an error if the directory walk fails
	if err != nil {
		logger.Errorf("Error walking through '%s' directory: %v\n", wikiDir, err)
	}
}

// processFile reads, processes, and writes back the file content
func processFile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		logger.Errorf("Error reading file %s: %v\n", path, err)
		return
	}

	text := string(content)
	text = replaceImages(text)
	text = replaceLinks(text)

	// write the modified content back to the file
	err = os.WriteFile(path, []byte(text), 0644)
	if err != nil {
		logger.Errorf("Error writing file %s: %v\n", path, err)
	}
}

// replaceImages updates all image paths in the markdown content to point to the
// flat .assets/ directory used by GitHub Wiki. Handles multiple images per line.
func replaceImages(text string) string {
	return imageRegex.ReplaceAllStringFunc(text, func(match string) string {
		groups := imageRegex.FindStringSubmatch(match)
		if len(groups) < 3 {
			return match
		}
		altText := groups[1]
		imagePath := groups[2]
		imageName := filepath.Base(imagePath)
		return fmt.Sprintf("![%s](.assets/%s)", altText, imageName)
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
