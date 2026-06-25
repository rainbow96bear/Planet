package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	BaseModel

	Username string `gorm:"uniqueIndex;not null;type:varchar(20)"`
	Nickname string `gorm:"uniqueIndex;not null;type:varchar(20)"`

	Password string

	Email string `gorm:"index;type:varchar(254)"`

	Provider   string `gorm:"default:'local';not null;type:varchar(20)"` // "local", "kakao", "naver"
	ProviderID string `gorm:"index"`

	UserStatus string `gorm:"default:'active';not null;type:varchar(20)"` // "active", "suspended", "banned"

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
