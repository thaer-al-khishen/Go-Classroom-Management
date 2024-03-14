package Models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"`
	Password string // This will store a hashed password
}

type RefreshTokenModel struct {
	gorm.Model
	Token     string    `gorm:"index;not null"`
	Username  string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
}
