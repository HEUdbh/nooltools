package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	// AppVersion is the local application version used for release comparison.
	AppVersion = "v1.2.0"

	githubRepoURL         = "https://github.com/HEUdbh/nooltools"
	githubTokenConfigName = "github_token.json"
	updateRequestTimeout  = 8 * time.Second
)

type UpdateCheckResult struct {
	HasUpdate      bool   `json:"has_update"`
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
	ReleaseName    string `json:"release_name"`
	ReleaseURL     string `json:"release_url"`
	PublishedAt    string `json:"published_at"`
	ReleaseNotes   string `json:"release_notes"`
	CheckedAt      string `json:"checked_at"`
	Message        string `json:"message"`
}

type githubTokenConfig struct {
	GitHubToken string `json:"github_token"`
}

type githubLatestRelease struct {
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	HTMLURL     string `json:"html_url"`
	PublishedAt string `json:"published_at"`
	Body        string `json:"body"`
}

func checkReleaseUpdate() (UpdateCheckResult, error) {
	result := UpdateCheckResult{
		CurrentVersion: AppVersion,
		CheckedAt:      time.Now().UTC().Format(time.RFC3339),
	}

	token, err := loadGitHubToken()
	if err != nil {
		return result, err
	}

	apiURL, err := buildLatestReleaseAPIURL(githubRepoURL)
	if err != nil {
		return result, err
	}

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return result, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "nooltools-update-checker")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: updateRequestTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return result, fmt.Errorf("failed to query github release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return result, fmt.Errorf("github api returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var release githubLatestRelease
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&release); err != nil {
		return result, fmt.Errorf("failed to decode github release response: %w", err)
	}

	latestVersion := strings.TrimSpace(release.TagName)
	if latestVersion == "" {
		latestVersion = strings.TrimSpace(release.Name)
	}
	if latestVersion == "" {
		return result, errors.New("latest release version is empty")
	}

	result.LatestVersion = latestVersion
	result.ReleaseName = strings.TrimSpace(release.Name)
	result.ReleaseURL = strings.TrimSpace(release.HTMLURL)
	result.PublishedAt = strings.TrimSpace(release.PublishedAt)
	result.ReleaseNotes = strings.TrimSpace(release.Body)

	cmp, err := compareVersions(AppVersion, latestVersion)
	if err != nil {
		return result, fmt.Errorf("failed to compare versions: %w", err)
	}

	if cmp < 0 {
		result.HasUpdate = true
		result.Message = "A new version is available."
	} else {
		result.HasUpdate = false
		result.Message = "You are using the latest version."
	}

	return result, nil
}

func getTokenConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(homeDir, ".nooltools", githubTokenConfigName), nil
}

func loadGitHubToken() (string, error) {
	configPath, err := getTokenConfigPath()
	if err != nil {
		return "", err
	}
	return loadGitHubTokenFromPath(configPath)
}

func loadGitHubTokenFromPath(configPath string) (string, error) {
	cfg, err := readOrCreateTokenConfig(configPath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(cfg.GitHubToken), nil
}

func readOrCreateTokenConfig(configPath string) (githubTokenConfig, error) {
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return githubTokenConfig{}, fmt.Errorf("failed to create config directory: %w", err)
	}

	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		cfg := githubTokenConfig{GitHubToken: ""}
		if err := writeTokenConfig(configPath, cfg); err != nil {
			return githubTokenConfig{}, err
		}
		return cfg, nil
	} else if err != nil {
		return githubTokenConfig{}, fmt.Errorf("failed to read token config: %w", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return githubTokenConfig{}, fmt.Errorf("failed to read token config file: %w", err)
	}

	var cfg githubTokenConfig
	if len(strings.TrimSpace(string(data))) == 0 {
		return githubTokenConfig{}, errors.New("token config file is empty")
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return githubTokenConfig{}, fmt.Errorf("failed to parse token config file: %w", err)
	}

	return cfg, nil
}

func writeTokenConfig(configPath string, cfg githubTokenConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode token config file: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(configPath, data, 0o600); err != nil {
		return fmt.Errorf("failed to write token config file: %w", err)
	}
	return nil
}

func buildLatestReleaseAPIURL(repoURL string) (string, error) {
	owner, repo, err := parseGitHubRepoURL(repoURL)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo), nil
}

func parseGitHubRepoURL(repoURL string) (string, string, error) {
	parsedURL, err := url.Parse(strings.TrimSpace(repoURL))
	if err != nil {
		return "", "", fmt.Errorf("invalid repository url: %w", err)
	}
	if !strings.EqualFold(parsedURL.Host, "github.com") {
		return "", "", errors.New("repository host must be github.com")
	}

	parts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(parts) < 2 {
		return "", "", errors.New("repository url must include owner and repo")
	}

	owner := strings.TrimSpace(parts[0])
	repo := strings.TrimSpace(parts[1])
	repo = strings.TrimSuffix(repo, ".git")
	if owner == "" || repo == "" {
		return "", "", errors.New("repository owner or repo is empty")
	}

	return owner, repo, nil
}

func compareVersions(current, latest string) (int, error) {
	currentParts, err := parseVersionParts(current)
	if err != nil {
		return 0, fmt.Errorf("invalid current version %q: %w", current, err)
	}
	latestParts, err := parseVersionParts(latest)
	if err != nil {
		return 0, fmt.Errorf("invalid latest version %q: %w", latest, err)
	}

	for i := 0; i < 3; i++ {
		if currentParts[i] < latestParts[i] {
			return -1, nil
		}
		if currentParts[i] > latestParts[i] {
			return 1, nil
		}
	}
	return 0, nil
}

func parseVersionParts(version string) ([3]int, error) {
	var parts [3]int
	normalized := normalizeVersion(version)
	if normalized == "" {
		return parts, errors.New("version is empty")
	}

	core := strings.SplitN(normalized, "-", 2)[0]
	core = strings.SplitN(core, "+", 2)[0]
	segments := strings.Split(core, ".")
	if len(segments) == 0 || len(segments) > 3 {
		return parts, errors.New("version must contain between 1 and 3 numeric segments")
	}

	for i := 0; i < len(segments); i++ {
		segment := strings.TrimSpace(segments[i])
		if segment == "" {
			return parts, errors.New("version segment is empty")
		}
		value, err := strconv.Atoi(segment)
		if err != nil {
			return parts, fmt.Errorf("version segment %q is not numeric", segment)
		}
		if value < 0 {
			return parts, errors.New("version segment cannot be negative")
		}
		parts[i] = value
	}

	return parts, nil
}

func normalizeVersion(version string) string {
	v := strings.TrimSpace(version)
	v = strings.TrimPrefix(v, "v")
	v = strings.TrimPrefix(v, "V")
	return strings.TrimSpace(v)
}
