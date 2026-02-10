package model

import "time"

// Notebook 笔记本 — 对应 Vault 中的一个文件夹
type Notebook struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Path        string    `json:"path" gorm:"not null;uniqueIndex"` // Vault 内相对目录名
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
