package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name    string
		current string
		latest  string
		want    int
		wantErr bool
	}{
		{name: "equal with prefix", current: "v1.2.0", latest: "1.2.0", want: 0},
		{name: "latest greater", current: "1.2.0", latest: "1.2.1", want: -1},
		{name: "current greater", current: "1.3.0", latest: "1.2.9", want: 1},
		{name: "semantic compare", current: "1.9.9", latest: "1.10.0", want: -1},
		{name: "prerelease ignored", current: "1.2.0", latest: "v1.2.0-beta", want: 0},
		{name: "invalid version", current: "1.2", latest: "abc", wantErr: true},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := compareVersions(tc.current, tc.latest)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("compareVersions returned error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("compareVersions() = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestParseGitHubRepoURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		owner   string
		repo    string
		wantErr bool
	}{
		{name: "normal", input: "https://github.com/HEUdbh/nooltools", owner: "HEUdbh", repo: "nooltools"},
		{name: "trailing slash", input: "https://github.com/HEUdbh/nooltools/", owner: "HEUdbh", repo: "nooltools"},
		{name: "git suffix", input: "https://github.com/HEUdbh/nooltools.git", owner: "HEUdbh", repo: "nooltools"},
		{name: "invalid host", input: "https://gitlab.com/HEUdbh/nooltools", wantErr: true},
		{name: "missing repo", input: "https://github.com/HEUdbh", wantErr: true},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			owner, repo, err := parseGitHubRepoURL(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("parseGitHubRepoURL returned error: %v", err)
			}
			if owner != tc.owner || repo != tc.repo {
				t.Fatalf("parseGitHubRepoURL() = (%q, %q), want (%q, %q)", owner, repo, tc.owner, tc.repo)
			}
		})
	}
}

func TestLoadGitHubTokenFromPath(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "github_token.json")

	token, err := loadGitHubTokenFromPath(configPath)
	if err != nil {
		t.Fatalf("loadGitHubTokenFromPath created file should not fail: %v", err)
	}
	if token != "" {
		t.Fatalf("expected empty token for first create, got %q", token)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read generated token config: %v", err)
	}
	if !strings.Contains(string(data), "\"github_token\"") {
		t.Fatalf("generated token config does not contain github_token field")
	}

	err = os.WriteFile(configPath, []byte("{\"github_token\":\"  test_token  \"}\n"), 0o600)
	if err != nil {
		t.Fatalf("failed to write token config: %v", err)
	}

	token, err = loadGitHubTokenFromPath(configPath)
	if err != nil {
		t.Fatalf("loadGitHubTokenFromPath existing file should not fail: %v", err)
	}
	if token != "test_token" {
		t.Fatalf("expected trimmed token test_token, got %q", token)
	}
}
