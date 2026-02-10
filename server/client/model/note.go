package model

import "time"

// Note 笔记 — 对应 Vault 中的一个 .md 文件
type Note struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	NotebookID uint      `json:"notebook_id" gorm:"index"`
	Title      string    `json:"title" gorm:"not null"`
	FilePath   string    `json:"file_path" gorm:"not null"` // .md 文件相对于 Vault 的路径
	Content    string    `json:"content" gorm:"-"`           // 从文件加载，不存数据库
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
