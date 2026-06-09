package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	BaseModel
	Username string `gorm:"uniqueIndex;not null"`
	Nickname string `gorm:"not null"`
	Password string

	Provider   string         `gorm:"default:'local';not null"`
	ProviderID string         `gorm:"index"`
	DeletedAt  gorm.DeletedAt `gorm:"index;type:timestamptz"`

	ProfileImage string

	TermsVersion   string
	PrivacyVersion string

	TermsAgreedAt   *time.Time
	PrivacyAgreedAt *time.Time
}
