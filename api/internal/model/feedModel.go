package model

type Feed struct {
	BaseModel
	UserID     string     `gorm:"type:char(36);not null;index"`
	Type       FeedType   `gorm:"type:varchar(50);not null"`
	TargetID   string     `gorm:"type:char(36);index"`
	TargetType TargetType `gorm:"type:varchar(50)"`
}

// UpdatedAt 불필요하므로 오버라이드
func (Feed) UpdatedAt() {}

type TargetType string

const (
	TargetTypeTask TargetType = "task"
)

type FeedType string

const (
	TaskCreated   FeedType = "task.created"
	TaskCompleted FeedType = "task.completed"
)
