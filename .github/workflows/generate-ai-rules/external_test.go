package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadExternalConfig(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectErr   bool
		expectCount int
	}{
		{
			name: "valid config with one source",
			content: `sources:
  - repo: 'example/repo'
    branch: 'main'
    agents:
      - source: 'plugins/foo/agents/bar.md'
        target: 'bar.md'
    commands:
      - source: 'plugins/foo/commands/baz.md'
        target: 'baz.md'
`,
			expectCount: 1,
		},
		{
			name:        "empty config",
			content:     "sources: []\n",
			expectCount: 0,
		},
		{
			name:      "invalid yaml",
			content:   "{{invalid",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.yaml")
			os.WriteFile(configPath, []byte(tt.content), 0644)

			// when
			config, err := loadExternalConfig(configPath)

			// then
			if tt.expectErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("loadExternalConfig() error: %v", err)
			}
			if len(config.Sources) != tt.expectCount {
				t.Errorf("expected %d sources, got %d", tt.expectCount, len(config.Sources))
			}
		})
	}
}

func TestLoadExternalConfigFileNotFound(t *testing.T) {
	// given
	path := "/nonexistent/config.yaml"

	// when
	_, err := loadExternalConfig(path)

	// then
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestFetchArtifact(t *testing.T) {
	// given
	agentContent := "---\nname: test-agent\n---\n\nYou are a test agent.\n"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/agents/test.md") {
			w.Write([]byte(agentContent))
		} else {
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	// override httpClient to point at test server
	originalClient := httpClient
	httpClient = server.Client()
	defer func() { httpClient = originalClient }()

	tmpDir := t.TempDir()
	artifact := ExternalArtifact{
		Source: "agents/test.md",
		Target: "test-agent.md",
	}

	// construct a URL that the test server will serve
	// We need to use the server URL as the "repo" base
	url := server.URL + "/" + artifact.Source

	// when - directly test httpGet and file writing
	body, err := httpGet(url)
	if err != nil {
		t.Fatalf("httpGet() error: %v", err)
	}

	dir := filepath.Join(tmpDir, ".ai", "claude", "agents")
	os.MkdirAll(dir, 0755)
	path := filepath.Join(dir, artifact.Target)
	os.WriteFile(path, body, 0644)

	// then
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading output file: %v", err)
	}
	if string(data) != agentContent {
		t.Errorf("file content\n  got:  %q\n  want: %q", string(data), agentContent)
	}
}

func TestHttpGetRetry(t *testing.T) {
	// given
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			http.Error(w, "temporary error", http.StatusInternalServerError)
			return
		}
		w.Write([]byte("success"))
	}))
	defer server.Close()

	originalClient := httpClient
	httpClient = server.Client()
	defer func() { httpClient = originalClient }()

	// when
	body, err := httpGet(server.URL + "/test")

	// then
	if err != nil {
		t.Fatalf("httpGet() should succeed on retry, got error: %v", err)
	}
	if string(body) != "success" {
		t.Errorf("expected 'success', got %q", string(body))
	}
	if callCount != 2 {
		t.Errorf("expected 2 calls (1 fail + 1 success), got %d", callCount)
	}
}

func TestHttpGetAllFail(t *testing.T) {
	// given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "always fail", http.StatusInternalServerError)
	}))
	defer server.Close()

	originalClient := httpClient
	httpClient = server.Client()
	defer func() { httpClient = originalClient }()

	// when
	_, err := httpGet(server.URL + "/test")

	// then
	if err == nil {
		t.Fatal("expected error when all retries fail")
	}
}

func TestFetchExternalSourcesEmptyConfig(t *testing.T) {
	// given
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	os.WriteFile(configPath, []byte("sources: []\n"), 0644)
	outputDir := t.TempDir()

	// when
	errorCount := fetchExternalSources(configPath, outputDir)

	// then
	if errorCount != 0 {
		t.Errorf("expected 0 errors for empty config, got %d", errorCount)
	}
}
