// Markdown编辑器处理器
package database

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// MarkdownFile Markdown文件信息
type MarkdownFile struct {
	Name    string `json:"name"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// GetMarkdownDir 获取markdown文件存储目录
func (d *Database) GetMarkdownDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("获取用户目录失败: %v", err)
	}

	markdownDir := filepath.Join(homeDir, ".nooltools", "markdown")
	if err := os.MkdirAll(markdownDir, 0755); err != nil {
		return "", fmt.Errorf("创建markdown目录失败: %v", err)
	}

	return markdownDir, nil
}

// GetMarkdownFiles 获取所有markdown文件列表
func (d *Database) GetMarkdownFiles() ([]MarkdownFile, error) {
	markdownDir, err := d.GetMarkdownDir()
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(markdownDir)
	if err != nil {
		return nil, fmt.Errorf("读取markdown目录失败: %v", err)
	}

	var markdownFiles []MarkdownFile
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.HasSuffix(file.Name(), ".md") {
			markdownFiles = append(markdownFiles, MarkdownFile{
				Name:  file.Name(),
				Title: strings.TrimSuffix(file.Name(), ".md"),
			})
		}
	}

	return markdownFiles, nil
}

// ReadMarkdownFile 读取markdown文件内容
func (d *Database) ReadMarkdownFile(filename string) (string, error) {
	markdownDir, err := d.GetMarkdownDir()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(markdownDir, filename)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %v", err)
	}

	return string(content), nil
}

// SaveMarkdownFile 保存markdown文件
func (d *Database) SaveMarkdownFile(title, content string) error {
	markdownDir, err := d.GetMarkdownDir()
	if err != nil {
		return err
	}

	// 确保文件名以.md结尾
	filename := title
	if !strings.HasSuffix(filename, ".md") {
		filename = filename + ".md"
	}

	filePath := filepath.Join(markdownDir, filename)

	// 写入文件
	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("保存文件失败: %v", err)
	}

	return nil
}

// DeleteMarkdownFile 删除markdown文件
func (d *Database) DeleteMarkdownFile(filename string) error {
	markdownDir, err := d.GetMarkdownDir()
	if err != nil {
		return err
	}

	// 确保文件名以.md结尾
	if !strings.HasSuffix(filename, ".md") {
		filename = filename + ".md"
	}

	filePath := filepath.Join(markdownDir, filename)

	err = os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("删除文件失败: %v", err)
	}

	return nil
}

// RenameMarkdownFile 重命名markdown文件
func (d *Database) RenameMarkdownFile(oldName, newName string) error {
	markdownDir, err := d.GetMarkdownDir()
	if err != nil {
		return err
	}

	// 确保文件名以.md结尾
	if !strings.HasSuffix(oldName, ".md") {
		oldName = oldName + ".md"
	}
	if !strings.HasSuffix(newName, ".md") {
		newName = newName + ".md"
	}

	oldPath := filepath.Join(markdownDir, oldName)
	newPath := filepath.Join(markdownDir, newName)

	err = os.Rename(oldPath, newPath)
	if err != nil {
		return fmt.Errorf("重命名文件失败: %v", err)
	}

	return nil
}
