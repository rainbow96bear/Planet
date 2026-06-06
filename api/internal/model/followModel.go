package model

type Follow struct {
	BaseModel
	FollowerID  string `gorm:"type:uuid;uniqueIndex:idx_follower_following"`
	FollowingID string `gorm:"type:uuid;uniqueIndex:idx_follower_following"`
}
