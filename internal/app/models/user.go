package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	gorm.Model
	UserName        string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email           string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash    string     `gorm:"type:varchar(255);not null" json:"-"`
	FirstName       string     `gorm:"type:varchar(100)" json:"first_name"`
	LastName        string     `gorm:"type:varchar(100)" json:"last_name"`
	AvatarURL       string     `gorm:"type:text" json:"avatar_url,omitempty"`
	DefaultCurrency string     `gorm:"type:varchar(10);not null;default:'USD'" json:"default_currency"`
	EmailVerified   bool       `gorm:"default:false" json:"email_verified"`
	LastLoginAt     *time.Time `gorm:"type:timestamp" json:"last_login_at,omitempty"`
	IsActive        bool       `gorm:"not null;default:true" json:"is_active"`
	IsSuperAdmin    bool       `gorm:"not null;default:false" json:"is_superadmin"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}
