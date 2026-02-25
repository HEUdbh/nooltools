package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"nooltools/apps/storage"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	// AppVersion is the local application version used for release comparison.
	AppVersion = "v1.4.0"

	githubRepoURL         = "https://github.com/HEUdbh/nooltools"
	githubTokenConfigName = "github_token.json"
	updateRequestTimeout  = 8 * time.Second

	releaseAssetPrimaryName  = "nooltools.exe"
	releaseAssetFallbackName = "noltools.exe"
)

type UpdateCheckResult struct {
	HasUpdate        bool   `json:"has_update"`
	CurrentVersion   string `json:"current_version"`
	LatestVersion    string `json:"latest_version"`
	ReleaseName      string `json:"release_name"`
	ReleaseURL       string `json:"release_url"`
	PublishedAt      string `json:"published_at"`
	ReleaseNotes     string `json:"release_notes"`
	CheckedAt        string `json:"checked_at"`
	Message          string `json:"message"`
	AssetName        string `json:"asset_name"`
	AssetSize        int64  `json:"asset_size"`
	CanAutoUpdate    bool   `json:"can_auto_update"`
	AutoUpdateReason string `json:"auto_update_reason"`
}

type githubTokenConfig struct {
	GitHubToken string `json:"github_token"`
}

type githubLatestRelease struct {
	TagName     string               `json:"tag_name"`
	Name        string               `json:"name"`
	HTMLURL     string               `json:"html_url"`
	PublishedAt string               `json:"published_at"`
	Body        string               `json:"body"`
	Assets      []githubReleaseAsset `json:"assets"`
}

type githubReleaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
	Digest             string `json:"digest"`
}

type updatePreparation struct {
	Result        UpdateCheckResult
	SelectedAsset *githubReleaseAsset
	SHA256Digest  string
}

func checkReleaseUpdate() (UpdateCheckResult, error) {
	token, err := loadGitHubToken()
	if err != nil {
		return UpdateCheckResult{}, err
	}

	prep, err := prepareLatestReleaseUpdateWithToken(token)
	if err != nil {
		return UpdateCheckResult{}, err
	}
	return prep.Result, nil
}

func getTokenConfigPath() (string, error) {
	dataDir, err := storage.ResolveDataDir()
	if err != nil {
		return "", fmt.Errorf("failed to resolve data directory: %w", err)
	}
	return filepath.Join(dataDir, githubTokenConfigName), nil
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

func prepareLatestReleaseUpdateWithToken(token string) (updatePreparation, error) {
	release, err := fetchLatestRelease(token)
	if err != nil {
		return updatePreparation{}, err
	}
	return buildUpdatePreparation(release)
}

func fetchLatestRelease(token string) (githubLatestRelease, error) {
	apiURL, err := buildLatestReleaseAPIURL(githubRepoURL)
	if err != nil {
		return githubLatestRelease{}, err
	}

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return githubLatestRelease{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "nooltools-update-checker")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: updateRequestTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return githubLatestRelease{}, fmt.Errorf("failed to query github release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return githubLatestRelease{}, fmt.Errorf("github api returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var release githubLatestRelease
	if err := json.NewDecoder(io.LimitReader(resp.Body, 2<<20)).Decode(&release); err != nil {
		return githubLatestRelease{}, fmt.Errorf("failed to decode github release response: %w", err)
	}
	return release, nil
}

func buildUpdatePreparation(release githubLatestRelease) (updatePreparation, error) {
	result := UpdateCheckResult{
		CurrentVersion: AppVersion,
		CheckedAt:      time.Now().UTC().Format(time.RFC3339),
	}

	latestVersion := strings.TrimSpace(release.TagName)
	if latestVersion == "" {
		latestVersion = strings.TrimSpace(release.Name)
	}
	if latestVersion == "" {
		return updatePreparation{}, errors.New("latest release version is empty")
	}

	result.LatestVersion = latestVersion
	result.ReleaseName = strings.TrimSpace(release.Name)
	result.ReleaseURL = strings.TrimSpace(release.HTMLURL)
	result.PublishedAt = strings.TrimSpace(release.PublishedAt)
	result.ReleaseNotes = strings.TrimSpace(release.Body)

	cmp, err := compareVersions(AppVersion, latestVersion)
	if err != nil {
		return updatePreparation{}, fmt.Errorf("failed to compare versions: %w", err)
	}
	if cmp >= 0 {
		result.HasUpdate = false
		result.Message = "You are using the latest version."
		return updatePreparation{Result: result}, nil
	}

	result.HasUpdate = true
	result.Message = "A new version is available."

	asset, ok := selectPreferredAsset(release.Assets)
	if !ok {
		result.CanAutoUpdate = false
		result.AutoUpdateReason = "No supported Windows release asset found."
		return updatePreparation{Result: result}, nil
	}
	result.AssetName = asset.Name
	result.AssetSize = asset.Size

	sha256Digest, err := parseSHA256Digest(asset.Digest)
	if err != nil {
		result.CanAutoUpdate = false
		result.AutoUpdateReason = "Missing valid SHA256 digest, automatic replacement is blocked."
		return updatePreparation{Result: result, SelectedAsset: &asset}, nil
	}

	result.CanAutoUpdate = true
	result.AutoUpdateReason = ""
	return updatePreparation{
		Result:        result,
		SelectedAsset: &asset,
		SHA256Digest:  sha256Digest,
	}, nil
}

func selectPreferredAsset(assets []githubReleaseAsset) (githubReleaseAsset, bool) {
	names := []string{releaseAssetPrimaryName, releaseAssetFallbackName}
	for _, expectedName := range names {
		for _, asset := range assets {
			if strings.EqualFold(strings.TrimSpace(asset.Name), expectedName) {
				return asset, true
			}
		}
	}
	return githubReleaseAsset{}, false
}

func parseSHA256Digest(digest string) (string, error) {
	trimmed := strings.TrimSpace(digest)
	if trimmed == "" {
		return "", errors.New("digest is empty")
	}

	prefix := "sha256:"
	if !strings.HasPrefix(strings.ToLower(trimmed), prefix) {
		return "", errors.New("digest must start with sha256:")
	}

	sum := strings.TrimSpace(trimmed[len(prefix):])
	if len(sum) != 64 {
		return "", errors.New("sha256 digest must be 64 hex chars")
	}
	if _, err := hex.DecodeString(sum); err != nil {
		return "", errors.New("sha256 digest contains non-hex chars")
	}
	return strings.ToLower(sum), nil
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
