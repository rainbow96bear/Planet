package model

import "time"

type User struct {
	ID         uint   `gorm:"primaryKey"`
	Username   string `gorm:"uniqueIndex;not null"`
	Nickname   string `gorm:"not null"`
	Password   string
	Provider   string `gorm:"default:'local';not null"`
	ProviderID string `gorm:"index"` // OAuth 사용자의 고유 ID
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
