package dto

type FollowRequest struct {
	FollowerID  uint
	FollowingID uint
}

type FollowResponse struct {
	IsFollowing bool `json:"is_following"`
}

type UnfollowRequest struct {
	FollowerID  uint
	FollowingID uint
}

type UnfollowResponse struct {
	IsFollowing bool `json:"is_following"`
}

type IsFollowingRequest struct {
	FollowerID  uint
	FollowingID uint
}

type IsFollowingResponse struct {
	IsFollowing bool `json:"is_following"`
}
