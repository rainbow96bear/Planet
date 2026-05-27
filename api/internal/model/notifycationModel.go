package model

type Notification struct {
	BaseModel
	ReceiverID string                 `gorm:"type:char(36);not null;index:idx_noti_unique,unique"`
	ActorID    string                 `gorm:"type:char(36);not null;index:idx_noti_unique,unique"`
	TargetID   string                 `gorm:"type:char(36);index:idx_noti_unique,unique"`
	TargetType NotificationTargetType `gorm:"type:varchar(50);index:idx_noti_unique,unique"`
	Type       NotificationType       `gorm:"type:varchar(50);not null;index:idx_noti_unique,unique"`
	IsRead     bool                   `gorm:"default:false"`
}

type NotificationType string

const (
	NotificationFollowed     NotificationType = "followed"
	NotificationTypeReaction NotificationType = "reaction"
	NotificationTypeComment  NotificationType = "comment"
)

type NotificationTargetType string

const (
	NotificationTargetTypeTask NotificationTargetType = "task"
	NotificationTargetTypeUser NotificationTargetType = "user"
)
