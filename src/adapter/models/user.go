package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Username     string
	PasswordHash string
	LastLoginAt  *time.Time
}

func (*User) TableName() string {
	return "users"
}
