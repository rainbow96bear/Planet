package repository

import (
	"planet/internal/model"

	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(tx *gorm.DB, n *model.Notification) error
	FindByReceiverID(receiverID string) ([]*model.Notification, error)
	CountUnread(receiverID string) (int64, error)
	MarkAllAsRead(receiverID string) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) getDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

func (r *notificationRepository) Create(tx *gorm.DB, n *model.Notification) error {
	return r.getDB(tx).Create(n).Error
}

func (r *notificationRepository) FindByReceiverID(receiverID string) ([]*model.Notification, error) {
	var notifications []*model.Notification
	err := r.db.
		Where("receiver_id = ?", receiverID).
		Order("created_at DESC").
		Limit(30).
		Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) CountUnread(receiverID string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Notification{}).
		Where("receiver_id = ? AND is_read = ?", receiverID, false).
		Count(&count).Error
	return count, err
}

func (r *notificationRepository) MarkAllAsRead(receiverID string) error {
	return r.db.Model(&model.Notification{}).
		Where("receiver_id = ? AND is_read = ?", receiverID, false).
		Update("is_read", true).Error
}
