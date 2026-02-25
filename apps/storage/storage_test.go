package storage

import (
	"path/filepath"
	"testing"
)

func TestResolveDataDir_DefaultAndCustom(t *testing.T) {
	tempHome := t.TempDir()
	originalUserHomeDir := userHomeDir
	userHomeDir = func() (string, error) {
		return tempHome, nil
	}
	t.Cleanup(func() {
		userHomeDir = originalUserHomeDir
	})

	defaultDir, err := DefaultDataDir()
	if err != nil {
		t.Fatalf("DefaultDataDir() failed: %v", err)
	}
	if defaultDir != filepath.Join(tempHome, DefaultDataDirName) {
		t.Fatalf("unexpected default dir: %q", defaultDir)
	}

	resolvedDir, err := ResolveDataDir()
	if err != nil {
		t.Fatalf("ResolveDataDir() failed: %v", err)
	}
	if !PathsEqual(resolvedDir, defaultDir) {
		t.Fatalf("expected default data dir, got %q", resolvedDir)
	}

	customDir := filepath.Join(tempHome, "custom-data")
	if err := SetCustomDataDir(customDir); err != nil {
		t.Fatalf("SetCustomDataDir() failed: %v", err)
	}

	resolvedDir, err = ResolveDataDir()
	if err != nil {
		t.Fatalf("ResolveDataDir() after custom failed: %v", err)
	}
	if !PathsEqual(resolvedDir, customDir) {
		t.Fatalf("expected custom data dir %q, got %q", customDir, resolvedDir)
	}
}

func TestBuildTargetDataDir(t *testing.T) {
	tempDir := t.TempDir()

	target, err := BuildTargetDataDir(tempDir)
	if err != nil {
		t.Fatalf("BuildTargetDataDir() failed: %v", err)
	}

	expected := filepath.Join(tempDir, MigratedDataDirName)
	if !PathsEqual(target, expected) {
		t.Fatalf("expected target %q, got %q", expected, target)
	}
}
