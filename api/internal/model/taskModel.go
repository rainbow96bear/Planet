package model

import "time"

type Task struct {
	ID          uint   `gorm:"primaryKey"`
	UserID      uint   `gorm:"not null;index"` // User 외래키
	User        User   `gorm:"foreignKey:UserID"`
	Title       string `gorm:"not null"`
	Description string
	Date        time.Time `gorm:"not null;index"` // 캘린더에서 표시할 날짜
	IsCompleted bool      `gorm:"default:false"`
	IsPublic    bool      `gorm:"default:true"` // SNS 공개 여부
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
