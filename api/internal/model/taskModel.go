package model

import (
	"time"
)

type Task struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"not null;index"`
	User        User      `gorm:"foreignKey:UserID"`
	Title       string    `gorm:"not null;size:255"`
	Description string    `gorm:"size:1000"`
	Date        time.Time `gorm:"not null;index"`
	IsCompleted bool      `gorm:"not null;default:false"`
	IsPublic    bool      `gorm:"not null;default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
