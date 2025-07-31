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
}
