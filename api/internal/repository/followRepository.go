package repository

import (
	"planet/internal/model"

	"gorm.io/gorm"
)

type FollowRepository interface {
	Follow(tx *gorm.DB, f *model.Follow) error
	Unfollow(tx *gorm.DB, followerID, followingID uint) error
	IsFollowing(tx *gorm.DB, followerID, followingID uint) (bool, error)
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

func (r *followRepository) Unfollow(tx *gorm.DB, followerID, followingID uint) error {
	return tx.Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&model.Follow{}).Error
}

func (r *followRepository) IsFollowing(tx *gorm.DB, followerID, followingID uint) (bool, error) {
	var count int64
	err := tx.Model(&model.Follow{}).Where("follower_id = ? AND following_id = ?", followerID, followingID).Count(&count).Error
	return count > 0, err
}
