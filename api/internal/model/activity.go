package model

type Activity struct {
	BaseModel
	UserID     string       `gorm:"type:char(36);not null;index"`
	Type       ActivityType `gorm:"type:varchar(50);not null"`
	TargetID   string       `gorm:"type:char(36);index"`
	TargetType TargetType   `gorm:"type:varchar(50)"`
}

// UpdatedAt 불필요하므로 오버라이드
func (Activity) UpdatedAt() {}

type TargetType string

const (
	TargetTypeTask TargetType = "task"
	TargetTypeUser TargetType = "user"
)

type ActivityType string

const (
	ActivityTaskCreated   ActivityType = "task.created"
	ActivityTaskCompleted ActivityType = "task.completed"
	ActivityFollowed      ActivityType = "user.followed"
)
