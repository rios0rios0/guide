package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// redirectTransport rewrites all requests to point at the given test server.
type redirectTransport struct {
	serverURL string
}

func (rt *redirectTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "http"
	req.URL.Host = strings.TrimPrefix(rt.serverURL, "http://")
	return http.DefaultTransport.RoundTrip(req)
}

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

func TestFetchExternalSourcesWithArtifacts(t *testing.T) {
	// given
	agentContent := "You are a test agent.\n"
	commandContent := "You are a test command.\n"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/agents/foo.md"):
			w.Write([]byte(agentContent))
		case strings.HasSuffix(r.URL.Path, "/commands/bar.md"):
			w.Write([]byte(commandContent))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	originalClient := httpClient
	httpClient = &http.Client{Transport: &redirectTransport{serverURL: server.URL}}
	defer func() { httpClient = originalClient }()

	outputDir := t.TempDir()
	configDir := t.TempDir()
	configPath := filepath.Join(configDir, "config.yaml")
	configContent := fmt.Sprintf(`sources:
  - repo: 'test/repo'
    branch: 'main'
    agents:
      - source: 'agents/foo.md'
        target: 'foo-agent.md'
    commands:
      - source: 'commands/bar.md'
        target: 'bar-command.md'
`)
	os.WriteFile(configPath, []byte(configContent), 0644)

	// when
	errorCount := fetchExternalSources(configPath, outputDir)

	// then
	if errorCount != 0 {
		t.Errorf("expected 0 errors, got %d", errorCount)
	}

	agentPath := filepath.Join(outputDir, ".ai", "claude", "agents", "foo-agent.md")
	data, err := os.ReadFile(agentPath)
	if err != nil {
		t.Fatalf("expected agent file at %s: %v", agentPath, err)
	}
	if string(data) != agentContent {
		t.Errorf("agent content\n  got:  %q\n  want: %q", string(data), agentContent)
	}

	cmdPath := filepath.Join(outputDir, ".ai", "claude", "commands", "bar-command.md")
	data, err = os.ReadFile(cmdPath)
	if err != nil {
		t.Fatalf("expected command file at %s: %v", cmdPath, err)
	}
	if string(data) != commandContent {
		t.Errorf("command content\n  got:  %q\n  want: %q", string(data), commandContent)
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
