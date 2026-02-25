package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"nooltools/apps/database"
	"nooltools/apps/storage"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type StorageSettings struct {
	CurrentDataDir string `json:"current_data_dir"`
	DefaultDataDir string `json:"default_data_dir"`
	IsCustom       bool   `json:"is_custom"`
	StartupNotice  string `json:"startup_notice"`
}

type StorageMigrationResult struct {
	FromDir            string   `json:"from_dir"`
	ToDir              string   `json:"to_dir"`
	BackedUpConflicts  []string `json:"backed_up_conflicts"`
	RestartRecommended bool     `json:"restart_recommended"`
}

func (a *app) GetStorageSettings() (StorageSettings, error) {
	currentDataDir, err := storage.ResolveDataDir()
	if err != nil {
		return StorageSettings{}, err
	}
	defaultDataDir, err := storage.DefaultDataDir()
	if err != nil {
		return StorageSettings{}, err
	}

	return StorageSettings{
		CurrentDataDir: currentDataDir,
		DefaultDataDir: defaultDataDir,
		IsCustom:       !storage.PathsEqual(currentDataDir, defaultDataDir),
		StartupNotice:  a.storageStartupNotice,
	}, nil
}

func (a *app) SelectStorageParentDirectory() (string, error) {
	if a.ctx == nil {
		return "", errors.New("application context is unavailable")
	}

	selectedPath, err := wailsRuntime.OpenDirectoryDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title:                "选择新的存储目录",
		CanCreateDirectories: true,
	})
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(selectedPath) == "" {
		return "", nil
	}
	return storage.NormalizePath(selectedPath)
}

func (a *app) MigrateStorageDirectory(parentDir string) (StorageMigrationResult, error) {
	fromDir, err := storage.ResolveDataDir()
	if err != nil {
		return StorageMigrationResult{}, err
	}
	toDir, err := storage.BuildTargetDataDir(parentDir)
	if err != nil {
		return StorageMigrationResult{}, err
	}

	if storage.PathsEqual(fromDir, toDir) {
		return StorageMigrationResult{}, errors.New("target storage directory is the same as the current directory")
	}
	if storage.IsSubPath(toDir, fromDir) {
		return StorageMigrationResult{}, errors.New("target storage directory cannot be inside current storage directory")
	}

	if err := os.MkdirAll(fromDir, 0o755); err != nil {
		return StorageMigrationResult{}, fmt.Errorf("failed to ensure source data directory: %w", err)
	}
	if err := os.MkdirAll(toDir, 0o755); err != nil {
		return StorageMigrationResult{}, fmt.Errorf("failed to create target data directory: %w", err)
	}

	configBefore, err := storage.LoadConfig()
	if err != nil {
		return StorageMigrationResult{}, err
	}

	if a.database != nil {
		if err := a.database.Close(); err != nil {
			return StorageMigrationResult{}, fmt.Errorf("failed to close current database: %w", err)
		}
	}
	a.database = nil

	backedUpConflicts, err := copyDirectoryWithConflictBackup(fromDir, toDir)
	if err != nil {
		if reopenErr := a.reopenDatabaseWithoutConfigChange(); reopenErr != nil {
			return StorageMigrationResult{}, fmt.Errorf("migration failed: %v; additionally failed to reopen database: %v", err, reopenErr)
		}
		return StorageMigrationResult{}, err
	}

	if err := storage.SetCustomDataDir(toDir); err != nil {
		if reopenErr := a.reopenDatabaseWithoutConfigChange(); reopenErr != nil {
			return StorageMigrationResult{}, fmt.Errorf("failed to switch storage config: %v; additionally failed to reopen database: %v", err, reopenErr)
		}
		return StorageMigrationResult{}, err
	}

	newDB, err := database.NewDatabase()
	if err != nil {
		rollbackErr := storage.SaveConfig(configBefore)
		reopenErr := a.reopenDatabaseWithoutConfigChange()
		return StorageMigrationResult{}, fmt.Errorf(
			"failed to initialize database in new directory: %v; rollback config error: %v; reopen old database error: %v",
			err,
			rollbackErr,
			reopenErr,
		)
	}
	a.database = newDB

	return StorageMigrationResult{
		FromDir:            fromDir,
		ToDir:              toDir,
		BackedUpConflicts:  backedUpConflicts,
		RestartRecommended: true,
	}, nil
}

