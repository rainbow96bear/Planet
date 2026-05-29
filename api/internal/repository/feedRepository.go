package repository

import (
	"errors"
	"fmt"
	"planet/internal/dto"
	"planet/internal/model"

	"gorm.io/gorm"
)

type FeedRepository interface {
	Create(tx *gorm.DB, a *model.Feed) error
	FindByUserID(userID string, limit int) ([]model.Feed, error)
	FindLatestByUserID(userID string) (*model.Feed, error)
	FindFeed(userID string, limit int) ([]*dto.GetFeedResponse, error)
	FindExploreFeed(limit int) ([]*dto.GetFeedResponse, error)
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

func (r *feedRepository) Create(tx *gorm.DB, a *model.Feed) error {
	return r.getDB(tx).Create(a).Error
}

func (r *feedRepository) FindByUserID(userID string, limit int) ([]model.Feed, error) {
	var feeds []model.Feed
	err := r.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&feeds).Error
	return feeds, err
}

func (r *feedRepository) FindLatestByUserID(userID string) (*model.Feed, error) {
	var feed model.Feed

	err := r.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(1).
		First(&feed).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &feed, nil
}

func (r *feedRepository) FindFeed(userID string, limit int) ([]*dto.GetFeedResponse, error) {
	var result []*dto.GetFeedResponse

	err := r.db.
		Table("feeds a").
		Select(`
            a.id           AS feed_id,
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
	fmt.Printf("result : %+v\n", result)
	return result, err
}

func (r *feedRepository) FindExploreFeed(limit int) ([]*dto.GetFeedResponse, error) {
	var result []*dto.GetFeedResponse

	err := r.db.
		Table("feeds a").
		Select(`
            a.id           AS feed_id,
            a.user_id      AS actor_id,
            u.nickname     AS actor_nickname,
            a.type,
            a.target_id,
            a.target_type,
            t.title        AS task_title,
            a.created_at
        `).
		Joins("JOIN users u ON u.id = a.user_id").
		Joins("JOIN tasks t ON t.id = a.target_id").
		Where("t.is_public = ?", true).
		Order("a.created_at DESC").
		Limit(limit).
		Scan(&result).Error

	return result, err
}
