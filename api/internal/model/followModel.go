package model

import "time"

type Follow struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	FollowerID  uint      `json:"follower_id" gorm:"uniqueIndex:idx_follower_following"`
	FollowingID uint      `json:"following_id" gorm:"uniqueIndex:idx_follower_following"`
	CreatedAt   time.Time `json:"created_at"`
}
