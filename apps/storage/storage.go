package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	goRuntime "runtime"
	"strings"
)

const (
	DefaultDataDirName  = ".nooltools"
	ConfigFileName      = "storage_config.json"
	MigratedDataDirName = "nooltools_data"
	DatabaseFileName    = "nooltools.db"
	MarkdownDirName     = "markdown"
	UpdatesDirName      = "updates"
	GitHubTokenFileName = "github_token.json"
)

var userHomeDir = os.UserHomeDir

type Config struct {
	CustomDataDir string `json:"custom_data_dir"`
}

func DefaultDataDir() (string, error) {
	homeDir, err := userHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Clean(filepath.Join(homeDir, DefaultDataDirName)), nil
}

func ConfigPath() (string, error) {
	defaultDir, err := DefaultDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(defaultDir, ConfigFileName), nil
}

func LoadConfig() (Config, error) {
	configPath, err := ConfigPath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(configPath)
	if errors.Is(err, os.ErrNotExist) {
		return Config{}, nil
	}
	if err != nil {
		return Config{}, fmt.Errorf("failed to read storage config: %w", err)
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return Config{}, nil
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse storage config: %w", err)
	}

	if strings.TrimSpace(cfg.CustomDataDir) != "" {
		normalized, err := NormalizePath(cfg.CustomDataDir)
		if err != nil {
			return Config{}, fmt.Errorf("invalid custom data dir in storage config: %w", err)
		}
		cfg.CustomDataDir = normalized
	}

	return cfg, nil
}

func SaveConfig(cfg Config) error {
	defaultDir, err := DefaultDataDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(defaultDir, 0o755); err != nil {
		return fmt.Errorf("failed to create default data dir: %w", err)
	}

	if strings.TrimSpace(cfg.CustomDataDir) != "" {
		normalized, err := NormalizePath(cfg.CustomDataDir)
		if err != nil {
			return fmt.Errorf("invalid custom data dir: %w", err)
		}
		cfg.CustomDataDir = normalized
	}

	configPath := filepath.Join(defaultDir, ConfigFileName)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode storage config: %w", err)
	}
	data = append(data, '\n')

	if err := os.WriteFile(configPath, data, 0o600); err != nil {
		return fmt.Errorf("failed to write storage config: %w", err)
	}
	return nil
}

func SetCustomDataDir(dataDir string) error {
	normalized, err := NormalizePath(dataDir)
	if err != nil {
		return err
	}

	cfg, err := LoadConfig()
	if err != nil {
		return err
	}
	cfg.CustomDataDir = normalized
	return SaveConfig(cfg)
}

func ClearCustomDataDir() error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}
	cfg.CustomDataDir = ""
	return SaveConfig(cfg)
}

func ResolveDataDir() (string, error) {
	defaultDir, err := DefaultDataDir()
	if err != nil {
		return "", err
	}

	cfg, err := LoadConfig()
	if err != nil {
		return "", err
	}

	if strings.TrimSpace(cfg.CustomDataDir) == "" {
		return defaultDir, nil
	}

	return cfg.CustomDataDir, nil
}

func BuildTargetDataDir(parentDir string) (string, error) {
	parent, err := NormalizePath(parentDir)
	if err != nil {
		return "", err
	}
	return filepath.Clean(filepath.Join(parent, MigratedDataDirName)), nil
}

func NormalizePath(path string) (string, error) {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return "", errors.New("path is empty")
	}

	abs, err := filepath.Abs(trimmed)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}
	return filepath.Clean(abs), nil
}

func PathsEqual(left, right string) bool {
	if goRuntime.GOOS == "windows" {
		return strings.EqualFold(filepath.Clean(left), filepath.Clean(right))
	}
	return filepath.Clean(left) == filepath.Clean(right)
}

func IsSubPath(path, parent string) bool {
	path = filepath.Clean(path)
	parent = filepath.Clean(parent)

	rel, err := filepath.Rel(parent, path)
	if err != nil {
		return false
	}
	if rel == "." {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator))
}
