package main

import (
	logger "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	wikiDir := "wiki"

	err := filepath.Walk(wikiDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			processFile(path)
		}
		return nil
	})

	if err != nil {
		logger.Errorf("Error walking through wiki directory: %v\n", err)
	}
}

func processFile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		logger.Errorf("Error reading file %s: %v\n", path, err)
		return
	}

	text := string(content)
	text = replaceImages(text)
	text = replaceLinks(text)

	err = os.WriteFile(path, []byte(text), 0644)
	if err != nil {
		logger.Errorf("Error writing file %s: %v\n", path, err)
	}
}

func replaceImages(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if strings.Contains(line, "!") && strings.Contains(line, "](") {
			start := strings.Index(line, "](") + 2
			end := strings.Index(line[start:], ")") + start
			if end > start {
				imagePath := line[start:end]
				imageName := filepath.Base(imagePath)
				lines[i] = strings.Replace(line, imagePath, ".assets/"+imageName, 1)
			}
		}
	}
	return strings.Join(lines, "\n")
}

func replaceLinks(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if strings.Contains(line, "") && strings.Contains(line, "](") {
			start := strings.Index(line, "](") + 2
			end := strings.Index(line[start:], ")") + start
			if end > start {
				linkPath := line[start:end]
				if !strings.HasPrefix(linkPath, "http") && strings.Contains(linkPath, "/") && strings.HasSuffix(linkPath, ".md") {
					linkName := filepath.Base(linkPath)
					linkName = strings.TrimSuffix(linkName, ".md")
					lines[i] = strings.Replace(line, linkPath, linkName, 1)
				}
			}
		}
	}
	return strings.Join(lines, "\n")
}
