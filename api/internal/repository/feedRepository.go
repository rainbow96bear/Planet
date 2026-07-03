package repository

import (
	"planet/internal/dto"
	"planet/internal/model"

	"gorm.io/gorm"
)

type FeedRepository interface {
	Create(tx *gorm.DB, feed *model.Feed) error
	DeleteByActorAndTask(tx *gorm.DB, actorID, taskID string, feedType model.FeedType) error
	DeleteByTaskID(tx *gorm.DB, taskID string) error
	FindFeed(userID string, limit int) ([]*dto.GetFeedResponse, error)
	FindExploreFeed(userID string, limit int) ([]*dto.GetFeedResponse, error)
}

type feedRepository struct {
	db *gorm.DB
}

func NewFeedRepository(db *gorm.DB) FeedRepository {
	return &feedRepository{db: db}
}

func (r *feedRepository) getDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

func (r *feedRepository) Create(tx *gorm.DB, feed *model.Feed) error {
	return r.getDB(tx).Create(feed).Error
}

func (r *feedRepository) DeleteByActorAndTask(tx *gorm.DB, actorID, taskID string, feedType model.FeedType) error {
	return r.getDB(tx).
		Where("actor_id = ? AND task_id = ? AND type = ?", actorID, taskID, feedType).
		Delete(&model.Feed{}).Error
}

func (r *feedRepository) DeleteByTaskID(tx *gorm.DB, taskID string) error {
	return r.getDB(tx).
		Where("task_id = ?", taskID).
		Delete(&model.Feed{}).Error
}

func (r *feedRepository) FindFeed(userID string, limit int) ([]*dto.GetFeedResponse, error) {
	var result []*dto.GetFeedResponse

	err := r.db.
		Table("feeds f").
		Select(`
			f.id                                                AS feed_id,
			f.actor_id,
			u.nickname                                          AS actor_nickname,
			f.type,
			f.task_id,
			t.title                                             AS task_title,
			COUNT(r.id) FILTER (WHERE r.type = 'like')          AS like_count,
			COUNT(r.id) FILTER (WHERE r.type = 'cheer')         AS cheer_count,
			BOOL_OR(r.type = 'like'  AND r.user_id = ?)         AS is_liked,
			BOOL_OR(r.type = 'cheer' AND r.user_id = ?)         AS is_cheered,
			f.created_at
		`, userID, userID).
		Joins("JOIN users u ON u.id = f.actor_id AND u.deleted_at IS NULL").
		Joins("JOIN tasks t ON t.id = f.task_id AND t.deleted_at IS NULL").
		Joins("LEFT JOIN reactions r ON r.task_id = f.task_id").
		Where("f.actor_id IN (?)",
			r.db.Table("follows").
				Select("following_id").
				Where("follower_id = ?", userID),
		).
		Where("(t.is_public = true OR f.actor_id = ?)", userID).
		Group("f.id, f.actor_id, u.nickname, f.type, f.task_id, t.title, f.created_at").
		Order("f.created_at DESC").
		Limit(limit).
		Scan(&result).Error

	return result, err
}

func (r *feedRepository) FindExploreFeed(userID string, limit int) ([]*dto.GetFeedResponse, error) {
	var result []*dto.GetFeedResponse

	err := r.db.
		Table("feeds f").
		Select(`
			f.id                                                AS feed_id,
			f.actor_id,
			u.nickname                                          AS actor_nickname,
			f.type,
			f.task_id,
			t.title                                             AS task_title,
			COUNT(r.id) FILTER (WHERE r.type = 'like')          AS like_count,
			COUNT(r.id) FILTER (WHERE r.type = 'cheer')         AS cheer_count,
			BOOL_OR(r.type = 'like'  AND r.user_id = ?)         AS is_liked,
			BOOL_OR(r.type = 'cheer' AND r.user_id = ?)         AS is_cheered,
			f.created_at
		`, userID, userID).
		Joins("JOIN users u ON u.id = f.actor_id AND u.deleted_at IS NULL").
		Joins("JOIN tasks t ON t.id = f.task_id AND t.deleted_at IS NULL").
		Joins("LEFT JOIN reactions r ON r.task_id = f.task_id").
		Where("t.is_public = ?", true).
		Group("f.id, f.actor_id, u.nickname, f.type, f.task_id, t.title, f.created_at").
		Order("f.created_at DESC").
		Limit(limit).
		Scan(&result).Error

	return result, err
}
