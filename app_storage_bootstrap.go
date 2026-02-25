package main

import (
	"errors"
	"fmt"
	"log"
	"nooltools/apps/database"
	"nooltools/apps/storage"
	"os"
	"path/filepath"
	"strings"
)

func (a *app) initializeStorageAndDatabase() error {
	a.storageStartupNotice = ""

	currentDataDir, err := storage.ResolveDataDir()
	if err != nil {
		return fmt.Errorf("failed to resolve current data directory: %w", err)
	}
	defaultDataDir, err := storage.DefaultDataDir()
	if err != nil {
		return fmt.Errorf("failed to resolve default data directory: %w", err)
	}

	cfg, err := storage.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to read storage config: %w", err)
	}
	hasCustomDir := strings.TrimSpace(cfg.CustomDataDir) != ""

	if err := a.bootstrapDataDir(currentDataDir); err == nil {
		return nil
	} else if !hasCustomDir {
		return err
	} else {
		a.storageStartupNotice = fmt.Sprintf(
			"检测到自定义存储目录不可用，已自动回退到默认目录。失败原因：%v",
			err,
		)
		log.Println(a.storageStartupNotice)
	}

	if err := storage.ClearCustomDataDir(); err != nil {
		return fmt.Errorf("failed to clear custom storage path after fallback: %w", err)
	}

	if err := a.bootstrapDataDir(defaultDataDir); err != nil {
		return fmt.Errorf("failed to initialize fallback default data directory: %w", err)
	}

	log.Printf("存储目录已回退至默认目录: %s", defaultDataDir)
	return nil
}

func (a *app) bootstrapDataDir(dataDir string) error {
	normalizedDir, err := storage.NormalizePath(dataDir)
	if err != nil {
		return fmt.Errorf("invalid data directory: %w", err)
	}

	dbPath := filepath.Join(normalizedDir, storage.DatabaseFileName)
	if _, err := os.Stat(dbPath); errors.Is(err, os.ErrNotExist) {
		log.Printf("检测到数据库文件缺失，开始自动初始化: %s", dbPath)
	} else if err != nil {
		return fmt.Errorf("failed to inspect database file: %w", err)
	}

	if err := ensureStorageLayout(normalizedDir); err != nil {
		return err
	}

	db, err := database.NewDatabaseAtDir(normalizedDir)
	if err != nil {
		return err
	}

	ok, err := db.CheckDatabase()
	if err != nil {
		_ = db.Close()
		return fmt.Errorf("database check failed after init: %w", err)
	}
	if !ok {
		_ = db.Close()
		return errors.New("database check returned abnormal state after init")
	}

	if a.database != nil {
		_ = a.database.Close()
	}
	a.database = db
	log.Printf("数据库初始化成功，当前数据目录: %s", normalizedDir)
	return nil
}

func ensureStorageLayout(dataDir string) error {
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	markdownDir := filepath.Join(dataDir, storage.MarkdownDirName)
	if err := os.MkdirAll(markdownDir, 0o755); err != nil {
		return fmt.Errorf("failed to create markdown directory: %w", err)
	}

	updatesDir := filepath.Join(dataDir, storage.UpdatesDirName)
	if err := os.MkdirAll(updatesDir, 0o755); err != nil {
		return fmt.Errorf("failed to create updates directory: %w", err)
	}

	if err := ensureGitHubTokenConfig(dataDir); err != nil {
		return err
	}

	return nil
}

func ensureGitHubTokenConfig(dataDir string) error {
	tokenPath := filepath.Join(dataDir, githubTokenConfigName)
	if _, err := readOrCreateTokenConfig(tokenPath); err != nil {
		return fmt.Errorf("failed to initialize github token config: %w", err)
	}
	return nil
}
