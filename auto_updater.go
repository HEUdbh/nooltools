package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	goRuntime "runtime"
	"strconv"
	"strings"
	"time"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	updateProgressEventName = "update:progress"
)

type updateProgressEvent struct {
	Stage   string `json:"stage"`
	Percent int    `json:"percent"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

func (a *app) StartAutoUpdate() (map[string]interface{}, error) {
	if goRuntime.GOOS != "windows" {
		return nil, errors.New("automatic replacement is only supported on Windows")
	}

	a.updateMu.Lock()
	if a.isUpdating {
		a.updateMu.Unlock()
		return nil, errors.New("update task is already running")
	}
	a.isUpdating = true
	a.updateMu.Unlock()

	go a.runAutoUpdateTask()

	return map[string]interface{}{
		"started": true,
		"message": "auto update task started",
	}, nil
}

func (a *app) runAutoUpdateTask() {
	defer func() {
		a.updateMu.Lock()
		a.isUpdating = false
		a.updateMu.Unlock()
	}()

	a.emitUpdateProgress("preparing", 0, "正在准备更新", "")

	token, err := loadGitHubToken()
	if err != nil {
		a.emitUpdateError(fmt.Errorf("failed to load github token config: %w", err))
		return
	}

	prep, err := prepareLatestReleaseUpdateWithToken(token)
	if err != nil {
		a.emitUpdateError(err)
		return
	}
	if !prep.Result.HasUpdate {
		a.emitUpdateError(errors.New("当前已是最新版本"))
		return
	}
	if !prep.Result.CanAutoUpdate {
		reason := prep.Result.AutoUpdateReason
		if reason == "" {
			reason = "自动更新条件不满足，请手动下载更新"
		}
		a.emitUpdateError(errors.New(reason))
		return
	}
	if prep.SelectedAsset == nil || prep.SelectedAsset.BrowserDownloadURL == "" || prep.SHA256Digest == "" {
		a.emitUpdateError(errors.New("release asset is incomplete, cannot run auto update"))
		return
	}

	updatesDir, err := getUpdatesDirectory()
	if err != nil {
		a.emitUpdateError(err)
		return
	}

	pendingPath := filepath.Join(updatesDir, fmt.Sprintf("pending_%s.exe", sanitizeVersionForFilename(prep.Result.LatestVersion)))
	if err := a.downloadUpdateAsset(token, prep.SelectedAsset.BrowserDownloadURL, pendingPath, prep.SelectedAsset.Size); err != nil {
		a.emitUpdateError(err)
		return
	}

	a.emitUpdateProgress("verifying", 90, "正在校验文件完整性", "")
	if err := verifyFileSHA256(pendingPath, prep.SHA256Digest); err != nil {
		a.emitUpdateError(err)
		return
	}

	targetPath, err := os.Executable()
	if err != nil {
		a.emitUpdateError(fmt.Errorf("failed to get executable path: %w", err))
		return
	}
	backupPath := targetPath + ".bak"
	scriptPath := filepath.Join(updatesDir, "apply_update.ps1")

	a.emitUpdateProgress("replacing", 96, "正在替换本地程序", "")
	if err := writeApplyUpdateScript(scriptPath); err != nil {
		a.emitUpdateError(err)
		return
	}
	if err := startApplyUpdateScript(scriptPath, os.Getpid(), targetPath, pendingPath, backupPath); err != nil {
		a.emitUpdateError(err)
		return
	}

	a.emitUpdateProgress("restarting", 100, "更新成功，正在重启应用", "")
	time.Sleep(250 * time.Millisecond)
	if a.ctx != nil {
		wailsRuntime.Quit(a.ctx)
	}
}

func (a *app) downloadUpdateAsset(token, downloadURL, outputPath string, expectedSize int64) error {
	a.emitUpdateProgress("downloading", 5, "开始下载更新包", "")

	req, err := http.NewRequest(http.MethodGet, downloadURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create download request: %w", err)
	}
	req.Header.Set("User-Agent", "nooltools-auto-updater")
	req.Header.Set("Accept", "application/octet-stream")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: 0}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download release asset: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("asset download failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	totalBytes := resp.ContentLength
	if totalBytes <= 0 {
		totalBytes = expectedSize
	}

	tempPath := outputPath + ".download"
	if err := os.RemoveAll(tempPath); err != nil {
		return fmt.Errorf("failed to cleanup temp download file: %w", err)
	}
	file, err := os.OpenFile(tempPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o755)
	if err != nil {
		return fmt.Errorf("failed to create temp download file: %w", err)
	}

	var downloaded int64
	buffer := make([]byte, 64*1024)
	for {
		n, readErr := resp.Body.Read(buffer)
		if n > 0 {
			if _, err := file.Write(buffer[:n]); err != nil {
				file.Close()
				return fmt.Errorf("failed to write update package: %w", err)
			}
			downloaded += int64(n)
			percent := calculateDownloadPercent(downloaded, totalBytes)
			a.emitUpdateProgress("downloading", percent, "正在下载更新包", fmt.Sprintf("%d/%d bytes", downloaded, totalBytes))
		}

		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				break
			}
			file.Close()
			return fmt.Errorf("failed during update download: %w", readErr)
		}
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("failed to finalize downloaded file: %w", err)
	}
	if err := os.RemoveAll(outputPath); err != nil {
		return fmt.Errorf("failed to cleanup previous pending file: %w", err)
	}
	if err := os.Rename(tempPath, outputPath); err != nil {
		return fmt.Errorf("failed to finalize pending update file: %w", err)
	}

	a.emitUpdateProgress("downloading", 85, "下载完成", fmt.Sprintf("已下载 %d bytes", downloaded))
	return nil
}

func calculateDownloadPercent(downloaded, total int64) int {
	if total <= 0 {
		return 50
	}
	percent := 5 + int(float64(downloaded)/float64(total)*80)
	if percent > 85 {
		return 85
	}
	if percent < 5 {
		return 5
	}
	return percent
}

func verifyFileSHA256(filePath, expectedDigest string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open downloaded file for hash check: %w", err)
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return fmt.Errorf("failed to hash downloaded file: %w", err)
	}

	actualDigest := hex.EncodeToString(hasher.Sum(nil))
	if !strings.EqualFold(actualDigest, expectedDigest) {
		return fmt.Errorf("sha256 mismatch, expected %s but got %s", expectedDigest, actualDigest)
	}
	return nil
}

func getUpdatesDirectory() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	updatesDir := filepath.Join(homeDir, ".nooltools", "updates")
	if err := os.MkdirAll(updatesDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create updates directory: %w", err)
	}
	return updatesDir, nil
}

func sanitizeVersionForFilename(version string) string {
	v := strings.TrimSpace(version)
	if v == "" {
		return "unknown"
	}

	var builder strings.Builder
	for _, ch := range v {
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') {
			builder.WriteRune(ch)
		} else {
			builder.WriteByte('_')
		}
	}
	result := strings.Trim(builder.String(), "_")
	if result == "" {
		return "unknown"
	}
	return result
}

func writeApplyUpdateScript(scriptPath string) error {
	content := `param(
  [int]$OldPid,
  [string]$TargetPath,
  [string]$NewFilePath,
  [string]$BackupPath
)

$ErrorActionPreference = "Stop"

$deadline = (Get-Date).AddSeconds(60)
while ((Get-Process -Id $OldPid -ErrorAction SilentlyContinue) -and (Get-Date) -lt $deadline) {
  Start-Sleep -Milliseconds 300
}

if (Get-Process -Id $OldPid -ErrorAction SilentlyContinue) {
  throw "Timed out waiting for target process exit."
}

if (Test-Path $BackupPath) {
  Remove-Item -Path $BackupPath -Force
}

$backupCreated = $false
if (Test-Path $TargetPath) {
  Copy-Item -Path $TargetPath -Destination $BackupPath -Force
  $backupCreated = $true
}

try {
  Copy-Item -Path $NewFilePath -Destination $TargetPath -Force
  Remove-Item -Path $NewFilePath -Force -ErrorAction SilentlyContinue
}
catch {
  if ($backupCreated -and (Test-Path $BackupPath)) {
    Copy-Item -Path $BackupPath -Destination $TargetPath -Force
  }
  throw
}

Start-Sleep -Milliseconds 300
Start-Process -FilePath $TargetPath
`

	if err := os.WriteFile(scriptPath, []byte(content), 0o700); err != nil {
		return fmt.Errorf("failed to write updater script: %w", err)
	}
	return nil
}

func startApplyUpdateScript(scriptPath string, oldPID int, targetPath, newFilePath, backupPath string) error {
	cmd := exec.Command(
		"powershell",
		"-NoProfile",
		"-ExecutionPolicy", "Bypass",
		"-File", scriptPath,
		"-OldPid", strconv.Itoa(oldPID),
		"-TargetPath", targetPath,
		"-NewFilePath", newFilePath,
		"-BackupPath", backupPath,
	)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start updater script: %w", err)
	}
	return nil
}

func (a *app) emitUpdateProgress(stage string, percent int, message, detail string) {
	if a.ctx == nil {
		return
	}
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	wailsRuntime.EventsEmit(a.ctx, updateProgressEventName, updateProgressEvent{
		Stage:   stage,
		Percent: percent,
		Message: message,
		Detail:  detail,
	})
}

func (a *app) emitUpdateError(err error) {
	if err == nil {
		return
	}
	a.emitUpdateProgress("error", 0, "自动更新失败", err.Error())
}
