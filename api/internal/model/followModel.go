package model

type Follow struct {
	BaseModel
	FollowerID  string `gorm:"type:char(36);uniqueIndex:idx_follower_following"`
	FollowingID string `gorm:"type:char(36);uniqueIndex:idx_follower_following"`
}
