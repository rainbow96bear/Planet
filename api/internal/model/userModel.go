package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	BaseModel

	Username string `gorm:"uniqueIndex;not null"`
	Nickname string `gorm:"uniqueIndex;not null"`

	Password string

	Email string `gorm:"index"`

	Provider   string `gorm:"default:'local';not null"`
	ProviderID string `gorm:"index"`

	UserStatus string `gorm:"default:'active';not null"`

	BanReason string
	BannedAt  *time.Time

	LastLoginAt *time.Time

	ProfileImage string

	TermsVersion   string
	PrivacyVersion string

	TermsAgreedAt   *time.Time
	PrivacyAgreedAt *time.Time

	DeletedAt gorm.DeletedAt `gorm:"index;type:timestamptz"`
}

const (
	UserStatusActive    = "active"
	UserStatusSuspended = "suspended"
	UserStatusBanned    = "banned"
)