func (a *app) RestartApplication() error {
	executablePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	cmd := exec.Command(executablePath, os.Args[1:]...)
	cmd.Dir = filepath.Dir(executablePath)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start new process: %w", err)
	}

	if a.ctx != nil {
		wailsRuntime.Quit(a.ctx)
	}
	return nil
}

func (a *app) reopenDatabaseWithoutConfigChange() error {
	db, err := database.NewDatabase()
	if err != nil {
		a.database = nil
		return err
	}
	a.database = db
	return nil
}

func copyDirectoryWithConflictBackup(sourceDir, targetDir string) ([]string, error) {
	var backedUp []string

	err := filepath.WalkDir(sourceDir, func(sourcePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		relPath, err := filepath.Rel(sourceDir, sourcePath)
		if err != nil {
			return fmt.Errorf("failed to resolve relative path for %q: %w", sourcePath, err)
		}
		if relPath == "." {
			return nil
		}

		targetPath := filepath.Join(targetDir, relPath)
		sourceInfo, err := entry.Info()
		if err != nil {
			return fmt.Errorf("failed to inspect source path %q: %w", sourcePath, err)
		}

		if sourceInfo.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("symlink is not supported for migration: %s", sourcePath)
		}

		if entry.IsDir() {
			if backupPath, err := ensureNoConflictForDirectory(targetPath); err != nil {
				return err
			} else if backupPath != "" {
				backedUp = append(backedUp, backupPath)
			}
			return os.MkdirAll(targetPath, sourceInfo.Mode().Perm())
		}

		if backupPath, err := backupIfPathExists(targetPath); err != nil {
			return err
		} else if backupPath != "" {
			backedUp = append(backedUp, backupPath)
		}
		return copySingleFile(sourcePath, targetPath, sourceInfo.Mode().Perm())
	})
	if err != nil {
		return nil, fmt.Errorf("failed to copy directory: %w", err)
	}

	return backedUp, nil
}

func ensureNoConflictForDirectory(path string) (string, error) {
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to check target path %q: %w", path, err)
	}
	if info.IsDir() {
		return "", nil
	}
	return backupIfPathExists(path)
}

func backupIfPathExists(path string) (string, error) {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to check path %q: %w", path, err)
	}

	backupPath, err := buildConflictBackupPath(path)
	if err != nil {
		return "", err
	}
	if err := os.Rename(path, backupPath); err != nil {
		return "", fmt.Errorf("failed to backup target conflict %q: %w", path, err)
	}
	return backupPath, nil
}

func buildConflictBackupPath(path string) (string, error) {
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	timestamp := time.Now().UTC().Format("20060102_150405")

	for idx := 0; idx < 1000; idx++ {
		suffix := ".backup_" + timestamp
		if idx > 0 {
			suffix += "_" + strconv.Itoa(idx)
		}

		candidate := filepath.Join(dir, base+suffix)
		if _, err := os.Stat(candidate); errors.Is(err, os.ErrNotExist) {
			return candidate, nil
		} else if err != nil {
			return "", fmt.Errorf("failed to inspect backup candidate %q: %w", candidate, err)
		}
	}

	return "", fmt.Errorf("failed to create conflict backup path for %q", path)
}

func copySingleFile(sourcePath, targetPath string, filePerm os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		return fmt.Errorf("failed to create target directory for %q: %w", targetPath, err)
	}

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file %q: %w", sourcePath, err)
	}
	defer sourceFile.Close()

	targetFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, filePerm)
	if err != nil {
		return fmt.Errorf("failed to create target file %q: %w", targetPath, err)
	}

	if _, err := io.Copy(targetFile, sourceFile); err != nil {
		targetFile.Close()
		return fmt.Errorf("failed to copy file from %q to %q: %w", sourcePath, targetPath, err)
	}
	if err := targetFile.Close(); err != nil {
		return fmt.Errorf("failed to finalize target file %q: %w", targetPath, err)
	}

	return nil
}
