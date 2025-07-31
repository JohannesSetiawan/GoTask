package database

import (
	"gorm.io/gorm"
)

type TaskStatus string

const (
	StatusCreated    TaskStatus = "Created"
	StatusInProgress TaskStatus = "In Progress"
	StatusDone       TaskStatus = "Done"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Password string `json:"password"`
}

type Task struct {
	gorm.Model
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	Deadline    *string    `json:"deadline,omitempty"`
	UserID      uint       `json:"user_id" gorm:"index"`
	User        User       `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Subtasks    []Subtask  `json:"subtasks,omitempty" gorm:"foreignKey:TaskID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Subtask struct {
	gorm.Model
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Order       int        `json:"order" gorm:"index:idx_subtask_task_order,priority:2"`
	TaskID      uint       `json:"task_id" gorm:"index;index:idx_subtask_task_order,priority:1"`
	Task        Task       `json:"task,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID      uint       `json:"user_id" gorm:"index"`
	User        User       `json:"user,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Deadline    *string    `json:"deadline,omitempty"`
	Status      TaskStatus `json:"status" gorm:"default:'Created'"`
}
