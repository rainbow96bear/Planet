package repository

import (
	"errors"
	"planet/internal/model"

	"gorm.io/gorm"
)

type FollowRepository interface {
	Follow(tx *gorm.DB, f *model.Follow) error
	Unfollow(tx *gorm.DB, followerID, followingID string) error
	IsFollowing(followerID, followingID string) (bool, error)
	CountFollowers(userID string) (int64, error)
	CountFollowing(userID string) (int64, error)
}

type followRepository struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) FollowRepository {
	return &followRepository{db: db}
}

func (r *followRepository) Follow(tx *gorm.DB, f *model.Follow) error {
	return tx.Create(f).Error
}

func (r *followRepository) Unfollow(tx *gorm.DB, followerID, followingID string) error {
	return tx.Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&model.Follow{}).Error
}

func (r *followRepository) IsFollowing(followerID, followingID string) (bool, error) {
	var follow model.Follow
	err := r.db.
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Limit(1).
		First(&follow).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return err == nil, err
}

func (r *followRepository) CountFollowers(userID string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Follow{}).Where("following_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *followRepository) CountFollowing(userID string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Follow{}).Where("follower_id = ?", userID).Count(&count).Error
	return count, err
}
