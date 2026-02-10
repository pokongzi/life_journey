// Package vault 管理 Obsidian 风格的 Markdown 文件库
// 笔记本 = 文件夹，笔记 = .md 文件
// 默认 Vault 目录: ~/LifeJourney/
package vault

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Dir 是 Vault 根目录的绝对路径
var Dir string

// Init 初始化 Vault 目录，返回绝对路径
func Init() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("无法获取用户主目录: %v", err)
	}
	Dir = filepath.Join(home, "LifeJourney")
	if err := os.MkdirAll(Dir, 0755); err != nil {
		log.Fatalf("无法创建 Vault 目录: %v", err)
	}
	log.Printf("Vault 目录: %s", Dir)
	return Dir
}

// ---------- 文件名工具 ----------

// illegalChars 文件名中不允许的字符
var illegalChars = regexp.MustCompile(`[/\\:*?"<>|]`)

// SanitizeName 清理文件/目录名中的非法字符
func SanitizeName(name string) string {
	name = illegalChars.ReplaceAllString(name, "-")
	name = strings.TrimSpace(name)
	if name == "" {
		name = "untitled"
	}
	return name
}

// ---------- 路径工具 ----------

// FullPath 将 Vault 内相对路径转为绝对路径
func FullPath(relPath string) string {
	return filepath.Join(Dir, relPath)
}

// NoteRelPath 生成笔记文件的相对路径: {notebookDir}/{title}.md
func NoteRelPath(notebookDir, title string) string {
	return filepath.Join(notebookDir, SanitizeName(title)+".md")
}

// UniqueNoteRelPath 如果路径已存在则追加数字后缀，保证不冲突
func UniqueNoteRelPath(notebookDir, title string) string {
	base := SanitizeName(title)
	rel := filepath.Join(notebookDir, base+".md")
	if _, err := os.Stat(FullPath(rel)); os.IsNotExist(err) {
		return rel
	}
	for i := 1; i < 10000; i++ {
		rel = filepath.Join(notebookDir, fmt.Sprintf("%s_%d.md", base, i))
		if _, err := os.Stat(FullPath(rel)); os.IsNotExist(err) {
			return rel
		}
	}
	return rel
}

// ---------- 笔记文件操作 ----------

// ReadNoteContent 从 .md 文件读取笔记内容
func ReadNoteContent(relPath string) (string, error) {
	data, err := os.ReadFile(FullPath(relPath))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteNoteContent 将内容写入 .md 文件（自动创建父目录）
func WriteNoteContent(relPath, content string) error {
	full := FullPath(relPath)
	if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
		return err
	}
	return os.WriteFile(full, []byte(content), 0644)
}

// DeleteNoteFile 删除笔记 .md 文件
func DeleteNoteFile(relPath string) error {
	err := os.Remove(FullPath(relPath))
	if os.IsNotExist(err) {
		return nil // 文件已不存在，不算错误
	}
	return err
}

// RenameNoteFile 重命名/移动笔记文件，返回新的相对路径
func RenameNoteFile(oldRelPath, newNotebookDir, newTitle string) (string, error) {
	newRel := UniqueNoteRelPath(newNotebookDir, newTitle)
	newFull := FullPath(newRel)
	if err := os.MkdirAll(filepath.Dir(newFull), 0755); err != nil {
		return "", err
	}
	if err := os.Rename(FullPath(oldRelPath), newFull); err != nil {
		return "", err
	}
	return newRel, nil
}

// ---------- 笔记本目录操作 ----------

// CreateNotebookDir 创建笔记本文件夹
func CreateNotebookDir(dirName string) error {
	return os.MkdirAll(FullPath(dirName), 0755)
}

// RenameNotebookDir 重命名笔记本文件夹
func RenameNotebookDir(oldDir, newDir string) error {
	oldFull := FullPath(oldDir)
	newFull := FullPath(newDir)
	if _, err := os.Stat(oldFull); os.IsNotExist(err) {
		// 旧目录不存在，直接创建新目录
		return os.MkdirAll(newFull, 0755)
	}
	return os.Rename(oldFull, newFull)
}

// DeleteNotebookDir 删除笔记本文件夹及其内所有文件
func DeleteNotebookDir(dirName string) error {
	return os.RemoveAll(FullPath(dirName))
}
