package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time `gorm:"type:timestamptz"`
	UpdatedAt time.Time `gorm:"type:timestamptz"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.NewString()
	}
	return nil
}
