package repository

import (
	"errors"
	"planet/internal/dto"
	"planet/internal/model"
	"time"

	"gorm.io/gorm"
)

type ActivityWithTask struct {
	ID        string
	Type      model.ActivityType
	TaskTitle *string
	CreatedAt time.Time
}

type ActivityRepository interface {
	Create(tx *gorm.DB, a *model.Activity) error
	FindByUserID(userID string, limit int) ([]model.Activity, error)
	FindLatestByUserID(userID string) (*model.Activity, error)
	FindFeed(userID string, limit int) ([]*dto.GetFeedResponse, error)
	FindExploreFeed(limit int) ([]*dto.GetFeedResponse, error)
}

type activityRepository struct {
	db *gorm.DB
}

func NewActivityRepository(db *gorm.DB) ActivityRepository {
	return &activityRepository{db: db}
}

func (r *activityRepository) getDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

func (r *activityRepository) Create(tx *gorm.DB, a *model.Activity) error {
	return r.getDB(tx).Create(a).Error
}

func (r *activityRepository) FindByUserID(userID string, limit int) ([]model.Activity, error) {
	var activities []model.Activity
	err := r.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&activities).Error
	return activities, err
}

func (r *activityRepository) FindLatestByUserID(userID string) (*model.Activity, error) {
	var activity model.Activity

	err := r.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(1).
		First(&activity).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &activity, nil
}

func (r *activityRepository) FindFeed(userID string, limit int) ([]*dto.GetFeedResponse, error) {
	var result []*dto.GetFeedResponse

	err := r.db.
		Table("activities a").
		Select(`
            a.id           AS activity_id,
            a.user_id      AS actor_id,
            u.nickname     AS actor_nickname,
            a.type,
            a.target_id,
            a.target_type,
            t.title        AS task_title,
            a.created_at
        `).
		Joins("JOIN users u ON u.id = a.user_id").
		Joins("LEFT JOIN tasks t ON t.id = a.target_id AND a.target_type = 'task'").
		Where("a.user_id IN (?)",
			r.db.Table("follows").
				Select("following_id").
				Where("follower_id = ?", userID),
		).
		Order("a.created_at DESC").
		Limit(limit).
		Scan(&result).Error

	return result, err
}

func (r *activityRepository) FindExploreFeed(limit int) ([]*dto.GetFeedResponse, error) {
	var result []*dto.GetFeedResponse

	err := r.db.
		Table("activities a").
		Select(`
            a.id           AS activity_id,
            a.user_id      AS actor_id,
            u.nickname     AS actor_nickname,
            a.type,
            a.target_id,
            a.target_type,
            t.title        AS task_title,
            a.created_at
        `).
		Joins("JOIN users u ON u.id = a.user_id").
		Joins("LEFT JOIN tasks t ON t.id = a.target_id AND a.target_type = 'task'").
		Where("t.is_public = ? OR a.target_type = ?", true, model.TargetTypeUser).
		Order("a.created_at DESC").
		Limit(limit).
		Scan(&result).Error

	return result, err
}
