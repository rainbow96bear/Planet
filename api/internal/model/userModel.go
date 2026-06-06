package model

import "gorm.io/gorm"

type User struct {
	BaseModel
	Username   string `gorm:"uniqueIndex;not null"`
	Nickname   string `gorm:"not null"`
	Password   string
	Provider   string         `gorm:"default:'local';not null"`
	ProviderID string         `gorm:"index"`
	DeletedAt  gorm.DeletedAt `gorm:"index;type:timestamptz"`
}
