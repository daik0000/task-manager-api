package model

import (
	"time"
)

type Task struct {
	ID          uint
	Title       string `gorm:"not null"`
	Description string
	Status      TaskStatus `gorm:"not null;default:StatusTodo"`
	DueDate     *time.Time
	OwnerID     uint `gorm:"not null"`
	Owner       User `gorm:"foreignKey:OwnerID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:"index"` // soft delete
}

type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)
