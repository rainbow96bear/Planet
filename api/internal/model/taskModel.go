package model

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	BaseModel
	UserID      string         `gorm:"type:uuid;not null;index"`
	User        User           `gorm:"foreignKey:UserID"`
	Title       string         `gorm:"not null;size:255"`
	Description string         `gorm:"size:1000"`
	Date        time.Time      `gorm:"not null;index"`
	IsCompleted bool           `gorm:"not null;default:false"`
	IsPublic    bool           `gorm:"not null;default:false"`
	Reactions   []Reaction     `gorm:"foreignKey:TaskID"`
	Feeds       []Feed         `gorm:"foreignKey:TaskID"`
	DeletedAt   gorm.DeletedAt `gorm:"index;type:timestamptz"`
}
