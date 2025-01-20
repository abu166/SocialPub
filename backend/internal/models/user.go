package models

import (
	"time"
)

type User struct {
	UserID    uint      `gorm:"primaryKey;column:user_id" json:"user_id"`
	UserName  string    `gorm:"column:user_name" json:"user_name"`
	UserEmail string    `gorm:"column:user_email;uniqueIndex" json:"user_email"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

type ResponseData struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginatedResponse struct {
	Status      string      `json:"status"`
	Message     string      `json:"message"`
	Data        interface{} `json:"data"`
	TotalItems  int64       `json:"total_items"`
	TotalPages  int         `json:"total_pages"`
	CurrentPage int         `json:"current_page"`
}
