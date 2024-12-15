package main

import (
  "fmt"
  logger "github.com/sirupsen/logrus"
  "os"
  "os/exec"
  "path/filepath"
  "strings"
)

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

// replaceImages updates image paths in the markdown content
func replaceImages(text string) string {
  lines := strings.Split(text, "\n")
  for i, line := range lines {
    // check if the line contains an image link
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

// replaceLinks updates internal Markdown links to remove the .md extension
func replaceLinks(text string) string {
  lines := strings.Split(text, "\n")
  for i, line := range lines {
    // check if the line contains a link
    if strings.Contains(line, "") && strings.Contains(line, "](") {
      start := strings.Index(line, "](") + 2
      end := strings.Index(line[start:], ")") + start
      if end > start {
        linkPath := line[start:end]
        // check if the link is not an external link and has a .md extension
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
