package repository

import (
	"planet/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type NotificationRepository interface {
	Upsert(tx *gorm.DB, n *model.Notification) error                                              // 추가
	DeleteByActorAndTask(tx *gorm.DB, actorID, taskID string, nType model.NotificationType) error // 추가
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

func (r *notificationRepository) Upsert(tx *gorm.DB, n *model.Notification) error {
	return r.getDB(tx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "receiver_id"},
				{Name: "actor_id"},
				{Name: "target_id"},
				{Name: "target_type"},
				{Name: "type"},
			},
			DoUpdates: clause.AssignmentColumns([]string{"updated_at", "is_read"}),
		}).
		Create(n).Error
}

func (r *notificationRepository) DeleteByActorAndTask(tx *gorm.DB, actorID, taskID string, nType model.NotificationType) error {
	return r.getDB(tx).
		Where("actor_id = ? AND target_id = ? AND type = ?", actorID, taskID, nType).
		Delete(&model.Notification{}).Error
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
