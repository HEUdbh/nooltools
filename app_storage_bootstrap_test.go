package main

import (
	"nooltools/apps/storage"
	"os"
	"path/filepath"
	"testing"
)

func setupTempHome(t *testing.T) string {
	t.Helper()
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	t.Setenv("USERPROFILE", tempHome)
	return tempHome
}

func TestInitializeStorageAndDatabase_CustomDirMissingDBCreatesLayout(t *testing.T) {
	tempHome := setupTempHome(t)
	customDir := filepath.Join(tempHome, "custom-data")
	if err := storage.SetCustomDataDir(customDir); err != nil {
		t.Fatalf("failed to set custom data dir: %v", err)
	}

	appInstance := NewApp()
	if err := appInstance.initializeStorageAndDatabase(); err != nil {
		t.Fatalf("initializeStorageAndDatabase() failed: %v", err)
	}
	t.Cleanup(func() {
		if appInstance.database != nil {
			_ = appInstance.database.Close()
		}
	})

	if appInstance.storageStartupNotice != "" {
		t.Fatalf("expected empty startup notice, got %q", appInstance.storageStartupNotice)
	}

	mustExist := []string{
		filepath.Join(customDir, storage.DatabaseFileName),
		filepath.Join(customDir, storage.MarkdownDirName),
		filepath.Join(customDir, storage.UpdatesDirName),
		filepath.Join(customDir, storage.GitHubTokenFileName),
	}
	for _, path := range mustExist {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected path to exist after init %q: %v", path, err)
		}
	}
}

func TestInitializeStorageAndDatabase_CustomDirInvalidFallsBackToDefault(t *testing.T) {
	tempHome := setupTempHome(t)
	defaultDir, err := storage.DefaultDataDir()
	if err != nil {
		t.Fatalf("DefaultDataDir() failed: %v", err)
	}

	invalidCustomPath := filepath.Join(tempHome, "invalid-custom-dir")
	if err := os.WriteFile(invalidCustomPath, []byte("not a directory"), 0o644); err != nil {
		t.Fatalf("failed to create invalid path file: %v", err)
	}
	if err := storage.SetCustomDataDir(invalidCustomPath); err != nil {
		t.Fatalf("failed to set invalid custom data dir: %v", err)
	}

	appInstance := NewApp()
	if err := appInstance.initializeStorageAndDatabase(); err != nil {
		t.Fatalf("initializeStorageAndDatabase() should fallback to default dir, got error: %v", err)
	}
	t.Cleanup(func() {
		if appInstance.database != nil {
			_ = appInstance.database.Close()
		}
	})

	if appInstance.storageStartupNotice == "" {
		t.Fatalf("expected non-empty startup notice after fallback")
	}

	currentDataDir, err := storage.ResolveDataDir()
	if err != nil {
		t.Fatalf("ResolveDataDir() failed: %v", err)
	}
	if !storage.PathsEqual(currentDataDir, defaultDir) {
		t.Fatalf("expected fallback data dir %q, got %q", defaultDir, currentDataDir)
	}

	cfg, err := storage.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}
	if cfg.CustomDataDir != "" {
		t.Fatalf("expected custom data dir to be cleared after fallback, got %q", cfg.CustomDataDir)
	}

	if _, err := os.Stat(filepath.Join(defaultDir, storage.DatabaseFileName)); err != nil {
		t.Fatalf("expected database file in default dir after fallback: %v", err)
	}
}

func TestInitializeStorageAndDatabase_DefaultDirNoNotice(t *testing.T) {
	setupTempHome(t)

	appInstance := NewApp()
	if err := appInstance.initializeStorageAndDatabase(); err != nil {
		t.Fatalf("initializeStorageAndDatabase() failed: %v", err)
	}
	t.Cleanup(func() {
		if appInstance.database != nil {
			_ = appInstance.database.Close()
		}
	})

	if appInstance.storageStartupNotice != "" {
		t.Fatalf("expected no startup notice for default dir init, got %q", appInstance.storageStartupNotice)
	}
}
