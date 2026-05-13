package model

type Notification struct {
	BaseModel
	ReceiverID string           `gorm:"type:char(36);not null;index"` // 받는 사람
	ActorID    string           `gorm:"type:char(36);not null"`       // 한 사람
	Type       NotificationType `gorm:"type:varchar(50);not null"`
	IsRead     bool             `gorm:"default:false"`
}

type NotificationType string

const (
	NotificationFollowed NotificationType = "followed"
)
