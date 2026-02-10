package model

import "time"

// Todo 待办事项
type Todo struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Title     string     `json:"title" gorm:"not null"`
	Content   string     `json:"content"`
	Status    string     `json:"status" gorm:"default:pending;index"`
	Priority  int        `json:"priority" gorm:"default:0"`
	DueDate   *time.Time `json:"due_date"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
