package main

import (
	"nooltools/apps/database"
	"nooltools/apps/storage"
	"os"
	"path/filepath"
	"testing"
)

func TestCopyDirectoryWithConflictBackup(t *testing.T) {
	sourceDir := filepath.Join(t.TempDir(), "src")
	targetDir := filepath.Join(t.TempDir(), "dst")

	if err := os.MkdirAll(sourceDir, 0o755); err != nil {
		t.Fatalf("failed to create source dir: %v", err)
	}
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatalf("failed to create target dir: %v", err)
	}

	sourceFile := filepath.Join(sourceDir, "demo.txt")
	targetFile := filepath.Join(targetDir, "demo.txt")
	if err := os.WriteFile(sourceFile, []byte("source"), 0o644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}
	if err := os.WriteFile(targetFile, []byte("target-old"), 0o644); err != nil {
		t.Fatalf("failed to write target file: %v", err)
	}

	backups, err := copyDirectoryWithConflictBackup(sourceDir, targetDir)
	if err != nil {
		t.Fatalf("copyDirectoryWithConflictBackup() failed: %v", err)
	}
	if len(backups) != 1 {
		t.Fatalf("expected 1 backup conflict, got %d", len(backups))
	}

	data, err := os.ReadFile(targetFile)
	if err != nil {
		t.Fatalf("failed to read copied target file: %v", err)
	}
	if string(data) != "source" {
		t.Fatalf("expected copied file content 'source', got %q", string(data))
	}

	backupData, err := os.ReadFile(backups[0])
	if err != nil {
		t.Fatalf("failed to read backup file: %v", err)
	}
	if string(backupData) != "target-old" {
		t.Fatalf("expected backup file content 'target-old', got %q", string(backupData))
	}
}

func TestMigrateStorageDirectory_SwitchesToNewDirAndKeepsSource(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	t.Setenv("USERPROFILE", tempHome)

	defaultDir, err := storage.DefaultDataDir()
	if err != nil {
		t.Fatalf("DefaultDataDir() failed: %v", err)
	}

	appInstance := NewApp()
	db, err := database.NewDatabase()
	if err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}
	appInstance.database = db
	t.Cleanup(func() {
		if appInstance.database != nil {
			_ = appInstance.database.Close()
		}
	})

	markdownDir := filepath.Join(defaultDir, storage.MarkdownDirName)
	if err := os.MkdirAll(markdownDir, 0o755); err != nil {
		t.Fatalf("failed to create markdown dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(markdownDir, "notes.md"), []byte("# demo"), 0o644); err != nil {
		t.Fatalf("failed to write markdown file: %v", err)
	}

	newParentDir := filepath.Join(tempHome, "new-storage-parent")
	result, err := appInstance.MigrateStorageDirectory(newParentDir)
	if err != nil {
		t.Fatalf("MigrateStorageDirectory() failed: %v", err)
	}

	expectedTarget := filepath.Join(newParentDir, storage.MigratedDataDirName)
	if !storage.PathsEqual(result.ToDir, expectedTarget) {
		t.Fatalf("expected target %q, got %q", expectedTarget, result.ToDir)
	}

	if _, err := os.Stat(filepath.Join(defaultDir, storage.DatabaseFileName)); err != nil {
		t.Fatalf("expected source database to remain, stat failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(result.ToDir, storage.DatabaseFileName)); err != nil {
		t.Fatalf("expected migrated database in target dir, stat failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(result.ToDir, storage.MarkdownDirName, "notes.md")); err != nil {
		t.Fatalf("expected migrated markdown file in target dir, stat failed: %v", err)
	}

	currentDir, err := storage.ResolveDataDir()
	if err != nil {
		t.Fatalf("ResolveDataDir() failed after migration: %v", err)
	}
	if !storage.PathsEqual(currentDir, result.ToDir) {
		t.Fatalf("expected current data dir to switch to %q, got %q", result.ToDir, currentDir)
	}
}

func TestMigrateStorageDirectory_RejectsTargetInsideSource(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	t.Setenv("USERPROFILE", tempHome)

	defaultDir, err := storage.DefaultDataDir()
	if err != nil {
		t.Fatalf("DefaultDataDir() failed: %v", err)
	}
	appInstance := NewApp()
	db, err := database.NewDatabase()
	if err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}
	appInstance.database = db
	t.Cleanup(func() {
		if appInstance.database != nil {
			_ = appInstance.database.Close()
		}
	})

	parentInsideSource := filepath.Join(defaultDir, "sub")
	if _, err := appInstance.MigrateStorageDirectory(parentInsideSource); err == nil {
		t.Fatalf("expected error when target is inside source directory")
	}
}
