package model

type Feed struct {
	BaseModel
	ActorID string   `gorm:"type:uuid;not null;index"`
	Actor   User     `gorm:"foreignKey:ActorID"`
	TaskID  string   `gorm:"type:char(36);not null;index"`
	Task    Task     `gorm:"foreignKey:TaskID"`
	Type    FeedType `gorm:"type:varchar(50);not null"`
}

func (Feed) UpdatedAt() {}

type FeedType string

const (
	TaskCreated   FeedType = "task.created"
	TaskCompleted FeedType = "task.completed"
)
