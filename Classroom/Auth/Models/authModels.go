package Models

import (
	"gorm.io/gorm"
	"time"
)

// UserRole represents the role of a User in the system.
type UserRole int

const (
	Student UserRole = iota // defaults to 0
	Teacher
	Admin
)

// User struct modified to include a role
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"`
	Password string
	Role     *UserRole `json:"role,omitempty"`
}

type RefreshTokenModel struct {
	gorm.Model
	Token     string    `gorm:"index;not null"`
	Username  string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
}
