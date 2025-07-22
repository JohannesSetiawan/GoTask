package database

import (
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Task struct{
	ID          uint `json:"id" gorm:"primary_key"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
	UserID      uint           `json:"user_id"`
	User        User           `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}