package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestScaffoldExamplesCompileAndRun(t *testing.T) {
	t.Parallel()
	// given
	source, err := os.ReadFile("commands/scaffold-go-project.md")
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	fence := strings.Repeat(string(rune(96)), 3)
	blocks := regexp.MustCompile("(?s)"+fence+"go\n// ([^\n]+\\.go)\n(.*?)"+fence).FindAllStringSubmatch(string(source), -1)
	if len(blocks) < 6 {
		t.Fatal("missing scaffold source examples")
	}
	for _, block := range blocks {
		file := filepath.Join(root, strings.ReplaceAll(block[1], "<app>", "app"))
		if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
			t.Fatal(err)
		}
		body := strings.ReplaceAll(block[2], "module/", "example.com/fixture/")
		if strings.Contains(body, "go.uber.org/dig") || strings.Contains(body, "github.com/google/wire") {
			t.Fatal("scaffold introduced a DI framework")
		}
		if err := os.WriteFile(file, []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/fixture\n\ngo 1.26.2\n\nrequire github.com/sirupsen/logrus v1.9.3\n"), 0644); err != nil {
		t.Fatal(err)
	}
	makefile := regexp.MustCompile("(?s)" + fence + "makefile\n(.*?)" + fence).FindStringSubmatch(string(source))
	if len(makefile) != 2 {
		t.Fatal("missing test targets")
	}
	if err := os.WriteFile(filepath.Join(root, "Makefile"), []byte(makefile[1]), 0644); err != nil {
		t.Fatal(err)
	}
	// when
	for _, args := range [][]string{{"go", "mod", "tidy"}, {"go", "build", "./..."}, {"make", "test"}, {"make", "test-unit"}, {"make", "test-integration"}} {
		command := exec.Command(args[0], args[1:]...)
		command.Dir = root
		output, err := command.CombinedOutput()
		// then
		if err != nil {
			t.Fatalf("%v: %v\n%s", args, err, output)
		}
	}
}
